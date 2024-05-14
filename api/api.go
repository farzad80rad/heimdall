package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"heimdall/config"
	"heimdall/errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Api interface {
	Proxy(ctx *gin.Context) error
	Ping(url string) bool
}

type api struct {
	cb    *gobreaker.CircuitBreaker
	proxy *httputil.ReverseProxy
	host  string
}

func NewApi(host string, config config.CircuitBreakerConfig) (Api, error) {
	h, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(h)

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:     host,
		Interval: config.ExamineWindow,
		Timeout:  config.QuarantineDuration,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			maxTolerance := uint32(10)
			if config.FierierToleranceCount != 0 {
				maxTolerance = config.FierierToleranceCount
			}
			return counts.ConsecutiveFailures > maxTolerance
		},
	})

	return &api{
		cb:    cb,
		host:  host,
		proxy: proxy,
	}, nil
}

func (a *api) Ping(url string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := doHTTPGetRequest(ctx, a.host+url)
	return err == nil
}

func (a *api) Proxy(c *gin.Context) error {
	_, err := a.cb.Execute(func() (interface{}, error) {
		var cbError error
		a.proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
			cbError = err
		}
		a.proxy.ServeHTTP(c.Writer, c.Request)
		return nil, cbError
	})
	if err == gobreaker.ErrOpenState {
		return errors.HostIsDown
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "service is currently unable to respond. please try again"})
	}
	return nil
}

func doHTTPGetRequest(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{
		Timeout: time.Minute * 1,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is not OK: %v", resp.StatusCode)
	}

	return nil
}

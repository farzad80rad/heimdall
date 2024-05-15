package proxyHttp

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"heimdall/internal/config"
	heimdallErrors "heimdall/internal/errors"
	"heimdall/internal/proxy"
	"heimdall/internal/utils"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type httpProxy struct {
	cb              *gobreaker.CircuitBreaker
	proxy           *httputil.ReverseProxy
	host            string
	bodyCheckConfig []config.HostMatchInfo
}

func New(host string, config config.CircuitBreakerConfig, checkConfig []config.HostMatchInfo) (proxy.Proxy, error) {
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
			if config.FailureToleranceCount != 0 {
				maxTolerance = config.FailureToleranceCount
			}
			return counts.ConsecutiveFailures > maxTolerance
		},
		IsSuccessful: func(err error) bool {
			return err == nil || errors.Is(err, heimdallErrors.BadRequest)
		},
	})

	return &httpProxy{
		cb:              cb,
		host:            host,
		proxy:           proxy,
		bodyCheckConfig: checkConfig,
	}, nil
}

func (a *httpProxy) Ping(url string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := utils.DoHTTPGetRequest(ctx, url)
	return err == nil
}

func (a *httpProxy) Proxy(c *gin.Context) error {

	if a.bodyCheckConfig != nil {
		if c.Request.Body != nil {
			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			err := utils.ValidateBody(c.Request.Method, body, a.bodyCheckConfig)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	}

	_, err := a.cb.Execute(func() (interface{}, error) {
		var cbError error
		a.proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
			cbError = err
		}
		a.proxy.ServeHTTP(c.Writer, c.Request)
		return nil, cbError
	})

	if err != nil {
		if err == gobreaker.ErrOpenState {
			c.JSON(http.StatusBadGateway, gin.H{"error": heimdallErrors.HostIsDown.Error()})
			return heimdallErrors.HostIsDown
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": heimdallErrors.ConnectionIssue.Error()})
		return heimdallErrors.ConnectionIssue
	}
	return nil
}

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"heimdall/config"
	"heimdall/errors"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Api interface {
	Proxy(ctx *gin.Context) error
}

type api struct {
	cb    *gobreaker.CircuitBreaker
	proxy *httputil.ReverseProxy
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
			fmt.Println("check ready to be tripped")
			maxTolerance := uint32(10)
			if config.FierierToleranceCount != 0 {
				maxTolerance = config.FierierToleranceCount
			}
			return counts.ConsecutiveFailures > maxTolerance
		},
	})

	return &api{
		cb:    cb,
		proxy: proxy,
	}, nil
}

func (api *api) Proxy(c *gin.Context) error {
	_, err := api.cb.Execute(func() (interface{}, error) {
		var cbError error
		api.proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
			cbError = err
		}
		api.proxy.ServeHTTP(c.Writer, c.Request)
		return nil, cbError
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "service is currently unable to respond. please try again"})
	}
	if err == gobreaker.ErrOpenState {
		return errors.HostIsDown
	}
	return err
}

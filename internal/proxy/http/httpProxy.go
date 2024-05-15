package proxyHttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"heimdall/internal/config"
	heimdallErrors "heimdall/internal/errors"
	"heimdall/internal/proxy"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"time"
)

type httpProxy struct {
	cb              *gobreaker.CircuitBreaker
	proxy           *httputil.ReverseProxy
	host            string
	bodyCheckConfig *config.RequestBodyCheckConfig
}

func New(host string, config config.CircuitBreakerConfig, checkConfig *config.RequestBodyCheckConfig) (proxy.Proxy, error) {
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
	err := proxy.DoHTTPGetRequest(ctx, url)
	return err == nil
}

func (a *httpProxy) Proxy(c *gin.Context) error {

	if a.bodyCheckConfig != nil {
		//req.Header = c.Request.Header
		// Handle body (for POST, PUT, etc.)
		if c.Request.Body != nil {
			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			var requestBodyMap map[string]interface{}
			err := json.Unmarshal(body, &requestBodyMap)
			if err != nil {
				return errors.Join(heimdallErrors.BadRequest, errors.New("not in json format"))
			}
			for _, feildInfo := range a.bodyCheckConfig.MandatoryFields {
				v, found := requestBodyMap[feildInfo.FieldName]
				if !(found && reflect.TypeOf(v).Kind() == feildInfo.Type) {
					return errors.Join(heimdallErrors.BadRequest, errors.New("missing required field "+feildInfo.FieldName))
				}
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
			return heimdallErrors.HostIsDown
		}
		return heimdallErrors.ConnectionIssue
	}
	return nil
}
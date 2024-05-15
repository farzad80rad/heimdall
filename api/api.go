package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"heimdall/config"
	heimdallErrors "heimdall/errors"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"time"
)

type Api interface {
	Proxy(ctx *gin.Context) error
	Ping(url string) bool
}

type api struct {
	cb              *gobreaker.CircuitBreaker
	proxy           *httputil.ReverseProxy
	host            string
	bodyCheckConfig *config.RequestBodyCheckConfig
}

func NewApi(host string, config config.CircuitBreakerConfig, checkConfig *config.RequestBodyCheckConfig) (Api, error) {
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

	return &api{
		cb:              cb,
		host:            host,
		proxy:           proxy,
		bodyCheckConfig: checkConfig,
	}, nil
}

func (a *api) Ping(url string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := doHTTPGetRequest(ctx, a.host+url)
	return err == nil
}

func (a *api) Proxy(c *gin.Context) error {

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

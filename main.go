package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type MatchPolicy struct {
	Url             string
	DestinationHost string
	HttpTypes       []string
}

type CircuitBreakerConfig struct {
	ExamineWindow         time.Duration
	QuarantineDuration    time.Duration
	FierierToleranceCount uint32
}

type ApiConfig struct {
	Match                MatchPolicy
	CircuitBreakerConfig CircuitBreakerConfig
}

func main() {

	apiConfigs := []ApiConfig{
		{Match: MatchPolicy{Url: "/api1/*any",
			DestinationHost: "http://localhost:23232",
			HttpTypes:       []string{http.MethodGet, http.MethodConnect, http.MethodOptions, http.MethodPost}},
			CircuitBreakerConfig: CircuitBreakerConfig{
				ExamineWindow:         60 * time.Second,
				QuarantineDuration:    10 * time.Second,
				FierierToleranceCount: 3,
			}},
		{Match: MatchPolicy{Url: "/api3", DestinationHost: "http://localhost:23232", HttpTypes: []string{http.MethodConnect, http.MethodPost, http.MethodOptions}}},
	}

	r := gin.Default()
	for _, apiConfig := range apiConfigs {
		targetURL, _ := url.Parse(apiConfig.Match.DestinationHost) // Replace with your target URL
		cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:     apiConfig.Match.Url + "-" + apiConfig.Match.DestinationHost,
			Interval: apiConfig.CircuitBreakerConfig.ExamineWindow,
			Timeout:  apiConfig.CircuitBreakerConfig.QuarantineDuration,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				maxTolerance := uint32(10)
				if apiConfig.CircuitBreakerConfig.FierierToleranceCount != 0 {
					maxTolerance = apiConfig.CircuitBreakerConfig.FierierToleranceCount
				}
				return counts.ConsecutiveFailures > maxTolerance
			},
		})
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		r.Match(apiConfig.Match.HttpTypes, apiConfig.Match.Url, func(c *gin.Context) {
			_, err := cb.Execute(func() (interface{}, error) {
				var cbError error
				proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
					cbError = err
				}
				proxy.ServeHTTP(c.Writer, c.Request)
				return nil, cbError
			})
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": "service is currently unable to respond. please try again"})
			}
		})
	}

	r.Run(":23982") // Run on port 8080
}

package main

import (
	"github.com/gin-gonic/gin"
	"heimdall/api"
	"heimdall/config"
	"heimdall/errors"
	"heimdall/loadBalancer"
	"net/http"
	"time"
)

func main() {

	apiConfigs := []config.ApiConfig{
		{Match: config.MatchPolicy{Url: "/api1/*any",
			HttpTypes: []string{http.MethodGet, http.MethodConnect, http.MethodOptions, http.MethodPost}},
			CircuitBreakerConfig: config.CircuitBreakerConfig{
				ExamineWindow:         60 * time.Second,
				QuarantineDuration:    10 * time.Second,
				FierierToleranceCount: 3,
			},
			HostInfo: config.HostInfo{
				HostAddress: []string{"http://localhost:23232", "http://localhost:23243"},
			},
		},
	}

	r := gin.Default()
	for _, apiConfig := range apiConfigs {
		proxyApi(apiConfig, r)
	}

	r.Run(":23982") // Run on port 8080
}

func proxyApi(apiConfig config.ApiConfig, r *gin.Engine) error {

	lb := loadBalancer.NewRoundRobin(apiConfig.HostInfo.HostAddress)
	hosts := make(map[string]api.Api, 3*len(apiConfig.HostInfo.HostAddress))
	for _, h := range apiConfig.HostInfo.HostAddress {
		p, err := api.NewApi(h, apiConfig.CircuitBreakerConfig)
		if err != nil {
			return err
		}
		hosts[h] = p
	}
	r.Match(apiConfig.Match.HttpTypes, apiConfig.Match.Url, func(c *gin.Context) {
		destination := lb.Next()
		if err := hosts[destination].Proxy(c); err == errors.HostIsDown {
			lb.DisableHost(destination, apiConfig.CircuitBreakerConfig.QuarantineDuration)
		}
	})
	return nil
}

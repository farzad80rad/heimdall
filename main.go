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
			HostInfo: config.HostLoadPolicy{
				LoadBalanceType: config.LoadBalanceType_WEIGHTED_ROUNDROBIN,
				HostUnits: []config.HostUnit{
					{
						Url:    "http://localhost:23232",
						Weight: 8,
					},
					{
						Url:    "http://localhost:23264",
						Weight: 4,
					},
				},
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

	var lb loadBalancer.LoadBalancer
	switch apiConfig.HostInfo.LoadBalanceType {
	case config.LoadBalanceType_WEIGHTED_ROUNDROBIN:
		lb = loadBalancer.NewWeightedRoundRobin(apiConfig.HostInfo.HostUnits)
	default:
		hosts := make([]string, len(apiConfig.HostInfo.HostUnits))
		for i, unit := range apiConfig.HostInfo.HostUnits {
			hosts[i] = unit.Url
		}
		lb = loadBalancer.NewRoundRobin(hosts)
	}

	hosts := make(map[string]api.Api, 3*len(apiConfig.HostInfo.HostUnits))
	for _, h := range apiConfig.HostInfo.HostUnits {
		p, err := api.NewApi(h.Url, apiConfig.CircuitBreakerConfig)
		if err != nil {
			return err
		}
		hosts[h.Url] = p
	}
	r.Match(apiConfig.Match.HttpTypes, apiConfig.Match.Url, func(c *gin.Context) {
		destination := lb.Next()
		if err := hosts[destination].Proxy(c); err == errors.HostIsDown {
			lb.DisableHost(destination, apiConfig.CircuitBreakerConfig.QuarantineDuration)
		}
	})
	return nil
}

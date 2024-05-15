package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"heimdall/api"
	"heimdall/config"
	heimdallErrors "heimdall/errors"
	"heimdall/loadBalancer"
	"net/http"
	"reflect"
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
						Host:   "http://localhost:23231",
						Weight: 8,
					},
					{
						Host:   "http://localhost:23235",
						Weight: 4,
					},
				},
			},
			HealthCheckConfig: &config.HealthCheckConfig{
				Path:              "/health",
				FailureThreshHold: 3,
				Interval:          5 * time.Second,
			},
			RequestBodyCheckConfig: &config.RequestBodyCheckConfig{
				MandatoryFields: []config.RequestValidationUnit{
					{FieldName: "f1", Type: reflect.String},
					{FieldName: "f2", Type: reflect.Bool},
					{FieldName: "f3", Type: reflect.Float64},
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
			hosts[i] = unit.Host
		}
		lb = loadBalancer.NewRoundRobin(hosts)
	}

	hosts := make(map[string]api.Api, 3*len(apiConfig.HostInfo.HostUnits))
	for _, h := range apiConfig.HostInfo.HostUnits {
		p, err := api.NewApi(h.Host, apiConfig.CircuitBreakerConfig, apiConfig.RequestBodyCheckConfig)
		if err != nil {
			return err
		}
		hosts[h.Host] = p

		if apiConfig.HealthCheckConfig != nil {
			go func(ap api.Api) {
				failureCount := 0
				for {
					if isActive := ap.Ping(apiConfig.HealthCheckConfig.Path); isActive {
						if failureCount > 0 {
							failureCount = 0
							lb.SetHostStatus(h.Host, true)
						}
					} else {
						failureCount++
					}

					if failureCount == apiConfig.HealthCheckConfig.FailureThreshHold {
						lb.SetHostStatus(h.Host, false)
					}

					sleepTime := apiConfig.HealthCheckConfig.Interval
					if sleepTime < time.Second {
						sleepTime = 5 * time.Second
					}
					time.Sleep(sleepTime)
				}
			}(p)
		}
	}

	r.Match(apiConfig.Match.HttpTypes, apiConfig.Match.Url, func(c *gin.Context) {
		destination := lb.Next()
		host := hosts[destination]
		err := host.Proxy(c)
		if err != nil {
			fmt.Println(err)
			if err == heimdallErrors.HostIsDown {
				if apiConfig.HealthCheckConfig == nil {
					lb.DisableHostForDuration(destination, apiConfig.CircuitBreakerConfig.QuarantineDuration)
				} else {
					lb.SetHostStatus(destination, false)
				}
			}
			if errors.Is(err, heimdallErrors.BadRequest) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		}
		return
	})
	return nil
}

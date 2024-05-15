package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"heimdall/config"
	heimdallErrors "heimdall/errors"
	"heimdall/loadBalancer"
	"heimdall/proxy"
	"heimdall/proxy/grpc"
	proxyHttp "heimdall/proxy/http"
	"net/http"
	"reflect"
	"time"
)

func main() {

	apiConfigs := []config.ApiConfig{
		{
			Match: config.MatchPolicy{
				ConnectionType: config.ConnectionType_GPRC,
				Name:           string(proxyGrpc.HeimdallGrpcService_MESSEGING),
				Url:            "/backend.messaging.v1.MessagingService/Echo",
				HttpTypes:      []string{http.MethodGet, http.MethodConnect, http.MethodOptions, http.MethodPost},
			},
			CircuitBreakerConfig: config.CircuitBreakerConfig{
				ExamineWindow:         60 * time.Second,
				QuarantineDuration:    time.Second,
				FierierToleranceCount: 3,
			},
			HostInfo: config.HostLoadPolicy{
				LoadBalanceType: config.LoadBalanceType_ROUNDROBIN,
				HostUnits: []config.HostUnit{
					{
						Host:   "localhost:50501",
						Weight: 8,
					},
					{
						Host:   "localhost:50502",
						Weight: 8,
					},
				},
			},
			/*HealthCheckConfig: &config.HealthCheckConfig{
				Path:              "/health",
				FailureThreshHold: 3,
				Interval:          5 * time.Second,
			},*/
			RequestBodyCheckConfig: &config.RequestBodyCheckConfig{
				MandatoryFields: []config.RequestValidationUnit{
					{FieldName: "f1", Type: reflect.String},
					{FieldName: "f2", Type: reflect.Bool},
					{FieldName: "f3", Type: reflect.Float64},
				},
			},
		},
		{
			Match: config.MatchPolicy{
				ConnectionType: config.ConnectionType_HTTP1,
				Url:            "/api1/*any",
				HttpTypes:      []string{http.MethodGet, http.MethodConnect, http.MethodOptions, http.MethodPost},
			},
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

// Example unary interceptor

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

	grpcMuxMap := make(map[string]*runtime.ServeMux)
	hosts := make(map[string]proxy.Proxy, 3*len(apiConfig.HostInfo.HostUnits))
	for _, h := range apiConfig.HostInfo.HostUnits {
		var p proxy.Proxy
		var err error
		if apiConfig.Match.ConnectionType == config.ConnectionType_GPRC {
			var mux *runtime.ServeMux
			if m, found := grpcMuxMap[h.Host]; found {
				mux = m
			} else {
				mux = runtime.NewServeMux()
				grpcMuxMap[h.Host] = mux
			}
			p, err = proxyGrpc.New(h.Host, apiConfig.CircuitBreakerConfig, apiConfig.RequestBodyCheckConfig, mux, proxyGrpc.HeimdallGrpcService(apiConfig.Match.Name))
		} else {
			p, err = proxyHttp.New(h.Host, apiConfig.CircuitBreakerConfig, apiConfig.RequestBodyCheckConfig)
		}
		if err != nil {
			return err
		}
		hosts[h.Host] = p

		if apiConfig.HealthCheckConfig != nil {
			go func(ap proxy.Proxy) {
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

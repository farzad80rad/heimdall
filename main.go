package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"heimdall/internal/config"
	"heimdall/internal/heimdall"
	"heimdall/internal/proxy/grpc"
	"log"
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
		h, err := heimdall.NewApiGateway(apiConfig, r)
		if err != nil {
			log.Println(fmt.Sprintf("creating gateway for url %s failed for reason: %v", apiConfig.Match.Url, err))
			continue
		}
		if err := h.Run(); err != nil {
			log.Println(fmt.Sprintf("running gateway for url %s failed for reason: %v", apiConfig.Match.Url, err))
		}
	}

	r.Run(":23982") // Run on port 8080
}

// Example unary interceptor

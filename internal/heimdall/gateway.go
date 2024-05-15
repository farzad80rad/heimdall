package heimdall

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"heimdall/internal/config"
	heimdallErrors "heimdall/internal/errors"
	loadBalancer2 "heimdall/internal/loadBalancer"
	"heimdall/internal/proxy"
	proxyGrpc "heimdall/internal/proxy/grpc"
	proxyHttp "heimdall/internal/proxy/http"
	"net/http"
	"time"
)

type ApiGateway interface {
	Run() error
}

type gateway struct {
	lb        loadBalancer2.LoadBalancer
	apiConfig config.ApiConfig
	hosts     map[string]proxy.Proxy
	r         *gin.Engine
}

func NewApiGateway(apiConfig config.ApiConfig, r *gin.Engine) (ApiGateway, error) {
	g := &gateway{
		lb:        createLoadBalancer(apiConfig),
		apiConfig: apiConfig,
		r:         r,
	}
	hosts, err := g.createHosts()
	if err != nil {
		return nil, err
	}
	g.hosts = hosts
	return g, nil
}

func (g *gateway) Run() error {
	g.r.Match(g.apiConfig.Match.HttpTypes, g.apiConfig.Match.Url, g.handleRequest)
	return nil
}

func (g *gateway) handleRequest(c *gin.Context) {
	destination := g.lb.Next()
	hostProxy := g.hosts[destination]

	err := hostProxy.Proxy(c)
	if err == nil {
		return
	}

	if err == heimdallErrors.HostIsDown {
		c.Status(http.StatusBadGateway)
		if g.apiConfig.HealthCheckConfig == nil {
			g.lb.DisableHostForDuration(destination, g.apiConfig.CircuitBreakerConfig.QuarantineDuration)
		} else {
			g.lb.SetHostStatus(destination, false)
		}
	}

	if errors.Is(err, heimdallErrors.BadRequest) {
		c.Status(http.StatusBadRequest)
		return
	}

}

func createLoadBalancer(apiConfig config.ApiConfig) loadBalancer2.LoadBalancer {
	switch apiConfig.HostInfo.LoadBalanceType {
	case config.LoadBalanceType_WEIGHTED_ROUNDROBIN:
		return loadBalancer2.NewWeightedRoundRobin(apiConfig.HostInfo.HostUnits)
	default:
		hosts := make([]string, len(apiConfig.HostInfo.HostUnits))
		for i, unit := range apiConfig.HostInfo.HostUnits {
			hosts[i] = unit.Host
		}
		return loadBalancer2.NewRoundRobin(hosts)
	}
}

func (g *gateway) createHosts() (map[string]proxy.Proxy, error) {
	hosts := make(map[string]proxy.Proxy, 3*len(g.apiConfig.HostInfo.HostUnits))
	muxPerHost := make(map[string]*runtime.ServeMux) // to prevent duplicate mux for a single grpc host
	for _, h := range g.apiConfig.HostInfo.HostUnits {
		var p proxy.Proxy
		var err error
		if g.apiConfig.Match.ConnectionType == config.ConnectionType_GPRC {
			var mux *runtime.ServeMux
			if m, found := muxPerHost[h.Host]; found {
				mux = m
			} else {
				mux = runtime.NewServeMux()
				muxPerHost[h.Host] = mux
			}
			p, err = proxyGrpc.New(h.Host, g.apiConfig.CircuitBreakerConfig, g.apiConfig.RequestBodyCheckConfig, mux, proxyGrpc.HeimdallGrpcService(g.apiConfig.Match.Name))
		} else {
			p, err = proxyHttp.New(h.Host, g.apiConfig.CircuitBreakerConfig, g.apiConfig.RequestBodyCheckConfig)
		}

		if err != nil {
			return nil, err
		}
		hosts[h.Host] = p

		if g.apiConfig.HealthCheckConfig != nil {
			go g.watchHealth(p, h.Host)
		}
	}
	return hosts, nil
}

func (g *gateway) watchHealth(ap proxy.Proxy, host string) {
	failureCount := 0
	for {
		if isActive := ap.Ping(host + g.apiConfig.HealthCheckConfig.Path); isActive {
			if failureCount > 0 {
				failureCount = 0
				g.lb.SetHostStatus(host, true)
			}
		} else {
			failureCount++
		}

		if failureCount == g.apiConfig.HealthCheckConfig.FailureThreshHold {
			g.lb.SetHostStatus(host, false)
		}

		sleepTime := g.apiConfig.HealthCheckConfig.Interval
		if sleepTime < time.Second {
			sleepTime = 5 * time.Second
		}
		time.Sleep(sleepTime)
	}
}
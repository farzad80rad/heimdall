package proxyGrpc

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	golang "heimdall/Proto/gen/Proto"
	"heimdall/internal/config"
	heimdallErrors "heimdall/internal/errors"
	"heimdall/internal/proxy"
	"heimdall/internal/utils"
	"io"
	"net/http"
	"time"
)

type grpcProxy struct {
	cb              *gobreaker.CircuitBreaker
	host            string
	bodyCheckConfig []config.HostMatchInfo
	mux             *runtime.ServeMux
	/* race condition may happen. but doesn't matter :).
	the purpose of this field is to check consecutive error happening.
	so if this error changes to nil or some other err rather than heimdallErrors.HostIsDown , then
	no consecutive  err has happened!*/
	lastError error
}

func New(host string, config config.CircuitBreakerConfig, checkConfig []config.HostMatchInfo,
	mux *runtime.ServeMux, grpcService HeimdallGrpcService) (proxy.Proxy, error) {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:     host,
		Interval: config.ExamineWindow,
		Timeout:  config.QuarantineDuration,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			maxTolerance := uint32(10)
			if config.FailureToleranceCount != 0 {
				maxTolerance = config.FailureToleranceCount
			}
			return counts.ConsecutiveFailures > maxTolerance
		},
		IsSuccessful: func(err error) bool {
			return err == nil || errors.Is(err, heimdallErrors.BadRequest)
		},
	})

	g := &grpcProxy{
		cb:              cb,
		host:            host,
		mux:             mux,
		bodyCheckConfig: checkConfig,
	}
	err := g.establishConnection(grpcService, host, mux)
	if err != nil {
		panic(err)
		return nil, err
	}

	return g, nil
}

func (g *grpcProxy) establishConnection(identifier HeimdallGrpcService, host string, mux *runtime.ServeMux) error {
	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(g.NewInterceptor()),
	}
	var err error
	switch identifier {
	case HeimdallGrpcService_MESSEGING:
		err = golang.RegisterMessagingServiceHandlerFromEndpoint(ctx, mux, "localhost:50501", opts)
	case HeimdallGrpcService_CARGO:
		err = golang.RegisterCargoServiceHandlerFromEndpoint(ctx, mux, host, opts)
	}
	return err
}

func (a *grpcProxy) Ping(url string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := utils.DoHTTPGetRequest(ctx, a.host+url)
	return err == nil
}

func (a *grpcProxy) Proxy(c *gin.Context) error {

	if a.bodyCheckConfig != nil {
		if c.Request.Body != nil {
			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			err := utils.ValidateBody(c.Request.Method, body, a.bodyCheckConfig)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	}

	_, err := a.cb.Execute(func() (interface{}, error) {
		a.mux.ServeHTTP(c.Writer, c.Request)
		return nil, a.lastError
	})
	if err == gobreaker.ErrOpenState {
		return heimdallErrors.HostIsDown
	}
	return err
}

func (a *grpcProxy) NewInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var err error
		err = invoker(ctx, method, req, resp, cc, opts...)
		if err != nil {
			errCode := status.Code(err)
			if errCode == codes.DeadlineExceeded || errCode == codes.Unavailable || errCode == codes.Unimplemented {
				err = heimdallErrors.HostIsDown
			}
			if errCode == codes.InvalidArgument || errCode == codes.PermissionDenied || errCode == codes.NotFound ||
				errCode == codes.Unauthenticated || errCode == codes.Canceled || errCode == codes.Aborted ||
				errCode == codes.AlreadyExists || errCode == codes.OutOfRange || errCode == codes.Unknown {
				// the error will still be showed to client, but should not lead to considering the host as down.
				err = nil
			}
		}
		a.lastError = err
		return err
	}
}

package server

import (
	"fmt"
	"net"
	"time"

	"achuala.in/rpc-bp/logger"
	"achuala.in/rpc-bp/util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server is the implementation of an API Server
type Server struct {
	// gRPC-server exposing both the API and health
	grpc *grpc.Server
	// Logger interface
	log logger.Logger
	//
	shutdown *util.ShutdownWaitGroup
}

// NewServer returns a new configured instance of Server
func NewGrpcServer(name string, keepAlive bool) *Server {
	return NewServerWithOpts(name, keepAlive, []grpc.UnaryServerInterceptor{}, []grpc.StreamServerInterceptor{})
}

// NewServerWithOpts returns a new configured instance of Server with additional interceptros specified
func NewServerWithOpts(name string, keepAlive bool, unaryServerInterceptors []grpc.UnaryServerInterceptor, streamServerInterceptors []grpc.StreamServerInterceptor) *Server {
	s := &Server{
		log:      logger.WithName(name),
		shutdown: util.NewShutdownWaitGroup(),
	}

	// Add default interceptors
	unaryServerInterceptors = append(unaryServerInterceptors,
		grpc_ctxtags.UnaryServerInterceptor(),
		logger.UnaryServerInterceptor(s.log),
		grpc_recovery.UnaryServerInterceptor(), // add recovery from panics
		// own wrapper is used to unpack nested messages
		//grpc_validator.UnaryServerInterceptor(), // add message validator
		//grpc_validator_wrapper.UnaryServerInterceptor(), // add message validator wrapper
	)
	streamServerInterceptors = append(streamServerInterceptors,
		grpc_recovery.StreamServerInterceptor(), // add recovery from panics
		// own wrapper is used to unpack nested messages
		//grpc_validator.StreamServerInterceptor(), // add message validator
		//grpc_validator_wrapper.StreamServerInterceptor(), // add message validator wrapper
	)

	// Configure gRPC server
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamServerInterceptors...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryServerInterceptors...)),
	}
	if keepAlive {
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Second,
		}))
	}
	s.grpc = grpc.NewServer(opts...)

	// Add grpc health check service
	healthpb.RegisterHealthServer(s.grpc, health.NewServer())

	// Enable reflection API
	reflection.Register(s.grpc)

	return s
}

// RegisterService registers your gRPC service implementation with the server
func (s *Server) RegisterService(f func(grpc.ServiceRegistrar)) {
	f(s.grpc)
}

// Serve starts the api listeners of the Server
func (s *Server) Serve(apiAddr string) error {
	// Setup grpc listener
	apiLis, err := net.Listen("tcp", apiAddr)
	if err != nil {
		return err
	}
	defer apiLis.Close()
	return s.ServeFromListener(apiLis)
}

// ServeFromListener starts the api listeners of the Server
func (s *Server) ServeFromListener(apiLis net.Listener) error {
	shutdown := s.shutdown

	// Start routine waiting for signals
	shutdown.RegisterSignalHandler(func() {
		//  gRPC server
		s.log.Info("grpc server stopping gracefully")
		s.grpc.GracefulStop()
	})

	s.log.Info("starting to serve grpc", "addr", apiLis.Addr())
	err := s.grpc.Serve(apiLis)
	s.log.Info("grpc server stopped")

	// Check if we are expecting shutdown
	if !shutdown.IsExpected() {
		panic(fmt.Sprintf("shutdown unexpected, grpc serve returned: %v", err))
	}
	// Wait for both shutdown signals and close the channel
	if ok := shutdown.WaitOrTimeout(30 * time.Second); !ok {
		panic("shutting down gracefully exceeded 30 seconds")
	}
	return err // Return the error, if grpc stopped gracefully there is no error
}

// Tell the server to shutdown
func (s *Server) Shutdown() {
	s.shutdown.Expect()
}

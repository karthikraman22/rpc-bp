package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"achuala.in/rpc-bp/logger"
	"achuala.in/rpc-bp/util"
	"github.com/gin-gonic/gin"
)

// Server is the implementation of an API Server
type RestServer struct {
	// Http Server
	httpServer *http.Server
	// Gin router
	router *gin.Engine
	// Logger interface
	log logger.Logger
	//
	shutdown *util.ShutdownWaitGroup
}

// ServiceRegistrar wraps a single method that supports service registration.
type RestServiceRegistrar interface {
	// RegisterService registers a service and its implementation to the
	// gin router
	RegisterRoutes(router *gin.Engine)
}

// NewServer returns a new configured instance of Server
func NewRestServer(name string, keepAlive bool) *RestServer {

	router := gin.Default()

	httpSrv := &http.Server{
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &RestServer{log: logger.WithName(name), httpServer: httpSrv, router: router, shutdown: util.NewShutdownWaitGroup()}
}

func (s *RestServer) RegisterService(f func(*gin.Engine)) {
	f(s.router)
}

// Serve starts the api listeners of the Server
func (s *RestServer) Serve(apiAddr string) error {
	// Setup grpc listener
	apiLis, err := net.Listen("tcp", apiAddr)
	if err != nil {
		return err
	}
	defer apiLis.Close()
	return s.ServeFromListener(apiLis)
}

// ServeFromListener starts the api listeners of the Server
func (s *RestServer) ServeFromListener(apiLis net.Listener) error {
	shutdown := s.shutdown

	// Start routine waiting for signals
	shutdown.RegisterSignalHandler(func() {
		//  gRPC server
		s.log.Info("rest server stopping gracefully")
		s.httpServer.Shutdown(context.Background())
	})

	s.log.Info("starting to serve rest", "addr", apiLis.Addr())
	err := s.httpServer.Serve(apiLis)
	s.log.Info("rest server stopped")

	// Check if we are expecting shutdown
	if !shutdown.IsExpected() {
		panic(fmt.Sprintf("shutdown unexpected, rest serve returned: %v", err))
	}
	// Wait for both shutdown signals and close the channel
	if ok := shutdown.WaitOrTimeout(30 * time.Second); !ok {
		panic("shutting down gracefully exceeded 30 seconds")
	}
	return err // Return the error, if grpc stopped gracefully there is no error
}

// Tell the server to shutdown
func (s *RestServer) Shutdown() {
	s.shutdown.Expect()
}

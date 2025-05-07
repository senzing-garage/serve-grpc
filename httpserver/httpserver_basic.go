package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicHTTPServer is the default implementation of the HttpServer interface.
type BasicHTTPServer struct {
	AvoidServing      bool
	EnableAll         bool
	EnableGRPC        bool
	GRPCRoutePrefix   string
	GRPCServer        *grpc.Server
	logger            logging.Logging
	LogLevelName      string
	ReadHeaderTimeout time.Duration
	ServerAddress     string
	ServerPort        int
}

const OptionCallerSkip = 3

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Handler method gets the http.ServeMux for the gRPC over HTTP service.

Input
  - ctx: A context to control lifecycle.

Output
  - httpServeMux - the Mux of the service.
*/
func (httpServer *BasicHTTPServer) Handler(ctx context.Context) *http.ServeMux {
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/", httpServer.grpcFunc(ctx))

	return rootMux
}

/*
The Serve method starts the HTTP server.

Input
  - ctx: A context to control lifecycle.

Output
  - Nothing is returned, except for an error.
*/
func (httpServer *BasicHTTPServer) Serve(ctx context.Context) error {
	var err error

	rootMux := http.NewServeMux()

	// Enable GRPC over HTTP.

	httpServer.registerGRPC(ctx, rootMux)

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	httpServer.log(2001, listenOnAddress)

	server := http.Server{
		ReadHeaderTimeout: httpServer.ReadHeaderTimeout,
		Addr:              listenOnAddress,
		Handler:           rootMux,
	}

	// Start a web browser.  Unless disabled.

	if !httpServer.AvoidServing {
		err = server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}

	return wraperror.Errorf(err, "httpserver error: %w", err)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (httpServer *BasicHTTPServer) grpcFunc(ctx context.Context) http.HandlerFunc {
	_ = ctx

	wrappedGrpc := grpcweb.WrapServer(httpServer.GRPCServer)

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			httpServer.log(1000, req)
			wrappedGrpc.ServeHTTP(resp, req)
		} else { // Fall back to other servers.
			httpServer.log(1001, req)
			http.DefaultServeMux.ServeHTTP(resp, req)
		}
	})
}

func (httpServer *BasicHTTPServer) registerGRPC(ctx context.Context, rootMux *http.ServeMux) {
	if httpServer.EnableAll || httpServer.EnableGRPC {
		pattern := fmt.Sprintf("/%s/", httpServer.GRPCRoutePrefix)
		prefix := "/" + httpServer.GRPCRoutePrefix
		grpcWebMux := httpServer.Handler(ctx)
		handler := http.StripPrefix(prefix, grpcWebMux)
		rootMux.Handle(pattern, handler)
		httpServer.log(2002, httpServer.ServerPort, httpServer.GRPCRoutePrefix)
	}
}

// --- Logging -------------------------------------------------------------------------

// Get the Logger singleton.
func (httpServer *BasicHTTPServer) getLogger() logging.Logging {
	var err error

	if httpServer.logger == nil {
		options := []interface{}{
			logging.OptionCallerSkip{Value: OptionCallerSkip},
			logging.OptionMessageFields{Value: []string{"id", "text", "reason", "errors", "details"}},
		}

		httpServer.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}

		err = httpServer.logger.SetLogLevel(httpServer.LogLevelName)
		if err != nil {
			panic(err)
		}
	}

	return httpServer.logger
}

// Log message.
func (httpServer *BasicHTTPServer) log(messageNumber int, details ...interface{}) {
	httpServer.getLogger().Log(messageNumber, details...)
}

package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/senzing-garage/go-helpers/wraperror"
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
	GRPCRoutePrefix   string // IMPROVE: Only works with "grpc"
	GRPCServer        *grpc.Server
	ReadHeaderTimeout time.Duration
	ServerAddress     string
	ServerPort        int
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Serve method starts the HTTP server.

Input
  - ctx: A context to control lifecycle.

Output
  - Nothing is returned, except for an error.
*/

func (httpServer *BasicHTTPServer) Serve(ctx context.Context) error {
	var (
		err          error
		userMessages []string
	)

	rootMux := http.NewServeMux()

	// Enable GRPC over HTTP.

	userMessages = httpServer.registerGRPC(ctx, rootMux, userMessages)

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	userMessages = append(userMessages, fmt.Sprintf("Starting server on interface:port '%s'...\n", listenOnAddress))
	outputUserMessages(userMessages)

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
			wrappedGrpc.ServeHTTP(resp, req)
		} else {
			// Fall back to other servers.
			http.DefaultServeMux.ServeHTTP(resp, req)
		}
	})
}

func (httpServer *BasicHTTPServer) registerGRPC(
	ctx context.Context,
	rootMux *http.ServeMux,
	userMessages []string,
) []string {
	result := userMessages

	if httpServer.EnableAll || httpServer.EnableGRPC {
		rootMux.HandleFunc(fmt.Sprintf("/%s/", httpServer.GRPCRoutePrefix), httpServer.grpcFunc(ctx))
		result = append(result,
			fmt.Sprintf(
				"Serving GRPC over HTTP at http://localhost:%d/%s",
				httpServer.ServerPort,
				httpServer.GRPCRoutePrefix,
			))
	}

	return result
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}

func outputUserMessages(messages []string) {
	for _, message := range messages {
		outputln(message)
	}
}

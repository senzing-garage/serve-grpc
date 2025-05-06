package httpserver

import (
	"context"
	"errors"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The HTTPServer interface...
type HTTPServer interface {
	Serve(ctx context.Context) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6204xxxx".
const ComponentID = 6042

// Log message prefix.
const Prefix = "serve-grpc.httpserver."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates.
var IDMessages = map[int]string{
	1000: "gRPC Web request: %+v",
	1001: "HTTP Web request: %+v",
	2001: "Starting HTTP server on interface:port '%s'",
	2002: "Serving GRPC over HTTP at http://localhost:%d/%s",
}

// Status strings for specific messages.
var IDStatuses = map[int]string{}

var errForPackage = errors.New("httpserver")

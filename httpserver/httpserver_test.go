package httpserver_test

import (
	"context"
	"testing"
	"time"

	"github.com/senzing-garage/serve-grpc/httpserver"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface methods
// ----------------------------------------------------------------------------

func TestBasicHTTPServer_Serve(test *testing.T) {
	ctx := test.Context()
	httpServer := getTestObject(ctx, test)
	err := httpServer.Serve(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, t *testing.T) *httpserver.BasicHTTPServer {
	t.Helper()

	_ = ctx
	result := &httpserver.BasicHTTPServer{
		GRPCRoutePrefix:   "api",
		AvoidServing:      true,
		EnableAll:         true,
		LogLevelName:      "INFO",
		ReadHeaderTimeout: 10 * time.Second,
	}

	return result
}

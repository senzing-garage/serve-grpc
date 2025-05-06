package httpserver

import (
	"context"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The HTTPServer interface...
type HTTPServer interface {
	Serve(ctx context.Context) error
}

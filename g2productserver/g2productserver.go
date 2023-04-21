package g2productserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go-base/g2product"
	"github.com/senzing/g2-sdk-go/g2api"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2product"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/observer"
)

var (
	g2productSingleton g2api.G2product
	g2productSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *G2ProductServer) getLogger() logging.LoggingInterface {
	var err error = nil
	if server.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 3},
		}
		server.logger, err = logging.NewSenzingToolsLogger(ProductId, IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return server.logger
}

// Trace method entry.
func (server *G2ProductServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *G2ProductServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (server *G2ProductServer) error(messageNumber int, details ...interface{}) error {
	return server.getLogger().Error(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2product() g2api.G2product {
	g2productSyncOnce.Do(func() {
		g2productSingleton = &g2sdk.G2product{}
	})
	return g2productSingleton
}

func GetSdkG2product() g2api.G2product {
	return getG2product()
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing/g2-sdk-go/g2product.G2product
// ----------------------------------------------------------------------------

func (server *G2ProductServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(3, request)
		defer func() { server.traceExit(4, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2product := getG2product()
	// err := g2product.Destroy(ctx)
	err = server.error(4001)
	response := g2pb.DestroyResponse{}
	return &response, err
}

func (server *G2ProductServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(9, request)
		defer func() { server.traceExit(10, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2product := getG2product()
	// err := g2product.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err = server.error(4002)
	response := g2pb.InitResponse{}
	return &response, err
}

func (server *G2ProductServer) License(ctx context.Context, request *g2pb.LicenseRequest) (*g2pb.LicenseResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(11, request)
		defer func() { server.traceExit(12, request, result, err, time.Since(entryTime)) }()
	}
	g2product := getG2product()
	result, err = g2product.License(ctx)
	response := g2pb.LicenseResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ProductServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, observer.GetObserverId(ctx))
		defer func() { server.traceExit(2, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2product := getG2product()
	return g2product.RegisterObserver(ctx, observer)
}

func (server *G2ProductServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(13, logLevelName)
		defer func() { server.traceExit(14, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	g2product := getG2product()
	g2product.SetLogLevel(ctx, logLevelName)
	server.getLogger().SetLogLevel(logLevelName)
	server.isTrace = (logLevelName == logging.LevelTraceName)
	return err
}

func (server *G2ProductServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(5, observer.GetObserverId(ctx))
		defer func() { server.traceExit(6, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2product := getG2product()
	return g2product.UnregisterObserver(ctx, observer)
}

func (server *G2ProductServer) ValidateLicenseFile(ctx context.Context, request *g2pb.ValidateLicenseFileRequest) (*g2pb.ValidateLicenseFileResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(14, request)
		defer func() { server.traceExit(16, request, result, err, time.Since(entryTime)) }()
	}
	g2product := getG2product()
	result, err = g2product.ValidateLicenseFile(ctx, request.GetLicenseFilePath())
	response := g2pb.ValidateLicenseFileResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ProductServer) ValidateLicenseStringBase64(ctx context.Context, request *g2pb.ValidateLicenseStringBase64Request) (*g2pb.ValidateLicenseStringBase64Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(17, request)
		defer func() { server.traceExit(18, request, result, err, time.Since(entryTime)) }()
	}
	g2product := getG2product()
	result, err = g2product.ValidateLicenseStringBase64(ctx, request.GetLicenseString())
	response := g2pb.ValidateLicenseStringBase64Response{
		Result: result,
	}
	return &response, err
}

func (server *G2ProductServer) Version(ctx context.Context, request *g2pb.VersionRequest) (*g2pb.VersionResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, result, err, time.Since(entryTime)) }()
	}
	g2product := getG2product()
	result, err = g2product.Version(ctx)
	response := g2pb.VersionResponse{
		Result: result,
	}
	return &response, err
}

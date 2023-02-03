package g2productserver

import (
	"context"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go/g2product"
	pb "github.com/senzing/g2-sdk-proto/go/g2product"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
)

var (
	g2productSingleton *g2sdk.G2productImpl
	g2productSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2product() *g2sdk.G2productImpl {
	g2productSyncOnce.Do(func() {
		g2productSingleton = &g2sdk.G2productImpl{}
	})
	return g2productSingleton
}

func GetSdkG2product() *g2sdk.G2productImpl {
	return getG2product()
}

// Get the Logger singleton.
func (server *G2ProductServer) getLogger() messagelogger.MessageLoggerInterface {
	if server.logger == nil {
		server.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return server.logger
}

// Trace method entry.
func (server *G2ProductServer) traceEntry(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (server *G2ProductServer) traceExit(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing/g2-sdk-go/g2product.G2product
// ----------------------------------------------------------------------------

func (server *G2ProductServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	if server.isTrace {
		server.traceEntry(3, request)
	}
	entryTime := time.Now()
	// g2product := getG2product()
	// err := g2product.Destroy(ctx)
	err := server.getLogger().Error(4001)
	response := pb.DestroyResponse{}
	if server.isTrace {
		defer server.traceExit(4, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ProductServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	if server.isTrace {
		server.traceEntry(9, request)
	}
	entryTime := time.Now()
	// g2product := getG2product()
	// err := g2product.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err := server.getLogger().Error(4002)
	response := pb.InitResponse{}
	if server.isTrace {
		defer server.traceExit(10, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ProductServer) License(ctx context.Context, request *pb.LicenseRequest) (*pb.LicenseResponse, error) {
	if server.isTrace {
		server.traceEntry(11, request)
	}
	entryTime := time.Now()
	g2product := getG2product()
	result, err := g2product.License(ctx)
	response := pb.LicenseResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(12, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ProductServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	g2product := getG2product()
	return g2product.RegisterObserver(ctx, observer)
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (server *G2ProductServer) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if server.isTrace {
		server.traceEntry(13, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	g2product := getG2product()
	g2product.SetLogLevel(ctx, logLevel)
	server.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	server.isTrace = (server.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if server.isTrace {
		defer server.traceExit(14, logLevel, err, time.Since(entryTime))
	}
	return err
}

func (server *G2ProductServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	g2product := getG2product()
	return g2product.UnregisterObserver(ctx, observer)
}

func (server *G2ProductServer) ValidateLicenseFile(ctx context.Context, request *pb.ValidateLicenseFileRequest) (*pb.ValidateLicenseFileResponse, error) {
	if server.isTrace {
		server.traceEntry(14, request)
	}
	entryTime := time.Now()
	g2product := getG2product()
	result, err := g2product.ValidateLicenseFile(ctx, request.GetLicenseFilePath())
	response := pb.ValidateLicenseFileResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(16, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ProductServer) ValidateLicenseStringBase64(ctx context.Context, request *pb.ValidateLicenseStringBase64Request) (*pb.ValidateLicenseStringBase64Response, error) {
	if server.isTrace {
		server.traceEntry(17, request)
	}
	entryTime := time.Now()
	g2product := getG2product()
	result, err := g2product.ValidateLicenseStringBase64(ctx, request.GetLicenseString())
	response := pb.ValidateLicenseStringBase64Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(18, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ProductServer) Version(ctx context.Context, request *pb.VersionRequest) (*pb.VersionResponse, error) {
	if server.isTrace {
		server.traceEntry(19, request)
	}
	entryTime := time.Now()
	g2product := getG2product()
	result, err := g2product.Version(ctx)
	response := pb.VersionResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(20, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

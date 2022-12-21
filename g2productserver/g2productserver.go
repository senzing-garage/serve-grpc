package g2productserver

import (
	"context"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go/g2product"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/go-servegrpc/protobuf/g2product"
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
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2ProductServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	g2product := getG2product()
	err := g2product.Destroy(ctx)
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2ProductServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	g2product := getG2product()
	err := g2product.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2ProductServer) License(ctx context.Context, request *pb.LicenseRequest) (*pb.LicenseResponse, error) {
	g2product := getG2product()
	result, err := g2product.License(ctx)
	response := pb.LicenseResponse{
		Result: result,
	}
	return &response, err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (server *G2ProductServer) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if server.isTrace {
		server.traceEntry(1, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	g2product := getG2product()
	g2product.SetLogLevel(ctx, logLevel)
	server.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	server.isTrace = (server.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if server.isTrace {
		defer server.traceExit(1, logLevel, err, time.Since(entryTime))
	}
	return err
}

func (server *G2ProductServer) ValidateLicenseFile(ctx context.Context, request *pb.ValidateLicenseFileRequest) (*pb.ValidateLicenseFileResponse, error) {
	g2product := getG2product()
	result, err := g2product.ValidateLicenseFile(ctx, request.GetLicenseFilePath())
	response := pb.ValidateLicenseFileResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ProductServer) ValidateLicenseStringBase64(ctx context.Context, request *pb.ValidateLicenseStringBase64Request) (*pb.ValidateLicenseStringBase64Response, error) {
	g2product := getG2product()
	result, err := g2product.ValidateLicenseStringBase64(ctx, request.GetLicenseString())
	response := pb.ValidateLicenseStringBase64Response{
		Result: result,
	}
	return &response, err
}

func (server *G2ProductServer) Version(ctx context.Context, request *pb.VersionRequest) (*pb.VersionResponse, error) {
	g2product := getG2product()
	result, err := g2product.Version(ctx)
	response := pb.VersionResponse{
		Result: result,
	}
	return &response, err
}

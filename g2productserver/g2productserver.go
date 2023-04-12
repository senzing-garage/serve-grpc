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

// Log message.
func (server *G2ProductServer) log(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (server *G2ProductServer) traceEntry(errorNumber int, details ...interface{}) {
	server.log(errorNumber, details...)
}

// Trace method exit.
func (server *G2ProductServer) traceExit(errorNumber int, details ...interface{}) {
	server.log(errorNumber, details...)
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
	if server.isTrace {
		server.traceEntry(3, request)
	}
	entryTime := time.Now()
	// g2product := getG2product()
	// err := g2product.Destroy(ctx)
	err := server.error(4001)
	response := g2pb.DestroyResponse{}
	if server.isTrace {
		defer server.traceExit(4, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ProductServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	if server.isTrace {
		server.traceEntry(9, request)
	}
	entryTime := time.Now()
	// g2product := getG2product()
	// err := g2product.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err := server.error(4002)
	response := g2pb.InitResponse{}
	if server.isTrace {
		defer server.traceExit(10, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ProductServer) License(ctx context.Context, request *g2pb.LicenseRequest) (*g2pb.LicenseResponse, error) {
	if server.isTrace {
		server.traceEntry(11, request)
	}
	entryTime := time.Now()
	g2product := getG2product()
	result, err := g2product.License(ctx)
	response := g2pb.LicenseResponse{
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

func (server *G2ProductServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	if server.isTrace {
		server.traceEntry(13, logLevelName)
	}
	entryTime := time.Now()
	var err error = nil
	if logging.IsValidLogLevelName(logLevelName) {
		g2product := getG2product()

		// TODO: Remove once g2configmgr.SetLogLevel(context.Context, string)
		logLevel := logging.TextToLoggerLevelMap[logLevelName]

		g2product.SetLogLevel(ctx, logLevel)
		server.getLogger().SetLogLevel(logLevelName)
		server.isTrace = (logLevelName == logging.LevelTraceName)
	} else {
		err = fmt.Errorf("invalid error level: %s", logLevelName)
	}
	if server.isTrace {
		defer server.traceExit(14, logLevelName, err, time.Since(entryTime))
	}
	return err
}

func (server *G2ProductServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	g2product := getG2product()
	return g2product.UnregisterObserver(ctx, observer)
}

func (server *G2ProductServer) ValidateLicenseFile(ctx context.Context, request *g2pb.ValidateLicenseFileRequest) (*g2pb.ValidateLicenseFileResponse, error) {
	if server.isTrace {
		server.traceEntry(14, request)
	}
	entryTime := time.Now()
	g2product := getG2product()
	result, err := g2product.ValidateLicenseFile(ctx, request.GetLicenseFilePath())
	response := g2pb.ValidateLicenseFileResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(16, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ProductServer) ValidateLicenseStringBase64(ctx context.Context, request *g2pb.ValidateLicenseStringBase64Request) (*g2pb.ValidateLicenseStringBase64Response, error) {
	if server.isTrace {
		server.traceEntry(17, request)
	}
	entryTime := time.Now()
	g2product := getG2product()
	result, err := g2product.ValidateLicenseStringBase64(ctx, request.GetLicenseString())
	response := g2pb.ValidateLicenseStringBase64Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(18, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ProductServer) Version(ctx context.Context, request *g2pb.VersionRequest) (*g2pb.VersionResponse, error) {
	if server.isTrace {
		server.traceEntry(19, request)
	}
	entryTime := time.Now()
	g2product := getG2product()
	result, err := g2product.Version(ctx)
	response := g2pb.VersionResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(20, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

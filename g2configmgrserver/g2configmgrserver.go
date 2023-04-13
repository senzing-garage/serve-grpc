package g2configmgrserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go-base/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2api"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/observer"
)

var (
	g2configmgrSingleton g2api.G2configmgr
	g2configmgrSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *G2ConfigmgrServer) getLogger() logging.LoggingInterface {
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
func (server *G2ConfigmgrServer) log(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (server *G2ConfigmgrServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *G2ConfigmgrServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (server *G2ConfigmgrServer) error(messageNumber int, details ...interface{}) error {
	return server.getLogger().Error(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

// Singleton pattern for g2configmgr.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2configmgr() g2api.G2configmgr {
	g2configmgrSyncOnce.Do(func() {
		g2configmgrSingleton = &g2sdk.G2configmgr{}
	})
	return g2configmgrSingleton
}

func GetSdkG2configmgr() g2api.G2configmgr {
	return getG2configmgr()
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing/g2-sdk-go/g2configmgr.G2configmgr
// ----------------------------------------------------------------------------

func (server *G2ConfigmgrServer) AddConfig(ctx context.Context, request *g2pb.AddConfigRequest) (*g2pb.AddConfigResponse, error) {
	if server.isTrace {
		server.traceEntry(1, request)
	}
	entryTime := time.Now()
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.AddConfig(ctx, request.GetConfigStr(), request.GetConfigComments())
	response := g2pb.AddConfigResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(2, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigmgrServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
	if server.isTrace {
		server.traceEntry(5, request)
	}
	entryTime := time.Now()
	// g2configmgr := getG2configmgr()
	// err := g2configmgr.Destroy(ctx)
	err := server.error(4001)
	response := g2pb.DestroyResponse{}
	if server.isTrace {
		defer server.traceExit(6, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfig(ctx context.Context, request *g2pb.GetConfigRequest) (*g2pb.GetConfigResponse, error) {
	if server.isTrace {
		server.traceEntry(7, request)
	}
	entryTime := time.Now()
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.GetConfig(ctx, request.GetConfigID())
	response := g2pb.GetConfigResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(8, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfigList(ctx context.Context, request *g2pb.GetConfigListRequest) (*g2pb.GetConfigListResponse, error) {
	if server.isTrace {
		server.traceEntry(9, request)
	}
	entryTime := time.Now()
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.GetConfigList(ctx)
	response := g2pb.GetConfigListResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(10, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetDefaultConfigID(ctx context.Context, request *g2pb.GetDefaultConfigIDRequest) (*g2pb.GetDefaultConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(11, request)
	}
	entryTime := time.Now()
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.GetDefaultConfigID(ctx)
	response := g2pb.GetDefaultConfigIDResponse{
		ConfigID: result,
	}
	if server.isTrace {
		defer server.traceExit(12, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigmgrServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	if server.isTrace {
		server.traceEntry(17, request)
	}
	entryTime := time.Now()
	// g2configmgr := getG2configmgr()
	// err := g2configmgr.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err := server.error(4002)
	response := g2pb.InitResponse{}
	if server.isTrace {
		defer server.traceExit(18, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigmgrServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	g2configmgr := getG2configmgr()
	return g2configmgr.RegisterObserver(ctx, observer)
}

func (server *G2ConfigmgrServer) ReplaceDefaultConfigID(ctx context.Context, request *g2pb.ReplaceDefaultConfigIDRequest) (*g2pb.ReplaceDefaultConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(19, request)
	}
	entryTime := time.Now()
	g2configmgr := getG2configmgr()
	err := g2configmgr.ReplaceDefaultConfigID(ctx, request.GetOldConfigID(), request.GetNewConfigID())
	response := g2pb.ReplaceDefaultConfigIDResponse{}
	if server.isTrace {
		defer server.traceExit(20, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigmgrServer) SetDefaultConfigID(ctx context.Context, request *g2pb.SetDefaultConfigIDRequest) (*g2pb.SetDefaultConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(21, request)
	}
	entryTime := time.Now()
	g2configmgr := getG2configmgr()
	err := g2configmgr.SetDefaultConfigID(ctx, request.GetConfigID())
	response := g2pb.SetDefaultConfigIDResponse{}
	if server.isTrace {
		defer server.traceExit(22, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigmgrServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	if server.isTrace {
		server.traceEntry(23, logLevelName)
	}
	entryTime := time.Now()
	var err error = nil
	if logging.IsValidLogLevelName(logLevelName) {
		g2configmgr := getG2configmgr()

		// TODO: Remove once g2configmgr.SetLogLevel(context.Context, string)
		logLevel := logging.TextToLoggerLevelMap[logLevelName]

		g2configmgr.SetLogLevel(ctx, logLevel)
		server.getLogger().SetLogLevel(logLevelName)
		server.isTrace = (logLevelName == logging.LevelTraceName)
	} else {
		err = fmt.Errorf("invalid error level: %s", logLevelName)
	}
	if server.isTrace {
		defer server.traceExit(24, logLevelName, err, time.Since(entryTime))
	}
	return err
}

func (server *G2ConfigmgrServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	g2configmgr := getG2configmgr()
	return g2configmgr.UnregisterObserver(ctx, observer)
}

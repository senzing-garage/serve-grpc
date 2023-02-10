package g2configmgrserver

import (
	"context"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go-base/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2api"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
)

var (
	g2configmgrSingleton g2api.G2configmgrInterface
	g2configmgrSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2configmgr.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2configmgr() g2api.G2configmgrInterface {
	g2configmgrSyncOnce.Do(func() {
		g2configmgrSingleton = &g2sdk.G2configmgr{}
	})
	return g2configmgrSingleton
}

func GetSdkG2configmgr() g2api.G2configmgrInterface {
	return getG2configmgr()
}

// Get the Logger singleton.
func (server *G2ConfigmgrServer) getLogger() messagelogger.MessageLoggerInterface {
	if server.logger == nil {
		server.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return server.logger
}

// Trace method entry.
func (server *G2ConfigmgrServer) traceEntry(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (server *G2ConfigmgrServer) traceExit(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
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
	err := server.getLogger().Error(4001)
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
	err := server.getLogger().Error(4002)
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

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (server *G2ConfigmgrServer) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if server.isTrace {
		server.traceEntry(23, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	g2configmgr := getG2configmgr()
	g2configmgr.SetLogLevel(ctx, logLevel)
	server.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	server.isTrace = (server.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if server.isTrace {
		defer server.traceExit(24, logLevel, err, time.Since(entryTime))
	}
	return err
}

func (server *G2ConfigmgrServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	g2configmgr := getG2configmgr()
	return g2configmgr.UnregisterObserver(ctx, observer)
}

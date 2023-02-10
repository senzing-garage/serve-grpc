package g2configserver

import (
	"context"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go-base/g2config"
	"github.com/senzing/g2-sdk-go/g2api"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
)

var (
	g2configSingleton g2api.G2configInterface
	g2configSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2config() g2api.G2configInterface {
	g2configSyncOnce.Do(func() {
		g2configSingleton = &g2sdk.G2config{}
	})
	return g2configSingleton
}

func GetSdkG2config() g2api.G2configInterface {
	return getG2config()
}

// Get the Logger singleton.
func (server *G2ConfigServer) getLogger() messagelogger.MessageLoggerInterface {
	if server.logger == nil {
		server.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return server.logger
}

// Trace method entry.
func (server *G2ConfigServer) traceEntry(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (server *G2ConfigServer) traceExit(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing/g2-sdk-go/g2config.G2config
// ----------------------------------------------------------------------------

func (server *G2ConfigServer) AddDataSource(ctx context.Context, request *g2pb.AddDataSourceRequest) (*g2pb.AddDataSourceResponse, error) {
	if server.isTrace {
		server.traceEntry(1, request)
	}
	entryTime := time.Now()
	g2config := getG2config()
	result, err := g2config.AddDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetInputJson())
	response := g2pb.AddDataSourceResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(2, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigServer) Close(ctx context.Context, request *g2pb.CloseRequest) (*g2pb.CloseResponse, error) {
	if server.isTrace {
		server.traceEntry(5, request)
	}
	entryTime := time.Now()
	g2config := getG2config()
	err := g2config.Close(ctx, uintptr(request.GetConfigHandle()))
	response := g2pb.CloseResponse{}
	if server.isTrace {
		defer server.traceExit(6, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigServer) Create(ctx context.Context, request *g2pb.CreateRequest) (*g2pb.CreateResponse, error) {
	if server.isTrace {
		server.traceEntry(7, request)
	}
	entryTime := time.Now()
	g2config := getG2config()
	result, err := g2config.Create(ctx)
	response := g2pb.CreateResponse{
		Result: int64(result),
	}
	if server.isTrace {
		defer server.traceExit(8, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigServer) DeleteDataSource(ctx context.Context, request *g2pb.DeleteDataSourceRequest) (*g2pb.DeleteDataSourceResponse, error) {
	if server.isTrace {
		server.traceEntry(9, request)
	}
	entryTime := time.Now()
	g2config := getG2config()
	err := g2config.DeleteDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetInputJson())
	response := g2pb.DeleteDataSourceResponse{}
	if server.isTrace {
		defer server.traceExit(10, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
	if server.isTrace {
		server.traceEntry(11, request)
	}
	entryTime := time.Now()
	// g2config := getG2config()
	// err := g2config.Destroy(ctx)
	err := server.getLogger().Error(4001)
	response := g2pb.DestroyResponse{}
	if server.isTrace {
		defer server.traceExit(12, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	if server.isTrace {
		server.traceEntry(17, request)
	}
	entryTime := time.Now()
	// g2config := getG2config()
	// err := g2config.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err := server.getLogger().Error(4002)
	response := g2pb.InitResponse{}
	if server.isTrace {
		defer server.traceExit(18, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigServer) ListDataSources(ctx context.Context, request *g2pb.ListDataSourcesRequest) (*g2pb.ListDataSourcesResponse, error) {
	if server.isTrace {
		server.traceEntry(19, request)
	}
	entryTime := time.Now()
	g2config := getG2config()
	result, err := g2config.ListDataSources(ctx, uintptr(request.GetConfigHandle()))
	response := g2pb.ListDataSourcesResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(20, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigServer) Load(ctx context.Context, request *g2pb.LoadRequest) (*g2pb.LoadResponse, error) {
	if server.isTrace {
		server.traceEntry(21, request)
	}
	entryTime := time.Now()
	g2config := getG2config()
	err := g2config.Load(ctx, uintptr(request.GetConfigHandle()), (request.GetJsonConfig()))
	response := g2pb.LoadResponse{}
	if server.isTrace {
		defer server.traceExit(22, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	g2config := getG2config()
	return g2config.RegisterObserver(ctx, observer)
}

func (server *G2ConfigServer) Save(ctx context.Context, request *g2pb.SaveRequest) (*g2pb.SaveResponse, error) {
	if server.isTrace {
		server.traceEntry(23, request)
	}
	entryTime := time.Now()
	g2config := getG2config()
	result, err := g2config.Save(ctx, uintptr(request.GetConfigHandle()))
	response := g2pb.SaveResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(24, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (server *G2ConfigServer) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if server.isTrace {
		server.traceEntry(25, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	g2config := getG2config()
	g2config.SetLogLevel(ctx, logLevel)
	server.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	server.isTrace = (server.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if server.isTrace {
		defer server.traceExit(26, logLevel, err, time.Since(entryTime))
	}
	return err
}

func (server *G2ConfigServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	g2config := getG2config()
	return g2config.UnregisterObserver(ctx, observer)
}

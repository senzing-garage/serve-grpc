package g2configserver

import (
	"context"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go/g2config"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/go-servegrpc/protobuf/g2config"
)

var (
	g2configSingleton *g2sdk.G2configImpl
	g2configSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2config() *g2sdk.G2configImpl {
	g2configSyncOnce.Do(func() {
		g2configSingleton = &g2sdk.G2configImpl{}
	})
	return g2configSingleton
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
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2ConfigServer) AddDataSource(ctx context.Context, request *pb.AddDataSourceRequest) (*pb.AddDataSourceResponse, error) {
	g2config := getG2config()
	result, err := g2config.AddDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetInputJson())
	response := pb.AddDataSourceResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigServer) Close(ctx context.Context, request *pb.CloseRequest) (*pb.CloseResponse, error) {
	g2config := getG2config()
	err := g2config.Close(ctx, uintptr(request.GetConfigHandle()))
	response := pb.CloseResponse{}
	return &response, err
}

func (server *G2ConfigServer) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	if server.isTrace {
		server.traceEntry(1, request)
	}
	entryTime := time.Now()
	g2config := getG2config()
	result, err := g2config.Create(ctx)
	response := pb.CreateResponse{
		Result: int64(result),
	}
	if server.isTrace {
		defer server.traceExit(6, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2ConfigServer) DeleteDataSource(ctx context.Context, request *pb.DeleteDataSourceRequest) (*pb.DeleteDataSourceResponse, error) {
	g2config := getG2config()
	err := g2config.DeleteDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetInputJson())
	response := pb.DeleteDataSourceResponse{}
	return &response, err
}

func (server *G2ConfigServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	g2config := getG2config()
	err := g2config.Destroy(ctx)
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2ConfigServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	g2config := getG2config()
	err := g2config.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2ConfigServer) ListDataSources(ctx context.Context, request *pb.ListDataSourcesRequest) (*pb.ListDataSourcesResponse, error) {
	g2config := getG2config()
	result, err := g2config.ListDataSources(ctx, uintptr(request.GetConfigHandle()))
	response := pb.ListDataSourcesResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigServer) Load(ctx context.Context, request *pb.LoadRequest) (*pb.LoadResponse, error) {
	g2config := getG2config()
	err := g2config.Load(ctx, uintptr(request.GetConfigHandle()), (request.GetJsonConfig()))
	response := pb.LoadResponse{}
	return &response, err
}

func (server *G2ConfigServer) Save(ctx context.Context, request *pb.SaveRequest) (*pb.SaveResponse, error) {
	g2config := getG2config()
	result, err := g2config.Save(ctx, uintptr(request.GetConfigHandle()))
	response := pb.SaveResponse{
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
func (server *G2ConfigServer) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if server.isTrace {
		server.traceEntry(1, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	server.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	server.isTrace = (server.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if server.isTrace {
		defer server.traceExit(1, logLevel, err, time.Since(entryTime))
	}
	return err
}

package g2configmgrserver

import (
	"context"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go/g2configmgr"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/go-servegrpc/protobuf/g2configmgr"
)

var (
	g2configmgrSingleton *g2sdk.G2configmgrImpl
	g2configmgrSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2configmgr.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2configmgr() *g2sdk.G2configmgrImpl {
	g2configmgrSyncOnce.Do(func() {
		g2configmgrSingleton = &g2sdk.G2configmgrImpl{}
	})
	return g2configmgrSingleton
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
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2ConfigmgrServer) AddConfig(ctx context.Context, request *pb.AddConfigRequest) (*pb.AddConfigResponse, error) {
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.AddConfig(ctx, request.GetConfigStr(), request.GetConfigComments())
	response := pb.AddConfigResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	g2configmgr := getG2configmgr()
	err := g2configmgr.Destroy(ctx)
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfig(ctx context.Context, request *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.GetConfig(ctx, request.GetConfigID())
	response := pb.GetConfigResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfigList(ctx context.Context, request *pb.GetConfigListRequest) (*pb.GetConfigListResponse, error) {
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.GetConfigList(ctx)
	response := pb.GetConfigListResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetDefaultConfigID(ctx context.Context, request *pb.GetDefaultConfigIDRequest) (*pb.GetDefaultConfigIDResponse, error) {
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.GetDefaultConfigID(ctx)
	response := pb.GetDefaultConfigIDResponse{
		ConfigID: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	g2configmgr := getG2configmgr()
	err := g2configmgr.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) ReplaceDefaultConfigID(ctx context.Context, request *pb.ReplaceDefaultConfigIDRequest) (*pb.ReplaceDefaultConfigIDResponse, error) {
	g2configmgr := getG2configmgr()
	err := g2configmgr.ReplaceDefaultConfigID(ctx, request.GetOldConfigID(), request.GetNewConfigID())
	response := pb.ReplaceDefaultConfigIDResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) SetDefaultConfigID(ctx context.Context, request *pb.SetDefaultConfigIDRequest) (*pb.SetDefaultConfigIDResponse, error) {
	g2configmgr := getG2configmgr()
	err := g2configmgr.SetDefaultConfigID(ctx, request.GetConfigID())
	response := pb.SetDefaultConfigIDResponse{}
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

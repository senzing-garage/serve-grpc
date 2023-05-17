package g2configserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go-base/g2config"
	"github.com/senzing/g2-sdk-go/g2api"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/observer"
)

var (
	g2configSingleton g2api.G2config
	g2configSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *G2ConfigServer) getLogger() logging.LoggingInterface {
	var err error = nil
	if server.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 3},
		}
		server.logger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return server.logger
}

// Trace method entry.
func (server *G2ConfigServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *G2ConfigServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (server *G2ConfigServer) error(messageNumber int, details ...interface{}) error {
	return server.getLogger().Error(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2config() g2api.G2config {
	g2configSyncOnce.Do(func() {
		g2configSingleton = &g2sdk.G2config{}
	})
	return g2configSingleton
}

func GetSdkG2config() g2api.G2config {
	return getG2config()
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing/g2-sdk-go/g2config.G2config
// ----------------------------------------------------------------------------

func (server *G2ConfigServer) AddDataSource(ctx context.Context, request *g2pb.AddDataSourceRequest) (*g2pb.AddDataSourceResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	result, err = g2config.AddDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetInputJson())
	response := g2pb.AddDataSourceResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigServer) Close(ctx context.Context, request *g2pb.CloseRequest) (*g2pb.CloseResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(5, request)
		defer func() { server.traceExit(6, request, err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	err = g2config.Close(ctx, uintptr(request.GetConfigHandle()))
	response := g2pb.CloseResponse{}
	return &response, err
}

func (server *G2ConfigServer) Create(ctx context.Context, request *g2pb.CreateRequest) (*g2pb.CreateResponse, error) {
	var err error = nil
	var result uintptr
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(7, request)
		defer func() { server.traceExit(8, request, result, err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	result, err = g2config.Create(ctx)
	response := g2pb.CreateResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *G2ConfigServer) DeleteDataSource(ctx context.Context, request *g2pb.DeleteDataSourceRequest) (*g2pb.DeleteDataSourceResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(9, request)
		defer func() { server.traceExit(10, request, err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	err = g2config.DeleteDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetInputJson())
	response := g2pb.DeleteDataSourceResponse{}
	return &response, err
}

func (server *G2ConfigServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(11, request)
		defer func() { server.traceExit(12, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2config := getG2config()
	// err := g2config.Destroy(ctx)
	err = server.error(4001)
	response := g2pb.DestroyResponse{}
	return &response, err
}

func (server *G2ConfigServer) GetObserverOrigin(ctx context.Context) string {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(27)
		defer func() { server.traceExit(28, err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	return g2config.GetObserverOrigin(ctx)
}

func (server *G2ConfigServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(17, request)
		defer func() { server.traceExit(18, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2config := getG2config()
	// err := g2config.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err = server.error(4002)
	response := g2pb.InitResponse{}
	return &response, err
}

func (server *G2ConfigServer) ListDataSources(ctx context.Context, request *g2pb.ListDataSourcesRequest) (*g2pb.ListDataSourcesResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, result, err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	result, err = g2config.ListDataSources(ctx, uintptr(request.GetConfigHandle()))
	response := g2pb.ListDataSourcesResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigServer) Load(ctx context.Context, request *g2pb.LoadRequest) (*g2pb.LoadResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(21, request)
		defer func() { server.traceExit(22, request, err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	err = g2config.Load(ctx, uintptr(request.GetConfigHandle()), (request.GetJsonConfig()))
	response := g2pb.LoadResponse{}
	return &response, err
}

func (server *G2ConfigServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(3, observer.GetObserverId(ctx))
		defer func() { server.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	return g2config.RegisterObserver(ctx, observer)
}

func (server *G2ConfigServer) Save(ctx context.Context, request *g2pb.SaveRequest) (*g2pb.SaveResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(23, request)
		defer func() { server.traceExit(24, request, result, err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	result, err = g2config.Save(ctx, uintptr(request.GetConfigHandle()))
	response := g2pb.SaveResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(25, logLevelName)
		defer func() { server.traceExit(26, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	g2config := getG2config()
	g2config.SetLogLevel(ctx, logLevelName)
	server.getLogger().SetLogLevel(logLevelName)
	server.isTrace = (logLevelName == logging.LevelTraceName)
	return err
}

func (server *G2ConfigServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(29, origin)
		defer func() { server.traceExit(30, origin, err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	g2config.SetObserverOrigin(ctx, origin)
}

func (server *G2ConfigServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(13, observer.GetObserverId(ctx))
		defer func() { server.traceExit(14, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2config := getG2config()
	return g2config.UnregisterObserver(ctx, observer)
}

package g2configmgrserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	g2sdk "github.com/senzing-garage/g2-sdk-go-base/g2configmgr"
	"github.com/senzing-garage/g2-sdk-go/g2api"
	g2pb "github.com/senzing-garage/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
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
		server.logger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return server.logger
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
	return server.getLogger().NewError(messageNumber, details...)
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
// Interface methods for github.com/senzing-garage/g2-sdk-go/g2configmgr.G2configmgr
// ----------------------------------------------------------------------------

func (server *G2ConfigmgrServer) AddConfig(ctx context.Context, request *g2pb.AddConfigRequest) (*g2pb.AddConfigResponse, error) {
	var err error = nil
	var result int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	result, err = g2configmgr.AddConfig(ctx, request.GetConfigStr(), request.GetConfigComments())
	response := g2pb.AddConfigResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(5, request)
		defer func() { server.traceExit(6, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2configmgr := getG2configmgr()
	// err := g2configmgr.Destroy(ctx)
	err = server.error(4001)
	response := g2pb.DestroyResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfig(ctx context.Context, request *g2pb.GetConfigRequest) (*g2pb.GetConfigResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(7, request)
		defer func() { server.traceExit(8, request, result, err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	result, err = g2configmgr.GetConfig(ctx, request.GetConfigID())
	response := g2pb.GetConfigResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfigList(ctx context.Context, request *g2pb.GetConfigListRequest) (*g2pb.GetConfigListResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(9, request)
		defer func() { server.traceExit(10, request, result, err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	result, err = g2configmgr.GetConfigList(ctx)
	response := g2pb.GetConfigListResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetDefaultConfigID(ctx context.Context, request *g2pb.GetDefaultConfigIDRequest) (*g2pb.GetDefaultConfigIDResponse, error) {
	var err error = nil
	var result int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(11, request)
		defer func() { server.traceExit(12, request, result, err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	result, err = g2configmgr.GetDefaultConfigID(ctx)
	response := g2pb.GetDefaultConfigIDResponse{
		ConfigID: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetObserverOrigin(ctx context.Context) string {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(25)
		defer func() { server.traceExit(26, err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	return g2configmgr.GetObserverOrigin(ctx)
}

func (server *G2ConfigmgrServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(17, request)
		defer func() { server.traceExit(18, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2configmgr := getG2configmgr()
	// err := g2configmgr.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err = server.error(4002)
	response := g2pb.InitResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(3, observer.GetObserverId(ctx))
		defer func() { server.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	return g2configmgr.RegisterObserver(ctx, observer)
}

func (server *G2ConfigmgrServer) ReplaceDefaultConfigID(ctx context.Context, request *g2pb.ReplaceDefaultConfigIDRequest) (*g2pb.ReplaceDefaultConfigIDResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	err = g2configmgr.ReplaceDefaultConfigID(ctx, request.GetOldConfigID(), request.GetNewConfigID())
	response := g2pb.ReplaceDefaultConfigIDResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) SetDefaultConfigID(ctx context.Context, request *g2pb.SetDefaultConfigIDRequest) (*g2pb.SetDefaultConfigIDResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(21, request)
		defer func() { server.traceExit(22, request, err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	err = g2configmgr.SetDefaultConfigID(ctx, request.GetConfigID())
	response := g2pb.SetDefaultConfigIDResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(23, logLevelName)
		defer func() { server.traceExit(24, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	g2configmgr := getG2configmgr()
	err = g2configmgr.SetLogLevel(ctx, logLevelName)
	if err != nil {
		return err
	}
	err = server.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		return err
	}
	server.isTrace = (logLevelName == logging.LevelTraceName)
	return err
}

func (server *G2ConfigmgrServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(27, origin)
		defer func() { server.traceExit(28, origin, err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	g2configmgr.SetObserverOrigin(ctx, origin)
}

func (server *G2ConfigmgrServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(13, observer.GetObserverId(ctx))
		defer func() { server.traceExit(14, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2configmgr := getG2configmgr()
	return g2configmgr.UnregisterObserver(ctx, observer)
}

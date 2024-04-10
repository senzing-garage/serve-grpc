package szconfigserver

import (
	"context"
	"sync"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	g2sdk "github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go/sz"
	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
)

var (
	szConfigSingleton sz.SzConfig
	szConfigSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *SzConfigServer) getLogger() logging.LoggingInterface {
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
func (server *SzConfigServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *SzConfigServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (server *SzConfigServer) error(messageNumber int, details ...interface{}) error {
	return server.getLogger().NewError(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

// Singleton pattern for szconfig.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getSzConfig() sz.SzConfig {
	szConfigSyncOnce.Do(func() {
		szConfigSingleton = &g2sdk.Szconfig{}
	})
	return szConfigSingleton
}

func GetSdkG2config() sz.SzConfig {
	return getSzConfig()
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/g2-sdk-go/szconfig.G2config
// ----------------------------------------------------------------------------

func (server *SzConfigServer) AddDataSource(ctx context.Context, request *g2pb.AddDataSourceRequest) (*g2pb.AddDataSourceResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err = szConfig.AddDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetDataSourceCode())
	response := g2pb.AddDataSourceResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzConfigServer) CloseConfig(ctx context.Context, request *g2pb.CloseConfigRequest) (*g2pb.CloseConfigResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(5, request)
		defer func() { server.traceExit(6, request, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	err = szConfig.CloseConfig(ctx, uintptr(request.GetConfigHandle()))
	response := g2pb.CloseConfigResponse{}
	return &response, err
}

func (server *SzConfigServer) CreateConfig(ctx context.Context, request *g2pb.CreateConfigRequest) (*g2pb.CreateConfigResponse, error) {
	var err error = nil
	var result uintptr
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(7, request)
		defer func() { server.traceExit(8, request, result, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err = szConfig.CreateConfig(ctx)
	response := g2pb.CreateConfigResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *SzConfigServer) DeleteDataSource(ctx context.Context, request *g2pb.DeleteDataSourceRequest) (*g2pb.DeleteDataSourceResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(9, request)
		defer func() { server.traceExit(10, request, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	err = szConfig.DeleteDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetDataSourceCode())
	response := g2pb.DeleteDataSourceResponse{}
	return &response, err
}

// func (server *SzConfigServer) GetObserverOrigin(ctx context.Context) string {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(27)
// 		defer func() { server.traceExit(28, err, time.Since(entryTime)) }()
// 	}
// 	szconfig := getSzConfig()
// 	return szconfig.GetObserverOrigin(ctx)
// }

func (server *SzConfigServer) GetDataSources(ctx context.Context, request *g2pb.GetDataSourcesRequest) (*g2pb.GetDataSourcesResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, result, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err = szConfig.GetDataSources(ctx, uintptr(request.GetConfigHandle()))
	response := g2pb.GetDataSourcesResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzConfigServer) ImportConfig(ctx context.Context, request *g2pb.ImportConfigRequest) (*g2pb.ImportConfigResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(21, request)
		defer func() { server.traceExit(22, request, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err := szConfig.ImportConfig(ctx, request.GetJsonConfig())
	response := g2pb.ImportConfigResponse{
		Result: int64(result),
	}
	return &response, err
}

// func (server *SzConfigServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(3, observer.GetObserverId(ctx))
// 		defer func() { server.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
// 	}
// 	szconfig := getSzConfig()
// 	return szconfig.RegisterObserver(ctx, observer)
// }

func (server *SzConfigServer) ExportConfig(ctx context.Context, request *g2pb.ExportConfigRequest) (*g2pb.ExportConfigResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(23, request)
		defer func() { server.traceExit(24, request, result, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err = szConfig.ExportConfig(ctx, uintptr(request.GetConfigHandle()))
	response := g2pb.ExportConfigResponse{
		Result: result,
	}
	return &response, err
}

// func (server *SzConfigServer) SetLogLevel(ctx context.Context, logLevelName string) error {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(25, logLevelName)
// 		defer func() { server.traceExit(26, logLevelName, err, time.Since(entryTime)) }()
// 	}
// 	if !logging.IsValidLogLevelName(logLevelName) {
// 		return fmt.Errorf("invalid error level: %s", logLevelName)
// 	}
// 	szconfig := getSzConfig()
// 	err = szconfig.SetLogLevel(ctx, logLevelName)
// 	if err != nil {
// 		return err
// 	}
// 	err = server.getLogger().SetLogLevel(logLevelName)
// 	if err != nil {
// 		return err
// 	}
// 	server.isTrace = (logLevelName == logging.LevelTraceName)
// 	return err
// }

// func (server *SzConfigServer) SetObserverOrigin(ctx context.Context, origin string) {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(29, origin)
// 		defer func() { server.traceExit(30, origin, err, time.Since(entryTime)) }()
// 	}
// 	szconfig := getSzConfig()
// 	szconfig.SetObserverOrigin(ctx, origin)
// }

// func (server *SzConfigServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(13, observer.GetObserverId(ctx))
// 		defer func() { server.traceExit(14, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
// 	}
// 	szconfig := getSzConfig()
// 	return szconfig.UnregisterObserver(ctx, observer)
// }

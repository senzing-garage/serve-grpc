package szconfigserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	szsdk "github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
)

var (
	szConfigSingleton *szsdk.Szconfig
	szConfigSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/sz-sdk-go/szconfig.SzCconfig
// ----------------------------------------------------------------------------

func (server *SzConfigServer) AddDataSource(ctx context.Context, request *szpb.AddDataSourceRequest) (*szpb.AddDataSourceResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err = szConfig.AddDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetDataSourceCode())
	response := szpb.AddDataSourceResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzConfigServer) CloseConfig(ctx context.Context, request *szpb.CloseConfigRequest) (*szpb.CloseConfigResponse, error) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(5, request)
		defer func() { server.traceExit(6, request, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	err = szConfig.CloseConfig(ctx, uintptr(request.GetConfigHandle()))
	response := szpb.CloseConfigResponse{}
	return &response, err
}

func (server *SzConfigServer) CreateConfig(ctx context.Context, request *szpb.CreateConfigRequest) (*szpb.CreateConfigResponse, error) {
	var err error
	var result uintptr
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(7, request)
		defer func() { server.traceExit(8, request, result, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err = szConfig.CreateConfig(ctx)
	response := szpb.CreateConfigResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *SzConfigServer) DeleteDataSource(ctx context.Context, request *szpb.DeleteDataSourceRequest) (*szpb.DeleteDataSourceResponse, error) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(9, request)
		defer func() { server.traceExit(10, request, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	err = szConfig.DeleteDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetDataSourceCode())
	response := szpb.DeleteDataSourceResponse{}
	return &response, err
}

func (server *SzConfigServer) ExportConfig(ctx context.Context, request *szpb.ExportConfigRequest) (*szpb.ExportConfigResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(23, request)
		defer func() { server.traceExit(24, request, result, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err = szConfig.ExportConfig(ctx, uintptr(request.GetConfigHandle()))
	response := szpb.ExportConfigResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzConfigServer) GetDataSources(ctx context.Context, request *szpb.GetDataSourcesRequest) (*szpb.GetDataSourcesResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, result, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err = szConfig.GetDataSources(ctx, uintptr(request.GetConfigHandle()))
	response := szpb.GetDataSourcesResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzConfigServer) ImportConfig(ctx context.Context, request *szpb.ImportConfigRequest) (*szpb.ImportConfigResponse, error) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(21, request)
		defer func() { server.traceExit(22, request, err, time.Since(entryTime)) }()
	}
	szConfig := getSzConfig()
	result, err := szConfig.ImportConfig(ctx, request.GetConfigDefinition())
	response := szpb.ImportConfigResponse{
		Result: int64(result),
	}
	return &response, err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *SzConfigServer) getLogger() logging.Logging {
	var err error
	if server.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 3},
		}
		server.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
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

func (server *SzConfigServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	_ = ctx
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(25, logLevelName)
		defer func() { server.traceExit(26, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	// szconfig := getSzConfig()
	// err = szconfig.SetLogLevel(ctx, logLevelName)
	// if err != nil {
	// 	return err
	// }
	err = server.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		return err
	}
	server.isTrace = (logLevelName == logging.LevelTraceName)
	return err
}

// --- Errors -----------------------------------------------------------------

// Create error.
// func (server *SzConfigServer) error(messageNumber int, details ...interface{}) error {
// 	return server.getLogger().NewError(messageNumber, details...)
// }

// --- Services ---------------------------------------------------------------

// Singleton pattern for szconfig.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getSzConfig() *szsdk.Szconfig {
	szConfigSyncOnce.Do(func() {
		szConfigSingleton = &szsdk.Szconfig{}
	})
	return szConfigSingleton
}

func GetSdkSzConfig() *szsdk.Szconfig {
	return getSzConfig()
}

func GetSdkSzConfigAsInterface() senzing.SzConfig {
	return getSzConfig()
}

// --- Observer ---------------------------------------------------------------

func (server *SzConfigServer) GetObserverOrigin(ctx context.Context) string {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(27)
		defer func() { server.traceExit(28, err, time.Since(entryTime)) }()
	}
	szconfig := getSzConfig()
	return szconfig.GetObserverOrigin(ctx)
}

func (server *SzConfigServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(3, observer.GetObserverID(ctx))
		defer func() { server.traceExit(4, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	szconfig := getSzConfig()
	return szconfig.RegisterObserver(ctx, observer)
}

func (server *SzConfigServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(29, origin)
		defer func() { server.traceExit(30, origin, err, time.Since(entryTime)) }()
	}
	szconfig := getSzConfig()
	szconfig.SetObserverOrigin(ctx, origin)
}

func (server *SzConfigServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(13, observer.GetObserverID(ctx))
		defer func() { server.traceExit(14, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	szconfig := getSzConfig()
	return szconfig.UnregisterObserver(ctx, observer)
}

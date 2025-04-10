package szconfigmanagerserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	szsdk "github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
)

var (
	szConfigManagerSingleton *szsdk.Szconfigmanager
	szConfigManagerSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/sz-sdk-go/szconfigmanager
// ----------------------------------------------------------------------------

func (server *SzConfigManagerServer) GetConfig(
	ctx context.Context,
	request *szpb.GetConfigRequest,
) (*szpb.GetConfigResponse, error) {
	var (
		err      error
		response *szpb.GetConfigResponse
		result   string
	)
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(7, request)
		defer func() { server.traceExit(8, request, result, err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, request.GetConfigId())
	if err != nil {
		return response, err
	}
	result, err = szConfig.Export(ctx)
	response = &szpb.GetConfigResponse{
		Result: result,
	}
	return response, err
}

func (server *SzConfigManagerServer) GetConfigs(
	ctx context.Context,
	request *szpb.GetConfigsRequest,
) (*szpb.GetConfigsResponse, error) {
	var (
		err      error
		response *szpb.GetConfigsResponse
		result   string
	)
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(9, request)
		defer func() { server.traceExit(10, request, result, err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	result, err = szConfigManager.GetConfigs(ctx)
	response = &szpb.GetConfigsResponse{
		Result: result,
	}
	return response, err
}

func (server *SzConfigManagerServer) GetDefaultConfigId( //revive:disable var-naming
	ctx context.Context,
	request *szpb.GetDefaultConfigIdRequest,
) (*szpb.GetDefaultConfigIdResponse, error) {
	var (
		err      error
		response *szpb.GetDefaultConfigIdResponse
		result   int64
	)
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(11, request)
		defer func() { server.traceExit(12, request, result, err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	result, err = szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		return response, err
	}
	response = &szpb.GetDefaultConfigIdResponse{
		Result: result,
	}
	return response, err
}

func (server *SzConfigManagerServer) GetTemplateConfig(
	ctx context.Context,
	request *szpb.GetTemplateConfigRequest,
) (*szpb.GetTemplateConfigResponse, error) {
	var (
		err      error
		result   string
		response *szpb.GetTemplateConfigResponse
	)
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(99, request)
		defer func() { server.traceExit(99, request, result, err, time.Since(entryTime)) }()
	}

	szConfigManager := getSzConfigManager()
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		return response, err
	}

	configDefinition, err := szConfig.Export(ctx)
	response = &szpb.GetTemplateConfigResponse{
		Result: configDefinition,
	}
	return response, err
}

func (server *SzConfigManagerServer) RegisterConfig(
	ctx context.Context,
	request *szpb.RegisterConfigRequest,
) (*szpb.RegisterConfigResponse, error) {
	var (
		err      error
		response *szpb.RegisterConfigResponse
		result   int64
	)
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}

	szConfigManager := getSzConfigManager()
	result, err = szConfigManager.RegisterConfig(ctx, request.GetConfigDefinition(), request.GetConfigComment())
	response = &szpb.RegisterConfigResponse{
		Result: result,
	}
	return response, err
}

func (server *SzConfigManagerServer) ReplaceDefaultConfigId(
	ctx context.Context,
	request *szpb.ReplaceDefaultConfigIdRequest,
) (*szpb.ReplaceDefaultConfigIdResponse, error) {
	var (
		err      error
		response *szpb.ReplaceDefaultConfigIdResponse
	)
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	err = szConfigManager.ReplaceDefaultConfigID(
		ctx,
		request.GetCurrentDefaultConfigId(),
		request.GetNewDefaultConfigId(),
	)
	response = &szpb.ReplaceDefaultConfigIdResponse{}
	return response, err
}

func (server *SzConfigManagerServer) SetDefaultConfig(
	ctx context.Context,
	request *szpb.SetDefaultConfigRequest,
) (*szpb.SetDefaultConfigResponse, error) {
	var (
		err      error
		response *szpb.SetDefaultConfigResponse
	)
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(21, request)
		defer func() { server.traceExit(22, request, err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	result, err := szConfigManager.SetDefaultConfig(ctx, request.GetConfigDefinition(), request.GetConfigComment())
	response = &szpb.SetDefaultConfigResponse{
		Result: result,
	}
	return response, err
}

func (server *SzConfigManagerServer) SetDefaultConfigId(
	ctx context.Context,
	request *szpb.SetDefaultConfigIdRequest,
) (*szpb.SetDefaultConfigIdResponse, error) {
	var (
		err      error
		response *szpb.SetDefaultConfigIdResponse
	)
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(21, request)
		defer func() { server.traceExit(22, request, err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	err = szConfigManager.SetDefaultConfigID(ctx, request.GetConfigId())
	response = &szpb.SetDefaultConfigIdResponse{}
	return response, err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *SzConfigManagerServer) getLogger() logging.Logging {
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
func (server *SzConfigManagerServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *SzConfigManagerServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

func (server *SzConfigManagerServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	_ = ctx
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(23, logLevelName)
		defer func() { server.traceExit(24, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	// szConfigManager := getSzConfigManager()
	// err = szConfigManager.SetLogLevel(ctx, logLevelName)
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
// func (server *SzConfigManagerServer) error(messageNumber int, details ...interface{}) error {
// 	return server.getLogger().NewError(messageNumber, details...)
// }

// --- Services ---------------------------------------------------------------

// Singleton pattern for szconfigmanager.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getSzConfigManager() *szsdk.Szconfigmanager {
	szConfigManagerSyncOnce.Do(func() {
		szConfigManagerSingleton = &szsdk.Szconfigmanager{}
	})
	return szConfigManagerSingleton
}

func GetSdkSzConfigManager() *szsdk.Szconfigmanager {
	return getSzConfigManager()
}

func GetSdkSzConfigManagerAsInterface() senzing.SzConfigManager {
	return getSzConfigManager()
}

// --- Observer ---------------------------------------------------------------

func (server *SzConfigManagerServer) GetObserverOrigin(ctx context.Context) string {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(25)
		defer func() { server.traceExit(26, err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	return szConfigManager.GetObserverOrigin(ctx)
}

func (server *SzConfigManagerServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(3, observer.GetObserverID(ctx))
		defer func() { server.traceExit(4, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	return szConfigManager.RegisterObserver(ctx, observer)
}

func (server *SzConfigManagerServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(27, origin)
		defer func() { server.traceExit(28, origin, err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	szConfigManager.SetObserverOrigin(ctx, origin)
}

func (server *SzConfigManagerServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(13, observer.GetObserverID(ctx))
		defer func() { server.traceExit(14, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	szConfigManager := getSzConfigManager()
	return szConfigManager.UnregisterObserver(ctx, observer)
}

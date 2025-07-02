package szconfigserver

import (
	"context"
	"sync"
	"time"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	szobserver "github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
)

const OptionCallerSkip = 3

var (
	szConfigManagerSingleton *szconfigmanager.Szconfigmanager
	szConfigManagerSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/sz-sdk-go/szconfig.SzConfig
// ----------------------------------------------------------------------------

func (server *SzConfigServer) RegisterDataSource(
	ctx context.Context,
	request *szpb.RegisterDataSourceRequest,
) (*szpb.RegisterDataSourceResponse, error) {
	var (
		err      error
		response *szpb.RegisterDataSourceResponse
		result   string
	)

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(1, request)

		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}

	szConfig, err := server.createSzConfig(ctx, request.GetConfigDefinition())
	if err != nil {
		return response, wraperror.Errorf(err, "createSzConfig")
	}

	result, err = szConfig.RegisterDataSource(ctx, request.GetDataSourceCode())
	if err != nil {
		return response, wraperror.Errorf(err, "RegisterDataSource: %s", request.GetDataSourceCode())
	}

	configDefinition, err := szConfig.Export(ctx)
	response = &szpb.RegisterDataSourceResponse{
		Result:           result,
		ConfigDefinition: configDefinition,
	}

	return response, wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzConfigServer) UnregisterDataSource(
	ctx context.Context,
	request *szpb.UnregisterDataSourceRequest,
) (*szpb.UnregisterDataSourceResponse, error) {
	var (
		err      error
		response *szpb.UnregisterDataSourceResponse
	)

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(9, request)

		defer func() { server.traceExit(10, request, err, time.Since(entryTime)) }()
	}

	szConfig, err := server.createSzConfig(ctx, request.GetConfigDefinition())
	if err != nil {
		return response, wraperror.Errorf(err, "createSzConfig")
	}

	result, err := szConfig.UnregisterDataSource(ctx, request.GetDataSourceCode())
	if err != nil {
		return response, wraperror.Errorf(err, "UnregisterDataSource: %s", request.GetDataSourceCode())
	}

	configDefinition, err := szConfig.Export(ctx)
	response = &szpb.UnregisterDataSourceResponse{
		Result:           result,
		ConfigDefinition: configDefinition,
	}

	return response, wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzConfigServer) GetDataSourceRegistry(
	ctx context.Context,
	request *szpb.GetDataSourceRegistryRequest,
) (*szpb.GetDataSourceRegistryResponse, error) {
	var (
		err      error
		response *szpb.GetDataSourceRegistryResponse
		result   string
	)

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(19, request)

		defer func() { server.traceExit(20, request, result, err, time.Since(entryTime)) }()
	}

	szConfig, err := server.createSzConfig(ctx, request.GetConfigDefinition())
	if err != nil {
		return response, err
	}

	result, err = szConfig.GetDataSourceRegistry(ctx)
	response = &szpb.GetDataSourceRegistryResponse{
		Result: result,
	}

	return response, wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzConfigServer) VerifyConfig(
	ctx context.Context,
	request *szpb.VerifyConfigRequest,
) (*szpb.VerifyConfigResponse, error) {
	var (
		err      error
		response *szpb.VerifyConfigResponse
		result   bool
	)

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(999, request)

		defer func() { server.traceExit(999, request, result, err, time.Since(entryTime)) }()
	}

	szConfig, err := server.createSzConfig(ctx, request.GetConfigDefinition())
	if err != nil {
		return response, err
	}

	result = true

	err = szConfig.VerifyConfigDefinition(ctx, request.GetConfigDefinition())
	if err != nil {
		result = false
	}

	response = &szpb.VerifyConfigResponse{
		Result: result,
	}

	return response, wraperror.Errorf(err, wraperror.NoMessage)
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
			&logging.OptionCallerSkip{Value: OptionCallerSkip},
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
		return wraperror.Errorf(errPackage, "invalid error level: %s", logLevelName)
	}

	server.logLevelName = logLevelName
	// szconfig := getSzConfig()
	// err = szconfig.SetLogLevel(ctx, logLevelName)
	// if err != nil {
	// 	return err
	// }
	err = server.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		return wraperror.Errorf(err, "SetLogLevel: %s", logLevelName)
	}

	server.isTrace = (logLevelName == logging.LevelTraceName)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// --- Services ---------------------------------------------------------------

func (server *SzConfigServer) createSzConfig(ctx context.Context, configDefinition string) (*szconfig.Szconfig, error) {
	szConfigManager := getSzConfigManager()

	result, err := szConfigManager.CreateConfigFromStringChoreography(ctx, configDefinition)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzConfigServer) GetSdkSzConfigAsInterface(
	ctx context.Context,
	configDefinition string,
) (senzing.SzConfig, error) {
	result, err := server.createSzConfig(ctx, configDefinition)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

// Singleton pattern for szconfigmanager.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getSzConfigManager() *szconfigmanager.Szconfigmanager {
	szConfigManagerSyncOnce.Do(func() {
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{}
	})

	return szConfigManagerSingleton
}

func GetSdkSzConfigManager() *szconfigmanager.Szconfigmanager {
	return getSzConfigManager()
}

func GetSdkSzConfigManagerAsInterface() senzing.SzConfigManager {
	return getSzConfigManager()
}

// --- Observer ---------------------------------------------------------------

func (server *SzConfigServer) GetObserverOrigin(ctx context.Context) string {
	var err error

	_ = ctx

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(27)

		defer func() { server.traceExit(28, err, time.Since(entryTime)) }()
	}

	return server.observerOrigin
}

func (server *SzConfigServer) RegisterObserver(ctx context.Context, observer szobserver.Observer) error {
	var err error

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(3, observer.GetObserverID(ctx))

		defer func() { server.traceExit(4, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}

	server.observers = append(server.observers, observer)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzConfigServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error

	_ = ctx

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(29, origin)

		defer func() { server.traceExit(30, origin, err, time.Since(entryTime)) }()
	}

	server.observerOrigin = origin
}

func (server *SzConfigServer) UnregisterObserver(ctx context.Context, observer szobserver.Observer) error {
	var err error

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(13, observer.GetObserverID(ctx))

		defer func() { server.traceExit(14, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}

	if len(server.observers) > 0 {
		result := make([]szobserver.Observer, 0, len(server.observers))

		for _, registeredObserver := range server.observers {
			if registeredObserver.GetObserverID(ctx) != observer.GetObserverID(ctx) {
				result = append(result, registeredObserver)
			}
		}

		server.observers = result
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

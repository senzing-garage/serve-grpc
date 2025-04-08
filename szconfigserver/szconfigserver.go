package szconfigserver

import (
	"context"
	"fmt"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	szobserver "github.com/senzing-garage/go-observing/observer"
	szsdk "github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
)

// var (
// 	szConfig         *szsdk.Szconfig
// )

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/sz-sdk-go/szconfig.SzConfig
// ----------------------------------------------------------------------------

func (server *SzConfigServer) AddDataSource(
	ctx context.Context,
	request *szpb.AddDataSourceRequest,
) (*szpb.AddDataSourceResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}

	szConfig := server.createSzConfig(ctx, request.GetConfigDefinition())
	result, err = szConfig.AddDataSource(ctx, request.GetDataSourceCode())
	response := szpb.AddDataSourceResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzConfigServer) DeleteDataSource(
	ctx context.Context,
	request *szpb.DeleteDataSourceRequest,
) (*szpb.DeleteDataSourceResponse, error) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(9, request)
		defer func() { server.traceExit(10, request, err, time.Since(entryTime)) }()
	}
	szConfig := server.createSzConfig(ctx, request.GetConfigDefinition())
	result, err := szConfig.DeleteDataSource(ctx, request.GetDataSourceCode())
	response := szpb.DeleteDataSourceResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzConfigServer) GetDataSources(
	ctx context.Context,
	request *szpb.GetDataSourcesRequest,
) (*szpb.GetDataSourcesResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, result, err, time.Since(entryTime)) }()
	}
	szConfig := server.createSzConfig(ctx, request.GetConfigDefinition())
	result, err = szConfig.GetDataSources(ctx)
	response := szpb.GetDataSourcesResponse{
		Result: result,
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

	server.logLevelName = logLevelName
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

// --- Services ---------------------------------------------------------------

func (server *SzConfigServer) createSzConfig(ctx context.Context, configDefinition string) *szsdk.Szconfig {
	szConfig := &szsdk.Szconfig{}
	szConfig.SetLogLevel(ctx, server.logLevelName)
	szConfig.Import(ctx, configDefinition)
	szConfig.SetObserverOrigin(ctx, server.observerOrigin)
	for _, x := range server.observers {
		szConfig.RegisterObserver(ctx, x)
	}
	return szConfig
}

func (server *SzConfigServer) GetSdkSzConfig() *szsdk.Szconfig {
	ctx := context.TODO()
	return server.createSzConfig(ctx, "")
}

func (server *SzConfigServer) GetSdkSzConfigAsInterface(ctx context.Context, configDefinition string) senzing.SzConfig {
	return server.createSzConfig(ctx, configDefinition)
}

// --- Observer ---------------------------------------------------------------

func (server *SzConfigServer) GetObserverOrigin(ctx context.Context) string {
	var err error
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
	return err
}

func (server *SzConfigServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error
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

	result := make([]szobserver.Observer, len(server.observers))
	for _, registeredObserver := range server.observers {
		if registeredObserver.GetObserverID(ctx) != observer.GetObserverID(ctx) {
			result = append(result, registeredObserver)
		}
	}
	server.observers = result
	return nil

}

package szdiagnosticserver

import (
	"context"
	"sync"
	"time"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	szsdk "github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
)

const OptionCallerSkip = 3

var (
	szDiagnosticSingleton *szsdk.Szdiagnostic
	szDiagnosticSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/sz-sdk-go/szdiagnostic.SzDdiagnostic
// ----------------------------------------------------------------------------

func (server *SzDiagnosticServer) CheckRepositoryPerformance(
	ctx context.Context,
	request *szpb.CheckRepositoryPerformanceRequest,
) (*szpb.CheckRepositoryPerformanceResponse, error) {
	var (
		err      error
		response *szpb.CheckRepositoryPerformanceResponse
		result   string
	)

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(1, request)

		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}

	szDiagnostic := getSzDiagnostic()
	result, err = szDiagnostic.CheckRepositoryPerformance(ctx, int(request.GetSecondsToRun()))
	response = &szpb.CheckRepositoryPerformanceResponse{
		Result: result,
	}

	return response, wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzDiagnosticServer) GetRepositoryInfo(
	ctx context.Context,
	request *szpb.GetRepositoryInfoRequest,
) (*szpb.GetRepositoryInfoResponse, error) {
	var (
		err      error
		response *szpb.GetRepositoryInfoResponse
		result   string
	)

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(1, request)

		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}

	szDiagnostic := getSzDiagnostic()
	result, err = szDiagnostic.GetRepositoryInfo(ctx)
	response = &szpb.GetRepositoryInfoResponse{
		Result: result,
	}

	return response, wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzDiagnosticServer) GetFeature(
	ctx context.Context,
	request *szpb.GetFeatureRequest,
) (*szpb.GetFeatureResponse, error) {
	var err error

	var result string

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(1, request)

		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}

	szDiagnostic := getSzDiagnostic()
	result, err = szDiagnostic.GetFeature(ctx, request.GetFeatureId())
	response := szpb.GetFeatureResponse{
		Result: result,
	}

	return &response, wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzDiagnosticServer) PurgeRepository(
	ctx context.Context,
	request *szpb.PurgeRepositoryRequest,
) (*szpb.PurgeRepositoryResponse, error) {
	var err error

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(117, request)

		defer func() { server.traceExit(118, request, err, time.Since(entryTime)) }()
	}

	szDiagnostic := getSzDiagnostic()
	err = szDiagnostic.PurgeRepository(ctx)
	response := szpb.PurgeRepositoryResponse{}

	return &response, wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzDiagnosticServer) Reinitialize(
	ctx context.Context,
	request *szpb.ReinitializeRequest,
) (*szpb.ReinitializeResponse, error) {
	var err error

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(51, request)

		defer func() { server.traceExit(52, request, err, time.Since(entryTime)) }()
	}

	szDiagnostic := getSzDiagnostic()
	err = szDiagnostic.Reinitialize(ctx, request.GetConfigId())
	response := szpb.ReinitializeResponse{}

	return &response, wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *SzDiagnosticServer) getLogger() logging.Logging {
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
func (server *SzDiagnosticServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *SzDiagnosticServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

func (server *SzDiagnosticServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	_ = ctx

	var err error

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(53, logLevelName)

		defer func() { server.traceExit(54, logLevelName, err, time.Since(entryTime)) }()
	}

	if !logging.IsValidLogLevelName(logLevelName) {
		return wraperror.Errorf(errPackage, "invalid error level: %s", logLevelName)
	}
	// szdiagnostic := getSzdiagnostic()
	// err = szdiagnostic.SetLogLevel(ctx, logLevelName)
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

// --- Errors -----------------------------------------------------------------

// Create error.
// func (server *SzDiagnosticServer) error(messageNumber int, details ...interface{}) error {
// 	return server.getLogger().NewError(messageNumber, details...)
// }

// --- Services ---------------------------------------------------------------

// Singleton pattern for szdiagnostic.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getSzDiagnostic() *szsdk.Szdiagnostic {
	szDiagnosticSyncOnce.Do(func() {
		szDiagnosticSingleton = &szsdk.Szdiagnostic{}
	})

	return szDiagnosticSingleton
}

func GetSdkSzDiagnostic() *szsdk.Szdiagnostic {
	return getSzDiagnostic()
}

func GetSdkSzDiagnosticAsInterface() senzing.SzDiagnostic {
	return getSzDiagnostic()
}

// --- Observer ---------------------------------------------------------------

func (server *SzDiagnosticServer) GetObserverOrigin(ctx context.Context) string {
	var err error

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(55)

		defer func() { server.traceExit(56, err, time.Since(entryTime)) }()
	}

	szDiagnostic := getSzDiagnostic()

	return szDiagnostic.GetObserverOrigin(ctx)
}

func (server *SzDiagnosticServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(3, observer.GetObserverID(ctx))

		defer func() { server.traceExit(4, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}

	szDiagnostic := getSzDiagnostic()

	err = szDiagnostic.RegisterObserver(ctx, observer)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

func (server *SzDiagnosticServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(57, origin)

		defer func() { server.traceExit(58, origin, err, time.Since(entryTime)) }()
	}

	szDiagnostic := getSzDiagnostic()
	szDiagnostic.SetObserverOrigin(ctx, origin)
}

func (server *SzDiagnosticServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if server.isTrace {
		entryTime := time.Now()

		server.traceEntry(31, observer.GetObserverID(ctx))

		defer func() { server.traceExit(32, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}

	szDiagnostic := getSzDiagnostic()

	err = szDiagnostic.UnregisterObserver(ctx, observer)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

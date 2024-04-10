package szdiagnosticserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	g2sdk "github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/sz"
	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
)

var (
	g2diagnosticSingleton sz.G2diagnostic
	g2diagnosticSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *SzDiagnosticServer) getLogger() logging.LoggingInterface {
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
func (server *SzDiagnosticServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *SzDiagnosticServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (server *SzDiagnosticServer) error(messageNumber int, details ...interface{}) error {
	return server.getLogger().NewError(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

// Singleton pattern for g2diagnostic.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2diagnostic() sz.G2diagnostic {
	g2diagnosticSyncOnce.Do(func() {
		g2diagnosticSingleton = &g2sdk.G2diagnostic{}
	})
	return g2diagnosticSingleton
}

func GetSdkSzDiagnostic() sz.G2diagnostic {
	return getG2diagnostic()
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/g2-sdk-go/g2diagnostic.G2diagnostic
// ----------------------------------------------------------------------------

func (server *SzDiagnosticServer) CheckDBPerf(ctx context.Context, request *g2pb.CheckDBPerfRequest) (*g2pb.CheckDBPerfResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.CheckDBPerf(ctx, int(request.GetSecondsToRun()))
	response := g2pb.CheckDBPerfResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzDiagnosticServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(7, request)
		defer func() { server.traceExit(8, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2diagnostic := getG2diagnostic()
	// err := g2diagnostic.Destroy(ctx)
	err = server.error(4001)
	response := g2pb.DestroyResponse{}
	return &response, err
}

func (server *SzDiagnosticServer) GetObserverOrigin(ctx context.Context) string {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(55)
		defer func() { server.traceExit(56, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	return g2diagnostic.GetObserverOrigin(ctx)
}

func (server *SzDiagnosticServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(47, request)
		defer func() { server.traceExit(48, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2diagnostic := getG2diagnostic()
	// err := g2diagnostic.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err = server.error(4002)
	response := g2pb.InitResponse{}

	return &response, err
}

func (server *SzDiagnosticServer) InitWithConfigID(ctx context.Context, request *g2pb.InitWithConfigIDRequest) (*g2pb.InitWithConfigIDResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(49, request)
		defer func() { server.traceExit(50, request, err, time.Since(entryTime)) }()
	}
	// g2diagnostic := getG2diagnostic()
	// err := g2diagnostic.InitWithConfigID(ctx, request.GetModuleName(), request.GetIniParams(), int64(request.GetInitConfigID()), int(request.GetVerboseLogging()))
	err = server.error(4003)
	response := g2pb.InitWithConfigIDResponse{}
	return &response, err
}

func (server *SzDiagnosticServer) PurgeRepository(ctx context.Context, request *g2pb.PurgeRepositoryRequest) (*g2pb.PurgeRepositoryResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(117, request)
		defer func() { server.traceExit(118, request, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	err = g2diagnostic.PurgeRepository(ctx)
	response := g2pb.PurgeRepositoryResponse{}
	return &response, err
}

func (server SzDiagnosticServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(3, observer.GetObserverId(ctx))
		defer func() { server.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	return g2diagnostic.RegisterObserver(ctx, observer)
}

func (server *SzDiagnosticServer) Reinit(ctx context.Context, request *g2pb.ReinitRequest) (*g2pb.ReinitResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(51, request)
		defer func() { server.traceExit(52, request, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	err = g2diagnostic.Reinit(ctx, int64(request.GetInitConfigID()))
	response := g2pb.ReinitResponse{}
	return &response, err
}

func (server *SzDiagnosticServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(53, logLevelName)
		defer func() { server.traceExit(54, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	g2diagnostic := getG2diagnostic()
	err = g2diagnostic.SetLogLevel(ctx, logLevelName)
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

func (server *SzDiagnosticServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(57, origin)
		defer func() { server.traceExit(58, origin, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	g2diagnostic.SetObserverOrigin(ctx, origin)
}

func (server *SzDiagnosticServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(31, observer.GetObserverId(ctx))
		defer func() { server.traceExit(32, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	return g2diagnostic.UnregisterObserver(ctx, observer)
}

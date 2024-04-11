package szdiagnosticserver

import (
	"context"
	"sync"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	szsdk "github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
)

var (
	szDiagnosticSingleton sz.SzDiagnostic
	szDiagnosticSyncOnce  sync.Once
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
func getG2diagnostic() sz.SzDiagnostic {
	szDiagnosticSyncOnce.Do(func() {
		szDiagnosticSingleton = &szsdk.Szdiagnostic{}
	})
	return szDiagnosticSingleton
}

func GetSdkSzDiagnostic() sz.SzDiagnostic {
	return getG2diagnostic()
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/g2-sdk-go/g2diagnostic.G2diagnostic
// ----------------------------------------------------------------------------

func (server *SzDiagnosticServer) CheckDatabasePerformance(ctx context.Context, request *szpb.CheckDatabasePerformanceRequest) (*szpb.CheckDatabasePerformanceResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.CheckDatabasePerformance(ctx, int(request.GetSecondsToRun()))
	response := szpb.CheckDatabasePerformanceResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzDiagnosticServer) PurgeRepository(ctx context.Context, request *szpb.PurgeRepositoryRequest) (*szpb.PurgeRepositoryResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(117, request)
		defer func() { server.traceExit(118, request, err, time.Since(entryTime)) }()
	}
	szDiagnostic := getG2diagnostic()
	err = szDiagnostic.PurgeRepository(ctx)
	response := szpb.PurgeRepositoryResponse{}
	return &response, err
}

// func (server *SzDiagnosticServer) GetObserverOrigin(ctx context.Context) string {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(55)
// 		defer func() { server.traceExit(56, err, time.Since(entryTime)) }()
// 	}
// 	g2diagnostic := getG2diagnostic()
// 	return g2diagnostic.GetObserverOrigin(ctx)
// }

// func (server SzDiagnosticServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(3, observer.GetObserverId(ctx))
// 		defer func() { server.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
// 	}
// 	g2diagnostic := getG2diagnostic()
// 	return g2diagnostic.RegisterObserver(ctx, observer)
// }

// func (server *SzDiagnosticServer) Reinit(ctx context.Context, request *g2pb.ReinitRequest) (*g2pb.ReinitResponse, error) {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(51, request)
// 		defer func() { server.traceExit(52, request, err, time.Since(entryTime)) }()
// 	}
// 	g2diagnostic := getG2diagnostic()
// 	err = g2diagnostic.Reinit(ctx, int64(request.GetInitConfigID()))
// 	response := g2pb.ReinitResponse{}
// 	return &response, err
// }

// func (server *SzDiagnosticServer) SetLogLevel(ctx context.Context, logLevelName string) error {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(53, logLevelName)
// 		defer func() { server.traceExit(54, logLevelName, err, time.Since(entryTime)) }()
// 	}
// 	if !logging.IsValidLogLevelName(logLevelName) {
// 		return fmt.Errorf("invalid error level: %s", logLevelName)
// 	}
// 	g2diagnostic := getG2diagnostic()
// 	err = g2diagnostic.SetLogLevel(ctx, logLevelName)
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

// func (server *SzDiagnosticServer) SetObserverOrigin(ctx context.Context, origin string) {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(57, origin)
// 		defer func() { server.traceExit(58, origin, err, time.Since(entryTime)) }()
// 	}
// 	g2diagnostic := getG2diagnostic()
// 	g2diagnostic.SetObserverOrigin(ctx, origin)
// }

// func (server *SzDiagnosticServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(31, observer.GetObserverId(ctx))
// 		defer func() { server.traceExit(32, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
// 	}
// 	g2diagnostic := getG2diagnostic()
// 	return g2diagnostic.UnregisterObserver(ctx, observer)
// }

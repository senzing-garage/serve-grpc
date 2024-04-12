package szproductserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	szsdk "github.com/senzing-garage/sz-sdk-go-core/szproduct"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
)

var (
	szProductSingleton sz.SzProduct
	szProductSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/sz-sdk-go/szproduct.SzProduct
// ----------------------------------------------------------------------------

func (server *SzProductServer) GetLicense(ctx context.Context, request *szpb.GetLicenseRequest) (*szpb.GetLicenseResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(11, request)
		defer func() { server.traceExit(12, request, result, err, time.Since(entryTime)) }()
	}
	szproduct := getSzProduct()
	result, err = szproduct.GetLicense(ctx)
	response := szpb.GetLicenseResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzProductServer) GetVersion(ctx context.Context, request *szpb.GetVersionRequest) (*szpb.GetVersionResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, result, err, time.Since(entryTime)) }()
	}
	szproduct := getSzProduct()
	result, err = szproduct.GetVersion(ctx)
	response := szpb.GetVersionResponse{
		Result: result,
	}
	return &response, err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *SzProductServer) getLogger() logging.LoggingInterface {
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
func (server *SzProductServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *SzProductServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

func (server *SzProductServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(13, logLevelName)
		defer func() { server.traceExit(14, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	// g2product := getG2product()
	// err = g2product.SetLogLevel(ctx, logLevelName)
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
// func (server *SzProductServer) error(messageNumber int, details ...interface{}) error {
// 	return server.getLogger().NewError(messageNumber, details...)
// }

// --- Services ---------------------------------------------------------------

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getSzProduct() sz.SzProduct {
	szProductSyncOnce.Do(func() {
		szProductSingleton = &szsdk.Szproduct{}
	})
	return szProductSingleton
}

func GetSdkSzProduct() sz.SzProduct {
	return getSzProduct()
}

// --- Observer ---------------------------------------------------------------

func (server *SzProductServer) GetObserverOrigin(ctx context.Context) string {
	// var err error = nil
	// if server.isTrace {
	// 	entryTime := time.Now()
	// 	server.traceEntry(21)
	// 	defer func() { server.traceExit(22, err, time.Since(entryTime)) }()
	// }
	// g2product := getG2product()
	// return g2product.GetObserverOrigin(ctx)
	return ""
}

func (server *SzProductServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	// var err error = nil
	// if server.isTrace {
	// 	entryTime := time.Now()
	// 	server.traceEntry(1, observer.GetObserverId(ctx))
	// 	defer func() { server.traceExit(2, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	// }
	// g2product := getSzProduct()
	// return g2product.RegisterObserver(ctx, observer)
	return nil
}

func (server *SzProductServer) SetObserverOrigin(ctx context.Context, origin string) {
	// var err error = nil
	// if server.isTrace {
	// 	entryTime := time.Now()
	// 	server.traceEntry(23, origin)
	// 	defer func() { server.traceExit(24, origin, err, time.Since(entryTime)) }()
	// }
	// g2product := getSzProduct()
	// g2product.SetObserverOrigin(ctx, origin)
}

func (server *SzProductServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	// var err error = nil
	// if server.isTrace {
	// 	entryTime := time.Now()
	// 	server.traceEntry(5, observer.GetObserverId(ctx))
	// 	defer func() { server.traceExit(6, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	// }
	// g2product := getSzProduct()
	// return g2product.UnregisterObserver(ctx, observer)
	return nil
}

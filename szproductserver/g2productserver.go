package szproductserver

import (
	"context"
	"sync"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	szsdk "github.com/senzing-garage/sz-sdk-go-core/szproduct"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
)

var (
	szProductSingleton sz.SzProduct
	szProductSyncOnce  sync.Once
)

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

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/g2-sdk-go/g2product.G2product
// ----------------------------------------------------------------------------

// func (server *SzProductServer) GetObserverOrigin(ctx context.Context) string {
// 	var err error = nil
// 	if server.isTrace {
// 		entryTime := time.Now()
// 		server.traceEntry(21)
// 		defer func() { server.traceExit(22, err, time.Since(entryTime)) }()
// 	}
// 	g2product := getG2product()
// 	return g2product.GetObserverOrigin(ctx)
// }

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

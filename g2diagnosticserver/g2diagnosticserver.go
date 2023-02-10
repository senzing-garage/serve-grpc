package g2diagnosticserver

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go-base/g2diagnostic"
	"github.com/senzing/g2-sdk-go/g2api"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
)

var (
	g2diagnosticSingleton g2api.G2diagnosticInterface
	g2diagnosticSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// func getLogger() messagelogger.MessageLoggerInterface {

// 	onceLogger.Do(func() {
// 		logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
// 	})
// 	return logger
// }

// Singleton pattern for g2diagnostic.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2diagnostic() g2api.G2diagnosticInterface {
	g2diagnosticSyncOnce.Do(func() {
		g2diagnosticSingleton = &g2sdk.G2diagnostic{}
	})
	return g2diagnosticSingleton
}

func GetSdkG2diagnostic() g2api.G2diagnosticInterface {
	return getG2diagnostic()
}

// Get the Logger singleton.
func (server *G2DiagnosticServer) getLogger() messagelogger.MessageLoggerInterface {
	if server.logger == nil {
		server.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return server.logger
}

// Trace method entry.
func (server *G2DiagnosticServer) traceEntry(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (server *G2DiagnosticServer) traceExit(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing/g2-sdk-go/g2diagnostic.G2diagnostic
// ----------------------------------------------------------------------------

func (server *G2DiagnosticServer) CheckDBPerf(ctx context.Context, request *g2pb.CheckDBPerfRequest) (*g2pb.CheckDBPerfResponse, error) {
	if server.isTrace {
		server.traceEntry(1, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.CheckDBPerf(ctx, int(request.GetSecondsToRun()))
	response := g2pb.CheckDBPerfResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(2, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) CloseEntityListBySize(ctx context.Context, request *g2pb.CloseEntityListBySizeRequest) (*g2pb.CloseEntityListBySizeResponse, error) {
	if server.isTrace {
		server.traceEntry(5, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	entityListBySizeHandleInt, err := strconv.ParseUint(request.GetEntityListBySizeHandle(), 10, 64)
	if err == nil {
		err = g2diagnostic.CloseEntityListBySize(ctx, uintptr(entityListBySizeHandleInt))
	}
	response := g2pb.CloseEntityListBySizeResponse{}
	if server.isTrace {
		defer server.traceExit(6, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
	if server.isTrace {
		server.traceEntry(7, request)
	}
	entryTime := time.Now()
	// g2diagnostic := getG2diagnostic()
	// err := g2diagnostic.Destroy(ctx)
	err := server.getLogger().Error(4001)
	response := g2pb.DestroyResponse{}
	if server.isTrace {
		defer server.traceExit(8, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) FetchNextEntityBySize(ctx context.Context, request *g2pb.FetchNextEntityBySizeRequest) (*g2pb.FetchNextEntityBySizeResponse, error) {
	var result string = ""
	if server.isTrace {
		server.traceEntry(9, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	entityListBySizeHandleInt, err := strconv.ParseUint(request.GetEntityListBySizeHandle(), 10, 64)
	if err == nil {
		result, err = g2diagnostic.FetchNextEntityBySize(ctx, uintptr(entityListBySizeHandleInt))

	}
	response := g2pb.FetchNextEntityBySizeResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(10, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) FindEntitiesByFeatureIDs(ctx context.Context, request *g2pb.FindEntitiesByFeatureIDsRequest) (*g2pb.FindEntitiesByFeatureIDsResponse, error) {
	if server.isTrace {
		server.traceEntry(11, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, request.GetFeatures())
	response := g2pb.FindEntitiesByFeatureIDsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(12, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetAvailableMemory(ctx context.Context, request *g2pb.GetAvailableMemoryRequest) (*g2pb.GetAvailableMemoryResponse, error) {
	if server.isTrace {
		server.traceEntry(13, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetAvailableMemory(ctx)
	response := g2pb.GetAvailableMemoryResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(14, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetDataSourceCounts(ctx context.Context, request *g2pb.GetDataSourceCountsRequest) (*g2pb.GetDataSourceCountsResponse, error) {
	if server.isTrace {
		server.traceEntry(15, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetDataSourceCounts(ctx)
	response := g2pb.GetDataSourceCountsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(16, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetDBInfo(ctx context.Context, request *g2pb.GetDBInfoRequest) (*g2pb.GetDBInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(17, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetDBInfo(ctx)
	response := g2pb.GetDBInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(18, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityDetails(ctx context.Context, request *g2pb.GetEntityDetailsRequest) (*g2pb.GetEntityDetailsResponse, error) {
	if server.isTrace {
		server.traceEntry(19, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntityDetails(ctx, request.GetEntityID(), int(request.GetIncludeInternalFeatures()))
	response := g2pb.GetEntityDetailsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(20, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityListBySize(ctx context.Context, request *g2pb.GetEntityListBySizeRequest) (*g2pb.GetEntityListBySizeResponse, error) {
	if server.isTrace {
		server.traceEntry(21, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntityListBySize(ctx, int(request.GetEntitySize()))
	response := g2pb.GetEntityListBySizeResponse{
		Result: fmt.Sprintf("%v", result),
	}
	if server.isTrace {
		defer server.traceExit(22, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityResume(ctx context.Context, request *g2pb.GetEntityResumeRequest) (*g2pb.GetEntityResumeResponse, error) {
	if server.isTrace {
		server.traceEntry(23, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntityResume(ctx, request.GetEntityID())
	response := g2pb.GetEntityResumeResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(24, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntitySizeBreakdown(ctx context.Context, request *g2pb.GetEntitySizeBreakdownRequest) (*g2pb.GetEntitySizeBreakdownResponse, error) {
	if server.isTrace {
		server.traceEntry(25, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntitySizeBreakdown(ctx, int(request.GetMinimumEntitySize()), int(request.GetIncludeInternalFeatures()))
	response := g2pb.GetEntitySizeBreakdownResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(26, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetFeature(ctx context.Context, request *g2pb.GetFeatureRequest) (*g2pb.GetFeatureResponse, error) {
	if server.isTrace {
		server.traceEntry(27, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetFeature(ctx, request.GetLibFeatID())
	response := g2pb.GetFeatureResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(28, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetGenericFeatures(ctx context.Context, request *g2pb.GetGenericFeaturesRequest) (*g2pb.GetGenericFeaturesResponse, error) {
	if server.isTrace {
		server.traceEntry(29, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetGenericFeatures(ctx, request.GetFeatureType(), int(request.GetMaximumEstimatedCount()))
	response := g2pb.GetGenericFeaturesResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(30, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetLogicalCores(ctx context.Context, request *g2pb.GetLogicalCoresRequest) (*g2pb.GetLogicalCoresResponse, error) {
	if server.isTrace {
		server.traceEntry(35, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetLogicalCores(ctx)
	response := g2pb.GetLogicalCoresResponse{
		Result: int32(result),
	}
	if server.isTrace {
		defer server.traceExit(36, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetMappingStatistics(ctx context.Context, request *g2pb.GetMappingStatisticsRequest) (*g2pb.GetMappingStatisticsResponse, error) {
	if server.isTrace {
		server.traceEntry(37, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetMappingStatistics(ctx, int(request.GetIncludeInternalFeatures()))
	response := g2pb.GetMappingStatisticsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(38, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetPhysicalCores(ctx context.Context, request *g2pb.GetPhysicalCoresRequest) (*g2pb.GetPhysicalCoresResponse, error) {
	if server.isTrace {
		server.traceEntry(39, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetPhysicalCores(ctx)
	response := g2pb.GetPhysicalCoresResponse{
		Result: int32(result),
	}
	if server.isTrace {
		defer server.traceExit(40, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetRelationshipDetails(ctx context.Context, request *g2pb.GetRelationshipDetailsRequest) (*g2pb.GetRelationshipDetailsResponse, error) {
	if server.isTrace {
		server.traceEntry(41, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetRelationshipDetails(ctx, request.GetRelationshipID(), int(request.GetIncludeInternalFeatures()))
	response := g2pb.GetRelationshipDetailsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(42, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetResolutionStatistics(ctx context.Context, request *g2pb.GetResolutionStatisticsRequest) (*g2pb.GetResolutionStatisticsResponse, error) {
	if server.isTrace {
		server.traceEntry(43, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetResolutionStatistics(ctx)
	response := g2pb.GetResolutionStatisticsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(44, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetTotalSystemMemory(ctx context.Context, request *g2pb.GetTotalSystemMemoryRequest) (*g2pb.GetTotalSystemMemoryResponse, error) {
	if server.isTrace {
		server.traceEntry(45, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetTotalSystemMemory(ctx)
	response := g2pb.GetTotalSystemMemoryResponse{
		Result: int64(result),
	}
	if server.isTrace {
		defer server.traceExit(46, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	if server.isTrace {
		server.traceEntry(47, request)
	}
	entryTime := time.Now()
	// g2diagnostic := getG2diagnostic()
	// err := g2diagnostic.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err := server.getLogger().Error(4002)
	response := g2pb.InitResponse{}
	if server.isTrace {
		defer server.traceExit(48, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2DiagnosticServer) InitWithConfigID(ctx context.Context, request *g2pb.InitWithConfigIDRequest) (*g2pb.InitWithConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(49, request)
	}
	entryTime := time.Now()
	// g2diagnostic := getG2diagnostic()
	// err := g2diagnostic.InitWithConfigID(ctx, request.GetModuleName(), request.GetIniParams(), int64(request.GetInitConfigID()), int(request.GetVerboseLogging()))
	err := server.getLogger().Error(4003)
	response := g2pb.InitWithConfigIDResponse{}
	if server.isTrace {
		defer server.traceExit(50, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server G2DiagnosticServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	g2diagnostic := getG2diagnostic()
	return g2diagnostic.RegisterObserver(ctx, observer)
}

func (server *G2DiagnosticServer) Reinit(ctx context.Context, request *g2pb.ReinitRequest) (*g2pb.ReinitResponse, error) {
	if server.isTrace {
		server.traceEntry(51, request)
	}
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.Reinit(ctx, int64(request.GetInitConfigID()))
	response := g2pb.ReinitResponse{}
	if server.isTrace {
		defer server.traceExit(52, request, err, time.Since(entryTime))
	}
	return &response, err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (server *G2DiagnosticServer) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if server.isTrace {
		server.traceEntry(53, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	g2diagnostic := getG2diagnostic()
	g2diagnostic.SetLogLevel(ctx, logLevel)
	server.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	server.isTrace = (server.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if server.isTrace {
		defer server.traceExit(54, logLevel, err, time.Since(entryTime))
	}
	return err
}

func (server *G2DiagnosticServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	g2diagnostic := getG2diagnostic()
	return g2diagnostic.UnregisterObserver(ctx, observer)
}

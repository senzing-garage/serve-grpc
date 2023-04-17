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
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/observer"
)

var (
	g2diagnosticSingleton g2api.G2diagnostic
	g2diagnosticSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *G2DiagnosticServer) getLogger() logging.LoggingInterface {
	var err error = nil
	if server.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 3},
		}
		server.logger, err = logging.NewSenzingToolsLogger(ProductId, IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return server.logger
}

// Trace method entry.
func (server *G2DiagnosticServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *G2DiagnosticServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (server *G2DiagnosticServer) error(messageNumber int, details ...interface{}) error {
	return server.getLogger().Error(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

// Singleton pattern for g2diagnostic.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2diagnostic() g2api.G2diagnostic {
	g2diagnosticSyncOnce.Do(func() {
		g2diagnosticSingleton = &g2sdk.G2diagnostic{}
	})
	return g2diagnosticSingleton
}

func GetSdkG2diagnostic() g2api.G2diagnostic {
	return getG2diagnostic()
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing/g2-sdk-go/g2diagnostic.G2diagnostic
// ----------------------------------------------------------------------------

func (server *G2DiagnosticServer) CheckDBPerf(ctx context.Context, request *g2pb.CheckDBPerfRequest) (*g2pb.CheckDBPerfResponse, error) {
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

func (server *G2DiagnosticServer) CloseEntityListBySize(ctx context.Context, request *g2pb.CloseEntityListBySizeRequest) (*g2pb.CloseEntityListBySizeResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(5, request)
		defer func() { server.traceExit(6, request, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	entityListBySizeHandleInt, err := strconv.ParseUint(request.GetEntityListBySizeHandle(), 10, 64)
	if err == nil {
		err = g2diagnostic.CloseEntityListBySize(ctx, uintptr(entityListBySizeHandleInt))
	}
	response := g2pb.CloseEntityListBySizeResponse{}
	return &response, err
}

func (server *G2DiagnosticServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
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

func (server *G2DiagnosticServer) FetchNextEntityBySize(ctx context.Context, request *g2pb.FetchNextEntityBySizeRequest) (*g2pb.FetchNextEntityBySizeResponse, error) {
	var err error = nil
	var result string = ""
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(9, request)
		defer func() { server.traceExit(10, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	entityListBySizeHandleInt, err := strconv.ParseUint(request.GetEntityListBySizeHandle(), 10, 64)
	if err == nil {
		result, err = g2diagnostic.FetchNextEntityBySize(ctx, uintptr(entityListBySizeHandleInt))
	}
	response := g2pb.FetchNextEntityBySizeResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) FindEntitiesByFeatureIDs(ctx context.Context, request *g2pb.FindEntitiesByFeatureIDsRequest) (*g2pb.FindEntitiesByFeatureIDsResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(11, request)
		defer func() { server.traceExit(12, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.FindEntitiesByFeatureIDs(ctx, request.GetFeatures())
	response := g2pb.FindEntitiesByFeatureIDsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetAvailableMemory(ctx context.Context, request *g2pb.GetAvailableMemoryRequest) (*g2pb.GetAvailableMemoryResponse, error) {
	var err error = nil
	var result int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(13, request)
		defer func() { server.traceExit(14, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetAvailableMemory(ctx)
	response := g2pb.GetAvailableMemoryResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetDataSourceCounts(ctx context.Context, request *g2pb.GetDataSourceCountsRequest) (*g2pb.GetDataSourceCountsResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(15, request)
		defer func() { server.traceExit(16, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetDataSourceCounts(ctx)
	response := g2pb.GetDataSourceCountsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetDBInfo(ctx context.Context, request *g2pb.GetDBInfoRequest) (*g2pb.GetDBInfoResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(17, request)
		defer func() { server.traceExit(18, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetDBInfo(ctx)
	response := g2pb.GetDBInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityDetails(ctx context.Context, request *g2pb.GetEntityDetailsRequest) (*g2pb.GetEntityDetailsResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetEntityDetails(ctx, request.GetEntityID(), int(request.GetIncludeInternalFeatures()))
	response := g2pb.GetEntityDetailsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityListBySize(ctx context.Context, request *g2pb.GetEntityListBySizeRequest) (*g2pb.GetEntityListBySizeResponse, error) {
	var err error = nil
	var result uintptr
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(21, request)
		defer func() { server.traceExit(22, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetEntityListBySize(ctx, int(request.GetEntitySize()))
	response := g2pb.GetEntityListBySizeResponse{
		Result: fmt.Sprintf("%v", result),
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityResume(ctx context.Context, request *g2pb.GetEntityResumeRequest) (*g2pb.GetEntityResumeResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(23, request)
		defer func() { server.traceExit(24, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetEntityResume(ctx, request.GetEntityID())
	response := g2pb.GetEntityResumeResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntitySizeBreakdown(ctx context.Context, request *g2pb.GetEntitySizeBreakdownRequest) (*g2pb.GetEntitySizeBreakdownResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(25, request)
		defer func() { server.traceExit(26, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetEntitySizeBreakdown(ctx, int(request.GetMinimumEntitySize()), int(request.GetIncludeInternalFeatures()))
	response := g2pb.GetEntitySizeBreakdownResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetFeature(ctx context.Context, request *g2pb.GetFeatureRequest) (*g2pb.GetFeatureResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(27, request)
		defer func() { server.traceExit(28, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetFeature(ctx, request.GetLibFeatID())
	response := g2pb.GetFeatureResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetGenericFeatures(ctx context.Context, request *g2pb.GetGenericFeaturesRequest) (*g2pb.GetGenericFeaturesResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(29, request)
		defer func() { server.traceExit(30, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetGenericFeatures(ctx, request.GetFeatureType(), int(request.GetMaximumEstimatedCount()))
	response := g2pb.GetGenericFeaturesResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetLogicalCores(ctx context.Context, request *g2pb.GetLogicalCoresRequest) (*g2pb.GetLogicalCoresResponse, error) {
	var err error = nil
	var result int
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(35, request)
		defer func() { server.traceExit(36, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetLogicalCores(ctx)
	response := g2pb.GetLogicalCoresResponse{
		Result: int32(result),
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetMappingStatistics(ctx context.Context, request *g2pb.GetMappingStatisticsRequest) (*g2pb.GetMappingStatisticsResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(37, request)
		defer func() { server.traceExit(38, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetMappingStatistics(ctx, int(request.GetIncludeInternalFeatures()))
	response := g2pb.GetMappingStatisticsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetPhysicalCores(ctx context.Context, request *g2pb.GetPhysicalCoresRequest) (*g2pb.GetPhysicalCoresResponse, error) {
	var err error = nil
	var result int
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(39, request)
		defer func() { server.traceExit(40, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetPhysicalCores(ctx)
	response := g2pb.GetPhysicalCoresResponse{
		Result: int32(result),
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetRelationshipDetails(ctx context.Context, request *g2pb.GetRelationshipDetailsRequest) (*g2pb.GetRelationshipDetailsResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(41, request)
		defer func() { server.traceExit(42, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetRelationshipDetails(ctx, request.GetRelationshipID(), int(request.GetIncludeInternalFeatures()))
	response := g2pb.GetRelationshipDetailsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetResolutionStatistics(ctx context.Context, request *g2pb.GetResolutionStatisticsRequest) (*g2pb.GetResolutionStatisticsResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(43, request)
		defer func() { server.traceExit(44, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetResolutionStatistics(ctx)
	response := g2pb.GetResolutionStatisticsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetTotalSystemMemory(ctx context.Context, request *g2pb.GetTotalSystemMemoryRequest) (*g2pb.GetTotalSystemMemoryResponse, error) {
	var err error = nil
	var result int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(45, request)
		defer func() { server.traceExit(46, request, result, err, time.Since(entryTime)) }()
	}
	g2diagnostic := getG2diagnostic()
	result, err = g2diagnostic.GetTotalSystemMemory(ctx)
	response := g2pb.GetTotalSystemMemoryResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *G2DiagnosticServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
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

func (server *G2DiagnosticServer) InitWithConfigID(ctx context.Context, request *g2pb.InitWithConfigIDRequest) (*g2pb.InitWithConfigIDResponse, error) {
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

func (server G2DiagnosticServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	g2diagnostic := getG2diagnostic()
	return g2diagnostic.RegisterObserver(ctx, observer)
}

func (server *G2DiagnosticServer) Reinit(ctx context.Context, request *g2pb.ReinitRequest) (*g2pb.ReinitResponse, error) {
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

func (server *G2DiagnosticServer) SetLogLevel(ctx context.Context, logLevelName string) error {
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

	// TODO: Remove once g2configmgr.SetLogLevel(context.Context, string)
	logLevel := logging.TextToLoggerLevelMap[logLevelName]

	g2diagnostic.SetLogLevel(ctx, logLevel)
	server.getLogger().SetLogLevel(logLevelName)
	server.isTrace = (logLevelName == logging.LevelTraceName)
	return err
}

func (server *G2DiagnosticServer) StreamEntityListBySize(request *g2pb.StreamEntityListBySizeRequest, stream g2pb.G2Diagnostic_StreamEntityListBySizeServer) (err error) {
	if server.isTrace {
		server.traceEntry(163, request)
	}
	ctx := stream.Context()
	entryTime := time.Now()
	g2diagnostic := getG2diagnostic()

	entitiesFetched := 0

	//get the query handle
	var queryHandle uintptr
	queryHandle, err = g2diagnostic.GetEntityListBySize(ctx, int(request.GetEntitySize()))
	if err != nil {
		return err
	}

	defer func() {
		err = g2diagnostic.CloseEntityListBySize(ctx, queryHandle)
		if server.isTrace {
			server.traceExit(165, request, entitiesFetched, err, time.Since(entryTime))
		}
	}()

	for {
		var fetchResult string
		fetchResult, err = g2diagnostic.FetchNextEntityBySize(ctx, queryHandle)
		if err != nil {
			return err
		}
		if len(fetchResult) == 0 {
			break
		}
		response := g2pb.StreamEntityListBySizeResponse{
			Result: fetchResult,
		}
		entitiesFetched += 1
		if err = stream.Send(&response); err != nil {
			return err
		}
		server.traceEntry(164, request, fetchResult)
	}

	err = nil
	return
}

func (server *G2DiagnosticServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	g2diagnostic := getG2diagnostic()
	return g2diagnostic.UnregisterObserver(ctx, observer)
}

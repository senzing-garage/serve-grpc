package g2diagnosticserver

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go/g2diagnostic"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/go-servegrpc/protobuf/g2diagnostic"
)

var (
	g2diagnosticSingleton *g2sdk.G2diagnosticImpl
	// logger           messagelogger.MessageLoggerInterface
	// onceLogger       sync.Once
	g2diagnosticSyncOnce sync.Once
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
func getG2diagnostic() *g2sdk.G2diagnosticImpl {
	g2diagnosticSyncOnce.Do(func() {
		g2diagnosticSingleton = &g2sdk.G2diagnosticImpl{}
	})
	return g2diagnosticSingleton
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
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2DiagnosticServer) CheckDBPerf(ctx context.Context, request *pb.CheckDBPerfRequest) (*pb.CheckDBPerfResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.CheckDBPerf(ctx, int(request.GetSecondsToRun()))
	response := pb.CheckDBPerfResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) CloseEntityListBySize(ctx context.Context, request *pb.CloseEntityListBySizeRequest) (*pb.CloseEntityListBySizeResponse, error) {
	g2diagnostic := getG2diagnostic()
	entityListBySizeHandleInt, err := strconv.ParseUint(request.GetEntityListBySizeHandle(), 10, 64)
	if err == nil {
		err = g2diagnostic.CloseEntityListBySize(ctx, uintptr(entityListBySizeHandleInt))
	}
	response := pb.CloseEntityListBySizeResponse{}
	return &response, err
}

func (server *G2DiagnosticServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.Destroy(ctx)
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2DiagnosticServer) FetchNextEntityBySize(ctx context.Context, request *pb.FetchNextEntityBySizeRequest) (*pb.FetchNextEntityBySizeResponse, error) {
	var result string = ""
	g2diagnostic := getG2diagnostic()
	entityListBySizeHandleInt, err := strconv.ParseUint(request.GetEntityListBySizeHandle(), 10, 64)
	if err == nil {
		result, err = g2diagnostic.FetchNextEntityBySize(ctx, uintptr(entityListBySizeHandleInt))

	}
	response := pb.FetchNextEntityBySizeResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) FindEntitiesByFeatureIDs(ctx context.Context, request *pb.FindEntitiesByFeatureIDsRequest) (*pb.FindEntitiesByFeatureIDsResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, request.GetFeatures())
	response := pb.FindEntitiesByFeatureIDsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetAvailableMemory(ctx context.Context, request *pb.GetAvailableMemoryRequest) (*pb.GetAvailableMemoryResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetAvailableMemory(ctx)
	response := pb.GetAvailableMemoryResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetDataSourceCounts(ctx context.Context, request *pb.GetDataSourceCountsRequest) (*pb.GetDataSourceCountsResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetDataSourceCounts(ctx)
	response := pb.GetDataSourceCountsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetDBInfo(ctx context.Context, request *pb.GetDBInfoRequest) (*pb.GetDBInfoResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetDBInfo(ctx)
	response := pb.GetDBInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityDetails(ctx context.Context, request *pb.GetEntityDetailsRequest) (*pb.GetEntityDetailsResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntityDetails(ctx, request.GetEntityID(), int(request.GetIncludeInternalFeatures()))
	response := pb.GetEntityDetailsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityListBySize(ctx context.Context, request *pb.GetEntityListBySizeRequest) (*pb.GetEntityListBySizeResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntityListBySize(ctx, int(request.GetEntitySize()))
	response := pb.GetEntityListBySizeResponse{
		Result: fmt.Sprintf("%v", result),
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityResume(ctx context.Context, request *pb.GetEntityResumeRequest) (*pb.GetEntityResumeResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntityResume(ctx, request.GetEntityID())
	response := pb.GetEntityResumeResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetEntitySizeBreakdown(ctx context.Context, request *pb.GetEntitySizeBreakdownRequest) (*pb.GetEntitySizeBreakdownResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntitySizeBreakdown(ctx, int(request.GetMinimumEntitySize()), int(request.GetIncludeInternalFeatures()))
	response := pb.GetEntitySizeBreakdownResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetFeature(ctx context.Context, request *pb.GetFeatureRequest) (*pb.GetFeatureResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetFeature(ctx, request.GetLibFeatID())
	response := pb.GetFeatureResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetGenericFeatures(ctx context.Context, request *pb.GetGenericFeaturesRequest) (*pb.GetGenericFeaturesResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetGenericFeatures(ctx, request.GetFeatureType(), int(request.GetMaximumEstimatedCount()))
	response := pb.GetGenericFeaturesResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetLogicalCores(ctx context.Context, request *pb.GetLogicalCoresRequest) (*pb.GetLogicalCoresResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetLogicalCores(ctx)
	response := pb.GetLogicalCoresResponse{
		Result: int32(result),
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetMappingStatistics(ctx context.Context, request *pb.GetMappingStatisticsRequest) (*pb.GetMappingStatisticsResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetMappingStatistics(ctx, int(request.GetIncludeInternalFeatures()))
	response := pb.GetMappingStatisticsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetPhysicalCores(ctx context.Context, request *pb.GetPhysicalCoresRequest) (*pb.GetPhysicalCoresResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetPhysicalCores(ctx)
	response := pb.GetPhysicalCoresResponse{
		Result: int32(result),
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetRelationshipDetails(ctx context.Context, request *pb.GetRelationshipDetailsRequest) (*pb.GetRelationshipDetailsResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetRelationshipDetails(ctx, request.GetRelationshipID(), int(request.GetIncludeInternalFeatures()))
	response := pb.GetRelationshipDetailsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetResolutionStatistics(ctx context.Context, request *pb.GetResolutionStatisticsRequest) (*pb.GetResolutionStatisticsResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetResolutionStatistics(ctx)
	response := pb.GetResolutionStatisticsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2DiagnosticServer) GetTotalSystemMemory(ctx context.Context, request *pb.GetTotalSystemMemoryRequest) (*pb.GetTotalSystemMemoryResponse, error) {
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetTotalSystemMemory(ctx)
	response := pb.GetTotalSystemMemoryResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *G2DiagnosticServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2DiagnosticServer) InitWithConfigID(ctx context.Context, request *pb.InitWithConfigIDRequest) (*pb.InitWithConfigIDResponse, error) {
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.InitWithConfigID(ctx, request.GetModuleName(), request.GetIniParams(), int64(request.GetInitConfigID()), int(request.GetVerboseLogging()))
	response := pb.InitWithConfigIDResponse{}
	return &response, err
}

func (server *G2DiagnosticServer) Reinit(ctx context.Context, request *pb.ReinitRequest) (*pb.ReinitResponse, error) {
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.Reinit(ctx, int64(request.GetInitConfigID()))
	response := pb.ReinitResponse{}
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
		server.traceEntry(1, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	g2diagnostic := getG2diagnostic()
	g2diagnostic.SetLogLevel(ctx, logLevel)
	server.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	server.isTrace = (server.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if server.isTrace {
		defer server.traceExit(1, logLevel, err, time.Since(entryTime))
	}
	return err
}

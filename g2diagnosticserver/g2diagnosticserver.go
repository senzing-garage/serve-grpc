package g2diagnosticserver

import (
	"context"
	"fmt"
	"sync"

	g2diagnosticsdk "github.com/senzing/g2-sdk-go/g2diagnostic"
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/go-servegrpc/g2diagnosticprotobuf"
)

var (
	g2diagnostic *g2diagnosticsdk.G2diagnosticImpl
	logger       messagelogger.MessageLoggerInterface
	once         sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods - names begin with lower case
// ----------------------------------------------------------------------------

func getLogger() messagelogger.MessageLoggerInterface {

	once.Do(func() {
		logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	})
	return logger
}

// Singleton pattern for g2diagnostic.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2diagnostic() *g2diagnosticsdk.G2diagnosticImpl {
	once.Do(func() {
		g2diagnostic = &g2diagnosticsdk.G2diagnosticImpl{}
	})
	return g2diagnostic
}

func traceEnter(messageNumber int, request interface{}) {
	if logger.IsTrace() {

	}
}

func traceExit(messageNumber int, response interface{}) {
	if logger.IsTrace() {

	}
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2DiagnosticServer) CheckDBPerf(ctx context.Context, request *pb.CheckDBPerfRequest) (*pb.CheckDBPerfResponse, error) {
	traceEnter(5100, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.CheckDBPerf(ctx, int(request.SecondsToRun))
	response := pb.CheckDBPerfResponse{
		Result: result,
	}
	traceExit(5101, response)
	return &response, err
}

func (server *G2DiagnosticServer) ClearLastException(ctx context.Context, request *pb.ClearLastExceptionRequest) (*pb.ClearLastExceptionResponse, error) {
	traceEnter(5102, request)
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.ClearLastException(ctx)
	response := pb.ClearLastExceptionResponse{}
	traceExit(5103, response)
	return &response, err
}

func (server *G2DiagnosticServer) CloseEntityListBySize(ctx context.Context, request *pb.CloseEntityListBySizeRequest) (*pb.CloseEntityListBySizeResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.CloseEntityListBySize(ctx, request.EntityListBySizeHandle)
	response := pb.CloseEntityListBySizeResponse{}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.Destroy(ctx)
	response := pb.DestroyResponse{}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) FetchNextEntityBySize(ctx context.Context, request *pb.FetchNextEntityBySizeRequest) (*pb.FetchNextEntityBySizeResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.FetchNextEntityBySize(ctx, request.EntityListBySizeHandle)
	response := pb.FetchNextEntityBySizeResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) FindEntitiesByFeatureIDs(ctx context.Context, request *pb.FindEntitiesByFeatureIDsRequest) (*pb.FindEntitiesByFeatureIDsResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, request.Features)
	response := pb.FindEntitiesByFeatureIDsResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetAvailableMemory(ctx context.Context, request *pb.GetAvailableMemoryRequest) (*pb.GetAvailableMemoryResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetAvailableMemory(ctx)
	response := pb.GetAvailableMemoryResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetDataSourceCounts(ctx context.Context, request *pb.GetDataSourceCountsRequest) (*pb.GetDataSourceCountsResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetDataSourceCounts(ctx)
	response := pb.GetDataSourceCountsResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetDBInfo(ctx context.Context, request *pb.GetDBInfoRequest) (*pb.GetDBInfoResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetDBInfo(ctx)
	response := pb.GetDBInfoResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityDetails(ctx context.Context, request *pb.GetEntityDetailsRequest) (*pb.GetEntityDetailsResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntityDetails(ctx, request.EntityID, int(request.IncludeInternalFeatures))
	response := pb.GetEntityDetailsResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityListBySize(ctx context.Context, request *pb.GetEntityListBySizeRequest) (*pb.GetEntityListBySizeResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntityListBySize(ctx, int(request.EntitySize))
	response := pb.GetEntityListBySizeResponse{
		Result: fmt.Sprintf("%v", result),
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetEntityResume(ctx context.Context, request *pb.GetEntityResumeRequest) (*pb.GetEntityResumeResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntityResume(ctx, request.EntityID)
	response := pb.GetEntityResumeResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetEntitySizeBreakdown(ctx context.Context, request *pb.GetEntitySizeBreakdownRequest) (*pb.GetEntitySizeBreakdownResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetEntitySizeBreakdown(ctx, int(request.MinimumEntitySize), int(request.IncludeInternalFeatures))
	response := pb.GetEntitySizeBreakdownResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetFeature(ctx context.Context, request *pb.GetFeatureRequest) (*pb.GetFeatureResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetFeature(ctx, request.LibFeatID)
	response := pb.GetFeatureResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetGenericFeatures(ctx context.Context, request *pb.GetGenericFeaturesRequest) (*pb.GetGenericFeaturesResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetGenericFeatures(ctx, request.FeatureType, int(request.MaximumEstimatedCount))
	response := pb.GetGenericFeaturesResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetLogicalCores(ctx context.Context, request *pb.GetLogicalCoresRequest) (*pb.GetLogicalCoresResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetLogicalCores(ctx)
	response := pb.GetLogicalCoresResponse{
		Result: int32(result),
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetMappingStatistics(ctx context.Context, request *pb.GetMappingStatisticsRequest) (*pb.GetMappingStatisticsResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetMappingStatistics(ctx, int(request.IncludeInternalFeatures))
	response := pb.GetMappingStatisticsResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetPhysicalCores(ctx context.Context, request *pb.GetPhysicalCoresRequest) (*pb.GetPhysicalCoresResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetPhysicalCores(ctx)
	response := pb.GetPhysicalCoresResponse{
		Result: int32(result),
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetRelationshipDetails(ctx context.Context, request *pb.GetRelationshipDetailsRequest) (*pb.GetRelationshipDetailsResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetRelationshipDetails(ctx, request.RelationshipID, int(request.IncludeInternalFeatures))
	response := pb.GetRelationshipDetailsResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetResolutionStatistics(ctx context.Context, request *pb.GetResolutionStatisticsRequest) (*pb.GetResolutionStatisticsResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetResolutionStatistics(ctx)
	response := pb.GetResolutionStatisticsResponse{
		Result: result,
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) GetTotalSystemMemory(ctx context.Context, request *pb.GetTotalSystemMemoryRequest) (*pb.GetTotalSystemMemoryResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	result, err := g2diagnostic.GetTotalSystemMemory(ctx)
	response := pb.GetTotalSystemMemoryResponse{
		Result: int64(result),
	}
	traceExit(5199, response)
	return &response, err
}

func (server *G2DiagnosticServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.Init(ctx, request.ModuleName, request.IniParams, int(request.VerboseLogging))
	response := pb.InitResponse{}
	traceExit(5156, response)
	return &response, err
}

func (server *G2DiagnosticServer) InitWithConfigID(ctx context.Context, request *pb.InitWithConfigIDRequest) (*pb.InitWithConfigIDResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.InitWithConfigID(ctx, request.ModuleName, request.IniParams, int64(request.InitConfigID), int(request.VerboseLogging))
	response := pb.InitWithConfigIDResponse{}
	traceExit(5156, response)
	return &response, err
}

func (server *G2DiagnosticServer) Reinit(ctx context.Context, request *pb.ReinitRequest) (*pb.ReinitResponse, error) {
	traceEnter(5198, request)
	g2diagnostic := getG2diagnostic()
	err := g2diagnostic.Reinit(ctx, int64(request.InitConfigID))
	response := pb.ReinitResponse{}
	traceExit(5156, response)
	return &response, err
}

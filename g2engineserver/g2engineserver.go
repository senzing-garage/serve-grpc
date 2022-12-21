package g2engineserver

import (
	"context"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go/g2engine"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/go-servegrpc/protobuf/g2engine"
)

var (
	g2engineSingleton *g2sdk.G2engineImpl
	g2engineSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2engine() *g2sdk.G2engineImpl {
	g2engineSyncOnce.Do(func() {
		g2engineSingleton = &g2sdk.G2engineImpl{}
	})
	return g2engineSingleton
}

// Get the Logger singleton.
func (server *G2EngineServer) getLogger() messagelogger.MessageLoggerInterface {
	if server.logger == nil {
		server.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return server.logger
}

// Trace method entry.
func (server *G2EngineServer) traceEntry(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (server *G2EngineServer) traceExit(errorNumber int, details ...interface{}) {
	server.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2EngineServer) AddRecord(ctx context.Context, request *pb.AddRecordRequest) (*pb.AddRecordResponse, error) {
	g2engine := getG2engine()
	err := g2engine.AddRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID())
	response := pb.AddRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithInfo(ctx context.Context, request *pb.AddRecordWithInfoRequest) (*pb.AddRecordWithInfoResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.AddRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := pb.AddRecordWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithInfoWithReturnedRecordID(ctx context.Context, request *pb.AddRecordWithInfoWithReturnedRecordIDRequest) (*pb.AddRecordWithInfoWithReturnedRecordIDResponse, error) {
	g2engine := getG2engine()
	result, recordId, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, request.GetDataSourceCode(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := pb.AddRecordWithInfoWithReturnedRecordIDResponse{
		RecordID: recordId,
		WithInfo: result,
	}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithReturnedRecordID(ctx context.Context, request *pb.AddRecordWithReturnedRecordIDRequest) (*pb.AddRecordWithReturnedRecordIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.AddRecordWithReturnedRecordID(ctx, request.GetDataSourceCode(), request.GetJsonData(), request.GetLoadID())
	response := pb.AddRecordWithReturnedRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) CheckRecord(ctx context.Context, request *pb.CheckRecordRequest) (*pb.CheckRecordResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.CheckRecord(ctx, request.GetRecord(), request.GetRecordQueryList())
	response := pb.CheckRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) CloseExport(ctx context.Context, request *pb.CloseExportRequest) (*pb.CloseExportResponse, error) {
	g2engine := getG2engine()
	err := g2engine.CloseExport(ctx, uintptr(request.GetResponseHandle()))
	response := pb.CloseExportResponse{}
	return &response, err
}

func (server *G2EngineServer) CountRedoRecords(ctx context.Context, request *pb.CountRedoRecordsRequest) (*pb.CountRedoRecordsResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.CountRedoRecords(ctx)
	response := pb.CountRedoRecordsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) DeleteRecord(ctx context.Context, request *pb.DeleteRecordRequest) (*pb.DeleteRecordResponse, error) {
	g2engine := getG2engine()
	err := g2engine.DeleteRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetLoadID())
	response := pb.DeleteRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) DeleteRecordWithInfo(ctx context.Context, request *pb.DeleteRecordWithInfoRequest) (*pb.DeleteRecordWithInfoResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.DeleteRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetLoadID(), request.GetFlags())
	response := pb.DeleteRecordWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	g2engine := getG2engine()
	err := g2engine.Destroy(ctx)
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2EngineServer) ExportConfig(ctx context.Context, request *pb.ExportConfigRequest) (*pb.ExportConfigResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ExportConfig(ctx)
	response := pb.ExportConfigResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) ExportConfigAndConfigID(ctx context.Context, request *pb.ExportConfigAndConfigIDRequest) (*pb.ExportConfigAndConfigIDResponse, error) {
	g2engine := getG2engine()
	result, configId, err := g2engine.ExportConfigAndConfigID(ctx)
	response := pb.ExportConfigAndConfigIDResponse{
		Config:   result,
		ConfigID: configId,
	}
	return &response, err
}

func (server *G2EngineServer) ExportCSVEntityReport(ctx context.Context, request *pb.ExportCSVEntityReportRequest) (*pb.ExportCSVEntityReportResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ExportCSVEntityReport(ctx, request.GetCsvColumnList(), request.GetFlags())
	response := pb.ExportCSVEntityReportResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *G2EngineServer) ExportJSONEntityReport(ctx context.Context, request *pb.ExportJSONEntityReportRequest) (*pb.ExportJSONEntityReportResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ExportJSONEntityReport(ctx, request.GetFlags())
	response := pb.ExportJSONEntityReportResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *G2EngineServer) FetchNext(ctx context.Context, request *pb.FetchNextRequest) (*pb.FetchNextResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FetchNext(ctx, uintptr(request.GetResponseHandle()))
	response := pb.FetchNextResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindInterestingEntitiesByEntityID(ctx context.Context, request *pb.FindInterestingEntitiesByEntityIDRequest) (*pb.FindInterestingEntitiesByEntityIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindInterestingEntitiesByEntityID(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.FindInterestingEntitiesByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindInterestingEntitiesByRecordID(ctx context.Context, request *pb.FindInterestingEntitiesByRecordIDRequest) (*pb.FindInterestingEntitiesByRecordIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindInterestingEntitiesByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.FindInterestingEntitiesByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByEntityID(ctx context.Context, request *pb.FindNetworkByEntityIDRequest) (*pb.FindNetworkByEntityIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByEntityID(ctx, request.GetEntityList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()))
	response := pb.FindNetworkByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByEntityID_V2(ctx context.Context, request *pb.FindNetworkByEntityID_V2Request) (*pb.FindNetworkByEntityID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByEntityID_V2(ctx, request.GetEntityList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()), request.GetFlags())
	response := pb.FindNetworkByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByRecordID(ctx context.Context, request *pb.FindNetworkByRecordIDRequest) (*pb.FindNetworkByRecordIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByRecordID(ctx, request.GetRecordList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()))
	response := pb.FindNetworkByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByRecordID_V2(ctx context.Context, request *pb.FindNetworkByRecordID_V2Request) (*pb.FindNetworkByRecordID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByRecordID_V2(ctx, request.GetRecordList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()), request.GetFlags())
	response := pb.FindNetworkByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByEntityID(ctx context.Context, request *pb.FindPathByEntityIDRequest) (*pb.FindPathByEntityIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()))
	response := pb.FindPathByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByEntityID_V2(ctx context.Context, request *pb.FindPathByEntityID_V2Request) (*pb.FindPathByEntityID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetFlags())
	response := pb.FindPathByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByRecordID(ctx context.Context, request *pb.FindPathByRecordIDRequest) (*pb.FindPathByRecordIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()))
	response := pb.FindPathByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByRecordID_V2(ctx context.Context, request *pb.FindPathByRecordID_V2Request) (*pb.FindPathByRecordID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetFlags())
	response := pb.FindPathByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByEntityID(ctx context.Context, request *pb.FindPathExcludingByEntityIDRequest) (*pb.FindPathExcludingByEntityIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities())
	response := pb.FindPathExcludingByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByEntityID_V2(ctx context.Context, request *pb.FindPathExcludingByEntityID_V2Request) (*pb.FindPathExcludingByEntityID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities(), request.GetFlags())
	response := pb.FindPathExcludingByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByRecordID(ctx context.Context, request *pb.FindPathExcludingByRecordIDRequest) (*pb.FindPathExcludingByRecordIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords())
	response := pb.FindPathExcludingByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByRecordID_V2(ctx context.Context, request *pb.FindPathExcludingByRecordID_V2Request) (*pb.FindPathExcludingByRecordID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords(), request.GetFlags())
	response := pb.FindPathExcludingByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByEntityID(ctx context.Context, request *pb.FindPathIncludingSourceByEntityIDRequest) (*pb.FindPathIncludingSourceByEntityIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities(), request.GetRequiredDsrcs())
	response := pb.FindPathIncludingSourceByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByEntityID_V2(ctx context.Context, request *pb.FindPathIncludingSourceByEntityID_V2Request) (*pb.FindPathIncludingSourceByEntityID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities(), request.GetRequiredDsrcs(), request.GetFlags())
	response := pb.FindPathIncludingSourceByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByRecordID(ctx context.Context, request *pb.FindPathIncludingSourceByRecordIDRequest) (*pb.FindPathIncludingSourceByRecordIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords(), request.GetRequiredDsrcs())
	response := pb.FindPathIncludingSourceByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByRecordID_V2(ctx context.Context, request *pb.FindPathIncludingSourceByRecordID_V2Request) (*pb.FindPathIncludingSourceByRecordID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords(), request.GetRequiredDsrcs(), request.GetFlags())
	response := pb.FindPathIncludingSourceByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetActiveConfigID(ctx context.Context, request *pb.GetActiveConfigIDRequest) (*pb.GetActiveConfigIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetActiveConfigID(ctx)
	response := pb.GetActiveConfigIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByEntityID(ctx context.Context, request *pb.GetEntityByEntityIDRequest) (*pb.GetEntityByEntityIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByEntityID(ctx, request.GetEntityID())
	response := pb.GetEntityByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByEntityID_V2(ctx context.Context, request *pb.GetEntityByEntityID_V2Request) (*pb.GetEntityByEntityID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.GetEntityByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByRecordID(ctx context.Context, request *pb.GetEntityByRecordIDRequest) (*pb.GetEntityByRecordIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := pb.GetEntityByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByRecordID_V2(ctx context.Context, request *pb.GetEntityByRecordID_V2Request) (*pb.GetEntityByRecordID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByRecordID_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.GetEntityByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetRecord(ctx context.Context, request *pb.GetRecordRequest) (*pb.GetRecordResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetRecord(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := pb.GetRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetRecord_V2(ctx context.Context, request *pb.GetRecord_V2Request) (*pb.GetRecord_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetRecord_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.GetRecord_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetRedoRecord(ctx context.Context, request *pb.GetRedoRecordRequest) (*pb.GetRedoRecordResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetRedoRecord(ctx)
	response := pb.GetRedoRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetRepositoryLastModifiedTime(ctx context.Context, request *pb.GetRepositoryLastModifiedTimeRequest) (*pb.GetRepositoryLastModifiedTimeResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetRepositoryLastModifiedTime(ctx)
	response := pb.GetRepositoryLastModifiedTimeResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetVirtualEntityByRecordID(ctx context.Context, request *pb.GetVirtualEntityByRecordIDRequest) (*pb.GetVirtualEntityByRecordIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetVirtualEntityByRecordID(ctx, request.GetRecordList())
	response := pb.GetVirtualEntityByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) GetVirtualEntityByRecordID_V2(ctx context.Context, request *pb.GetVirtualEntityByRecordID_V2Request) (*pb.GetVirtualEntityByRecordID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, request.GetRecordList(), request.GetFlags())
	response := pb.GetVirtualEntityByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) HowEntityByEntityID(ctx context.Context, request *pb.HowEntityByEntityIDRequest) (*pb.HowEntityByEntityIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.HowEntityByEntityID(ctx, request.GetEntityID())
	response := pb.HowEntityByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) HowEntityByEntityID_V2(ctx context.Context, request *pb.HowEntityByEntityID_V2Request) (*pb.HowEntityByEntityID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.HowEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.HowEntityByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	g2engine := getG2engine()
	err := g2engine.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2EngineServer) InitWithConfigID(ctx context.Context, request *pb.InitWithConfigIDRequest) (*pb.InitWithConfigIDResponse, error) {
	g2engine := getG2engine()
	err := g2engine.InitWithConfigID(ctx, request.GetModuleName(), request.GetIniParams(), request.GetInitConfigID(), int(request.GetVerboseLogging()))
	response := pb.InitWithConfigIDResponse{}
	return &response, err
}

func (server *G2EngineServer) PrimeEngine(ctx context.Context, request *pb.PrimeEngineRequest) (*pb.PrimeEngineResponse, error) {
	g2engine := getG2engine()
	err := g2engine.PrimeEngine(ctx)
	response := pb.PrimeEngineResponse{}
	return &response, err
}

func (server *G2EngineServer) Process(ctx context.Context, request *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	g2engine := getG2engine()
	err := g2engine.Process(ctx, request.GetRecord())
	response := pb.ProcessResponse{}
	return &response, err
}

func (server *G2EngineServer) ProcessRedoRecord(ctx context.Context, request *pb.ProcessRedoRecordRequest) (*pb.ProcessRedoRecordResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ProcessRedoRecord(ctx)
	response := pb.ProcessRedoRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) ProcessRedoRecordWithInfo(ctx context.Context, request *pb.ProcessRedoRecordWithInfoRequest) (*pb.ProcessRedoRecordWithInfoResponse, error) {
	g2engine := getG2engine()
	result, withInfo, err := g2engine.ProcessRedoRecordWithInfo(ctx, request.GetFlags())
	response := pb.ProcessRedoRecordWithInfoResponse{
		Result:   result,
		WithInfo: withInfo,
	}
	return &response, err
}

func (server *G2EngineServer) ProcessWithInfo(ctx context.Context, request *pb.ProcessWithInfoRequest) (*pb.ProcessWithInfoResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ProcessWithInfo(ctx, request.GetRecord(), request.GetFlags())
	response := pb.ProcessWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) ProcessWithResponse(ctx context.Context, request *pb.ProcessWithResponseRequest) (*pb.ProcessWithResponseResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ProcessWithResponse(ctx, request.GetRecord())
	response := pb.ProcessWithResponseResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) ProcessWithResponseResize(ctx context.Context, request *pb.ProcessWithResponseResizeRequest) (*pb.ProcessWithResponseResizeResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ProcessWithResponseResize(ctx, request.GetRecord())
	response := pb.ProcessWithResponseResizeResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) PurgeRepository(ctx context.Context, request *pb.PurgeRepositoryRequest) (*pb.PurgeRepositoryResponse, error) {
	g2engine := getG2engine()
	err := g2engine.PurgeRepository(ctx)
	response := pb.PurgeRepositoryResponse{}
	return &response, err
}

func (server *G2EngineServer) ReevaluateEntity(ctx context.Context, request *pb.ReevaluateEntityRequest) (*pb.ReevaluateEntityResponse, error) {
	g2engine := getG2engine()
	err := g2engine.ReevaluateEntity(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.ReevaluateEntityResponse{}
	return &response, err
}

func (server *G2EngineServer) ReevaluateEntityWithInfo(ctx context.Context, request *pb.ReevaluateEntityWithInfoRequest) (*pb.ReevaluateEntityWithInfoResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ReevaluateEntityWithInfo(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.ReevaluateEntityWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) ReevaluateRecord(ctx context.Context, request *pb.ReevaluateRecordRequest) (*pb.ReevaluateRecordResponse, error) {
	g2engine := getG2engine()
	err := g2engine.ReevaluateRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.ReevaluateRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) ReevaluateRecordWithInfo(ctx context.Context, request *pb.ReevaluateRecordWithInfoRequest) (*pb.ReevaluateRecordWithInfoResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ReevaluateRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.ReevaluateRecordWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) Reinit(ctx context.Context, request *pb.ReinitRequest) (*pb.ReinitResponse, error) {
	g2engine := getG2engine()
	err := g2engine.Reinit(ctx, request.GetInitConfigID())
	response := pb.ReinitResponse{}
	return &response, err
}

func (server *G2EngineServer) ReplaceRecord(ctx context.Context, request *pb.ReplaceRecordRequest) (*pb.ReplaceRecordResponse, error) {
	g2engine := getG2engine()
	err := g2engine.ReplaceRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID())
	response := pb.ReplaceRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) ReplaceRecordWithInfo(ctx context.Context, request *pb.ReplaceRecordWithInfoRequest) (*pb.ReplaceRecordWithInfoResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.ReplaceRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := pb.ReplaceRecordWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) SearchByAttributes(ctx context.Context, request *pb.SearchByAttributesRequest) (*pb.SearchByAttributesResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.SearchByAttributes(ctx, request.GetJsonData())
	response := pb.SearchByAttributesResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) SearchByAttributes_V2(ctx context.Context, request *pb.SearchByAttributes_V2Request) (*pb.SearchByAttributes_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.SearchByAttributes_V2(ctx, request.GetJsonData(), request.GetFlags())
	response := pb.SearchByAttributes_V2Response{
		Result: result,
	}
	return &response, err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (server *G2EngineServer) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if server.isTrace {
		server.traceEntry(1, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	server.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	server.isTrace = (server.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if server.isTrace {
		defer server.traceExit(1, logLevel, err, time.Since(entryTime))
	}
	return err
}

func (server *G2EngineServer) Stats(ctx context.Context, request *pb.StatsRequest) (*pb.StatsResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.Stats(ctx)
	response := pb.StatsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntities(ctx context.Context, request *pb.WhyEntitiesRequest) (*pb.WhyEntitiesResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.WhyEntities(ctx, request.GetEntityID1(), request.GetEntityID2())
	response := pb.WhyEntitiesResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntities_V2(ctx context.Context, request *pb.WhyEntities_V2Request) (*pb.WhyEntities_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.WhyEntities_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetFlags())
	response := pb.WhyEntities_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByEntityID(ctx context.Context, request *pb.WhyEntityByEntityIDRequest) (*pb.WhyEntityByEntityIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByEntityID(ctx, request.GetEntityID())
	response := pb.WhyEntityByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByEntityID_V2(ctx context.Context, request *pb.WhyEntityByEntityID_V2Request) (*pb.WhyEntityByEntityID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.WhyEntityByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByRecordID(ctx context.Context, request *pb.WhyEntityByRecordIDRequest) (*pb.WhyEntityByRecordIDResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := pb.WhyEntityByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByRecordID_V2(ctx context.Context, request *pb.WhyEntityByRecordID_V2Request) (*pb.WhyEntityByRecordID_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByRecordID_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.WhyEntityByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) WhyRecords(ctx context.Context, request *pb.WhyRecordsRequest) (*pb.WhyRecordsResponse, error) {
	g2engine := getG2engine()
	result, err := g2engine.WhyRecords(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2())
	response := pb.WhyRecordsResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2EngineServer) WhyRecords_V2(ctx context.Context, request *pb.WhyRecords_V2Request) (*pb.WhyRecords_V2Response, error) {
	g2engine := getG2engine()
	result, err := g2engine.WhyRecords_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetFlags())
	response := pb.WhyRecords_V2Response{
		Result: result,
	}
	return &response, err
}

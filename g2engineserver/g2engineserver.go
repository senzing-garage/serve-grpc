package g2engineserver

import (
	"context"
	"sync"

	g2sdk "github.com/senzing/g2-sdk-go/g2engine"
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
func getG2config() *g2sdk.G2engineImpl {
	g2engineSyncOnce.Do(func() {
		g2engineSingleton = &g2sdk.G2engineImpl{}
	})
	return g2engineSingleton
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2EngineServer) AddRecord(ctx context.Context, request *pb.AddRecordRequest) (*pb.AddRecordResponse, error) {
	var err error = nil
	response := pb.AddRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithInfo(ctx context.Context, request *pb.AddRecordWithInfoRequest) (*pb.AddRecordWithInfoResponse, error) {
	var err error = nil
	response := pb.AddRecordWithInfoResponse{}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithInfoWithReturnedRecordID(ctx context.Context, request *pb.AddRecordWithInfoWithReturnedRecordIDRequest) (*pb.AddRecordWithInfoWithReturnedRecordIDResponse, error) {
	var err error = nil
	response := pb.AddRecordWithInfoWithReturnedRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithReturnedRecordID(ctx context.Context, request *pb.AddRecordWithReturnedRecordIDRequest) (*pb.AddRecordWithReturnedRecordIDResponse, error) {
	var err error = nil
	response := pb.AddRecordWithReturnedRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) CheckRecord(ctx context.Context, request *pb.CheckRecordRequest) (*pb.CheckRecordResponse, error) {
	var err error = nil
	response := pb.CheckRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) CloseExport(ctx context.Context, request *pb.CloseExportRequest) (*pb.CloseExportResponse, error) {
	var err error = nil
	response := pb.CloseExportResponse{}
	return &response, err
}

func (server *G2EngineServer) CountRedoRecords(ctx context.Context, request *pb.CountRedoRecordsRequest) (*pb.CountRedoRecordsResponse, error) {
	var err error = nil
	response := pb.CountRedoRecordsResponse{}
	return &response, err
}

func (server *G2EngineServer) DeleteRecord(ctx context.Context, request *pb.DeleteRecordRequest) (*pb.DeleteRecordResponse, error) {
	var err error = nil
	response := pb.DeleteRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) DeleteRecordWithInfo(ctx context.Context, request *pb.DeleteRecordWithInfoRequest) (*pb.DeleteRecordWithInfoResponse, error) {
	var err error = nil
	response := pb.DeleteRecordWithInfoResponse{}
	return &response, err
}

func (server *G2EngineServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	var err error = nil
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2EngineServer) ExportConfig(ctx context.Context, request *pb.ExportConfigRequest) (*pb.ExportConfigResponse, error) {
	var err error = nil
	response := pb.ExportConfigResponse{}
	return &response, err
}

func (server *G2EngineServer) ExportConfigAndConfigID(ctx context.Context, request *pb.ExportConfigAndConfigIDRequest) (*pb.ExportConfigAndConfigIDResponse, error) {
	var err error = nil
	response := pb.ExportConfigAndConfigIDResponse{}
	return &response, err
}

func (server *G2EngineServer) ExportCSVEntityReport(ctx context.Context, request *pb.ExportCSVEntityReportRequest) (*pb.ExportCSVEntityReportResponse, error) {
	var err error = nil
	response := pb.ExportCSVEntityReportResponse{}
	return &response, err
}

func (server *G2EngineServer) ExportJSONEntityReport(ctx context.Context, request *pb.ExportJSONEntityReportRequest) (*pb.ExportJSONEntityReportResponse, error) {
	var err error = nil
	response := pb.ExportJSONEntityReportResponse{}
	return &response, err
}

func (server *G2EngineServer) FetchNext(ctx context.Context, request *pb.FetchNextRequest) (*pb.FetchNextResponse, error) {
	var err error = nil
	response := pb.FetchNextResponse{}
	return &response, err
}

func (server *G2EngineServer) FindInterestingEntitiesByEntityID(ctx context.Context, request *pb.FindInterestingEntitiesByEntityIDRequest) (*pb.FindInterestingEntitiesByEntityIDResponse, error) {
	var err error = nil
	response := pb.FindInterestingEntitiesByEntityIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindInterestingEntitiesByRecordID(ctx context.Context, request *pb.FindInterestingEntitiesByRecordIDRequest) (*pb.FindInterestingEntitiesByRecordIDResponse, error) {
	var err error = nil
	response := pb.FindInterestingEntitiesByRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByEntityID(ctx context.Context, request *pb.FindNetworkByEntityIDRequest) (*pb.FindNetworkByEntityIDResponse, error) {
	var err error = nil
	response := pb.FindNetworkByEntityIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByEntityID_V2(ctx context.Context, request *pb.FindNetworkByEntityID_V2Request) (*pb.FindNetworkByEntityID_V2Response, error) {
	var err error = nil
	response := pb.FindNetworkByEntityID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByRecordID(ctx context.Context, request *pb.FindNetworkByRecordIDRequest) (*pb.FindNetworkByRecordIDResponse, error) {
	var err error = nil
	response := pb.FindNetworkByRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByRecordID_V2(ctx context.Context, request *pb.FindNetworkByRecordID_V2Request) (*pb.FindNetworkByRecordID_V2Response, error) {
	var err error = nil
	response := pb.FindNetworkByRecordID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) FindPathByEntityID(ctx context.Context, request *pb.FindPathByEntityIDRequest) (*pb.FindPathByEntityIDResponse, error) {
	var err error = nil
	response := pb.FindPathByEntityIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindPathByEntityID_V2(ctx context.Context, request *pb.FindPathByEntityID_V2Request) (*pb.FindPathByEntityID_V2Response, error) {
	var err error = nil
	response := pb.FindPathByEntityID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) FindPathByRecordID(ctx context.Context, request *pb.FindPathByRecordIDRequest) (*pb.FindPathByRecordIDResponse, error) {
	var err error = nil
	response := pb.FindPathByRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindPathByRecordID_V2(ctx context.Context, request *pb.FindPathByRecordID_V2Request) (*pb.FindPathByRecordID_V2Response, error) {
	var err error = nil
	response := pb.FindPathByRecordID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByEntityID(ctx context.Context, request *pb.FindPathExcludingByEntityIDRequest) (*pb.FindPathExcludingByEntityIDResponse, error) {
	var err error = nil
	response := pb.FindPathExcludingByEntityIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByEntityID_V2(ctx context.Context, request *pb.FindPathExcludingByEntityID_V2Request) (*pb.FindPathExcludingByEntityID_V2Response, error) {
	var err error = nil
	response := pb.FindPathExcludingByEntityID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByRecordID(ctx context.Context, request *pb.FindPathExcludingByRecordIDRequest) (*pb.FindPathExcludingByRecordIDResponse, error) {
	var err error = nil
	response := pb.FindPathExcludingByRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByRecordID_V2(ctx context.Context, request *pb.FindPathExcludingByRecordID_V2Request) (*pb.FindPathExcludingByRecordID_V2Response, error) {
	var err error = nil
	response := pb.FindPathExcludingByRecordID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByEntityID(ctx context.Context, request *pb.FindPathIncludingSourceByEntityIDRequest) (*pb.FindPathIncludingSourceByEntityIDResponse, error) {
	var err error = nil
	response := pb.FindPathIncludingSourceByEntityIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByEntityID_V2(ctx context.Context, request *pb.FindPathIncludingSourceByEntityID_V2Request) (*pb.FindPathIncludingSourceByEntityID_V2Response, error) {
	var err error = nil
	response := pb.FindPathIncludingSourceByEntityID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByRecordID(ctx context.Context, request *pb.FindPathIncludingSourceByRecordIDRequest) (*pb.FindPathIncludingSourceByRecordIDResponse, error) {
	var err error = nil
	response := pb.FindPathIncludingSourceByRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByRecordID_V2(ctx context.Context, request *pb.FindPathIncludingSourceByRecordID_V2Request) (*pb.FindPathIncludingSourceByRecordID_V2Response, error) {
	var err error = nil
	response := pb.FindPathIncludingSourceByRecordID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) GetActiveConfigID(ctx context.Context, request *pb.GetActiveConfigIDRequest) (*pb.GetActiveConfigIDResponse, error) {
	var err error = nil
	response := pb.GetActiveConfigIDResponse{}
	return &response, err
}

func (server *G2EngineServer) GetEntityByEntityID(ctx context.Context, request *pb.GetEntityByEntityIDRequest) (*pb.GetEntityByEntityIDResponse, error) {
	var err error = nil
	response := pb.GetEntityByEntityIDResponse{}
	return &response, err
}

func (server *G2EngineServer) GetEntityByEntityID_V2(ctx context.Context, request *pb.GetEntityByEntityID_V2Request) (*pb.GetEntityByEntityID_V2Response, error) {
	var err error = nil
	response := pb.GetEntityByEntityID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) GetEntityByRecordID(ctx context.Context, request *pb.GetEntityByRecordIDRequest) (*pb.GetEntityByRecordIDResponse, error) {
	var err error = nil
	response := pb.GetEntityByRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) GetEntityByRecordID_V2(ctx context.Context, request *pb.GetEntityByRecordID_V2Request) (*pb.GetEntityByRecordID_V2Response, error) {
	var err error = nil
	response := pb.GetEntityByRecordID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) GetRecord(ctx context.Context, request *pb.GetRecordRequest) (*pb.GetRecordResponse, error) {
	var err error = nil
	response := pb.GetRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) GetRecord_V2(ctx context.Context, request *pb.GetRecord_V2Request) (*pb.GetRecord_V2Response, error) {
	var err error = nil
	response := pb.GetRecord_V2Response{}
	return &response, err
}

func (server *G2EngineServer) GetRedoRecord(ctx context.Context, request *pb.GetRedoRecordRequest) (*pb.GetRedoRecordResponse, error) {
	var err error = nil
	response := pb.GetRedoRecordResponse{}
	return &response, err

}

func (server *G2EngineServer) GetRepositoryLastModifiedTime(ctx context.Context, request *pb.GetRepositoryLastModifiedTimeRequest) (*pb.GetRepositoryLastModifiedTimeResponse, error) {
	var err error = nil
	response := pb.GetRepositoryLastModifiedTimeResponse{}
	return &response, err
}

func (server *G2EngineServer) GetVirtualEntityByRecordID(ctx context.Context, request *pb.GetVirtualEntityByRecordIDRequest) (*pb.GetVirtualEntityByRecordIDResponse, error) {
	var err error = nil
	response := pb.GetVirtualEntityByRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) GetVirtualEntityByRecordID_V2(ctx context.Context, request *pb.GetVirtualEntityByRecordID_V2Request) (*pb.GetVirtualEntityByRecordID_V2Response, error) {
	var err error = nil
	response := pb.GetVirtualEntityByRecordID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) HowEntityByEntityID(ctx context.Context, request *pb.HowEntityByEntityIDRequest) (*pb.HowEntityByEntityIDResponse, error) {
	var err error = nil
	response := pb.HowEntityByEntityIDResponse{}
	return &response, err
}

func (server *G2EngineServer) HowEntityByEntityID_V2(ctx context.Context, request *pb.HowEntityByEntityID_V2Request) (*pb.HowEntityByEntityID_V2Response, error) {
	var err error = nil
	response := pb.HowEntityByEntityID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	var err error = nil
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2EngineServer) InitWithConfigID(ctx context.Context, request *pb.InitWithConfigIDRequest) (*pb.InitWithConfigIDResponse, error) {
	var err error = nil
	response := pb.InitWithConfigIDResponse{}
	return &response, err
}

func (server *G2EngineServer) PrimeEngine(ctx context.Context, request *pb.PrimeEngineRequest) (*pb.PrimeEngineResponse, error) {
	var err error = nil
	response := pb.PrimeEngineResponse{}
	return &response, err
}

func (server *G2EngineServer) Process(ctx context.Context, request *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	var err error = nil
	response := pb.ProcessResponse{}
	return &response, err
}

func (server *G2EngineServer) ProcessRedoRecord(ctx context.Context, request *pb.ProcessRedoRecordRequest) (*pb.ProcessRedoRecordResponse, error) {
	var err error = nil
	response := pb.ProcessRedoRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) ProcessRedoRecordWithInfo(ctx context.Context, request *pb.ProcessRedoRecordWithInfoRequest) (*pb.ProcessRedoRecordWithInfoResponse, error) {
	var err error = nil
	response := pb.ProcessRedoRecordWithInfoResponse{}
	return &response, err
}

func (server *G2EngineServer) ProcessWithInfo(ctx context.Context, request *pb.ProcessWithInfoRequest) (*pb.ProcessWithInfoResponse, error) {
	var err error = nil
	response := pb.ProcessWithInfoResponse{}
	return &response, err
}

func (server *G2EngineServer) ProcessWithResponse(ctx context.Context, request *pb.ProcessWithResponseRequest) (*pb.ProcessWithResponseResponse, error) {
	var err error = nil
	response := pb.ProcessWithResponseResponse{}
	return &response, err
}

func (server *G2EngineServer) ProcessWithResponseResize(ctx context.Context, request *pb.ProcessWithResponseResizeRequest) (*pb.ProcessWithResponseResizeResponse, error) {
	var err error = nil
	response := pb.ProcessWithResponseResizeResponse{}
	return &response, err
}

func (server *G2EngineServer) PurgeRepository(ctx context.Context, request *pb.PurgeRepositoryRequest) (*pb.PurgeRepositoryResponse, error) {
	var err error = nil
	response := pb.PurgeRepositoryResponse{}
	return &response, err
}

func (server *G2EngineServer) ReevaluateEntity(ctx context.Context, request *pb.ReevaluateEntityRequest) (*pb.ReevaluateEntityResponse, error) {
	var err error = nil
	response := pb.ReevaluateEntityResponse{}
	return &response, err
}

func (server *G2EngineServer) ReevaluateEntityWithInfo(ctx context.Context, request *pb.ReevaluateEntityWithInfoRequest) (*pb.ReevaluateEntityWithInfoResponse, error) {
	var err error = nil
	response := pb.ReevaluateEntityWithInfoResponse{}
	return &response, err
}

func (server *G2EngineServer) ReevaluateRecord(ctx context.Context, request *pb.ReevaluateRecordRequest) (*pb.ReevaluateRecordResponse, error) {
	var err error = nil
	response := pb.ReevaluateRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) ReevaluateRecordWithInfo(ctx context.Context, request *pb.ReevaluateRecordWithInfoRequest) (*pb.ReevaluateRecordWithInfoResponse, error) {
	var err error = nil
	response := pb.ReevaluateRecordWithInfoResponse{}
	return &response, err
}

func (server *G2EngineServer) Reinit(ctx context.Context, request *pb.ReinitRequest) (*pb.ReinitResponse, error) {
	var err error = nil
	response := pb.ReinitResponse{}
	return &response, err
}

func (server *G2EngineServer) ReplaceRecord(ctx context.Context, request *pb.ReplaceRecordRequest) (*pb.ReplaceRecordResponse, error) {
	var err error = nil
	response := pb.ReplaceRecordResponse{}
	return &response, err
}

func (server *G2EngineServer) ReplaceRecordWithInfo(ctx context.Context, request *pb.ReplaceRecordWithInfoRequest) (*pb.ReplaceRecordWithInfoResponse, error) {
	var err error = nil
	response := pb.ReplaceRecordWithInfoResponse{}
	return &response, err
}

func (server *G2EngineServer) SearchByAttributes(ctx context.Context, request *pb.SearchByAttributesRequest) (*pb.SearchByAttributesResponse, error) {
	var err error = nil
	response := pb.SearchByAttributesResponse{}
	return &response, err
}

func (server *G2EngineServer) SearchByAttributes_V2(ctx context.Context, request *pb.SearchByAttributes_V2Request) (*pb.SearchByAttributes_V2Response, error) {
	var err error = nil
	response := pb.SearchByAttributes_V2Response{}
	return &response, err
}

func (server *G2EngineServer) Stats(ctx context.Context, request *pb.StatsRequest) (*pb.StatsResponse, error) {
	var err error = nil
	response := pb.StatsResponse{}
	return &response, err
}

func (server *G2EngineServer) WhyEntities(ctx context.Context, request *pb.WhyEntitiesRequest) (*pb.WhyEntitiesResponse, error) {
	var err error = nil
	response := pb.WhyEntitiesResponse{}
	return &response, err
}

func (server *G2EngineServer) WhyEntities_V2(ctx context.Context, request *pb.WhyEntities_V2Request) (*pb.WhyEntities_V2Response, error) {
	var err error = nil
	response := pb.WhyEntities_V2Response{}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByEntityID(ctx context.Context, request *pb.WhyEntityByEntityIDRequest) (*pb.WhyEntityByEntityIDResponse, error) {
	var err error = nil
	response := pb.WhyEntityByEntityIDResponse{}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByEntityID_V2(ctx context.Context, request *pb.WhyEntityByEntityID_V2Request) (*pb.WhyEntityByEntityID_V2Response, error) {
	var err error = nil
	response := pb.WhyEntityByEntityID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByRecordID(ctx context.Context, request *pb.WhyEntityByRecordIDRequest) (*pb.WhyEntityByRecordIDResponse, error) {
	var err error = nil
	response := pb.WhyEntityByRecordIDResponse{}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByRecordID_V2(ctx context.Context, request *pb.WhyEntityByRecordID_V2Request) (*pb.WhyEntityByRecordID_V2Response, error) {
	var err error = nil
	response := pb.WhyEntityByRecordID_V2Response{}
	return &response, err
}

func (server *G2EngineServer) WhyRecords(ctx context.Context, request *pb.WhyRecordsRequest) (*pb.WhyRecordsResponse, error) {
	var err error = nil
	response := pb.WhyRecordsResponse{}
	return &response, err
}

func (server *G2EngineServer) WhyRecords_V2(ctx context.Context, request *pb.WhyRecords_V2Request) (*pb.WhyRecords_V2Response, error) {
	var err error = nil
	response := pb.WhyRecords_V2Response{}
	return &response, err
}

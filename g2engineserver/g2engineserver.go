package g2engineserver

import (
	"context"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go-base/g2engine"
	"github.com/senzing/g2-sdk-go/g2api"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2engine"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
)

var (
	g2engineSingleton g2api.G2engine
	g2engineSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2engine() g2api.G2engine {
	g2engineSyncOnce.Do(func() {
		g2engineSingleton = &g2sdk.G2engine{}
	})
	return g2engineSingleton
}

func GetSdkG2engine() g2api.G2engine {
	return getG2engine()
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
// Interface methods for github.com/senzing/g2-sdk-go/g2engine.G2engine
// ----------------------------------------------------------------------------

func (server *G2EngineServer) AddRecord(ctx context.Context, request *g2pb.AddRecordRequest) (*g2pb.AddRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(1, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.AddRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID())
	response := g2pb.AddRecordResponse{}
	if server.isTrace {
		defer server.traceExit(2, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithInfo(ctx context.Context, request *g2pb.AddRecordWithInfoRequest) (*g2pb.AddRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(3, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.AddRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := g2pb.AddRecordWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(4, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithInfoWithReturnedRecordID(ctx context.Context, request *g2pb.AddRecordWithInfoWithReturnedRecordIDRequest) (*g2pb.AddRecordWithInfoWithReturnedRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(5, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, recordId, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, request.GetDataSourceCode(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := g2pb.AddRecordWithInfoWithReturnedRecordIDResponse{
		RecordID: recordId,
		WithInfo: result,
	}
	if server.isTrace {
		defer server.traceExit(6, request, result, recordId, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithReturnedRecordID(ctx context.Context, request *g2pb.AddRecordWithReturnedRecordIDRequest) (*g2pb.AddRecordWithReturnedRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(7, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.AddRecordWithReturnedRecordID(ctx, request.GetDataSourceCode(), request.GetJsonData(), request.GetLoadID())
	response := g2pb.AddRecordWithReturnedRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(8, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) CheckRecord(ctx context.Context, request *g2pb.CheckRecordRequest) (*g2pb.CheckRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(9, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.CheckRecord(ctx, request.GetRecord(), request.GetRecordQueryList())
	response := g2pb.CheckRecordResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(10, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) CloseExport(ctx context.Context, request *g2pb.CloseExportRequest) (*g2pb.CloseExportResponse, error) {
	if server.isTrace {
		server.traceEntry(13, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.CloseExport(ctx, uintptr(request.GetResponseHandle()))
	response := g2pb.CloseExportResponse{}
	if server.isTrace {
		defer server.traceExit(14, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) CountRedoRecords(ctx context.Context, request *g2pb.CountRedoRecordsRequest) (*g2pb.CountRedoRecordsResponse, error) {
	if server.isTrace {
		server.traceEntry(15, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.CountRedoRecords(ctx)
	response := g2pb.CountRedoRecordsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(16, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) DeleteRecord(ctx context.Context, request *g2pb.DeleteRecordRequest) (*g2pb.DeleteRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(17, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.DeleteRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetLoadID())
	response := g2pb.DeleteRecordResponse{}
	if server.isTrace {
		defer server.traceExit(18, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) DeleteRecordWithInfo(ctx context.Context, request *g2pb.DeleteRecordWithInfoRequest) (*g2pb.DeleteRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(19, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.DeleteRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetLoadID(), request.GetFlags())
	response := g2pb.DeleteRecordWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(20, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
	if server.isTrace {
		server.traceEntry(21, request)
	}
	entryTime := time.Now()
	// g2engine := getG2engine()
	// err := g2engine.Destroy(ctx)
	err := server.getLogger().Error(4001)
	response := g2pb.DestroyResponse{}
	if server.isTrace {
		defer server.traceExit(22, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ExportConfig(ctx context.Context, request *g2pb.ExportConfigRequest) (*g2pb.ExportConfigResponse, error) {
	if server.isTrace {
		server.traceEntry(23, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ExportConfig(ctx)
	response := g2pb.ExportConfigResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(24, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ExportConfigAndConfigID(ctx context.Context, request *g2pb.ExportConfigAndConfigIDRequest) (*g2pb.ExportConfigAndConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(25, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, configId, err := g2engine.ExportConfigAndConfigID(ctx)
	response := g2pb.ExportConfigAndConfigIDResponse{
		Config:   result,
		ConfigID: configId,
	}
	if server.isTrace {
		defer server.traceExit(26, request, result, configId, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ExportCSVEntityReport(ctx context.Context, request *g2pb.ExportCSVEntityReportRequest) (*g2pb.ExportCSVEntityReportResponse, error) {
	if server.isTrace {
		server.traceEntry(27, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ExportCSVEntityReport(ctx, request.GetCsvColumnList(), request.GetFlags())
	response := g2pb.ExportCSVEntityReportResponse{
		Result: int64(result),
	}
	if server.isTrace {
		defer server.traceExit(28, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ExportJSONEntityReport(ctx context.Context, request *g2pb.ExportJSONEntityReportRequest) (*g2pb.ExportJSONEntityReportResponse, error) {
	if server.isTrace {
		server.traceEntry(29, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ExportJSONEntityReport(ctx, request.GetFlags())
	response := g2pb.ExportJSONEntityReportResponse{
		Result: int64(result),
	}
	if server.isTrace {
		defer server.traceExit(30, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FetchNext(ctx context.Context, request *g2pb.FetchNextRequest) (*g2pb.FetchNextResponse, error) {
	if server.isTrace {
		server.traceEntry(31, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FetchNext(ctx, uintptr(request.GetResponseHandle()))
	response := g2pb.FetchNextResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(32, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindInterestingEntitiesByEntityID(ctx context.Context, request *g2pb.FindInterestingEntitiesByEntityIDRequest) (*g2pb.FindInterestingEntitiesByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(33, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindInterestingEntitiesByEntityID(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.FindInterestingEntitiesByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(34, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindInterestingEntitiesByRecordID(ctx context.Context, request *g2pb.FindInterestingEntitiesByRecordIDRequest) (*g2pb.FindInterestingEntitiesByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(35, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindInterestingEntitiesByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.FindInterestingEntitiesByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(36, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByEntityID(ctx context.Context, request *g2pb.FindNetworkByEntityIDRequest) (*g2pb.FindNetworkByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(37, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByEntityID(ctx, request.GetEntityList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()))
	response := g2pb.FindNetworkByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(38, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByEntityID_V2(ctx context.Context, request *g2pb.FindNetworkByEntityID_V2Request) (*g2pb.FindNetworkByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(39, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByEntityID_V2(ctx, request.GetEntityList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()), request.GetFlags())
	response := g2pb.FindNetworkByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(40, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByRecordID(ctx context.Context, request *g2pb.FindNetworkByRecordIDRequest) (*g2pb.FindNetworkByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(41, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByRecordID(ctx, request.GetRecordList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()))
	response := g2pb.FindNetworkByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(42, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByRecordID_V2(ctx context.Context, request *g2pb.FindNetworkByRecordID_V2Request) (*g2pb.FindNetworkByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(43, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByRecordID_V2(ctx, request.GetRecordList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()), request.GetFlags())
	response := g2pb.FindNetworkByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(44, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByEntityID(ctx context.Context, request *g2pb.FindPathByEntityIDRequest) (*g2pb.FindPathByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(45, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()))
	response := g2pb.FindPathByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(46, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByEntityID_V2(ctx context.Context, request *g2pb.FindPathByEntityID_V2Request) (*g2pb.FindPathByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(47, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetFlags())
	response := g2pb.FindPathByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(48, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByRecordID(ctx context.Context, request *g2pb.FindPathByRecordIDRequest) (*g2pb.FindPathByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(49, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()))
	response := g2pb.FindPathByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(50, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByRecordID_V2(ctx context.Context, request *g2pb.FindPathByRecordID_V2Request) (*g2pb.FindPathByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(51, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetFlags())
	response := g2pb.FindPathByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(52, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByEntityID(ctx context.Context, request *g2pb.FindPathExcludingByEntityIDRequest) (*g2pb.FindPathExcludingByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(53, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities())
	response := g2pb.FindPathExcludingByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(54, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByEntityID_V2(ctx context.Context, request *g2pb.FindPathExcludingByEntityID_V2Request) (*g2pb.FindPathExcludingByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(55, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities(), request.GetFlags())
	response := g2pb.FindPathExcludingByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(56, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByRecordID(ctx context.Context, request *g2pb.FindPathExcludingByRecordIDRequest) (*g2pb.FindPathExcludingByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(57, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords())
	response := g2pb.FindPathExcludingByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(58, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByRecordID_V2(ctx context.Context, request *g2pb.FindPathExcludingByRecordID_V2Request) (*g2pb.FindPathExcludingByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(59, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords(), request.GetFlags())
	response := g2pb.FindPathExcludingByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(60, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByEntityID(ctx context.Context, request *g2pb.FindPathIncludingSourceByEntityIDRequest) (*g2pb.FindPathIncludingSourceByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(61, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities(), request.GetRequiredDsrcs())
	response := g2pb.FindPathIncludingSourceByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(62, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByEntityID_V2(ctx context.Context, request *g2pb.FindPathIncludingSourceByEntityID_V2Request) (*g2pb.FindPathIncludingSourceByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(63, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities(), request.GetRequiredDsrcs(), request.GetFlags())
	response := g2pb.FindPathIncludingSourceByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(64, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByRecordID(ctx context.Context, request *g2pb.FindPathIncludingSourceByRecordIDRequest) (*g2pb.FindPathIncludingSourceByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(65, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords(), request.GetRequiredDsrcs())
	response := g2pb.FindPathIncludingSourceByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(66, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByRecordID_V2(ctx context.Context, request *g2pb.FindPathIncludingSourceByRecordID_V2Request) (*g2pb.FindPathIncludingSourceByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(67, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords(), request.GetRequiredDsrcs(), request.GetFlags())
	response := g2pb.FindPathIncludingSourceByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(68, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetActiveConfigID(ctx context.Context, request *g2pb.GetActiveConfigIDRequest) (*g2pb.GetActiveConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(69, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetActiveConfigID(ctx)
	response := g2pb.GetActiveConfigIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(70, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByEntityID(ctx context.Context, request *g2pb.GetEntityByEntityIDRequest) (*g2pb.GetEntityByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(71, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByEntityID(ctx, request.GetEntityID())
	response := g2pb.GetEntityByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(72, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByEntityID_V2(ctx context.Context, request *g2pb.GetEntityByEntityID_V2Request) (*g2pb.GetEntityByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(73, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.GetEntityByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(74, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByRecordID(ctx context.Context, request *g2pb.GetEntityByRecordIDRequest) (*g2pb.GetEntityByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(75, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := g2pb.GetEntityByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(76, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByRecordID_V2(ctx context.Context, request *g2pb.GetEntityByRecordID_V2Request) (*g2pb.GetEntityByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(77, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByRecordID_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.GetEntityByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(78, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetRecord(ctx context.Context, request *g2pb.GetRecordRequest) (*g2pb.GetRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(83, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetRecord(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := g2pb.GetRecordResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(84, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetRecord_V2(ctx context.Context, request *g2pb.GetRecord_V2Request) (*g2pb.GetRecord_V2Response, error) {
	if server.isTrace {
		server.traceEntry(85, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetRecord_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.GetRecord_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(86, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetRedoRecord(ctx context.Context, request *g2pb.GetRedoRecordRequest) (*g2pb.GetRedoRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(87, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetRedoRecord(ctx)
	response := g2pb.GetRedoRecordResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(88, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetRepositoryLastModifiedTime(ctx context.Context, request *g2pb.GetRepositoryLastModifiedTimeRequest) (*g2pb.GetRepositoryLastModifiedTimeResponse, error) {
	if server.isTrace {
		server.traceEntry(89, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetRepositoryLastModifiedTime(ctx)
	response := g2pb.GetRepositoryLastModifiedTimeResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(90, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetVirtualEntityByRecordID(ctx context.Context, request *g2pb.GetVirtualEntityByRecordIDRequest) (*g2pb.GetVirtualEntityByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(91, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetVirtualEntityByRecordID(ctx, request.GetRecordList())
	response := g2pb.GetVirtualEntityByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(92, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetVirtualEntityByRecordID_V2(ctx context.Context, request *g2pb.GetVirtualEntityByRecordID_V2Request) (*g2pb.GetVirtualEntityByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(93, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, request.GetRecordList(), request.GetFlags())
	response := g2pb.GetVirtualEntityByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(94, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) HowEntityByEntityID(ctx context.Context, request *g2pb.HowEntityByEntityIDRequest) (*g2pb.HowEntityByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(95, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.HowEntityByEntityID(ctx, request.GetEntityID())
	response := g2pb.HowEntityByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(96, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) HowEntityByEntityID_V2(ctx context.Context, request *g2pb.HowEntityByEntityID_V2Request) (*g2pb.HowEntityByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(97, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.HowEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.HowEntityByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(98, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	if server.isTrace {
		server.traceEntry(99, request)
	}
	entryTime := time.Now()
	// g2engine := getG2engine()
	// err := g2engine.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err := server.getLogger().Error(4002)
	response := g2pb.InitResponse{}
	if server.isTrace {
		defer server.traceExit(100, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) InitWithConfigID(ctx context.Context, request *g2pb.InitWithConfigIDRequest) (*g2pb.InitWithConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(101, request)
	}
	entryTime := time.Now()
	// g2engine := getG2engine()
	// err := g2engine.InitWithConfigID(ctx, request.GetModuleName(), request.GetIniParams(), request.GetInitConfigID(), int(request.GetVerboseLogging()))
	err := server.getLogger().Error(4003)
	response := g2pb.InitWithConfigIDResponse{}
	if server.isTrace {
		defer server.traceExit(102, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) PrimeEngine(ctx context.Context, request *g2pb.PrimeEngineRequest) (*g2pb.PrimeEngineResponse, error) {
	if server.isTrace {
		server.traceEntry(103, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.PrimeEngine(ctx)
	response := g2pb.PrimeEngineResponse{}
	if server.isTrace {
		defer server.traceExit(104, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) Process(ctx context.Context, request *g2pb.ProcessRequest) (*g2pb.ProcessResponse, error) {
	if server.isTrace {
		server.traceEntry(105, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.Process(ctx, request.GetRecord())
	response := g2pb.ProcessResponse{}
	if server.isTrace {
		defer server.traceExit(106, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessRedoRecord(ctx context.Context, request *g2pb.ProcessRedoRecordRequest) (*g2pb.ProcessRedoRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(107, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ProcessRedoRecord(ctx)
	response := g2pb.ProcessRedoRecordResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(108, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessRedoRecordWithInfo(ctx context.Context, request *g2pb.ProcessRedoRecordWithInfoRequest) (*g2pb.ProcessRedoRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(109, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, withInfo, err := g2engine.ProcessRedoRecordWithInfo(ctx, request.GetFlags())
	response := g2pb.ProcessRedoRecordWithInfoResponse{
		Result:   result,
		WithInfo: withInfo,
	}
	if server.isTrace {
		defer server.traceExit(110, request, result, withInfo, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessWithInfo(ctx context.Context, request *g2pb.ProcessWithInfoRequest) (*g2pb.ProcessWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(111, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ProcessWithInfo(ctx, request.GetRecord(), request.GetFlags())
	response := g2pb.ProcessWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(112, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessWithResponse(ctx context.Context, request *g2pb.ProcessWithResponseRequest) (*g2pb.ProcessWithResponseResponse, error) {
	if server.isTrace {
		server.traceEntry(113, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ProcessWithResponse(ctx, request.GetRecord())
	response := g2pb.ProcessWithResponseResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(114, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessWithResponseResize(ctx context.Context, request *g2pb.ProcessWithResponseResizeRequest) (*g2pb.ProcessWithResponseResizeResponse, error) {
	if server.isTrace {
		server.traceEntry(115, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ProcessWithResponseResize(ctx, request.GetRecord())
	response := g2pb.ProcessWithResponseResizeResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(116, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) PurgeRepository(ctx context.Context, request *g2pb.PurgeRepositoryRequest) (*g2pb.PurgeRepositoryResponse, error) {
	if server.isTrace {
		server.traceEntry(117, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.PurgeRepository(ctx)
	response := g2pb.PurgeRepositoryResponse{}
	if server.isTrace {
		defer server.traceExit(118, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReevaluateEntity(ctx context.Context, request *g2pb.ReevaluateEntityRequest) (*g2pb.ReevaluateEntityResponse, error) {
	if server.isTrace {
		server.traceEntry(119, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.ReevaluateEntity(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.ReevaluateEntityResponse{}
	if server.isTrace {
		defer server.traceExit(120, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReevaluateEntityWithInfo(ctx context.Context, request *g2pb.ReevaluateEntityWithInfoRequest) (*g2pb.ReevaluateEntityWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(121, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ReevaluateEntityWithInfo(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.ReevaluateEntityWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(122, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReevaluateRecord(ctx context.Context, request *g2pb.ReevaluateRecordRequest) (*g2pb.ReevaluateRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(123, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.ReevaluateRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.ReevaluateRecordResponse{}
	if server.isTrace {
		defer server.traceExit(124, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReevaluateRecordWithInfo(ctx context.Context, request *g2pb.ReevaluateRecordWithInfoRequest) (*g2pb.ReevaluateRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(125, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ReevaluateRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.ReevaluateRecordWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(126, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	g2engine := getG2engine()
	return g2engine.RegisterObserver(ctx, observer)
}

func (server *G2EngineServer) Reinit(ctx context.Context, request *g2pb.ReinitRequest) (*g2pb.ReinitResponse, error) {
	if server.isTrace {
		server.traceEntry(127, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.Reinit(ctx, request.GetInitConfigID())
	response := g2pb.ReinitResponse{}
	if server.isTrace {
		defer server.traceExit(128, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReplaceRecord(ctx context.Context, request *g2pb.ReplaceRecordRequest) (*g2pb.ReplaceRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(129, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.ReplaceRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID())
	response := g2pb.ReplaceRecordResponse{}
	if server.isTrace {
		defer server.traceExit(130, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReplaceRecordWithInfo(ctx context.Context, request *g2pb.ReplaceRecordWithInfoRequest) (*g2pb.ReplaceRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(131, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ReplaceRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := g2pb.ReplaceRecordWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(132, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) SearchByAttributes(ctx context.Context, request *g2pb.SearchByAttributesRequest) (*g2pb.SearchByAttributesResponse, error) {
	if server.isTrace {
		server.traceEntry(133, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.SearchByAttributes(ctx, request.GetJsonData())
	response := g2pb.SearchByAttributesResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(134, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) SearchByAttributes_V2(ctx context.Context, request *g2pb.SearchByAttributes_V2Request) (*g2pb.SearchByAttributes_V2Response, error) {
	if server.isTrace {
		server.traceEntry(135, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.SearchByAttributes_V2(ctx, request.GetJsonData(), request.GetFlags())
	response := g2pb.SearchByAttributes_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(136, request, result, err, time.Since(entryTime))
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
		server.traceEntry(137, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	g2engine := getG2engine()
	g2engine.SetLogLevel(ctx, logLevel)
	server.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	server.isTrace = (server.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if server.isTrace {
		defer server.traceExit(138, logLevel, err, time.Since(entryTime))
	}
	return err
}

func (server *G2EngineServer) Stats(ctx context.Context, request *g2pb.StatsRequest) (*g2pb.StatsResponse, error) {
	if server.isTrace {
		server.traceEntry(139, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.Stats(ctx)
	response := g2pb.StatsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(140, request, result, err, time.Since(entryTime))
	}
	return &response, err
}


func (server *G2EngineServer) StreamExportCSVEntityReport(request *g2pb.StreamExportCSVEntityReportRequest, stream g2pb.G2Engine_StreamExportCSVEntityReportServer) error {
	if server.isTrace {
		server.traceEntry(157, request)
	}
    ctx := stream.Context()
	entryTime := time.Now()
	g2engine := getG2engine()

	rowsFetched := 0

	//get the query handle
	queryHandle, err := g2engine.ExportCSVEntityReport(ctx, request.GetCsvColumnList(), request.GetFlags())
	if err != nil {
	    return err
	}

    for {
	    fetchResult, err := g2engine.FetchNext(ctx, queryHandle)
		if err != nil {
			return err
		}
		if len(fetchResult) == 0 {
			break
		}
		response := g2pb.StreamExportCSVEntityReportResponse{
			Result: fetchResult,
		}
		if err = stream.Send(&response); err != nil {
			return err
		}
		server.traceEntry(158, request, fetchResult)
		rowsFetched += 1
     }

	err = g2engine.CloseExport(ctx, queryHandle)
	if err != nil {
		return err
	}

	if server.isTrace {
		defer server.traceExit(159, request, rowsFetched, err, time.Since(entryTime))
	}
	return nil
}

func (server *G2EngineServer) StreamExportJSONEntityReport(request *g2pb.StreamExportJSONEntityReportRequest, stream g2pb.G2Engine_StreamExportJSONEntityReportServer) error {
	if server.isTrace {
		server.traceEntry(160, request)
	}
    ctx := stream.Context()
	entryTime := time.Now()
	g2engine := getG2engine()

	rowsFetched := 0

	//get the query handle
	queryHandle, err := g2engine.ExportJSONEntityReport(ctx, request.GetFlags())
	if err != nil {
	    return err
	}

    for {
	    fetchResult, err := g2engine.FetchNext(ctx, queryHandle)
		if err != nil {
			return err
		}
		if len(fetchResult) == 0 {
			break
		}
		response := g2pb.StreamExportJSONEntityReportResponse{
			Result: fetchResult,
		}
		if err = stream.Send(&response); err != nil {
			return err
		}
		server.traceEntry(161, request, fetchResult)
		rowsFetched += 1
     }

	err = g2engine.CloseExport(ctx, queryHandle)
	if err != nil {
		return err
	}

	if server.isTrace {
		defer server.traceExit(162, request, rowsFetched, err, time.Since(entryTime))
	}
	return nil
}

func (server *G2EngineServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	g2engine := getG2engine()
	return g2engine.UnregisterObserver(ctx, observer)
}

func (server *G2EngineServer) WhyEntities(ctx context.Context, request *g2pb.WhyEntitiesRequest) (*g2pb.WhyEntitiesResponse, error) {
	if server.isTrace {
		server.traceEntry(141, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntities(ctx, request.GetEntityID1(), request.GetEntityID2())
	response := g2pb.WhyEntitiesResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(142, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntities_V2(ctx context.Context, request *g2pb.WhyEntities_V2Request) (*g2pb.WhyEntities_V2Response, error) {
	if server.isTrace {
		server.traceEntry(143, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntities_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetFlags())
	response := g2pb.WhyEntities_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(144, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByEntityID(ctx context.Context, request *g2pb.WhyEntityByEntityIDRequest) (*g2pb.WhyEntityByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(145, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByEntityID(ctx, request.GetEntityID())
	response := g2pb.WhyEntityByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(146, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByEntityID_V2(ctx context.Context, request *g2pb.WhyEntityByEntityID_V2Request) (*g2pb.WhyEntityByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(147, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.WhyEntityByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(148, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByRecordID(ctx context.Context, request *g2pb.WhyEntityByRecordIDRequest) (*g2pb.WhyEntityByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(149, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := g2pb.WhyEntityByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(150, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByRecordID_V2(ctx context.Context, request *g2pb.WhyEntityByRecordID_V2Request) (*g2pb.WhyEntityByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(151, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByRecordID_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.WhyEntityByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(152, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyRecords(ctx context.Context, request *g2pb.WhyRecordsRequest) (*g2pb.WhyRecordsResponse, error) {
	if server.isTrace {
		server.traceEntry(153, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyRecords(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2())
	response := g2pb.WhyRecordsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(154, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyRecords_V2(ctx context.Context, request *g2pb.WhyRecords_V2Request) (*g2pb.WhyRecords_V2Response, error) {
	if server.isTrace {
		server.traceEntry(155, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyRecords_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetFlags())
	response := g2pb.WhyRecords_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(156, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

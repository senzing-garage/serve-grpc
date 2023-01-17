package g2engineserver

import (
	"context"
	"sync"
	"time"

	g2sdk "github.com/senzing/g2-sdk-go/g2engine"
	pb "github.com/senzing/g2-sdk-proto/go/g2engine"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
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

func GetSdkG2engine() *g2sdk.G2engineImpl {
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

func (server *G2EngineServer) AddRecord(ctx context.Context, request *pb.AddRecordRequest) (*pb.AddRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(1, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.AddRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID())
	response := pb.AddRecordResponse{}
	if server.isTrace {
		defer server.traceExit(2, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithInfo(ctx context.Context, request *pb.AddRecordWithInfoRequest) (*pb.AddRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(3, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.AddRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := pb.AddRecordWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(4, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithInfoWithReturnedRecordID(ctx context.Context, request *pb.AddRecordWithInfoWithReturnedRecordIDRequest) (*pb.AddRecordWithInfoWithReturnedRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(5, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, recordId, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, request.GetDataSourceCode(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := pb.AddRecordWithInfoWithReturnedRecordIDResponse{
		RecordID: recordId,
		WithInfo: result,
	}
	if server.isTrace {
		defer server.traceExit(6, request, result, recordId, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) AddRecordWithReturnedRecordID(ctx context.Context, request *pb.AddRecordWithReturnedRecordIDRequest) (*pb.AddRecordWithReturnedRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(7, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.AddRecordWithReturnedRecordID(ctx, request.GetDataSourceCode(), request.GetJsonData(), request.GetLoadID())
	response := pb.AddRecordWithReturnedRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(8, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) CheckRecord(ctx context.Context, request *pb.CheckRecordRequest) (*pb.CheckRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(9, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.CheckRecord(ctx, request.GetRecord(), request.GetRecordQueryList())
	response := pb.CheckRecordResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(10, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) CloseExport(ctx context.Context, request *pb.CloseExportRequest) (*pb.CloseExportResponse, error) {
	if server.isTrace {
		server.traceEntry(13, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.CloseExport(ctx, uintptr(request.GetResponseHandle()))
	response := pb.CloseExportResponse{}
	if server.isTrace {
		defer server.traceExit(14, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) CountRedoRecords(ctx context.Context, request *pb.CountRedoRecordsRequest) (*pb.CountRedoRecordsResponse, error) {
	if server.isTrace {
		server.traceEntry(15, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.CountRedoRecords(ctx)
	response := pb.CountRedoRecordsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(16, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) DeleteRecord(ctx context.Context, request *pb.DeleteRecordRequest) (*pb.DeleteRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(17, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.DeleteRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetLoadID())
	response := pb.DeleteRecordResponse{}
	if server.isTrace {
		defer server.traceExit(18, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) DeleteRecordWithInfo(ctx context.Context, request *pb.DeleteRecordWithInfoRequest) (*pb.DeleteRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(19, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.DeleteRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetLoadID(), request.GetFlags())
	response := pb.DeleteRecordWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(20, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	if server.isTrace {
		server.traceEntry(21, request)
	}
	entryTime := time.Now()
	// g2engine := getG2engine()
	// err := g2engine.Destroy(ctx)
	err := server.getLogger().Error(4001)
	response := pb.DestroyResponse{}
	if server.isTrace {
		defer server.traceExit(22, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ExportConfig(ctx context.Context, request *pb.ExportConfigRequest) (*pb.ExportConfigResponse, error) {
	if server.isTrace {
		server.traceEntry(23, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ExportConfig(ctx)
	response := pb.ExportConfigResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(24, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ExportConfigAndConfigID(ctx context.Context, request *pb.ExportConfigAndConfigIDRequest) (*pb.ExportConfigAndConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(25, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, configId, err := g2engine.ExportConfigAndConfigID(ctx)
	response := pb.ExportConfigAndConfigIDResponse{
		Config:   result,
		ConfigID: configId,
	}
	if server.isTrace {
		defer server.traceExit(26, request, result, configId, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ExportCSVEntityReport(ctx context.Context, request *pb.ExportCSVEntityReportRequest) (*pb.ExportCSVEntityReportResponse, error) {
	if server.isTrace {
		server.traceEntry(27, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ExportCSVEntityReport(ctx, request.GetCsvColumnList(), request.GetFlags())
	response := pb.ExportCSVEntityReportResponse{
		Result: int64(result),
	}
	if server.isTrace {
		defer server.traceExit(28, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ExportJSONEntityReport(ctx context.Context, request *pb.ExportJSONEntityReportRequest) (*pb.ExportJSONEntityReportResponse, error) {
	if server.isTrace {
		server.traceEntry(29, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ExportJSONEntityReport(ctx, request.GetFlags())
	response := pb.ExportJSONEntityReportResponse{
		Result: int64(result),
	}
	if server.isTrace {
		defer server.traceExit(30, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FetchNext(ctx context.Context, request *pb.FetchNextRequest) (*pb.FetchNextResponse, error) {
	if server.isTrace {
		server.traceEntry(31, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FetchNext(ctx, uintptr(request.GetResponseHandle()))
	response := pb.FetchNextResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(32, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindInterestingEntitiesByEntityID(ctx context.Context, request *pb.FindInterestingEntitiesByEntityIDRequest) (*pb.FindInterestingEntitiesByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(33, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindInterestingEntitiesByEntityID(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.FindInterestingEntitiesByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(34, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindInterestingEntitiesByRecordID(ctx context.Context, request *pb.FindInterestingEntitiesByRecordIDRequest) (*pb.FindInterestingEntitiesByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(35, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindInterestingEntitiesByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.FindInterestingEntitiesByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(36, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByEntityID(ctx context.Context, request *pb.FindNetworkByEntityIDRequest) (*pb.FindNetworkByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(37, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByEntityID(ctx, request.GetEntityList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()))
	response := pb.FindNetworkByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(38, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByEntityID_V2(ctx context.Context, request *pb.FindNetworkByEntityID_V2Request) (*pb.FindNetworkByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(39, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByEntityID_V2(ctx, request.GetEntityList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()), request.GetFlags())
	response := pb.FindNetworkByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(40, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByRecordID(ctx context.Context, request *pb.FindNetworkByRecordIDRequest) (*pb.FindNetworkByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(41, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByRecordID(ctx, request.GetRecordList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()))
	response := pb.FindNetworkByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(42, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindNetworkByRecordID_V2(ctx context.Context, request *pb.FindNetworkByRecordID_V2Request) (*pb.FindNetworkByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(43, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindNetworkByRecordID_V2(ctx, request.GetRecordList(), int(request.GetMaxDegree()), int(request.GetBuildOutDegree()), int(request.GetMaxEntities()), request.GetFlags())
	response := pb.FindNetworkByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(44, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByEntityID(ctx context.Context, request *pb.FindPathByEntityIDRequest) (*pb.FindPathByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(45, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()))
	response := pb.FindPathByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(46, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByEntityID_V2(ctx context.Context, request *pb.FindPathByEntityID_V2Request) (*pb.FindPathByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(47, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetFlags())
	response := pb.FindPathByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(48, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByRecordID(ctx context.Context, request *pb.FindPathByRecordIDRequest) (*pb.FindPathByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(49, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()))
	response := pb.FindPathByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(50, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathByRecordID_V2(ctx context.Context, request *pb.FindPathByRecordID_V2Request) (*pb.FindPathByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(51, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetFlags())
	response := pb.FindPathByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(52, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByEntityID(ctx context.Context, request *pb.FindPathExcludingByEntityIDRequest) (*pb.FindPathExcludingByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(53, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities())
	response := pb.FindPathExcludingByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(54, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByEntityID_V2(ctx context.Context, request *pb.FindPathExcludingByEntityID_V2Request) (*pb.FindPathExcludingByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(55, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities(), request.GetFlags())
	response := pb.FindPathExcludingByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(56, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByRecordID(ctx context.Context, request *pb.FindPathExcludingByRecordIDRequest) (*pb.FindPathExcludingByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(57, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords())
	response := pb.FindPathExcludingByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(58, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathExcludingByRecordID_V2(ctx context.Context, request *pb.FindPathExcludingByRecordID_V2Request) (*pb.FindPathExcludingByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(59, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords(), request.GetFlags())
	response := pb.FindPathExcludingByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(60, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByEntityID(ctx context.Context, request *pb.FindPathIncludingSourceByEntityIDRequest) (*pb.FindPathIncludingSourceByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(61, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities(), request.GetRequiredDsrcs())
	response := pb.FindPathIncludingSourceByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(62, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByEntityID_V2(ctx context.Context, request *pb.FindPathIncludingSourceByEntityID_V2Request) (*pb.FindPathIncludingSourceByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(63, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), int(request.GetMaxDegree()), request.GetExcludedEntities(), request.GetRequiredDsrcs(), request.GetFlags())
	response := pb.FindPathIncludingSourceByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(64, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByRecordID(ctx context.Context, request *pb.FindPathIncludingSourceByRecordIDRequest) (*pb.FindPathIncludingSourceByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(65, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords(), request.GetRequiredDsrcs())
	response := pb.FindPathIncludingSourceByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(66, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) FindPathIncludingSourceByRecordID_V2(ctx context.Context, request *pb.FindPathIncludingSourceByRecordID_V2Request) (*pb.FindPathIncludingSourceByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(67, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), int(request.GetMaxDegree()), request.GetExcludedRecords(), request.GetRequiredDsrcs(), request.GetFlags())
	response := pb.FindPathIncludingSourceByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(68, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetActiveConfigID(ctx context.Context, request *pb.GetActiveConfigIDRequest) (*pb.GetActiveConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(69, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetActiveConfigID(ctx)
	response := pb.GetActiveConfigIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(70, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByEntityID(ctx context.Context, request *pb.GetEntityByEntityIDRequest) (*pb.GetEntityByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(71, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByEntityID(ctx, request.GetEntityID())
	response := pb.GetEntityByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(72, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByEntityID_V2(ctx context.Context, request *pb.GetEntityByEntityID_V2Request) (*pb.GetEntityByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(73, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.GetEntityByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(74, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByRecordID(ctx context.Context, request *pb.GetEntityByRecordIDRequest) (*pb.GetEntityByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(75, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := pb.GetEntityByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(76, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetEntityByRecordID_V2(ctx context.Context, request *pb.GetEntityByRecordID_V2Request) (*pb.GetEntityByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(77, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetEntityByRecordID_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.GetEntityByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(78, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetRecord(ctx context.Context, request *pb.GetRecordRequest) (*pb.GetRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(83, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetRecord(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := pb.GetRecordResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(84, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetRecord_V2(ctx context.Context, request *pb.GetRecord_V2Request) (*pb.GetRecord_V2Response, error) {
	if server.isTrace {
		server.traceEntry(85, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetRecord_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.GetRecord_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(86, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetRedoRecord(ctx context.Context, request *pb.GetRedoRecordRequest) (*pb.GetRedoRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(87, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetRedoRecord(ctx)
	response := pb.GetRedoRecordResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(88, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetRepositoryLastModifiedTime(ctx context.Context, request *pb.GetRepositoryLastModifiedTimeRequest) (*pb.GetRepositoryLastModifiedTimeResponse, error) {
	if server.isTrace {
		server.traceEntry(89, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetRepositoryLastModifiedTime(ctx)
	response := pb.GetRepositoryLastModifiedTimeResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(90, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetVirtualEntityByRecordID(ctx context.Context, request *pb.GetVirtualEntityByRecordIDRequest) (*pb.GetVirtualEntityByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(91, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetVirtualEntityByRecordID(ctx, request.GetRecordList())
	response := pb.GetVirtualEntityByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(92, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) GetVirtualEntityByRecordID_V2(ctx context.Context, request *pb.GetVirtualEntityByRecordID_V2Request) (*pb.GetVirtualEntityByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(93, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, request.GetRecordList(), request.GetFlags())
	response := pb.GetVirtualEntityByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(94, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) HowEntityByEntityID(ctx context.Context, request *pb.HowEntityByEntityIDRequest) (*pb.HowEntityByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(95, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.HowEntityByEntityID(ctx, request.GetEntityID())
	response := pb.HowEntityByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(96, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) HowEntityByEntityID_V2(ctx context.Context, request *pb.HowEntityByEntityID_V2Request) (*pb.HowEntityByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(97, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.HowEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.HowEntityByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(98, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	if server.isTrace {
		server.traceEntry(99, request)
	}
	entryTime := time.Now()
	// g2engine := getG2engine()
	// err := g2engine.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err := server.getLogger().Error(4002)
	response := pb.InitResponse{}
	if server.isTrace {
		defer server.traceExit(100, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) InitWithConfigID(ctx context.Context, request *pb.InitWithConfigIDRequest) (*pb.InitWithConfigIDResponse, error) {
	if server.isTrace {
		server.traceEntry(101, request)
	}
	entryTime := time.Now()
	// g2engine := getG2engine()
	// err := g2engine.InitWithConfigID(ctx, request.GetModuleName(), request.GetIniParams(), request.GetInitConfigID(), int(request.GetVerboseLogging()))
	err := server.getLogger().Error(4003)
	response := pb.InitWithConfigIDResponse{}
	if server.isTrace {
		defer server.traceExit(102, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) PrimeEngine(ctx context.Context, request *pb.PrimeEngineRequest) (*pb.PrimeEngineResponse, error) {
	if server.isTrace {
		server.traceEntry(103, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.PrimeEngine(ctx)
	response := pb.PrimeEngineResponse{}
	if server.isTrace {
		defer server.traceExit(104, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) Process(ctx context.Context, request *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	if server.isTrace {
		server.traceEntry(105, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.Process(ctx, request.GetRecord())
	response := pb.ProcessResponse{}
	if server.isTrace {
		defer server.traceExit(106, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessRedoRecord(ctx context.Context, request *pb.ProcessRedoRecordRequest) (*pb.ProcessRedoRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(107, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ProcessRedoRecord(ctx)
	response := pb.ProcessRedoRecordResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(108, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessRedoRecordWithInfo(ctx context.Context, request *pb.ProcessRedoRecordWithInfoRequest) (*pb.ProcessRedoRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(109, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, withInfo, err := g2engine.ProcessRedoRecordWithInfo(ctx, request.GetFlags())
	response := pb.ProcessRedoRecordWithInfoResponse{
		Result:   result,
		WithInfo: withInfo,
	}
	if server.isTrace {
		defer server.traceExit(110, request, result, withInfo, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessWithInfo(ctx context.Context, request *pb.ProcessWithInfoRequest) (*pb.ProcessWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(111, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ProcessWithInfo(ctx, request.GetRecord(), request.GetFlags())
	response := pb.ProcessWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(112, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessWithResponse(ctx context.Context, request *pb.ProcessWithResponseRequest) (*pb.ProcessWithResponseResponse, error) {
	if server.isTrace {
		server.traceEntry(113, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ProcessWithResponse(ctx, request.GetRecord())
	response := pb.ProcessWithResponseResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(114, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ProcessWithResponseResize(ctx context.Context, request *pb.ProcessWithResponseResizeRequest) (*pb.ProcessWithResponseResizeResponse, error) {
	if server.isTrace {
		server.traceEntry(115, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ProcessWithResponseResize(ctx, request.GetRecord())
	response := pb.ProcessWithResponseResizeResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(116, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) PurgeRepository(ctx context.Context, request *pb.PurgeRepositoryRequest) (*pb.PurgeRepositoryResponse, error) {
	if server.isTrace {
		server.traceEntry(117, request)
	}
	entryTime := time.Now()
	// g2engine := getG2engine()
	// err := g2engine.PurgeRepository(ctx)
	err := server.getLogger().Error(4004)
	response := pb.PurgeRepositoryResponse{}
	if server.isTrace {
		defer server.traceExit(118, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReevaluateEntity(ctx context.Context, request *pb.ReevaluateEntityRequest) (*pb.ReevaluateEntityResponse, error) {
	if server.isTrace {
		server.traceEntry(119, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.ReevaluateEntity(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.ReevaluateEntityResponse{}
	if server.isTrace {
		defer server.traceExit(120, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReevaluateEntityWithInfo(ctx context.Context, request *pb.ReevaluateEntityWithInfoRequest) (*pb.ReevaluateEntityWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(121, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ReevaluateEntityWithInfo(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.ReevaluateEntityWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(122, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReevaluateRecord(ctx context.Context, request *pb.ReevaluateRecordRequest) (*pb.ReevaluateRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(123, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.ReevaluateRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.ReevaluateRecordResponse{}
	if server.isTrace {
		defer server.traceExit(124, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReevaluateRecordWithInfo(ctx context.Context, request *pb.ReevaluateRecordWithInfoRequest) (*pb.ReevaluateRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(125, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ReevaluateRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.ReevaluateRecordWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(126, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) Reinit(ctx context.Context, request *pb.ReinitRequest) (*pb.ReinitResponse, error) {
	if server.isTrace {
		server.traceEntry(127, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.Reinit(ctx, request.GetInitConfigID())
	response := pb.ReinitResponse{}
	if server.isTrace {
		defer server.traceExit(128, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReplaceRecord(ctx context.Context, request *pb.ReplaceRecordRequest) (*pb.ReplaceRecordResponse, error) {
	if server.isTrace {
		server.traceEntry(129, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	err := g2engine.ReplaceRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID())
	response := pb.ReplaceRecordResponse{}
	if server.isTrace {
		defer server.traceExit(130, request, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) ReplaceRecordWithInfo(ctx context.Context, request *pb.ReplaceRecordWithInfoRequest) (*pb.ReplaceRecordWithInfoResponse, error) {
	if server.isTrace {
		server.traceEntry(131, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.ReplaceRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := pb.ReplaceRecordWithInfoResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(132, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) SearchByAttributes(ctx context.Context, request *pb.SearchByAttributesRequest) (*pb.SearchByAttributesResponse, error) {
	if server.isTrace {
		server.traceEntry(133, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.SearchByAttributes(ctx, request.GetJsonData())
	response := pb.SearchByAttributesResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(134, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) SearchByAttributes_V2(ctx context.Context, request *pb.SearchByAttributes_V2Request) (*pb.SearchByAttributes_V2Response, error) {
	if server.isTrace {
		server.traceEntry(135, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.SearchByAttributes_V2(ctx, request.GetJsonData(), request.GetFlags())
	response := pb.SearchByAttributes_V2Response{
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

func (server *G2EngineServer) Stats(ctx context.Context, request *pb.StatsRequest) (*pb.StatsResponse, error) {
	if server.isTrace {
		server.traceEntry(139, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.Stats(ctx)
	response := pb.StatsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(140, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntities(ctx context.Context, request *pb.WhyEntitiesRequest) (*pb.WhyEntitiesResponse, error) {
	if server.isTrace {
		server.traceEntry(141, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntities(ctx, request.GetEntityID1(), request.GetEntityID2())
	response := pb.WhyEntitiesResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(142, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntities_V2(ctx context.Context, request *pb.WhyEntities_V2Request) (*pb.WhyEntities_V2Response, error) {
	if server.isTrace {
		server.traceEntry(143, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntities_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetFlags())
	response := pb.WhyEntities_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(144, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByEntityID(ctx context.Context, request *pb.WhyEntityByEntityIDRequest) (*pb.WhyEntityByEntityIDResponse, error) {
	if server.isTrace {
		server.traceEntry(145, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByEntityID(ctx, request.GetEntityID())
	response := pb.WhyEntityByEntityIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(146, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByEntityID_V2(ctx context.Context, request *pb.WhyEntityByEntityID_V2Request) (*pb.WhyEntityByEntityID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(147, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := pb.WhyEntityByEntityID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(148, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByRecordID(ctx context.Context, request *pb.WhyEntityByRecordIDRequest) (*pb.WhyEntityByRecordIDResponse, error) {
	if server.isTrace {
		server.traceEntry(149, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := pb.WhyEntityByRecordIDResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(150, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyEntityByRecordID_V2(ctx context.Context, request *pb.WhyEntityByRecordID_V2Request) (*pb.WhyEntityByRecordID_V2Response, error) {
	if server.isTrace {
		server.traceEntry(151, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyEntityByRecordID_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := pb.WhyEntityByRecordID_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(152, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyRecords(ctx context.Context, request *pb.WhyRecordsRequest) (*pb.WhyRecordsResponse, error) {
	if server.isTrace {
		server.traceEntry(153, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyRecords(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2())
	response := pb.WhyRecordsResponse{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(154, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

func (server *G2EngineServer) WhyRecords_V2(ctx context.Context, request *pb.WhyRecords_V2Request) (*pb.WhyRecords_V2Response, error) {
	if server.isTrace {
		server.traceEntry(155, request)
	}
	entryTime := time.Now()
	g2engine := getG2engine()
	result, err := g2engine.WhyRecords_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetFlags())
	response := pb.WhyRecords_V2Response{
		Result: result,
	}
	if server.isTrace {
		defer server.traceExit(156, request, result, err, time.Since(entryTime))
	}
	return &response, err
}

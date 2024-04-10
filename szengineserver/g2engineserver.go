package szengineserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	g2sdk "github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/sz"
	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
)

var (
	g2engineSingleton sz.SzEngine
	g2engineSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *SzEngineServer) getLogger() logging.LoggingInterface {
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
func (server *SzEngineServer) traceEntry(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (server *SzEngineServer) traceExit(messageNumber int, details ...interface{}) {
	server.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (server *SzEngineServer) error(messageNumber int, details ...interface{}) error {
	return server.getLogger().NewError(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2engine() sz.SzEngine {
	g2engineSyncOnce.Do(func() {
		g2engineSingleton = &g2sdk.G2engine{}
	})
	return g2engineSingleton
}

func GetSdkG2engine() sz.SzEngine {
	return getG2engine()
}

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/g2-sdk-go/g2engine.G2engine
// ----------------------------------------------------------------------------

func (server *SzEngineServer) AddRecord(ctx context.Context, request *g2pb.AddRecordRequest) (*g2pb.AddRecordResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	err = g2engine.AddRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID())
	response := g2pb.AddRecordResponse{}
	return &response, err
}

func (server *SzEngineServer) AddRecordWithInfo(ctx context.Context, request *g2pb.AddRecordWithInfoRequest) (*g2pb.AddRecordWithInfoResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(3, request)
		defer func() { server.traceExit(4, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.AddRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := g2pb.AddRecordWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) CloseExport(ctx context.Context, request *g2pb.CloseExportRequest) (*g2pb.CloseExportResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(13, request)
		defer func() { server.traceExit(14, request, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	err = g2engine.CloseExport(ctx, uintptr(request.GetResponseHandle()))
	response := g2pb.CloseExportResponse{}
	return &response, err
}

func (server *SzEngineServer) CountRedoRecords(ctx context.Context, request *g2pb.CountRedoRecordsRequest) (*g2pb.CountRedoRecordsResponse, error) {
	var err error = nil
	var result int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(15, request)
		defer func() { server.traceExit(16, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.CountRedoRecords(ctx)
	response := g2pb.CountRedoRecordsResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) DeleteRecord(ctx context.Context, request *g2pb.DeleteRecordRequest) (*g2pb.DeleteRecordResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(17, request)
		defer func() { server.traceExit(18, request, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	err = g2engine.DeleteRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetLoadID())
	response := g2pb.DeleteRecordResponse{}
	return &response, err
}

func (server *SzEngineServer) DeleteRecordWithInfo(ctx context.Context, request *g2pb.DeleteRecordWithInfoRequest) (*g2pb.DeleteRecordWithInfoResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(19, request)
		defer func() { server.traceExit(20, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.DeleteRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetLoadID(), request.GetFlags())
	response := g2pb.DeleteRecordWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) Destroy(ctx context.Context, request *g2pb.DestroyRequest) (*g2pb.DestroyResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(21, request)
		defer func() { server.traceExit(22, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2engine := getG2engine()
	// err := g2engine.Destroy(ctx)
	err = server.error(4001)
	response := g2pb.DestroyResponse{}
	return &response, err
}

func (server *SzEngineServer) ExportConfig(ctx context.Context, request *g2pb.ExportConfigRequest) (*g2pb.ExportConfigResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(23, request)
		defer func() { server.traceExit(24, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.ExportConfig(ctx)
	response := g2pb.ExportConfigResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) ExportConfigAndConfigID(ctx context.Context, request *g2pb.ExportConfigAndConfigIDRequest) (*g2pb.ExportConfigAndConfigIDResponse, error) {
	var err error = nil
	var result string
	var configId int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(25, request)
		defer func() { server.traceExit(26, request, result, configId, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, configId, err = g2engine.ExportConfigAndConfigID(ctx)
	response := g2pb.ExportConfigAndConfigIDResponse{
		Config:   result,
		ConfigID: configId,
	}
	return &response, err
}

func (server *SzEngineServer) ExportCSVEntityReport(ctx context.Context, request *g2pb.ExportCSVEntityReportRequest) (*g2pb.ExportCSVEntityReportResponse, error) {
	var err error = nil
	var result uintptr
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(27, request)
		defer func() { server.traceExit(28, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.ExportCSVEntityReport(ctx, request.GetCsvColumnList(), request.GetFlags())
	response := g2pb.ExportCSVEntityReportResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *SzEngineServer) ExportJSONEntityReport(ctx context.Context, request *g2pb.ExportJSONEntityReportRequest) (*g2pb.ExportJSONEntityReportResponse, error) {
	var err error = nil
	var result uintptr
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(29, request)
		defer func() { server.traceExit(30, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.ExportJSONEntityReport(ctx, request.GetFlags())
	response := g2pb.ExportJSONEntityReportResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *SzEngineServer) FetchNext(ctx context.Context, request *g2pb.FetchNextRequest) (*g2pb.FetchNextResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(31, request)
		defer func() { server.traceExit(32, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FetchNext(ctx, uintptr(request.GetResponseHandle()))
	response := g2pb.FetchNextResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindInterestingEntitiesByEntityID(ctx context.Context, request *g2pb.FindInterestingEntitiesByEntityIDRequest) (*g2pb.FindInterestingEntitiesByEntityIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(33, request)
		defer func() { server.traceExit(34, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindInterestingEntitiesByEntityID(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.FindInterestingEntitiesByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindInterestingEntitiesByRecordID(ctx context.Context, request *g2pb.FindInterestingEntitiesByRecordIDRequest) (*g2pb.FindInterestingEntitiesByRecordIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(35, request)
		defer func() { server.traceExit(36, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindInterestingEntitiesByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.FindInterestingEntitiesByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindNetworkByEntityID(ctx context.Context, request *g2pb.FindNetworkByEntityIDRequest) (*g2pb.FindNetworkByEntityIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(37, request)
		defer func() { server.traceExit(38, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindNetworkByEntityID(ctx, request.GetEntityList(), request.GetMaxDegree(), request.GetBuildOutDegree(), request.GetMaxEntities())
	response := g2pb.FindNetworkByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindNetworkByEntityID_V2(ctx context.Context, request *g2pb.FindNetworkByEntityID_V2Request) (*g2pb.FindNetworkByEntityID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(39, request)
		defer func() { server.traceExit(40, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindNetworkByEntityID_V2(ctx, request.GetEntityList(), request.GetMaxDegree(), request.GetBuildOutDegree(), request.GetMaxEntities(), request.GetFlags())
	response := g2pb.FindNetworkByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindNetworkByRecordID(ctx context.Context, request *g2pb.FindNetworkByRecordIDRequest) (*g2pb.FindNetworkByRecordIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(41, request)
		defer func() { server.traceExit(42, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindNetworkByRecordID(ctx, request.GetRecordList(), request.GetMaxDegree(), request.GetBuildOutDegree(), request.GetMaxEntities())
	response := g2pb.FindNetworkByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindNetworkByRecordID_V2(ctx context.Context, request *g2pb.FindNetworkByRecordID_V2Request) (*g2pb.FindNetworkByRecordID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(43, request)
		defer func() { server.traceExit(44, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindNetworkByRecordID_V2(ctx, request.GetRecordList(), request.GetMaxDegree(), request.GetBuildOutDegree(), request.GetMaxEntities(), request.GetFlags())
	response := g2pb.FindNetworkByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathByEntityID(ctx context.Context, request *g2pb.FindPathByEntityIDRequest) (*g2pb.FindPathByEntityIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(45, request)
		defer func() { server.traceExit(46, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetMaxDegree())
	response := g2pb.FindPathByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathByEntityID_V2(ctx context.Context, request *g2pb.FindPathByEntityID_V2Request) (*g2pb.FindPathByEntityID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(47, request)
		defer func() { server.traceExit(48, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetMaxDegree(), request.GetFlags())
	response := g2pb.FindPathByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathByRecordID(ctx context.Context, request *g2pb.FindPathByRecordIDRequest) (*g2pb.FindPathByRecordIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(49, request)
		defer func() { server.traceExit(50, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetMaxDegree())
	response := g2pb.FindPathByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathByRecordID_V2(ctx context.Context, request *g2pb.FindPathByRecordID_V2Request) (*g2pb.FindPathByRecordID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(51, request)
		defer func() { server.traceExit(52, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetMaxDegree(), request.GetFlags())
	response := g2pb.FindPathByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathExcludingByEntityID(ctx context.Context, request *g2pb.FindPathExcludingByEntityIDRequest) (*g2pb.FindPathExcludingByEntityIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(53, request)
		defer func() { server.traceExit(54, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathExcludingByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetMaxDegree(), request.GetExcludedEntities())
	response := g2pb.FindPathExcludingByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathExcludingByEntityID_V2(ctx context.Context, request *g2pb.FindPathExcludingByEntityID_V2Request) (*g2pb.FindPathExcludingByEntityID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(55, request)
		defer func() { server.traceExit(56, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathExcludingByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetMaxDegree(), request.GetExcludedEntities(), request.GetFlags())
	response := g2pb.FindPathExcludingByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathExcludingByRecordID(ctx context.Context, request *g2pb.FindPathExcludingByRecordIDRequest) (*g2pb.FindPathExcludingByRecordIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(57, request)
		defer func() { server.traceExit(58, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathExcludingByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetMaxDegree(), request.GetExcludedRecords())
	response := g2pb.FindPathExcludingByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathExcludingByRecordID_V2(ctx context.Context, request *g2pb.FindPathExcludingByRecordID_V2Request) (*g2pb.FindPathExcludingByRecordID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(59, request)
		defer func() { server.traceExit(60, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathExcludingByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetMaxDegree(), request.GetExcludedRecords(), request.GetFlags())
	response := g2pb.FindPathExcludingByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathIncludingSourceByEntityID(ctx context.Context, request *g2pb.FindPathIncludingSourceByEntityIDRequest) (*g2pb.FindPathIncludingSourceByEntityIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(61, request)
		defer func() { server.traceExit(62, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathIncludingSourceByEntityID(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetMaxDegree(), request.GetExcludedEntities(), request.GetRequiredDsrcs())
	response := g2pb.FindPathIncludingSourceByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathIncludingSourceByEntityID_V2(ctx context.Context, request *g2pb.FindPathIncludingSourceByEntityID_V2Request) (*g2pb.FindPathIncludingSourceByEntityID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(63, request)
		defer func() { server.traceExit(64, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetMaxDegree(), request.GetExcludedEntities(), request.GetRequiredDsrcs(), request.GetFlags())
	response := g2pb.FindPathIncludingSourceByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathIncludingSourceByRecordID(ctx context.Context, request *g2pb.FindPathIncludingSourceByRecordIDRequest) (*g2pb.FindPathIncludingSourceByRecordIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(65, request)
		defer func() { server.traceExit(66, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathIncludingSourceByRecordID(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetMaxDegree(), request.GetExcludedRecords(), request.GetRequiredDsrcs())
	response := g2pb.FindPathIncludingSourceByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathIncludingSourceByRecordID_V2(ctx context.Context, request *g2pb.FindPathIncludingSourceByRecordID_V2Request) (*g2pb.FindPathIncludingSourceByRecordID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(67, request)
		defer func() { server.traceExit(68, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetMaxDegree(), request.GetExcludedRecords(), request.GetRequiredDsrcs(), request.GetFlags())
	response := g2pb.FindPathIncludingSourceByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetActiveConfigID(ctx context.Context, request *g2pb.GetActiveConfigIDRequest) (*g2pb.GetActiveConfigIDResponse, error) {
	var err error = nil
	var result int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(69, request)
		defer func() { server.traceExit(70, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetActiveConfigID(ctx)
	response := g2pb.GetActiveConfigIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetEntityByEntityID(ctx context.Context, request *g2pb.GetEntityByEntityIDRequest) (*g2pb.GetEntityByEntityIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(71, request)
		defer func() { server.traceExit(72, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetEntityByEntityID(ctx, request.GetEntityID())
	response := g2pb.GetEntityByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetEntityByEntityID_V2(ctx context.Context, request *g2pb.GetEntityByEntityID_V2Request) (*g2pb.GetEntityByEntityID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(73, request)
		defer func() { server.traceExit(74, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.GetEntityByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetEntityByRecordID(ctx context.Context, request *g2pb.GetEntityByRecordIDRequest) (*g2pb.GetEntityByRecordIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(75, request)
		defer func() { server.traceExit(76, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetEntityByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := g2pb.GetEntityByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetEntityByRecordID_V2(ctx context.Context, request *g2pb.GetEntityByRecordID_V2Request) (*g2pb.GetEntityByRecordID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(77, request)
		defer func() { server.traceExit(78, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetEntityByRecordID_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.GetEntityByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetObserverOrigin(ctx context.Context) string {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(161)
		defer func() { server.traceExit(162, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	return g2engine.GetObserverOrigin(ctx)
}

func (server *SzEngineServer) GetRecord(ctx context.Context, request *g2pb.GetRecordRequest) (*g2pb.GetRecordResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(83, request)
		defer func() { server.traceExit(84, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetRecord(ctx, request.GetDataSourceCode(), request.GetRecordID())
	response := g2pb.GetRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetRecord_V2(ctx context.Context, request *g2pb.GetRecord_V2Request) (*g2pb.GetRecord_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(85, request)
		defer func() { server.traceExit(86, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetRecord_V2(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.GetRecord_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetRedoRecord(ctx context.Context, request *g2pb.GetRedoRecordRequest) (*g2pb.GetRedoRecordResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(87, request)
		defer func() { server.traceExit(88, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetRedoRecord(ctx)
	response := g2pb.GetRedoRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetRepositoryLastModifiedTime(ctx context.Context, request *g2pb.GetRepositoryLastModifiedTimeRequest) (*g2pb.GetRepositoryLastModifiedTimeResponse, error) {
	var err error = nil
	var result int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(89, request)
		defer func() { server.traceExit(90, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetRepositoryLastModifiedTime(ctx)
	response := g2pb.GetRepositoryLastModifiedTimeResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetVirtualEntityByRecordID(ctx context.Context, request *g2pb.GetVirtualEntityByRecordIDRequest) (*g2pb.GetVirtualEntityByRecordIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(91, request)
		defer func() { server.traceExit(92, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetVirtualEntityByRecordID(ctx, request.GetRecordList())
	response := g2pb.GetVirtualEntityByRecordIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetVirtualEntityByRecordID_V2(ctx context.Context, request *g2pb.GetVirtualEntityByRecordID_V2Request) (*g2pb.GetVirtualEntityByRecordID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(93, request)
		defer func() { server.traceExit(94, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.GetVirtualEntityByRecordID_V2(ctx, request.GetRecordList(), request.GetFlags())
	response := g2pb.GetVirtualEntityByRecordID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) HowEntityByEntityID(ctx context.Context, request *g2pb.HowEntityByEntityIDRequest) (*g2pb.HowEntityByEntityIDResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(95, request)
		defer func() { server.traceExit(96, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.HowEntityByEntityID(ctx, request.GetEntityID())
	response := g2pb.HowEntityByEntityIDResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) HowEntityByEntityID_V2(ctx context.Context, request *g2pb.HowEntityByEntityID_V2Request) (*g2pb.HowEntityByEntityID_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(97, request)
		defer func() { server.traceExit(98, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.HowEntityByEntityID_V2(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.HowEntityByEntityID_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) Init(ctx context.Context, request *g2pb.InitRequest) (*g2pb.InitResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(99, request)
		defer func() { server.traceExit(100, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2engine := getG2engine()
	// err := g2engine.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	err = server.error(4002)
	response := g2pb.InitResponse{}
	return &response, err
}

func (server *SzEngineServer) InitWithConfigID(ctx context.Context, request *g2pb.InitWithConfigIDRequest) (*g2pb.InitWithConfigIDResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(101, request)
		defer func() { server.traceExit(102, request, err, time.Since(entryTime)) }()
	}
	// Not allowed by gRPC server
	// g2engine := getG2engine()
	// err := g2engine.InitWithConfigID(ctx, request.GetModuleName(), request.GetIniParams(), request.GetInitConfigID(), int(request.GetVerboseLogging()))
	err = server.error(4003)
	response := g2pb.InitWithConfigIDResponse{}
	return &response, err
}

func (server *SzEngineServer) PrimeEngine(ctx context.Context, request *g2pb.PrimeEngineRequest) (*g2pb.PrimeEngineResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(103, request)
		defer func() { server.traceExit(104, request, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	err = g2engine.PrimeEngine(ctx)
	response := g2pb.PrimeEngineResponse{}
	return &response, err
}

func (server *SzEngineServer) ReevaluateEntity(ctx context.Context, request *g2pb.ReevaluateEntityRequest) (*g2pb.ReevaluateEntityResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(119, request)
		defer func() { server.traceExit(120, request, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	err = g2engine.ReevaluateEntity(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.ReevaluateEntityResponse{}
	return &response, err
}

func (server *SzEngineServer) ReevaluateEntityWithInfo(ctx context.Context, request *g2pb.ReevaluateEntityWithInfoRequest) (*g2pb.ReevaluateEntityWithInfoResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(121, request)
		defer func() { server.traceExit(122, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.ReevaluateEntityWithInfo(ctx, request.GetEntityID(), request.GetFlags())
	response := g2pb.ReevaluateEntityWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) ReevaluateRecord(ctx context.Context, request *g2pb.ReevaluateRecordRequest) (*g2pb.ReevaluateRecordResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(123, request)
		defer func() { server.traceExit(124, request, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	err = g2engine.ReevaluateRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.ReevaluateRecordResponse{}
	return &response, err
}

func (server *SzEngineServer) ReevaluateRecordWithInfo(ctx context.Context, request *g2pb.ReevaluateRecordWithInfoRequest) (*g2pb.ReevaluateRecordWithInfoResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(125, request)
		defer func() { server.traceExit(126, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.ReevaluateRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetFlags())
	response := g2pb.ReevaluateRecordWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(11, observer.GetObserverId(ctx))
		defer func() { server.traceExit(12, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	return g2engine.RegisterObserver(ctx, observer)
}

func (server *SzEngineServer) Reinit(ctx context.Context, request *g2pb.ReinitRequest) (*g2pb.ReinitResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(127, request)
		defer func() { server.traceExit(128, request, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	err = g2engine.Reinit(ctx, request.GetInitConfigID())
	response := g2pb.ReinitResponse{}
	return &response, err
}

func (server *SzEngineServer) ReplaceRecord(ctx context.Context, request *g2pb.ReplaceRecordRequest) (*g2pb.ReplaceRecordResponse, error) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(129, request)
		defer func() { server.traceExit(130, request, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	err = g2engine.ReplaceRecord(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID())
	response := g2pb.ReplaceRecordResponse{}
	return &response, err
}

func (server *SzEngineServer) ReplaceRecordWithInfo(ctx context.Context, request *g2pb.ReplaceRecordWithInfoRequest) (*g2pb.ReplaceRecordWithInfoResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(131, request)
		defer func() { server.traceExit(132, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.ReplaceRecordWithInfo(ctx, request.GetDataSourceCode(), request.GetRecordID(), request.GetJsonData(), request.GetLoadID(), request.GetFlags())
	response := g2pb.ReplaceRecordWithInfoResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) SearchByAttributes(ctx context.Context, request *g2pb.SearchByAttributesRequest) (*g2pb.SearchByAttributesResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(133, request)
		defer func() { server.traceExit(134, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.SearchByAttributes(ctx, request.GetJsonData())
	response := g2pb.SearchByAttributesResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) SearchByAttributes_V2(ctx context.Context, request *g2pb.SearchByAttributes_V2Request) (*g2pb.SearchByAttributes_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(135, request)
		defer func() { server.traceExit(136, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.SearchByAttributes_V2(ctx, request.GetJsonData(), request.GetFlags())
	response := g2pb.SearchByAttributes_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(137, logLevelName)
		defer func() { server.traceExit(138, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	g2engine := getG2engine()
	err = g2engine.SetLogLevel(ctx, logLevelName)
	if err != nil {
		return err
	}
	err = server.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		return err
	}
	server.isTrace = (logLevelName == logging.LevelTraceName)
	return err
}

func (server *SzEngineServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(163, origin)
		defer func() { server.traceExit(164, origin, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	g2engine.SetObserverOrigin(ctx, origin)
}

func (server *SzEngineServer) Stats(ctx context.Context, request *g2pb.StatsRequest) (*g2pb.StatsResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(139, request)
		defer func() { server.traceExit(140, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.Stats(ctx)
	response := g2pb.StatsResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) StreamExportCSVEntityReport(request *g2pb.StreamExportCSVEntityReportRequest, stream g2pb.G2Engine_StreamExportCSVEntityReportServer) (err error) {
	if server.isTrace {
		server.traceEntry(157, request)
	}
	ctx := stream.Context()
	entryTime := time.Now()
	g2engine := getG2engine()
	rowsFetched := 0

	// Get the query handle.

	var queryHandle uintptr
	queryHandle, err = g2engine.ExportCSVEntityReport(ctx, request.GetCsvColumnList(), request.GetFlags())
	if err != nil {
		return err
	}

	// Defer the CloseExport in case we exit early for any reason.

	defer func() {
		err = g2engine.CloseExport(ctx, queryHandle)
		if server.isTrace {
			server.traceExit(158, request, rowsFetched, err, time.Since(entryTime))
		}
	}()

	// Stream the results.

	for {
		var fetchResult string
		fetchResult, err = g2engine.FetchNext(ctx, queryHandle)
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
		server.traceEntry(601, request, fetchResult)
		rowsFetched += 1
	}

	err = nil
	return
}

func (server *SzEngineServer) StreamExportJSONEntityReport(request *g2pb.StreamExportJSONEntityReportRequest, stream g2pb.G2Engine_StreamExportJSONEntityReportServer) (err error) {
	if server.isTrace {
		server.traceEntry(159, request)
	}
	ctx := stream.Context()
	entryTime := time.Now()
	g2engine := getG2engine()
	rowsFetched := 0

	// Get the query handle.

	var queryHandle uintptr
	queryHandle, err = g2engine.ExportJSONEntityReport(ctx, request.GetFlags())
	if err != nil {
		return err
	}

	// Defer the CloseExport in case we exit early for any reason.

	defer func() {
		err = g2engine.CloseExport(ctx, queryHandle)
		if server.isTrace {
			server.traceExit(160, request, rowsFetched, err, time.Since(entryTime))
		}
	}()

	// Stream the results.

	for {
		var fetchResult string
		fetchResult, err = g2engine.FetchNext(ctx, queryHandle)
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
		server.traceEntry(602, request, fetchResult)
		rowsFetched += 1
	}

	err = nil
	return
}

func (server *SzEngineServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(79, observer.GetObserverId(ctx))
		defer func() { server.traceExit(80, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	return g2engine.UnregisterObserver(ctx, observer)
}

func (server *SzEngineServer) WhyEntities(ctx context.Context, request *g2pb.WhyEntitiesRequest) (*g2pb.WhyEntitiesResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(141, request)
		defer func() { server.traceExit(142, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.WhyEntities(ctx, request.GetEntityID1(), request.GetEntityID2())
	response := g2pb.WhyEntitiesResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) WhyEntities_V2(ctx context.Context, request *g2pb.WhyEntities_V2Request) (*g2pb.WhyEntities_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(143, request)
		defer func() { server.traceExit(144, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.WhyEntities_V2(ctx, request.GetEntityID1(), request.GetEntityID2(), request.GetFlags())
	response := g2pb.WhyEntities_V2Response{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) WhyRecords(ctx context.Context, request *g2pb.WhyRecordsRequest) (*g2pb.WhyRecordsResponse, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(153, request)
		defer func() { server.traceExit(154, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.WhyRecords(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2())
	response := g2pb.WhyRecordsResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) WhyRecords_V2(ctx context.Context, request *g2pb.WhyRecords_V2Request) (*g2pb.WhyRecords_V2Response, error) {
	var err error = nil
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(155, request)
		defer func() { server.traceExit(156, request, result, err, time.Since(entryTime)) }()
	}
	g2engine := getG2engine()
	result, err = g2engine.WhyRecords_V2(ctx, request.GetDataSourceCode1(), request.GetRecordID1(), request.GetDataSourceCode2(), request.GetRecordID2(), request.GetFlags())
	response := g2pb.WhyRecords_V2Response{
		Result: result,
	}
	return &response, err
}

package szengineserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	szsdk "github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
)

var (
	szEngineSingleton *szsdk.Szengine
	szEngineSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Interface methods for github.com/senzing-garage/sz-sdk-go/szengine.SzEngine
// ----------------------------------------------------------------------------

func (server *SzEngineServer) AddRecord(ctx context.Context, request *szpb.AddRecordRequest) (*szpb.AddRecordResponse, error) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(1, request)
		defer func() { server.traceExit(2, request, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err := szEngine.AddRecord(ctx, request.GetDataSourceCode(), request.GetRecordId(), request.GetRecordDefinition(), request.GetFlags())
	response := szpb.AddRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) CloseExport(ctx context.Context, request *szpb.CloseExportRequest) (*szpb.CloseExportResponse, error) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(13, request)
		defer func() { server.traceExit(14, request, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	err = szEngine.CloseExport(ctx, uintptr(request.GetExportHandle()))
	response := szpb.CloseExportResponse{}
	return &response, err
}

func (server *SzEngineServer) CountRedoRecords(ctx context.Context, request *szpb.CountRedoRecordsRequest) (*szpb.CountRedoRecordsResponse, error) {
	var err error
	var result int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(15, request)
		defer func() { server.traceExit(16, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.CountRedoRecords(ctx)
	response := szpb.CountRedoRecordsResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) DeleteRecord(ctx context.Context, request *szpb.DeleteRecordRequest) (*szpb.DeleteRecordResponse, error) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(17, request)
		defer func() { server.traceExit(18, request, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err := szEngine.DeleteRecord(ctx, request.GetDataSourceCode(), request.GetRecordId(), request.GetFlags())
	response := szpb.DeleteRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) ExportCsvEntityReport(ctx context.Context, request *szpb.ExportCsvEntityReportRequest) (*szpb.ExportCsvEntityReportResponse, error) {
	var err error
	var result uintptr
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(27, request)
		defer func() { server.traceExit(28, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.ExportCsvEntityReport(ctx, request.GetCsvColumnList(), request.GetFlags())
	response := szpb.ExportCsvEntityReportResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *SzEngineServer) ExportJsonEntityReport(ctx context.Context, request *szpb.ExportJsonEntityReportRequest) (*szpb.ExportJsonEntityReportResponse, error) { //revive:disable-line var-naming
	var err error
	var result uintptr
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(29, request)
		defer func() { server.traceExit(30, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.ExportJSONEntityReport(ctx, request.GetFlags())
	response := szpb.ExportJsonEntityReportResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *SzEngineServer) FetchNext(ctx context.Context, request *szpb.FetchNextRequest) (*szpb.FetchNextResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(31, request)
		defer func() { server.traceExit(32, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.FetchNext(ctx, uintptr(request.GetExportHandle()))
	response := szpb.FetchNextResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindInterestingEntitiesByEntityId(ctx context.Context, request *szpb.FindInterestingEntitiesByEntityIdRequest) (*szpb.FindInterestingEntitiesByEntityIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(33, request)
		defer func() { server.traceExit(34, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.FindInterestingEntitiesByEntityID(ctx, request.GetEntityId(), request.GetFlags())
	response := szpb.FindInterestingEntitiesByEntityIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindInterestingEntitiesByRecordId(ctx context.Context, request *szpb.FindInterestingEntitiesByRecordIdRequest) (*szpb.FindInterestingEntitiesByRecordIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(35, request)
		defer func() { server.traceExit(36, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.FindInterestingEntitiesByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordId(), request.GetFlags())
	response := szpb.FindInterestingEntitiesByRecordIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindNetworkByEntityId(ctx context.Context, request *szpb.FindNetworkByEntityIdRequest) (*szpb.FindNetworkByEntityIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(37, request)
		defer func() { server.traceExit(38, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.FindNetworkByEntityID(ctx, request.GetEntityIds(), request.GetMaxDegrees(), request.GetBuildOutDegree(), request.GetBuildOutMaxEntities(), request.GetFlags())
	response := szpb.FindNetworkByEntityIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindNetworkByRecordId(ctx context.Context, request *szpb.FindNetworkByRecordIdRequest) (*szpb.FindNetworkByRecordIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(41, request)
		defer func() { server.traceExit(42, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.FindNetworkByRecordID(ctx, request.GetRecordKeys(), request.GetMaxDegrees(), request.GetBuildOutDegree(), request.GetBuildOutMaxEntities(), request.GetFlags())
	response := szpb.FindNetworkByRecordIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathByEntityId(ctx context.Context, request *szpb.FindPathByEntityIdRequest) (*szpb.FindPathByEntityIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(45, request)
		defer func() { server.traceExit(46, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.FindPathByEntityID(ctx, request.GetStartEntityId(), request.GetEndEntityId(), request.GetMaxDegrees(), request.GetAvoidEntityIds(), request.GetRequiredDataSources(), request.GetFlags())
	response := szpb.FindPathByEntityIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) FindPathByRecordId(ctx context.Context, request *szpb.FindPathByRecordIdRequest) (*szpb.FindPathByRecordIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(49, request)
		defer func() { server.traceExit(50, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.FindPathByRecordID(ctx, request.GetStartDataSourceCode(), request.GetStartRecordId(), request.GetEndDataSourceCode(), request.GetEndRecordId(), request.GetMaxDegrees(), request.GetAvoidRecordKeys(), request.GetRequiredDataSources(), request.GetFlags())
	response := szpb.FindPathByRecordIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetActiveConfigId(ctx context.Context, request *szpb.GetActiveConfigIdRequest) (*szpb.GetActiveConfigIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result int64
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(69, request)
		defer func() { server.traceExit(70, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.GetActiveConfigID(ctx)
	response := szpb.GetActiveConfigIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetEntityByEntityId(ctx context.Context, request *szpb.GetEntityByEntityIdRequest) (*szpb.GetEntityByEntityIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(71, request)
		defer func() { server.traceExit(72, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.GetEntityByEntityID(ctx, request.GetEntityId(), request.GetFlags())
	response := szpb.GetEntityByEntityIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetEntityByRecordId(ctx context.Context, request *szpb.GetEntityByRecordIdRequest) (*szpb.GetEntityByRecordIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(75, request)
		defer func() { server.traceExit(76, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.GetEntityByRecordID(ctx, request.GetDataSourceCode(), request.GetRecordId(), request.GetFlags())
	response := szpb.GetEntityByRecordIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetRecord(ctx context.Context, request *szpb.GetRecordRequest) (*szpb.GetRecordResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(83, request)
		defer func() { server.traceExit(84, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.GetRecord(ctx, request.GetDataSourceCode(), request.GetRecordId(), request.GetFlags())
	response := szpb.GetRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetRedoRecord(ctx context.Context, request *szpb.GetRedoRecordRequest) (*szpb.GetRedoRecordResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(87, request)
		defer func() { server.traceExit(88, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.GetRedoRecord(ctx)
	response := szpb.GetRedoRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetStats(ctx context.Context, request *szpb.GetStatsRequest) (*szpb.GetStatsResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(139, request)
		defer func() { server.traceExit(140, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.GetStats(ctx)
	response := szpb.GetStatsResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) GetVirtualEntityByRecordId(ctx context.Context, request *szpb.GetVirtualEntityByRecordIdRequest) (*szpb.GetVirtualEntityByRecordIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(91, request)
		defer func() { server.traceExit(92, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.GetVirtualEntityByRecordID(ctx, request.GetRecordKeys(), request.GetFlags())
	response := szpb.GetVirtualEntityByRecordIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) HowEntityByEntityId(ctx context.Context, request *szpb.HowEntityByEntityIdRequest) (*szpb.HowEntityByEntityIdResponse, error) { //revive:disable-line var-naming
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(95, request)
		defer func() { server.traceExit(96, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.HowEntityByEntityID(ctx, request.GetEntityId(), request.GetFlags())
	response := szpb.HowEntityByEntityIdResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) PrimeEngine(ctx context.Context, request *szpb.PrimeEngineRequest) (*szpb.PrimeEngineResponse, error) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(103, request)
		defer func() { server.traceExit(104, request, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	err = szEngine.PrimeEngine(ctx)
	response := szpb.PrimeEngineResponse{}
	return &response, err
}

func (server *SzEngineServer) ProcessRedoRecord(ctx context.Context, request *szpb.ProcessRedoRecordRequest) (*szpb.ProcessRedoRecordResponse, error) {
	// TODO: Fix trace IDs.
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(999, request)
		defer func() { server.traceExit(999, request, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err := szEngine.ProcessRedoRecord(ctx, request.GetRedoRecord(), request.GetFlags())
	response := szpb.ProcessRedoRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) ReevaluateEntity(ctx context.Context, request *szpb.ReevaluateEntityRequest) (*szpb.ReevaluateEntityResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(119, request)
		defer func() { server.traceExit(120, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.ReevaluateEntity(ctx, request.GetEntityId(), request.GetFlags())
	response := szpb.ReevaluateEntityResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) ReevaluateRecord(ctx context.Context, request *szpb.ReevaluateRecordRequest) (*szpb.ReevaluateRecordResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(123, request)
		defer func() { server.traceExit(124, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.ReevaluateRecord(ctx, request.GetDataSourceCode(), request.GetRecordId(), request.GetFlags())
	response := szpb.ReevaluateRecordResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) Reinitialize(ctx context.Context, request *szpb.ReinitializeRequest) (*szpb.ReinitializeResponse, error) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(127, request)
		defer func() { server.traceExit(128, request, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	err = szEngine.Reinitialize(ctx, request.GetConfigId())
	response := szpb.ReinitializeResponse{}
	return &response, err
}

func (server *SzEngineServer) SearchByAttributes(ctx context.Context, request *szpb.SearchByAttributesRequest) (*szpb.SearchByAttributesResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(133, request)
		defer func() { server.traceExit(134, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.SearchByAttributes(ctx, request.GetAttributes(), request.GetSearchProfile(), request.GetFlags())
	response := szpb.SearchByAttributesResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) StreamExportCsvEntityReport(request *szpb.StreamExportCsvEntityReportRequest, stream szpb.SzEngine_StreamExportCsvEntityReportServer) (err error) {
	if server.isTrace {
		server.traceEntry(157, request)
	}
	ctx := stream.Context()
	entryTime := time.Now()
	szEngine := getSzEngine()
	rowsFetched := 0

	// Get the query handle.

	var queryHandle uintptr
	queryHandle, err = szEngine.ExportCsvEntityReport(ctx, request.GetCsvColumnList(), request.GetFlags())
	if err != nil {
		return err
	}

	// Defer the CloseExport in case we exit early for any reason.

	defer func() {
		err = szEngine.CloseExport(ctx, queryHandle)
		if server.isTrace {
			server.traceExit(158, request, rowsFetched, err, time.Since(entryTime))
		}
	}()

	// Stream the results.

	for {
		var fetchResult string
		fetchResult, err = szEngine.FetchNext(ctx, queryHandle)
		if err != nil {
			return err
		}
		if len(fetchResult) == 0 {
			break
		}
		response := szpb.StreamExportCsvEntityReportResponse{
			Result: fetchResult,
		}
		if err = stream.Send(&response); err != nil {
			return err
		}
		server.traceEntry(601, request, fetchResult)
		rowsFetched++
	}

	err = nil
	return
}

func (server *SzEngineServer) StreamExportJsonEntityReport(request *szpb.StreamExportJsonEntityReportRequest, stream szpb.SzEngine_StreamExportJsonEntityReportServer) (err error) { //revive:disable-line var-naming
	if server.isTrace {
		server.traceEntry(159, request)
	}
	ctx := stream.Context()
	entryTime := time.Now()
	szEngine := getSzEngine()
	rowsFetched := 0

	// Get the query handle.

	var queryHandle uintptr
	queryHandle, err = szEngine.ExportJSONEntityReport(ctx, request.GetFlags())
	if err != nil {
		return err
	}

	// Defer the CloseExport in case we exit early for any reason.

	defer func() {
		err = szEngine.CloseExport(ctx, queryHandle)
		if server.isTrace {
			server.traceExit(160, request, rowsFetched, err, time.Since(entryTime))
		}
	}()

	// Stream the results.

	for {
		var fetchResult string
		fetchResult, err = szEngine.FetchNext(ctx, queryHandle)
		if err != nil {
			return err
		}
		if len(fetchResult) == 0 {
			break
		}
		response := szpb.StreamExportJsonEntityReportResponse{
			Result: fetchResult,
		}
		if err = stream.Send(&response); err != nil {
			return err
		}
		server.traceEntry(602, request, fetchResult)
		rowsFetched++
	}

	err = nil
	return
}

func (server *SzEngineServer) WhyEntities(ctx context.Context, request *szpb.WhyEntitiesRequest) (*szpb.WhyEntitiesResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(141, request)
		defer func() { server.traceExit(142, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.WhyEntities(ctx, request.GetEntityId1(), request.GetEntityId2(), request.GetFlags())
	response := szpb.WhyEntitiesResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) WhyRecordInEntity(ctx context.Context, request *szpb.WhyRecordInEntityRequest) (*szpb.WhyRecordInEntityResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(153, request)
		defer func() { server.traceExit(154, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.WhyRecordInEntity(ctx, request.GetDataSourceCode(), request.GetRecordId(), request.GetFlags())
	response := szpb.WhyRecordInEntityResponse{
		Result: result,
	}
	return &response, err
}

func (server *SzEngineServer) WhyRecords(ctx context.Context, request *szpb.WhyRecordsRequest) (*szpb.WhyRecordsResponse, error) {
	var err error
	var result string
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(153, request)
		defer func() { server.traceExit(154, request, result, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	result, err = szEngine.WhyRecords(ctx, request.GetDataSourceCode1(), request.GetRecordId1(), request.GetDataSourceCode2(), request.GetRecordId2(), request.GetFlags())
	response := szpb.WhyRecordsResponse{
		Result: result,
	}
	return &response, err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (server *SzEngineServer) getLogger() logging.Logging {
	var err error
	if server.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 3},
		}
		server.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
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

func (server *SzEngineServer) SetLogLevel(ctx context.Context, logLevelName string) error {
	_ = ctx
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(137, logLevelName)
		defer func() { server.traceExit(138, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	// szengine := getSzengine()
	// err = szengine.SetLogLevel(ctx, logLevelName)
	// if err != nil {
	// 	return err
	// }
	err = server.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		return err
	}
	server.isTrace = (logLevelName == logging.LevelTraceName)
	return err
}

// --- Errors -----------------------------------------------------------------

// Create error.
// func (server *SzEngineServer) error(messageNumber int, details ...interface{}) error {
// 	return server.getLogger().NewError(messageNumber, details...)
// }

// --- Services ---------------------------------------------------------------

// Singleton pattern for szconfig.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getSzEngine() *szsdk.Szengine {
	szEngineSyncOnce.Do(func() {
		szEngineSingleton = &szsdk.Szengine{}
	})
	return szEngineSingleton
}

func GetSdkSzEngine() *szsdk.Szengine {
	return getSzEngine()
}

func GetSdkSzEngineAsInterface() senzing.SzEngine {
	return getSzEngine()
}

// --- Observer ---------------------------------------------------------------

func (server *SzEngineServer) GetObserverOrigin(ctx context.Context) string {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(161)
		defer func() { server.traceExit(162, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	return szEngine.GetObserverOrigin(ctx)
}

func (server *SzEngineServer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(11, observer.GetObserverID(ctx))
		defer func() { server.traceExit(12, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	return szEngine.RegisterObserver(ctx, observer)
}

func (server *SzEngineServer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(163, origin)
		defer func() { server.traceExit(164, origin, err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	szEngine.SetObserverOrigin(ctx, origin)
}

func (server *SzEngineServer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if server.isTrace {
		entryTime := time.Now()
		server.traceEntry(79, observer.GetObserverID(ctx))
		defer func() { server.traceExit(80, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	szEngine := getSzEngine()
	return szEngine.UnregisterObserver(ctx, observer)
}

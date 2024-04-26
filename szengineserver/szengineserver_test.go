package szengineserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

type GetEntityByRecordIdResponse struct {
	ResolvedEntity struct {
		EntityId int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

var (
	szEngineTestSingleton *SzEngineServer
	localLogger           logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzEngineServer_AddRecord(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	request1 := &szpb.AddRecordRequest{
		DataSourceCode:   record1.DataSource,
		Flags:            sz.SZ_WITH_INFO,
		RecordDefinition: record1.Json,
		RecordId:         record1.Id,
	}
	response1, err := szEngineServer.AddRecord(ctx, request1)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response1.GetResult())
	request2 := &szpb.AddRecordRequest{
		DataSourceCode:   record2.DataSource,
		Flags:            sz.SZ_WITH_INFO,
		RecordDefinition: record2.Json,
		RecordId:         record2.Id,
	}
	response2, err := szEngineServer.AddRecord(ctx, request2)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response2.GetResult())
}

func TestSzEngineServer_AddRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	request := &szpb.AddRecordRequest{
		DataSourceCode:   record.DataSource,
		Flags:            sz.SZ_WITH_INFO,
		RecordDefinition: record.Json,
		RecordId:         record.Id,
	}
	response, err := szEngineServer.AddRecord(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_CountRedoRecords(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.CountRedoRecordsRequest{}
	response, err := szEngineServer.CountRedoRecords(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_ExportJsonEntityReport(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	flags := sz.SZ_NO_FLAGS
	request := &szpb.ExportJsonEntityReportRequest{
		Flags: flags,
	}
	response, err := szEngineServer.ExportJsonEntityReport(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_ExportCsvEntityReport(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.ExportCsvEntityReportRequest{
		CsvColumnList: "",
		Flags:         sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.ExportCsvEntityReport(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_FindNetworkByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}, {"ENTITY_ID": ` + getEntityIdString(record2) + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := sz.SZ_NO_FLAGS
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegree: buildOutDegree,
		EntityList:     entityList,
		Flags:          flags,
		MaxDegrees:     maxDegrees,
		MaxEntities:    maxEntities,
	}
	response, err := szEngineServer.FindNetworkByEntityId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_FindNetworkByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.Id + `"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := sz.SZ_NO_FLAGS
	request := &szpb.FindNetworkByRecordIdRequest{
		BuildOutDegree: buildOutDegree,
		Flags:          flags,
		MaxDegrees:     maxDegrees,
		MaxEntities:    maxEntities,
		RecordList:     recordList,
	}
	response, err := szEngineServer.FindNetworkByRecordId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_FindPathByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	startEntityId := getEntityId(truthset.CustomerRecords["1001"])
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	flags := sz.SZ_NO_FLAGS
	request := &szpb.FindPathByEntityIdRequest{
		EndEntityId:   endEntityId,
		Flags:         flags,
		MaxDegrees:    maxDegrees,
		StartEntityId: startEntityId,
	}
	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_FindPathByEntityId_exclusions(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	startEntityId := getEntityId(record1)
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	flags := sz.SZ_NO_FLAGS
	request := &szpb.FindPathByEntityIdRequest{
		EndEntityId:   endEntityId,
		Exclusions:    exclusions,
		Flags:         flags,
		MaxDegrees:    maxDegrees,
		StartEntityId: startEntityId,
	}
	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_FindPathByEntityId_inclusions(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	startEntityId := getEntityId(record1)
	endEntityId := getEntityId(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	request := &szpb.FindPathByEntityIdRequest{
		EndEntityId:         endEntityId,
		Exclusions:          exclusions,
		MaxDegrees:          maxDegrees,
		RequiredDataSources: requiredDataSources,
		StartEntityId:       startEntityId,
	}
	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_FindPathByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegrees := int64(1)
	flags := sz.SZ_NO_FLAGS
	request := &szpb.FindPathByRecordIdRequest{
		EndDataSourceCode:   record2.DataSource,
		EndRecordId:         record2.Id,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		StartDataSourceCode: record1.DataSource,
		StartRecordId:       record1.Id,
	}
	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_FindPathByRecordId_exclusions(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegrees := int64(1)
	exclusions := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	flags := sz.SZ_NO_FLAGS
	request := &szpb.FindPathByRecordIdRequest{
		EndDataSourceCode:   record2.DataSource,
		EndRecordId:         record2.Id,
		Exclusions:          exclusions,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		StartDataSourceCode: record1.DataSource,
		StartRecordId:       record1.Id,
	}
	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_FindPathByRecordId_inclusions(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegrees := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := sz.SZ_NO_FLAGS
	request := &szpb.FindPathByRecordIdRequest{
		EndDataSourceCode:   record2.DataSource,
		EndRecordId:         record1.Id,
		Exclusions:          exclusions,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		RequiredDataSources: requiredDataSources,
		StartDataSourceCode: record1.DataSource,
		StartRecordId:       record1.Id,
	}
	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_GetActiveConfigId(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.GetActiveConfigIdRequest{}
	response, err := szEngineServer.GetActiveConfigId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_GetEntityByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_NO_FLAGS
	request := &szpb.GetEntityByEntityIdRequest{
		EntityId: entityId,
		Flags:    flags,
	}
	response, err := szEngineServer.GetEntityByEntityId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_GetEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_NO_FLAGS
	request := &szpb.GetEntityByRecordIdRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.Id,
	}
	response, err := szEngineServer.GetEntityByRecordId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_GetRecord(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_NO_FLAGS
	request := &szpb.GetRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.Id,
	}
	response, err := szEngineServer.GetRecord(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_GetRedoRecord(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.GetRedoRecordRequest{}
	response, err := szEngineServer.GetRedoRecord(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_GetRepositoryLastModifiedTime(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.GetRepositoryLastModifiedTimeRequest{}
	response, err := szEngineServer.GetRepositoryLastModifiedTime(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_GetVirtualEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}]}`
	flags := sz.SZ_NO_FLAGS
	request := &szpb.GetVirtualEntityByRecordIdRequest{
		Flags:      flags,
		RecordList: recordList,
	}
	response, err := szEngineServer.GetVirtualEntityByRecordId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_HowEntityByEntityId(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_NO_FLAGS
	request := &szpb.HowEntityByEntityIdRequest{
		EntityId: entityId,
		Flags:    flags,
	}
	response, err := szEngineServer.HowEntityByEntityId(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_PrimeEngine(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.PrimeEngineRequest{}
	response, err := szEngineServer.PrimeEngine(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response)
}

func TestSzEngineServer_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_WITHOUT_INFO
	request := &szpb.ReevaluateEntityRequest{
		EntityId: entityId,
		Flags:    flags,
	}
	response, err := szEngineServer.ReevaluateEntity(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_ReevaluateEntity_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := sz.SZ_WITH_INFO
	request := &szpb.ReevaluateEntityRequest{
		EntityId: entityId,
		Flags:    flags,
	}
	response, err := szEngineServer.ReevaluateEntity(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_ReevaluateRecord(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_WITHOUT_INFO
	request := &szpb.ReevaluateRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.Id,
	}
	response, err := szEngineServer.ReevaluateRecord(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_ReevaluateRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := sz.SZ_WITH_INFO
	request := &szpb.ReevaluateRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.Id,
	}
	response, err := szEngineServer.ReevaluateRecord(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)

	requestToGetActiveConfigId := &szpb.GetActiveConfigIdRequest{}
	responseFromGetActiveConfigId, err := szEngineServer.GetActiveConfigId(ctx, requestToGetActiveConfigId)
	testError(test, ctx, szEngineServer, err)

	request := &szpb.ReinitializeRequest{
		ConfigId: responseFromGetActiveConfigId.GetResult(),
	}
	response, err := szEngineServer.Reinitialize(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response)
}

func TestSzEngineServer_SearchByAttributes(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	flags := sz.SZ_NO_FLAGS
	request := &szpb.SearchByAttributesRequest{
		Attributes: attributes,
		Flags:      flags,
	}
	response, err := szEngineServer.SearchByAttributes(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_SearchByAttributes_searchProfile(test *testing.T) {
	// TODO:  Use actual searchProfile
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile := "SEARCH"
	flags := sz.SZ_NO_FLAGS
	request := &szpb.SearchByAttributesRequest{
		Attributes:    attributes,
		Flags:         flags,
		SearchProfile: searchProfile,
	}
	response, err := szEngineServer.SearchByAttributes(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_Stats(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.GetStatsRequest{}
	response, err := szEngineServer.GetStats(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_WhyEntities(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	entityId1 := getEntityId(truthset.CustomerRecords["1001"])
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	flags := sz.SZ_NO_FLAGS
	request := &szpb.WhyEntitiesRequest{
		EntityId1: entityId1,
		EntityId2: entityId2,
		Flags:     flags,
	}
	response, err := szEngineServer.WhyEntities(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_WhyRecordInEntity(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	flags := sz.SZ_NO_FLAGS
	request := &szpb.WhyRecordInEntityRequest{
		DataSourceCode: record1.DataSource,
		Flags:          flags,
		RecordId:       record1.Id,
	}
	response, err := szEngineServer.WhyRecordInEntity(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := sz.SZ_NO_FLAGS
	request := &szpb.WhyRecordsRequest{
		DataSourceCode1: record1.DataSource,
		DataSourceCode2: record2.DataSource,
		Flags:           flags,
		RecordId1:       record1.Id,
		RecordId2:       record2.Id,
	}
	response, err := szEngineServer.WhyRecords(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_DeleteRecord(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	flags := sz.SZ_WITHOUT_INFO
	request := &szpb.DeleteRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.Id,
	}
	response, err := szEngineServer.DeleteRecord(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

func TestSzEngineServer_DeleteRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	flags := sz.SZ_WITH_INFO
	request := &szpb.DeleteRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.Id,
	}
	response, err := szEngineServer.DeleteRecord(ctx, request)
	testError(test, ctx, szEngineServer, err)
	printResponse(test, response.GetResult())
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(localLogger.NewError(errorId, err), err)
}

func getEntityId(record record.Record) int64 {
	return getEntityIdForRecord(record.DataSource, record.Id)
}

func getEntityIdForRecord(datasource string, id string) int64 {
	ctx := context.TODO()
	var result int64 = 0
	szEngine := getSzEngineServer(ctx)
	request := &szpb.GetEntityByRecordIdRequest{
		DataSourceCode: datasource,
		RecordId:       id,
	}
	response, err := szEngine.GetEntityByRecordId(ctx, request)
	if err != nil {
		return result
	}

	getEntityByRecordIdResponse := &GetEntityByRecordIdResponse{}
	err = json.Unmarshal([]byte(response.Result), &getEntityByRecordIdResponse)
	if err != nil {
		return result
	}
	return getEntityByRecordIdResponse.ResolvedEntity.EntityId
}

func getEntityIdString(record record.Record) string {
	entityId := getEntityId(record)
	return strconv.FormatInt(entityId, 10)
}

func getEntityIdStringForRecord(datasource string, id string) string {
	entityId := getEntityIdForRecord(datasource, id)
	return strconv.FormatInt(entityId, 10)
}

func getSzEngineServer(ctx context.Context) SzEngineServer {
	if szEngineTestSingleton == nil {
		szEngineTestSingleton = &SzEngineServer{}
		instanceName := "Test name"
		verboseLogging := sz.SZ_NO_LOGGING
		configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
		setting, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkSzEngine().Initialize(ctx, instanceName, setting, configId, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *szEngineTestSingleton
}

func getTestObject(ctx context.Context, test *testing.T) SzEngineServer {
	if szEngineTestSingleton == nil {
		szEngineTestSingleton = &SzEngineServer{}
		instanceName := "Test name"
		verboseLogging := sz.SZ_NO_LOGGING
		configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
		settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkSzEngine().Initialize(ctx, instanceName, settings, configId, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *szEngineTestSingleton
}

func printResponse(test *testing.T, response interface{}) {
	printResult(test, "Response", response)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func testError(test *testing.T, ctx context.Context, szEngine SzEngineServer, err error) {
	_ = ctx
	_ = szEngine
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		if szerror.Is(err, szerror.SzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzBadInput) {
			fmt.Printf("\nBad user input error detected. \n\n")
		}
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setupSenzingConfig(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	now := time.Now()

	szConfig := &szconfig.Szconfig{}
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range datasourceNames {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return createError(5908, err)
		}
	}

	configStr, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	err = szConfig.Destroy(ctx)
	if err != nil {
		return createError(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComments := fmt.Sprintf("Created by szengine_test at %s", now.UTC())
	configId, err := szConfigManager.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return createError(5914, err)
	}

	err = szConfigManager.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, instanceName string, settings string, verboseLogging int64, configId int64) error {
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	err := szDiagnostic.Initialize(ctx, instanceName, settings, configId, verboseLogging)
	if err != nil {
		return createError(5903, err)
	}

	err = szDiagnostic.PurgeRepository(ctx)
	if err != nil {
		return createError(5904, err)
	}

	err = szDiagnostic.Destroy(ctx)
	if err != nil {
		return createError(5905, err)
	}
	return err
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	instanceName := "Test name"
	verboseLogging := sz.SZ_NO_LOGGING
	configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
	localLogger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages)
	if err != nil {
		panic(err)
	}

	settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		return createError(5902, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5920, err)
	}

	// Purge repository.

	err = setupPurgeRepository(ctx, instanceName, settings, verboseLogging, configId)
	if err != nil {
		return createError(5921, err)
	}
	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestBuildSimpleSystemConfigurationJsonUsingEnvVars(test *testing.T) {
	actual, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printResponse(test, actual)
}

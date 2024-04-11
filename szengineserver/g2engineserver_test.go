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
	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	loadId            = "G2Engine_test"
	printResults      = false
)

type GetEntityByRecordIdResponse struct {
	ResolvedEntity struct {
		EntityId int64 `json:"ENTITY_Id"`
	} `json:"RESOLVED_ENTITY"`
}

var (
	g2engineTestSingleton *SzEngineServer
	localLogger           logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(localLogger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) SzEngineServer {
	if g2engineTestSingleton == nil {
		g2engineTestSingleton = &SzEngineServer{}
		instanceName := "Test name"
		verboseLogging := sz.SZ_NO_LOGGING
		configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
		settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkG2engine().Initialize(ctx, instanceName, settings, verboseLogging, configId)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *g2engineTestSingleton
}

func getG2EngineServer(ctx context.Context) SzEngineServer {
	if g2engineTestSingleton == nil {
		g2engineTestSingleton = &SzEngineServer{}
		instanceName := "Test name"
		verboseLogging := sz.SZ_NO_LOGGING
		configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
		setting, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkG2engine().Initialize(ctx, instanceName, setting, verboseLogging, configId)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *g2engineTestSingleton
}

func getEntityIdForRecord(datasource string, id string) int64 {
	ctx := context.TODO()
	var result int64 = 0
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetEntityByRecordIdRequest{
		DataSourceCode: datasource,
		RecordId:       id,
	}
	response, err := g2engine.GetEntityByRecordId(ctx, request)
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

func getEntityIdStringForRecord(datasource string, id string) string {
	entityId := getEntityIdForRecord(datasource, id)
	return strconv.FormatInt(entityId, 10)
}

func getEntityId(record record.Record) int64 {
	return getEntityIdForRecord(record.DataSource, record.Id)
}

func getEntityIdString(record record.Record) string {
	entityId := getEntityId(record)
	return strconv.FormatInt(entityId, 10)
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printResponse(test *testing.T, response interface{}) {
	printResult(test, "Response", response)
}

func testError(test *testing.T, ctx context.Context, g2engine SzEngineServer, err error) {
	_ = ctx
	_ = g2engine
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2engine SzEngineServer, err error, messageId string) {
	_ = ctx
	_ = g2engine
	if err != nil {
		var dictionary map[string]interface{}
		unmarshalErr := json.Unmarshal([]byte(err.Error()), &dictionary)
		if unmarshalErr != nil {
			test.Log("Unmarshal Error:", unmarshalErr.Error())
		}
		assert.Equal(test, messageId, dictionary["id"].(string))
	} else {
		assert.FailNow(test, "Should have failed with", messageId)
	}
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
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := szConfig.AddDataSource(ctx, configHandle, datasource.Json)
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

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
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
	err := szDiagnostic.Initialize(ctx, instanceName, settings, verboseLogging, configId)
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

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2engineServer_AddRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	request1 := &g2pb.AddRecordRequest{
		DataSourceCode: record1.DataSource,
		RecordId:       record1.Id,
		JsonData:       record1.Json,
		LoadId:         loadId,
	}
	response1, err := g2engine.AddRecord(ctx, request1)
	testError(test, ctx, g2engine, err)
	printResponse(test, response1)
	request2 := &g2pb.AddRecordRequest{
		DataSourceCode: record2.DataSource,
		RecordId:       record2.Id,
		JsonData:       record2.Json,
		LoadId:         loadId,
	}
	response2, err := g2engine.AddRecord(ctx, request2)
	testError(test, ctx, g2engine, err)
	printResponse(test, response2)
}

func TestG2engineServer_AddRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	request := &g2pb.AddRecordWithInfoRequest{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
		JsonData:       record.Json,
		LoadId:         loadId,
	}
	response, err := g2engine.AddRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_CountRedoRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.CountRedoRecordsRequest{}
	response, err := g2engine.CountRedoRecords(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ExportJSONEntityReport(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	flags := int64(0)
	request := &g2pb.ExportJSONEntityReportRequest{
		Flags: flags,
	}
	response, err := g2engine.ExportJSONEntityReport(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ExportConfigAndConfigId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.ExportConfigAndConfigIdRequest{}
	response, err := g2engine.ExportConfigAndConfigId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.ExportConfigRequest{}
	response, err := g2engine.ExportConfig(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ExportCSVEntityReport(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.ExportCSVEntityReportRequest{
		CsvColumnList: "",
		Flags:         0,
	}
	response, err := g2engine.ExportCSVEntityReport(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindInterestingEntitiesByEntityId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := int64(0)
	request := &g2pb.FindInterestingEntitiesByEntityIdRequest{
		EntityId: entityId,
		Flags:    flags,
	}
	response, err := g2engine.FindInterestingEntitiesByEntityId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindInterestingEntitiesByRecordId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := int64(0)
	request := &g2pb.FindInterestingEntitiesByRecordIdRequest{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
		Flags:          flags,
	}
	response, err := g2engine.FindInterestingEntitiesByRecordId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByEntityId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_Id": ` + getEntityIdString(record1) + `}, {"ENTITY_Id": ` + getEntityIdString(record2) + `}]}`
	maxDegree := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	request := &g2pb.FindNetworkByEntityIdRequest{
		EntityList:     entityList,
		MaxDegree:      maxDegree,
		BuildOutDegree: buildOutDegree,
		MaxEntities:    maxEntities,
	}
	response, err := g2engine.FindNetworkByEntityId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByEntityId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_Id": ` + getEntityIdString(record1) + `}, {"ENTITY_Id": ` + getEntityIdString(record2) + `}]}`
	maxDegree := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := int64(0)
	request := &g2pb.FindNetworkByEntityId_V2Request{
		EntityList:     entityList,
		MaxDegree:      maxDegree,
		BuildOutDegree: buildOutDegree,
		MaxEntities:    maxEntities,
		Flags:          flags,
	}
	response, err := g2engine.FindNetworkByEntityId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByRecordId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_Id": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_Id": "` + record2.Id + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_Id": "` + record3.Id + `"}]}`
	maxDegree := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	request := &g2pb.FindNetworkByRecordIdRequest{
		RecordList:     recordList,
		MaxDegree:      maxDegree,
		BuildOutDegree: buildOutDegree,
		MaxEntities:    maxEntities,
	}
	response, err := g2engine.FindNetworkByRecordId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByRecordId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_Id": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_Id": "` + record2.Id + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_Id": "` + record3.Id + `"}]}`
	maxDegree := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := int64(0)
	request := &g2pb.FindNetworkByRecordId_V2Request{
		RecordList:     recordList,
		MaxDegree:      maxDegree,
		BuildOutDegree: buildOutDegree,
		MaxEntities:    maxEntities,
		Flags:          flags,
	}
	response, err := g2engine.FindNetworkByRecordId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByEntityId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId1 := getEntityId(truthset.CustomerRecords["1001"])
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	request := &g2pb.FindPathByEntityIdRequest{
		EntityId1: entityId1,
		EntityId2: entityId2,
		MaxDegree: maxDegree,
	}
	response, err := g2engine.FindPathByEntityId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByEntityId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId1 := getEntityId(truthset.CustomerRecords["1001"])
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	flags := int64(0)
	request := &g2pb.FindPathByEntityId_V2Request{
		EntityId1: entityId1,
		EntityId2: entityId2,
		MaxDegree: maxDegree,
		Flags:     flags,
	}
	response, err := g2engine.FindPathByEntityId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByRecordId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	request := &g2pb.FindPathByRecordIdRequest{
		DataSourceCode1: record1.DataSource,
		RecordId1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordId2:       record2.Id,
		MaxDegree:       maxDegree,
	}
	response, err := g2engine.FindPathByRecordId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByRecordId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	flags := int64(0)
	request := &g2pb.FindPathByRecordId_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordId1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordId2:       record2.Id,
		MaxDegree:       maxDegree,
		Flags:           flags,
	}
	response, err := g2engine.FindPathByRecordId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByEntityId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityId1 := getEntityId(record1)
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_Id": ` + getEntityIdString(record1) + `}]}`
	request := &g2pb.FindPathExcludingByEntityIdRequest{
		EntityId1:        entityId1,
		EntityId2:        entityId2,
		MaxDegree:        maxDegree,
		ExcludedEntities: excludedEntities,
	}
	response, err := g2engine.FindPathExcludingByEntityId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByEntityId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityId1 := getEntityId(record1)
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_Id": ` + getEntityIdString(record1) + `}]}`
	flags := int64(0)
	request := &g2pb.FindPathExcludingByEntityId_V2Request{
		EntityId1:        entityId1,
		EntityId2:        entityId2,
		MaxDegree:        maxDegree,
		ExcludedEntities: excludedEntities,
		Flags:            flags,
	}
	response, err := g2engine.FindPathExcludingByEntityId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByRecordId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_Id": "` + record1.Id + `"}]}`
	request := &g2pb.FindPathExcludingByRecordIdRequest{
		DataSourceCode1: record1.DataSource,
		RecordId1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordId2:       record2.Id,
		MaxDegree:       maxDegree,
		ExcludedRecords: excludedRecords,
	}
	response, err := g2engine.FindPathExcludingByRecordId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByRecordId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_Id": "` + record1.Id + `"}]}`
	flags := int64(0)
	request := &g2pb.FindPathExcludingByRecordId_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordId1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordId2:       record2.Id,
		MaxDegree:       maxDegree,
		ExcludedRecords: excludedRecords,
		Flags:           flags,
	}
	response, err := g2engine.FindPathExcludingByRecordId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByEntityId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityId1 := getEntityId(record1)
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_Id": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	request := &g2pb.FindPathIncludingSourceByEntityIdRequest{
		EntityId1:        entityId1,
		EntityId2:        entityId2,
		MaxDegree:        maxDegree,
		ExcludedEntities: excludedEntities,
		RequiredDsrcs:    requiredDsrcs,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByEntityId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityId1 := getEntityId(record1)
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_Id": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := int64(0)
	request := &g2pb.FindPathIncludingSourceByEntityId_V2Request{
		EntityId1:        entityId1,
		EntityId2:        entityId2,
		MaxDegree:        maxDegree,
		ExcludedEntities: excludedEntities,
		RequiredDsrcs:    requiredDsrcs,
		Flags:            flags,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByRecordId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_Id": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	request := &g2pb.FindPathIncludingSourceByRecordIdRequest{
		DataSourceCode1: record1.DataSource,
		RecordId1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordId2:       record1.Id,
		MaxDegree:       maxDegree,
		ExcludedRecords: excludedEntities,
		RequiredDsrcs:   requiredDsrcs,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByRecordId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_Id": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := int64(0)
	request := &g2pb.FindPathIncludingSourceByRecordId_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordId1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordId2:       record1.Id,
		MaxDegree:       maxDegree,
		ExcludedRecords: excludedEntities,
		RequiredDsrcs:   requiredDsrcs,
		Flags:           flags,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetActiveConfigId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.GetActiveConfigIdRequest{}
	response, err := g2engine.GetActiveConfigId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByEntityId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	request := &g2pb.GetEntityByEntityIdRequest{
		EntityId: entityId,
	}
	response, err := g2engine.GetEntityByEntityId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByEntityId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	request := &g2pb.GetEntityByEntityId_V2Request{
		EntityId: entityId,
		Flags:    0,
	}
	response, err := g2engine.GetEntityByEntityId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	request := &g2pb.GetEntityByRecordIdRequest{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
	}
	response, err := g2engine.GetEntityByRecordId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByRecordId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := int64(0)
	request := &g2pb.GetEntityByRecordId_V2Request{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
		Flags:          flags,
	}
	response, err := g2engine.GetEntityByRecordId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	request := &g2pb.GetRecordRequest{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
	}
	response, err := g2engine.GetRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetRecord_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := int64(0)
	request := &g2pb.GetRecord_V2Request{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
		Flags:          flags,
	}
	response, err := g2engine.GetRecord_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetRedoRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.GetRedoRecordRequest{}
	response, err := g2engine.GetRedoRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetRepositoryLastModifiedTime(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.GetRepositoryLastModifiedTimeRequest{}
	response, err := g2engine.GetRepositoryLastModifiedTime(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetVirtualEntityByRecordId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_Id": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_Id": "` + record2.Id + `"}]}`
	request := &g2pb.GetVirtualEntityByRecordIdRequest{
		RecordList: recordList,
	}
	response, err := g2engine.GetVirtualEntityByRecordId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetVirtualEntityByRecordId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_Id": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_Id": "` + record2.Id + `"}]}`
	flags := int64(0)
	request := &g2pb.GetVirtualEntityByRecordId_V2Request{
		RecordList: recordList,
		Flags:      flags,
	}
	response, err := g2engine.GetVirtualEntityByRecordId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_HowEntityByEntityId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	request := &g2pb.HowEntityByEntityIdRequest{
		EntityId: entityId,
	}
	response, err := g2engine.HowEntityByEntityId(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_HowEntityByEntityId_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	request := &g2pb.HowEntityByEntityId_V2Request{
		EntityId: entityId,
		Flags:    0,
	}
	response, err := g2engine.HowEntityByEntityId_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_PrimeEngine(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.PrimeEngineRequest{}
	response, err := g2engine.PrimeEngine(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := int64(0)
	request := &g2pb.ReevaluateEntityRequest{
		EntityId: entityId,
		Flags:    flags,
	}
	response, err := g2engine.ReevaluateEntity(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateEntityWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId := getEntityId(truthset.CustomerRecords["1001"])
	flags := int64(0)
	request := &g2pb.ReevaluateEntityWithInfoRequest{
		EntityId: entityId,
		Flags:    flags,
	}
	response, err := g2engine.ReevaluateEntityWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := int64(0)
	request := &g2pb.ReevaluateRecordRequest{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
		Flags:          flags,
	}
	response, err := g2engine.ReevaluateRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := int64(0)
	request := &g2pb.ReevaluateRecordWithInfoRequest{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
		Flags:          flags,
	}
	response, err := g2engine.ReevaluateRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

// FIXME: Remove after GDEV-3576 is fixed
// func TestG2engineServer_ReplaceRecord(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	request := &g2pb.ReplaceRecordRequest{
// 		DataSourceCode: "CUSTOMERS",
// 		RecordId:       "1001",
// 		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_Id": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
// 		LoadId:         "TEST",
// 	}
// 	response, err := g2engine.ReplaceRecord(ctx, request)
// 	testError(test, ctx, g2engine, err)
// 	printResponse(test, response)
// }

// FIXME: Remove after GDEV-3576 is fixed
// func TestG2engineServer_ReplaceRecordWithInfo(test *testing.T) {
// 	ctx := context.TODO()
// 	g2engine := getTestObject(ctx, test)
// 	request := &g2pb.ReplaceRecordWithInfoRequest{
// 		DataSourceCode: "CUSTOMERS",
// 		RecordId:       "1001",
// 		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_Id": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
// 		LoadId:         "TEST",
// 		Flags:          0,
// 	}
// 	response, err := g2engine.ReplaceRecordWithInfo(ctx, request)
// 	testError(test, ctx, g2engine, err)
// 	printResponse(test, response)
// }

func TestG2engineServer_SearchByAttributes(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	request := &g2pb.SearchByAttributesRequest{
		JsonData: jsonData,
	}
	response, err := g2engine.SearchByAttributes(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_SearchByAttributes_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	flags := int64(0)
	request := &g2pb.SearchByAttributes_V2Request{
		JsonData: jsonData,
		Flags:    flags,
	}
	response, err := g2engine.SearchByAttributes_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_Stats(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.StatsRequest{}
	response, err := g2engine.Stats(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntities(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId1 := getEntityId(truthset.CustomerRecords["1001"])
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	request := &g2pb.WhyEntitiesRequest{
		EntityId1: entityId1,
		EntityId2: entityId2,
	}
	response, err := g2engine.WhyEntities(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntities_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityId1 := getEntityId(truthset.CustomerRecords["1001"])
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	flags := int64(0)
	request := &g2pb.WhyEntities_V2Request{
		EntityId1: entityId1,
		EntityId2: entityId2,
		Flags:     flags,
	}
	response, err := g2engine.WhyEntities_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	request := &g2pb.WhyRecordsRequest{
		DataSourceCode1: record1.DataSource,
		RecordId1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordId2:       record2.Id,
	}
	response, err := g2engine.WhyRecords(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyRecords_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := int64(0)
	request := &g2pb.WhyRecords_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordId1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordId2:       record2.Id,
		Flags:           flags,
	}
	response, err := g2engine.WhyRecords_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_Init(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	moduleName := "Test module name"
	verboseLogging := int64(0) // 0 for no Senzing logging; 1 for logging
	iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: verboseLogging,
	}
	response, err := g2engine.Init(ctx, request)
	expectError(test, ctx, g2engine, err, "senzing-60144002")
	printResponse(test, response)
}

func TestG2engineServer_InitWithConfigId(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	moduleName := "Test module name"
	var initConfigId int64 = 1
	verboseLogging := int64(0) // 0 for no Senzing logging; 1 for logging
	iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitWithConfigIdRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		InitConfigId:   initConfigId,
		VerboseLogging: verboseLogging,
	}
	response, err := g2engine.InitWithConfigId(ctx, request)
	expectError(test, ctx, g2engine, err, "senzing-60144003")
	printResponse(test, response)
}

func TestG2engineServer_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)

	requestToGetActiveConfigId := &g2pb.GetActiveConfigIdRequest{}
	responseFromGetActiveConfigId, err := g2engine.GetActiveConfigId(ctx, requestToGetActiveConfigId)
	testError(test, ctx, g2engine, err)

	request := &g2pb.ReinitRequest{
		InitConfigId: responseFromGetActiveConfigId.GetResult(),
	}
	response, err := g2engine.Reinit(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_DeleteRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	request := &g2pb.DeleteRecordRequest{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
		LoadId:         loadId,
	}
	response, err := g2engine.DeleteRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_DeleteRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	flags := int64(0)
	request := &g2pb.DeleteRecordWithInfoRequest{
		DataSourceCode: record.DataSource,
		RecordId:       record.Id,
		LoadId:         loadId,
		Flags:          flags,
	}
	response, err := g2engine.DeleteRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.DestroyRequest{}
	response, err := g2engine.Destroy(ctx, request)
	expectError(test, ctx, g2engine, err, "senzing-60144001")
	printResponse(test, response)
	g2engineTestSingleton = nil
}

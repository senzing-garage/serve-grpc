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
	"github.com/senzing-garage/sz-sdk-go/szerror"
	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	loadId            = "G2Engine_test"
	printResults      = false
)

type GetEntityByRecordIDResponse struct {
	ResolvedEntity struct {
		EntityId int64 `json:"ENTITY_ID"`
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
		moduleName := "Test module name"
		verboseLogging := int64(0)
		iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkG2engine().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *g2engineTestSingleton
}

func getG2EngineServer(ctx context.Context) SzEngineServer {
	if g2engineTestSingleton == nil {
		g2engineTestSingleton = &SzEngineServer{}
		moduleName := "Test module name"
		verboseLogging := int64(0)
		iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkG2engine().Init(ctx, moduleName, iniParams, verboseLogging)
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
	request := &g2pb.GetEntityByRecordIDRequest{
		DataSourceCode: datasource,
		RecordID:       id,
	}
	response, err := g2engine.GetEntityByRecordID(ctx, request)
	if err != nil {
		return result
	}

	getEntityByRecordIDResponse := &GetEntityByRecordIDResponse{}
	err = json.Unmarshal([]byte(response.Result), &getEntityByRecordIDResponse)
	if err != nil {
		return result
	}
	return getEntityByRecordIDResponse.ResolvedEntity.EntityId
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

	aG2config := &szconfig.Szconfig{}
	err := aG2config.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}

	configHandle, err := aG2config.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := aG2config.AddDataSource(ctx, configHandle, datasource.Json)
		if err != nil {
			return createError(5908, err)
		}
	}

	configStr, err := aG2config.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = aG2config.CloseConfig(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return createError(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &szconfigmanager.Szconfigmanager{}
	err = aG2configmgr.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return createError(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigId(ctx, configID)
	if err != nil {
		return createError(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	aG2diagnostic := &szdiagnostic.Szdiagnostic{}
	err := aG2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5903, err)
	}

	err = aG2diagnostic.PurgeRepository(ctx)
	if err != nil {
		return createError(5904, err)
	}

	err = aG2diagnostic.Destroy(ctx)
	if err != nil {
		return createError(5905, err)
	}
	return err
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	moduleName := "Test module name"
	verboseLogging := int64(0)
	localLogger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages)
	if err != nil {
		panic(err)
	}

	iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		return createError(5902, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5920, err)
	}

	// Purge repository.

	err = setupPurgeRepository(ctx, moduleName, iniParams, verboseLogging)
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
		RecordID:       record1.Id,
		JsonData:       record1.Json,
		LoadID:         loadId,
	}
	response1, err := g2engine.AddRecord(ctx, request1)
	testError(test, ctx, g2engine, err)
	printResponse(test, response1)
	request2 := &g2pb.AddRecordRequest{
		DataSourceCode: record2.DataSource,
		RecordID:       record2.Id,
		JsonData:       record2.Json,
		LoadID:         loadId,
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
		RecordID:       record.Id,
		JsonData:       record.Json,
		LoadID:         loadId,
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

func TestG2engineServer_ExportConfigAndConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.ExportConfigAndConfigIDRequest{}
	response, err := g2engine.ExportConfigAndConfigID(ctx, request)
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

func TestG2engineServer_FindInterestingEntitiesByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	flags := int64(0)
	request := &g2pb.FindInterestingEntitiesByEntityIDRequest{
		EntityID: entityID,
		Flags:    flags,
	}
	response, err := g2engine.FindInterestingEntitiesByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindInterestingEntitiesByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := int64(0)
	request := &g2pb.FindInterestingEntitiesByRecordIDRequest{
		DataSourceCode: record.DataSource,
		RecordID:       record.Id,
		Flags:          flags,
	}
	response, err := g2engine.FindInterestingEntitiesByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}, {"ENTITY_ID": ` + getEntityIdString(record2) + `}]}`
	maxDegree := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	request := &g2pb.FindNetworkByEntityIDRequest{
		EntityList:     entityList,
		MaxDegree:      maxDegree,
		BuildOutDegree: buildOutDegree,
		MaxEntities:    maxEntities,
	}
	response, err := g2engine.FindNetworkByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}, {"ENTITY_ID": ` + getEntityIdString(record2) + `}]}`
	maxDegree := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := int64(0)
	request := &g2pb.FindNetworkByEntityID_V2Request{
		EntityList:     entityList,
		MaxDegree:      maxDegree,
		BuildOutDegree: buildOutDegree,
		MaxEntities:    maxEntities,
		Flags:          flags,
	}
	response, err := g2engine.FindNetworkByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.Id + `"}]}`
	maxDegree := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	request := &g2pb.FindNetworkByRecordIDRequest{
		RecordList:     recordList,
		MaxDegree:      maxDegree,
		BuildOutDegree: buildOutDegree,
		MaxEntities:    maxEntities,
	}
	response, err := g2engine.FindNetworkByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.Id + `"}]}`
	maxDegree := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := int64(0)
	request := &g2pb.FindNetworkByRecordID_V2Request{
		RecordList:     recordList,
		MaxDegree:      maxDegree,
		BuildOutDegree: buildOutDegree,
		MaxEntities:    maxEntities,
		Flags:          flags,
	}
	response, err := g2engine.FindNetworkByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	request := &g2pb.FindPathByEntityIDRequest{
		EntityID1: entityID1,
		EntityID2: entityID2,
		MaxDegree: maxDegree,
	}
	response, err := g2engine.FindPathByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	flags := int64(0)
	request := &g2pb.FindPathByEntityID_V2Request{
		EntityID1: entityID1,
		EntityID2: entityID2,
		MaxDegree: maxDegree,
		Flags:     flags,
	}
	response, err := g2engine.FindPathByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	request := &g2pb.FindPathByRecordIDRequest{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
		MaxDegree:       maxDegree,
	}
	response, err := g2engine.FindPathByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	flags := int64(0)
	request := &g2pb.FindPathByRecordID_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
		MaxDegree:       maxDegree,
		Flags:           flags,
	}
	response, err := g2engine.FindPathByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityID1 := getEntityId(record1)
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	request := &g2pb.FindPathExcludingByEntityIDRequest{
		EntityID1:        entityID1,
		EntityID2:        entityID2,
		MaxDegree:        maxDegree,
		ExcludedEntities: excludedEntities,
	}
	response, err := g2engine.FindPathExcludingByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityID1 := getEntityId(record1)
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	flags := int64(0)
	request := &g2pb.FindPathExcludingByEntityID_V2Request{
		EntityID1:        entityID1,
		EntityID2:        entityID2,
		MaxDegree:        maxDegree,
		ExcludedEntities: excludedEntities,
		Flags:            flags,
	}
	response, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	request := &g2pb.FindPathExcludingByRecordIDRequest{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
		MaxDegree:       maxDegree,
		ExcludedRecords: excludedRecords,
	}
	response, err := g2engine.FindPathExcludingByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	flags := int64(0)
	request := &g2pb.FindPathExcludingByRecordID_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
		MaxDegree:       maxDegree,
		ExcludedRecords: excludedRecords,
		Flags:           flags,
	}
	response, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityID1 := getEntityId(record1)
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	request := &g2pb.FindPathIncludingSourceByEntityIDRequest{
		EntityID1:        entityID1,
		EntityID2:        entityID2,
		MaxDegree:        maxDegree,
		ExcludedEntities: excludedEntities,
		RequiredDsrcs:    requiredDsrcs,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	entityID1 := getEntityId(record1)
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := int64(0)
	request := &g2pb.FindPathIncludingSourceByEntityID_V2Request{
		EntityID1:        entityID1,
		EntityID2:        entityID2,
		MaxDegree:        maxDegree,
		ExcludedEntities: excludedEntities,
		RequiredDsrcs:    requiredDsrcs,
		Flags:            flags,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	request := &g2pb.FindPathIncludingSourceByRecordIDRequest{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record1.Id,
		MaxDegree:       maxDegree,
		ExcludedRecords: excludedEntities,
		RequiredDsrcs:   requiredDsrcs,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := int64(0)
	request := &g2pb.FindPathIncludingSourceByRecordID_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record1.Id,
		MaxDegree:       maxDegree,
		ExcludedRecords: excludedEntities,
		RequiredDsrcs:   requiredDsrcs,
		Flags:           flags,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetActiveConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &g2pb.GetActiveConfigIDRequest{}
	response, err := g2engine.GetActiveConfigID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	request := &g2pb.GetEntityByEntityIDRequest{
		EntityID: entityID,
	}
	response, err := g2engine.GetEntityByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	request := &g2pb.GetEntityByEntityID_V2Request{
		EntityID: entityID,
		Flags:    0,
	}
	response, err := g2engine.GetEntityByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	request := &g2pb.GetEntityByRecordIDRequest{
		DataSourceCode: record.DataSource,
		RecordID:       record.Id,
	}
	response, err := g2engine.GetEntityByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := int64(0)
	request := &g2pb.GetEntityByRecordID_V2Request{
		DataSourceCode: record.DataSource,
		RecordID:       record.Id,
		Flags:          flags,
	}
	response, err := g2engine.GetEntityByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	request := &g2pb.GetRecordRequest{
		DataSourceCode: record.DataSource,
		RecordID:       record.Id,
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
		RecordID:       record.Id,
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

func TestG2engineServer_GetVirtualEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}]}`
	request := &g2pb.GetVirtualEntityByRecordIDRequest{
		RecordList: recordList,
	}
	response, err := g2engine.GetVirtualEntityByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetVirtualEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.Id + `"}]}`
	flags := int64(0)
	request := &g2pb.GetVirtualEntityByRecordID_V2Request{
		RecordList: recordList,
		Flags:      flags,
	}
	response, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_HowEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	request := &g2pb.HowEntityByEntityIDRequest{
		EntityID: entityID,
	}
	response, err := g2engine.HowEntityByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_HowEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	request := &g2pb.HowEntityByEntityID_V2Request{
		EntityID: entityID,
		Flags:    0,
	}
	response, err := g2engine.HowEntityByEntityID_V2(ctx, request)
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
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	flags := int64(0)
	request := &g2pb.ReevaluateEntityRequest{
		EntityID: entityID,
		Flags:    flags,
	}
	response, err := g2engine.ReevaluateEntity(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateEntityWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	flags := int64(0)
	request := &g2pb.ReevaluateEntityWithInfoRequest{
		EntityID: entityID,
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
		RecordID:       record.Id,
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
		RecordID:       record.Id,
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
// 		RecordID:       "1001",
// 		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
// 		LoadID:         "TEST",
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
// 		RecordID:       "1001",
// 		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
// 		LoadID:         "TEST",
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
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	request := &g2pb.WhyEntitiesRequest{
		EntityID1: entityID1,
		EntityID2: entityID2,
	}
	response, err := g2engine.WhyEntities(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntities_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	flags := int64(0)
	request := &g2pb.WhyEntities_V2Request{
		EntityID1: entityID1,
		EntityID2: entityID2,
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
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
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
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
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

func TestG2engineServer_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	moduleName := "Test module name"
	var initConfigID int64 = 1
	verboseLogging := int64(0) // 0 for no Senzing logging; 1 for logging
	iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitWithConfigIDRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		InitConfigID:   initConfigID,
		VerboseLogging: verboseLogging,
	}
	response, err := g2engine.InitWithConfigID(ctx, request)
	expectError(test, ctx, g2engine, err, "senzing-60144003")
	printResponse(test, response)
}

func TestG2engineServer_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)

	requestToGetActiveConfigID := &g2pb.GetActiveConfigIDRequest{}
	responseFromGetActiveConfigID, err := g2engine.GetActiveConfigID(ctx, requestToGetActiveConfigID)
	testError(test, ctx, g2engine, err)

	request := &g2pb.ReinitRequest{
		InitConfigID: responseFromGetActiveConfigID.GetResult(),
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
		RecordID:       record.Id,
		LoadID:         loadId,
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
		RecordID:       record.Id,
		LoadID:         loadId,
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

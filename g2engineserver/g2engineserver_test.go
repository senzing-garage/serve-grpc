package g2engineserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go/g2config"
	"github.com/senzing/g2-sdk-go/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2engine"
	pb "github.com/senzing/g2-sdk-proto/go/g2engine"
	"github.com/senzing/go-common/record"
	"github.com/senzing/go-common/truthset"
	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	"github.com/senzing/go-logging/messagelogger"
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
	g2engineTestSingleton *G2EngineServer
	localLogger           messagelogger.MessageLoggerInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) G2EngineServer {
	if g2engineTestSingleton == nil {
		g2engineTestSingleton = &G2EngineServer{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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

func getG2EngineServer(ctx context.Context) G2EngineServer {
	if g2engineTestSingleton == nil {
		g2engineTestSingleton = &G2EngineServer{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
	request := &pb.GetEntityByRecordIDRequest{
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

func testError(test *testing.T, ctx context.Context, g2engine G2EngineServer, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2engine G2EngineServer, err error, messageId string) {
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

func testErrorNoFail(test *testing.T, ctx context.Context, g2engine G2EngineServer, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
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

func setupSenzingConfig(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	now := time.Now()

	aG2config := &g2config.G2configImpl{}
	err := aG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
	if err != nil {
		return localLogger.Error(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := aG2config.AddDataSource(ctx, configHandle, datasource.Json)
		if err != nil {
			return localLogger.Error(5908, err)
		}
	}

	configStr, err := aG2config.Save(ctx, configHandle)
	if err != nil {
		return localLogger.Error(5909, err)
	}

	err = aG2config.Close(ctx, configHandle)
	if err != nil {
		return localLogger.Error(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return localLogger.Error(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &g2configmgr.G2configmgrImpl{}
	err = aG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return localLogger.Error(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return localLogger.Error(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return localLogger.Error(5915, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	aG2engine := &g2engine.G2engineImpl{}
	err := aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5903, err)
	}

	err = aG2engine.PurgeRepository(ctx)
	if err != nil {
		return localLogger.Error(5904, err)
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return localLogger.Error(5905, err)
	}
	return err
}

func setup() error {
	ctx := context.TODO()
	var err error = nil
	moduleName := "Test module name"
	verboseLogging := 0
	localLogger, err = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	if err != nil {
		return localLogger.Error(5901, err)
	}

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		return localLogger.Error(5902, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5920, err)
	}

	// Purge repository.

	err = setupPurgeRepository(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5921, err)
	}
	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestBuildSimpleSystemConfigurationJson(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
	request1 := &pb.AddRecordRequest{
		DataSourceCode: record1.DataSource,
		RecordID:       record1.Id,
		JsonData:       record1.Json,
		LoadID:         loadId,
	}
	response1, err := g2engine.AddRecord(ctx, request1)
	testError(test, ctx, g2engine, err)
	printResponse(test, response1)
	request2 := &pb.AddRecordRequest{
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
	request := &pb.AddRecordWithInfoRequest{
		DataSourceCode: record.DataSource,
		RecordID:       record.Id,
		JsonData:       record.Json,
		LoadID:         loadId,
	}
	response, err := g2engine.AddRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_AddRecordWithInfoWithReturnedRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.TestRecordsWithoutRecordId[0]
	request := &pb.AddRecordWithInfoWithReturnedRecordIDRequest{
		DataSourceCode: record.DataSource,
		JsonData:       record.Json,
		LoadID:         loadId,
	}
	response, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_AddRecordWithReturnedRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.TestRecordsWithoutRecordId[1]
	request := &pb.AddRecordWithReturnedRecordIDRequest{
		DataSourceCode: record.DataSource,
		JsonData:       record.Json,
		LoadID:         loadId,
	}
	response, err := g2engine.AddRecordWithReturnedRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_CheckRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	recordQueryList := `{"RECORDS": [{"DATA_SOURCE": "` + record.DataSource + `","RECORD_ID": "` + record.Id + `"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "123456789"}]}`
	request := &pb.CheckRecordRequest{
		Record:          record.Json,
		RecordQueryList: recordQueryList,
	}
	response, err := g2engine.CheckRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_CountRedoRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.CountRedoRecordsRequest{}
	response, err := g2engine.CountRedoRecords(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ExportJSONEntityReport(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	flags := int64(0)
	request := &pb.ExportJSONEntityReportRequest{
		Flags: flags,
	}
	response, err := g2engine.ExportJSONEntityReport(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ExportConfigAndConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ExportConfigAndConfigIDRequest{}
	response, err := g2engine.ExportConfigAndConfigID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ExportConfigRequest{}
	response, err := g2engine.ExportConfig(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ExportCSVEntityReport(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ExportCSVEntityReportRequest{
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
	var flags int64 = 0
	request := &pb.FindInterestingEntitiesByEntityIDRequest{
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
	var flags int64 = 0
	request := &pb.FindInterestingEntitiesByRecordIDRequest{
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
	maxDegree := 2
	buildOutDegree := 1
	maxEntities := 10
	request := &pb.FindNetworkByEntityIDRequest{
		EntityList:     entityList,
		MaxDegree:      int32(maxDegree),
		BuildOutDegree: int32(buildOutDegree),
		MaxEntities:    int32(maxEntities),
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
	maxDegree := 2
	buildOutDegree := 1
	maxEntities := 10
	var flags int64 = 0
	request := &pb.FindNetworkByEntityID_V2Request{
		EntityList:     entityList,
		MaxDegree:      int32(maxDegree),
		BuildOutDegree: int32(buildOutDegree),
		MaxEntities:    int32(maxEntities),
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
	maxDegree := 1
	buildOutDegree := 2
	maxEntities := 10
	request := &pb.FindNetworkByRecordIDRequest{
		RecordList:     recordList,
		MaxDegree:      int32(maxDegree),
		BuildOutDegree: int32(buildOutDegree),
		MaxEntities:    int32(maxEntities),
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
	maxDegree := 1
	buildOutDegree := 2
	maxEntities := 10
	var flags int64 = 0
	request := &pb.FindNetworkByRecordID_V2Request{
		RecordList:     recordList,
		MaxDegree:      int32(maxDegree),
		BuildOutDegree: int32(buildOutDegree),
		MaxEntities:    int32(maxEntities),
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
	maxDegree := 1
	request := &pb.FindPathByEntityIDRequest{
		EntityID1: entityID1,
		EntityID2: entityID2,
		MaxDegree: int32(maxDegree),
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
	maxDegree := 1
	var flags int64 = 0
	request := &pb.FindPathByEntityID_V2Request{
		EntityID1: entityID1,
		EntityID2: entityID2,
		MaxDegree: int32(maxDegree),
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
	maxDegree := 1
	request := &pb.FindPathByRecordIDRequest{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
		MaxDegree:       int32(maxDegree),
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
	maxDegree := 1
	var flags int64 = 0
	request := &pb.FindPathByRecordID_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
		MaxDegree:       int32(maxDegree),
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
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	request := &pb.FindPathExcludingByEntityIDRequest{
		EntityID1:        entityID1,
		EntityID2:        entityID2,
		MaxDegree:        int32(maxDegree),
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
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	var flags int64 = 0
	request := &pb.FindPathExcludingByEntityID_V2Request{
		EntityID1:        entityID1,
		EntityID2:        entityID2,
		MaxDegree:        int32(maxDegree),
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
	maxDegree := 1
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	request := &pb.FindPathExcludingByRecordIDRequest{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
		MaxDegree:       int32(maxDegree),
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
	maxDegree := 1
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.Id + `"}]}`
	var flags int64 = 0
	request := &pb.FindPathExcludingByRecordID_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record2.Id,
		MaxDegree:       int32(maxDegree),
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
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	request := &pb.FindPathIncludingSourceByEntityIDRequest{
		EntityID1:        entityID1,
		EntityID2:        entityID2,
		MaxDegree:        int32(maxDegree),
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
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	var flags int64 = 0
	request := &pb.FindPathIncludingSourceByEntityID_V2Request{
		EntityID1:        entityID1,
		EntityID2:        entityID2,
		MaxDegree:        int32(maxDegree),
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
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	request := &pb.FindPathIncludingSourceByRecordIDRequest{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record1.Id,
		MaxDegree:       int32(maxDegree),
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
	maxDegree := 1
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdString(record1) + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	var flags int64 = 0
	request := &pb.FindPathIncludingSourceByRecordID_V2Request{
		DataSourceCode1: record1.DataSource,
		RecordID1:       record1.Id,
		DataSourceCode2: record2.DataSource,
		RecordID2:       record1.Id,
		MaxDegree:       int32(maxDegree),
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
	request := &pb.GetActiveConfigIDRequest{}
	response, err := g2engine.GetActiveConfigID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	request := &pb.GetEntityByEntityIDRequest{
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
	request := &pb.GetEntityByEntityID_V2Request{
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
	request := &pb.GetEntityByRecordIDRequest{
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
	var flags int64 = 0
	request := &pb.GetEntityByRecordID_V2Request{
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
	request := &pb.GetRecordRequest{
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
	var flags int64 = 0
	request := &pb.GetRecord_V2Request{
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
	request := &pb.GetRedoRecordRequest{}
	response, err := g2engine.GetRedoRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetRepositoryLastModifiedTime(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetRepositoryLastModifiedTimeRequest{}
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
	request := &pb.GetVirtualEntityByRecordIDRequest{
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
	var flags int64 = 0
	request := &pb.GetVirtualEntityByRecordID_V2Request{
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
	request := &pb.HowEntityByEntityIDRequest{
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
	request := &pb.HowEntityByEntityID_V2Request{
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
	request := &pb.PrimeEngineRequest{}
	response, err := g2engine.PrimeEngine(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_Process(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	request := &pb.ProcessRequest{
		Record: record.Json,
	}
	response, err := g2engine.Process(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ProcessRedoRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessRedoRecordRequest{}
	response, err := g2engine.ProcessRedoRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ProcessRedoRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	var flags int64 = 0
	request := &pb.ProcessRedoRecordWithInfoRequest{
		Flags: flags,
	}
	response, err := g2engine.ProcessRedoRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ProcessWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	var flags int64 = 0
	request := &pb.ProcessWithInfoRequest{
		Record: record.Json,
		Flags:  flags,
	}
	response, err := g2engine.ProcessWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ProcessWithResponse(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	request := &pb.ProcessWithResponseRequest{
		Record: record.Json,
	}
	response, err := g2engine.ProcessWithResponse(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ProcessWithResponseResize(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	request := &pb.ProcessWithResponseResizeRequest{
		Record: record.Json,
	}
	response, err := g2engine.ProcessWithResponseResize(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	var flags int64 = 0
	request := &pb.ReevaluateEntityRequest{
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
	var flags int64 = 0
	request := &pb.ReevaluateEntityWithInfoRequest{
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
	var flags int64 = 0
	request := &pb.ReevaluateRecordRequest{
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
	var flags int64 = 0
	request := &pb.ReevaluateRecordWithInfoRequest{
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
// 	request := &pb.ReplaceRecordRequest{
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
// 	request := &pb.ReplaceRecordWithInfoRequest{
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
	request := &pb.SearchByAttributesRequest{
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
	var flags int64 = 0
	request := &pb.SearchByAttributes_V2Request{
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
	request := &pb.StatsRequest{}
	response, err := g2engine.Stats(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntities(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	request := &pb.WhyEntitiesRequest{
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
	var flags int64 = 0
	request := &pb.WhyEntities_V2Request{
		EntityID1: entityID1,
		EntityID2: entityID2,
		Flags:     flags,
	}
	response, err := g2engine.WhyEntities_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	request := &pb.WhyEntityByEntityIDRequest{
		EntityID: entityID,
	}
	response, err := g2engine.WhyEntityByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	entityID := getEntityId(truthset.CustomerRecords["1001"])
	var flags int64 = 0
	request := &pb.WhyEntityByEntityID_V2Request{
		EntityID: entityID,
		Flags:    flags,
	}
	response, err := g2engine.WhyEntityByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	request := &pb.WhyEntityByRecordIDRequest{
		DataSourceCode: record.DataSource,
		RecordID:       record.Id,
	}
	response, err := g2engine.WhyEntityByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	var flags int64 = 0
	request := &pb.WhyEntityByRecordID_V2Request{
		DataSourceCode: record.DataSource,
		RecordID:       record.Id,
		Flags:          flags,
	}
	response, err := g2engine.WhyEntityByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	request := &pb.WhyRecordsRequest{
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
	var flags int64 = 0
	request := &pb.WhyRecords_V2Request{
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
	verboseLogging := 0 // 0 for no Senzing logging; 1 for logging
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
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
	verboseLogging := 0 // 0 for no Senzing logging; 1 for logging
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &pb.InitWithConfigIDRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		InitConfigID:   initConfigID,
		VerboseLogging: int32(verboseLogging),
	}
	response, err := g2engine.InitWithConfigID(ctx, request)
	expectError(test, ctx, g2engine, err, "senzing-60144003")
	printResponse(test, response)
}

func TestG2engineServer_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)

	requestToGetActiveConfigID := &pb.GetActiveConfigIDRequest{}
	responseFromGetActiveConfigID, err := g2engine.GetActiveConfigID(ctx, requestToGetActiveConfigID)
	testError(test, ctx, g2engine, err)

	request := &pb.ReinitRequest{
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
	request := &pb.DeleteRecordRequest{
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
	var flags int64 = 0
	request := &pb.DeleteRecordWithInfoRequest{
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
	request := &pb.DestroyRequest{}
	response, err := g2engine.Destroy(ctx, request)
	expectError(test, ctx, g2engine, err, "senzing-60144001")
	printResponse(test, response)
	g2engineTestSingleton = nil
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2EngineServer_AddRecord() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_AddRecord_secondRecord() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1002",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`,
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_AddRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordWithInfoRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1003",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`,
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.AddRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_AddRecordWithInfoWithReturnedRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordWithInfoWithReturnedRecordIDRequest{
		DataSourceCode: "CUSTOMERS",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Kellar", "PRIMARY_NAME_FIRST": "Candace", "ADDR_LINE1": "1824 AspenOak Way", "ADDR_CITY": "Elmwood Park", "ADDR_STATE": "CA", "ADDR_POSTAL_CODE": "95865", "EMAIL_ADDRESS": "info@ca-state.gov"}`,
		LoadID:         "G2Engine_test",
		Flags:          0,
	}
	response, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetWithInfo(), 42))
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2EngineServer_AddRecordWithReturnedRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordWithReturnedRecordIDRequest{
		DataSourceCode: "CUSTOMERS",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Sanders", "PRIMARY_NAME_FIRST": "Sandy", "ADDR_LINE1": "1376 BlueBell Rd", "ADDR_CITY": "Sacramento", "ADDR_STATE": "CA", "ADDR_POSTAL_CODE": "95823", "EMAIL_ADDRESS": "info@ca-state.gov"}`,
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.AddRecordWithReturnedRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Length of record identifier is %d hexadecimal characters.\n", len(response.GetResult()))
	// Output: Length of record identifier is 40 hexadecimal characters.
}

func ExampleG2EngineServer_CheckRecord() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.CheckRecordRequest{
		Record:          `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		RecordQueryList: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "123456789"}]}`,
	}
	response, err := g2engine.CheckRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"CHECK_RECORD_RESPONSE":[{"DSRC_CODE":"CUSTOMERS","RECORD_ID":"1001","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","MATCH_KEY":"","ERRULE_CODE":"","ERRULE_ID":0,"CANDIDATE_MATCH":"N","NON_GENERIC_CANDIDATE_MATCH":"N"}]}
}

func ExampleG2EngineServer_CloseExport() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)

	// Create a handle for the example.
	requestToExportJSONEntityReport := &pb.ExportJSONEntityReportRequest{
		Flags: 0,
	}
	responseFromExportJSONEntityReport, err := g2engine.ExportJSONEntityReport(ctx, requestToExportJSONEntityReport)

	// Example
	request := &pb.CloseExportRequest{
		ResponseHandle: responseFromExportJSONEntityReport.GetResult(),
	}
	response, err := g2engine.CloseExport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_CountRedoRecords() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.CountRedoRecordsRequest{}
	response, err := g2engine.CountRedoRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: 1
}

func ExampleG2EngineServer_ExportCSVEntityReport() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ExportCSVEntityReportRequest{
		CsvColumnList: "",
		Flags:         0,
	}
	response, err := g2engine.ExportCSVEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_ExportConfig() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ExportConfigRequest{}
	response, err := g2engine.ExportConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 42))
	// Output: {"G2_CONFIG":{"CFG_ETYPE":[{"ETYPE_ID":...
}

func ExampleG2EngineServer_ExportConfigAndConfigID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ExportConfigAndConfigIDRequest{}
	response, err := g2engine.ExportConfigAndConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetConfigID() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_ExportJSONEntityReport() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ExportJSONEntityReportRequest{
		Flags: 0,
	}
	response, err := g2engine.ExportJSONEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_FetchNext() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)

	// Create a handle for the example.
	requestToExportJSONEntityReport := &pb.ExportJSONEntityReportRequest{
		Flags: 0,
	}
	responseFromExportJSONEntityReport, err := g2engine.ExportJSONEntityReport(ctx, requestToExportJSONEntityReport)

	// Example
	request := &pb.FetchNextRequest{
		ResponseHandle: responseFromExportJSONEntityReport.GetResult(),
	}
	response, err := g2engine.FetchNext(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(response.GetResult()) >= 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_FindInterestingEntitiesByEntityID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindInterestingEntitiesByEntityIDRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.FindInterestingEntitiesByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_FindInterestingEntitiesByRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindInterestingEntitiesByRecordIDRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.FindInterestingEntitiesByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_FindNetworkByEntityID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	request := &pb.FindNetworkByEntityIDRequest{
		EntityList:     entityList,
		MaxDegree:      2,
		BuildOutDegree: 1,
		MaxEntities:    10,
	}
	response, err := g2engine.FindNetworkByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleG2EngineServer_FindNetworkByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	request := &pb.FindNetworkByEntityID_V2Request{
		EntityList:     entityList,
		MaxDegree:      2,
		BuildOutDegree: 1,
		MaxEntities:    10,
		Flags:          0,
	}
	response, err := g2engine.FindNetworkByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindNetworkByRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByRecordIDRequest{
		RecordList:     `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`,
		MaxDegree:      1,
		BuildOutDegree: 2,
		MaxEntities:    10,
	}
	response, err := g2engine.FindNetworkByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleG2EngineServer_FindNetworkByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByRecordID_V2Request{
		RecordList:     `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`,
		MaxDegree:      1,
		BuildOutDegree: 2,
		MaxEntities:    10,
		Flags:          0,
	}
	response, err := g2engine.FindNetworkByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathByEntityID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByEntityIDRequest{
		EntityID1: getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2: getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree: 1,
	}
	response, err := g2engine.FindPathByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByEntityID_V2Request{
		EntityID1: getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2: getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree: 1,
		Flags:     0,
	}
	response, err := g2engine.FindPathByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathByRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByRecordIDRequest{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
	}
	response, err := g2engine.FindPathByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 87))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":...
}

func ExampleG2EngineServer_FindPathByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByRecordID_V2Request{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		Flags:           0,
	}
	response, err := g2engine.FindPathByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathExcludingByEntityID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &pb.FindPathExcludingByEntityIDRequest{
		EntityID1:        getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2:        getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree:        1,
		ExcludedEntities: excludedEntities,
	}
	response, err := g2engine.FindPathExcludingByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathExcludingByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &pb.FindPathExcludingByEntityID_V2Request{
		EntityID1:        getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2:        getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree:        1,
		ExcludedEntities: excludedEntities,
		Flags:            0,
	}
	response, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathExcludingByRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByRecordIDRequest{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		ExcludedRecords: `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`,
	}
	response, err := g2engine.FindPathExcludingByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathExcludingByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByRecordID_V2Request{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		ExcludedRecords: `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`,
		Flags:           0,
	}
	response, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathIncludingSourceByEntityID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &pb.FindPathIncludingSourceByEntityIDRequest{
		EntityID1:        getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2:        getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree:        1,
		ExcludedEntities: excludedEntities,
		RequiredDsrcs:    `{"DATA_SOURCES": ["CUSTOMERS"]}`,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 106))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathIncludingSourceByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &pb.FindPathIncludingSourceByEntityID_V2Request{
		EntityID1:        getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2:        getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree:        1,
		ExcludedEntities: excludedEntities,
		RequiredDsrcs:    `{"DATA_SOURCES": ["CUSTOMERS"]}`,
		Flags:            0,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathIncludingSourceByRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByRecordIDRequest{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		ExcludedRecords: `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`,
		RequiredDsrcs:   `{"DATA_SOURCES": ["CUSTOMERS"]}`,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 119))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2EngineServer_FindPathIncludingSourceByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByRecordID_V2Request{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		ExcludedRecords: `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`,
		RequiredDsrcs:   `{"DATA_SOURCES": ["CUSTOMERS"]}`,
		Flags:           0,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_GetActiveConfigID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetActiveConfigIDRequest{}
	response, err := g2engine.GetActiveConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_GetEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByEntityIDRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
	}
	response, err := g2engine.GetEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_GetEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByEntityID_V2Request{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.GetEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2EngineServer_GetEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByRecordIDRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
	}
	response, err := g2engine.GetEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 35))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2EngineServer_GetEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByRecordID_V2Request{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.GetEntityByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2EngineServer_GetRecord() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
	}
	response, err := g2engine.GetRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","JSON_DATA":{"RECORD_TYPE":"PERSON","PRIMARY_NAME_LAST":"Smith","PRIMARY_NAME_FIRST":"Robert","DATE_OF_BIRTH":"12/11/1978","ADDR_TYPE":"MAILING","ADDR_LINE1":"123 Main Street, Las Vegas NV 89132","PHONE_TYPE":"HOME","PHONE_NUMBER":"702-919-1300","EMAIL_ADDRESS":"bsmith@work.com","DATE":"1/2/18","STATUS":"Active","AMOUNT":"100","DATA_SOURCE":"CUSTOMERS","ENTITY_TYPE":"GENERIC","DSRC_ACTION":"A","RECORD_ID":"1001"}}
}

func ExampleG2EngineServer_GetRecord_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRecord_V2Request{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.GetRecord_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}
}

func ExampleG2EngineServer_GetRedoRecord() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRedoRecordRequest{}
	response, err := g2engine.GetRedoRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","ENTITY_TYPE":"GENERIC","DSRC_ACTION":"X"}
}

func ExampleG2EngineServer_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRepositoryLastModifiedTimeRequest{}
	response, err := g2engine.GetRepositoryLastModifiedTime(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_GetVirtualEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetVirtualEntityByRecordIDRequest{
		RecordList: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`,
	}
	response, err := g2engine.GetVirtualEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_GetVirtualEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetVirtualEntityByRecordID_V2Request{
		RecordList: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`,
		Flags:      0,
	}
	response, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2EngineServer_HowEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.HowEntityByEntityIDRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
	}
	response, err := g2engine.HowEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL","FEATURE_SCORES":{"ADDRESS":[{"INBOUND_FEAT_ID":20,"INBOUND_FEAT":"1515 Adela Lane Las Vegas NV 89111","INBOUND_FEAT_USAGE_TYPE":"HOME","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":42,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":19,"INBOUND_FEAT":"11/12/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":95,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"FMES"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":1,"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":97,"GNR_SN":100,"GNR_GN":95,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"702-919-1300","INBOUND_FEAT_USAGE_TYPE":"MOBILE","CANDIDATE_FEAT_ID":4,"CANDIDATE_FEAT":"702-919-1300","CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1-S1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V100001","MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","RESULT_VIRTUAL_ENTITY_ID":"V1-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"DOB":[{"INBOUND_FEAT_ID":2,"INBOUND_FEAT":"12/11/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FMES"}],"EMAIL":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":93,"GNR_SN":100,"GNR_GN":93,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"},{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"Robert Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":90,"GNR_SN":100,"GNR_GN":88,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}}],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1-S2","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleG2EngineServer_HowEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.HowEntityByEntityID_V2Request{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.HowEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL"}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1-S1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V100001","MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","RESULT_VIRTUAL_ENTITY_ID":"V1-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB"}}],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1-S2","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleG2EngineServer_PrimeEngine() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.PrimeEngineRequest{}
	response, err := g2engine.PrimeEngine(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_SearchByAttributes() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.SearchByAttributesRequest{
		JsonData: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`,
	}
	response, err := g2engine.SearchByAttributes(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 1962))
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1","FEATURE_SCORES":{"EMAIL":[{"INBOUND_FEAT":"bsmith@work.com","CANDIDATE_FEAT":"bsmith@work.com","FULL_SCORE":100}],"NAME":[{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Bob J Smith","GNR_FN":83,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1},{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Robert Smith","GNR_FN":88,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1}]}},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","FEATURES":{"ADDRESS":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20}]},{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3,"USAGE_TYPE":"MAILING","FEAT_DESC_VALUES":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3}]}],"DOB":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2},{"FEAT_DESC":"11/12/1978","LIB_FEAT_ID":19}]}],"EMAIL":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5}]}],"NAME":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1,"USAGE_TYPE":"PRIMARY","FEAT_DESC_VALUES":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1},{"FEAT_DESC":"Bob J Smith","LIB_FEAT_ID":32},{"FEAT_DESC":"Bob Smith","LIB_FEAT_ID":18}]}],"PHONE":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]},{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"MOBILE","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]}],"RECORD_TYPE":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleG2EngineServer_SearchByAttributes_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.SearchByAttributes_V2Request{
		JsonData: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`,
		Flags:    0,
	}
	response, err := g2engine.SearchByAttributes_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1"},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1}}}]}
}

// func ExampleG2engineImpl_SetLogLevel() {
// }

func ExampleG2EngineServer_Stats() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.StatsRequest{}
	response, err := g2engine.Stats(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 138))
	// Output: { "workload": { "loadedRecords": 5,  "addedRecords": 5,  "deletedRecords": 1,  "reevaluations": 0,  "repairedEntities": 0,  "duration":...
}

func ExampleG2EngineServer_WhyEntities() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntitiesRequest{
		EntityID1: getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2: getEntityIdForRecord("CUSTOMERS", "1002"),
	}
	response, err := g2engine.WhyEntities(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 74))
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":1,"MATCH_INFO":{"WHY_KEY":...
}

func ExampleG2EngineServer_WhyEntities_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntities_V2Request{
		EntityID1: getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2: getEntityIdForRecord("CUSTOMERS", "1002"),
		Flags:     0,
	}
	response, err := g2engine.WhyEntities_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":1,"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+ADDRESS+PHONE+EMAIL","WHY_ERRULE_CODE":"SF1_SNAME_CFF_CSTAB","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_WhyEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByEntityIDRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
	}
	response, err := g2engine.WhyEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2EngineServer_WhyEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByEntityID_V2Request{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.WhyEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2EngineServer_WhyEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByRecordIDRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
	}
	response, err := g2engine.WhyEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2EngineServer_WhyEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByRecordID_V2Request{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.WhyEntityByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2EngineServer_WhyRecords() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyRecordsRequest{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
	}
	response, err := g2engine.WhyRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 115))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],...
}

func ExampleG2EngineServer_WhyRecords_V2() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyRecords_V2Request{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		Flags:           0,
	}
	response, err := g2engine.WhyRecords_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":1,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}],"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+PHONE","WHY_ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_Process() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessRequest{
		Record: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
	}
	response, err := g2engine.Process(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ProcessRedoRecord() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessRedoRecordRequest{}
	response, err := g2engine.ProcessRedoRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ProcessRedoRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessRedoRecordWithInfoRequest{
		Flags: 0,
	}
	response, err := g2engine.ProcessRedoRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ProcessWithInfo() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessWithInfoRequest{
		Record: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		Flags:  0,
	}
	response, err := g2engine.ProcessWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ProcessWithResponse() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessWithResponseRequest{
		Record: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
	}
	response, err := g2engine.ProcessWithResponse(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
}

func ExampleG2EngineServer_ProcessWithResponseResize() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessWithResponseResizeRequest{
		Record: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
	}
	response, err := g2engine.ProcessWithResponseResize(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
}

func ExampleG2EngineServer_ReevaluateEntity() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateEntityRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.ReevaluateEntity(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ReevaluateEntityWithInfo() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateEntityWithInfoRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.ReevaluateEntityWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ReevaluateRecord() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.ReevaluateRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ReevaluateRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateRecordWithInfoRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.ReevaluateRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ReplaceRecord() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReplaceRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.ReplaceRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ReplaceRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReplaceRecordWithInfoRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		LoadID:         "G2Engine_test",
		Flags:          0,
	}
	response, err := g2engine.ReplaceRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_DeleteRecord() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.DeleteRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1003",
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.DeleteRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_DeleteRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.DeleteRecordWithInfoRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1003",
		LoadID:         "G2Engine_test",
		Flags:          0,
	}
	response, err := g2engine.DeleteRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_Init() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)

	}
	request := &pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: 0,
	}
	response, err := g2engine.Init(ctx, request)
	if err != nil {
		// This should produce a "senzing-60144002" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_InitWithConfigID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)

	}
	request := &pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   1,
		VerboseLogging: 0,
	}
	response, err := g2engine.InitWithConfigID(ctx, request)
	if err != nil {
		// This should produce a "senzing-60144003" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_Reinit() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)

	// Get a Senzing configuration ID for testing.
	requestToGetActiveConfigID := &pb.GetActiveConfigIDRequest{}
	responseFromGetActiveConfigID, err := g2engine.GetActiveConfigID(ctx, requestToGetActiveConfigID)

	// Example
	request := &pb.ReinitRequest{
		InitConfigID: responseFromGetActiveConfigID.GetResult(),
	}
	response, err := g2engine.Reinit(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_PurgeRepository() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.PurgeRepositoryRequest{}
	response, err := g2engine.PurgeRepository(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_Destroy() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.DestroyRequest{}
	response, err := g2engine.Destroy(ctx, request)
	if err != nil {
		// This should produce a "senzing-60164001" error.
	}
	fmt.Println(response)
	// Output:
}

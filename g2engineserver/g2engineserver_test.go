package g2engineserver

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go/g2config"
	"github.com/senzing/g2-sdk-go/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2engine"
	"github.com/senzing/g2-sdk-go/testhelpers"
	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/go-servegrpc/protobuf/g2engine"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2engineTestSingleton *G2EngineServer
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) G2EngineServer {
	if g2engineTestSingleton == nil {
		g2engineTestSingleton = &G2EngineServer{}

		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if jsonErr != nil {
			test.Logf("Cannot construct system configuration. Error: %v", jsonErr)
		}

		request := &pb.InitRequest{
			ModuleName:     moduleName,
			IniParams:      iniParams,
			VerboseLogging: int32(verboseLogging),
		}
		g2engineTestSingleton.Init(ctx, request)
	}
	return *g2engineTestSingleton
}

func getG2EngineServer(ctx context.Context) G2EngineServer {
	G2EngineServer := &G2EngineServer{}
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	request := &pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	G2EngineServer.Init(ctx, request)
	return *G2EngineServer
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

func setup() error {
	ctx := context.TODO()
	now := time.Now()
	moduleName := "Test module name"
	verboseLogging := 0
	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	// if err != nil {
	// 	return logger.Error(5901, err)
	// }

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		return logger.Error(5902, err)
	}

	// Add Data Sources to in-memory Senzing configuration.

	aG2config := &g2config.G2configImpl{}
	err = aG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return logger.Error(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
	if err != nil {
		return logger.Error(5907, err)
	}

	for _, testDataSource := range testhelpers.TestDataSources {
		_, err := aG2config.AddDataSource(ctx, configHandle, testDataSource.Data)
		if err != nil {
			return logger.Error(5908, err)
		}
	}

	configStr, err := aG2config.Save(ctx, configHandle)
	if err != nil {
		return logger.Error(5909, err)
	}

	err = aG2config.Close(ctx, configHandle)
	if err != nil {
		return logger.Error(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return logger.Error(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &g2configmgr.G2configmgrImpl{}
	err = aG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return logger.Error(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return logger.Error(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return logger.Error(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return logger.Error(5915, err)
	}

	// Purge repository.

	aG2engine := &g2engine.G2engineImpl{}
	err = aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return logger.Error(5903, err)
	}

	err = aG2engine.PurgeRepository(ctx)
	if err != nil {
		return logger.Error(5904, err)
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return logger.Error(5905, err)
	}

	// Add records.

	aG2engine = &g2engine.G2engineImpl{}
	err = aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return logger.Error(5916, err)
	}

	for _, testRecord := range testhelpers.TestRecords {
		err := aG2engine.AddRecord(ctx, testRecord.DataSource, testRecord.Id, testRecord.Data, testRecord.LoadId)
		if err != nil {
			return logger.Error(5917, err)
		}
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return logger.Error(5918, err)
	}

	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestG2engineserver_BuildSimpleSystemConfigurationJson(test *testing.T) {
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

// PurgeRepository() is first to start with a clean database.
func TestG2engineServer_PurgeRepository(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.PurgeRepositoryRequest{}
	response, err := g2engine.PurgeRepository(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_AddRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request1 := &pb.AddRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "111", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response1, err := g2engine.AddRecord(ctx, request1)
	testError(test, ctx, g2engine, err)
	printResponse(test, response1)

	request2 := &pb.AddRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "222",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "222", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response2, err := g2engine.AddRecord(ctx, request2)
	testError(test, ctx, g2engine, err)
	printResponse(test, response2)

}

func TestG2engineServer_AddRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.AddRecordWithInfoRequest{
		DataSourceCode: "TEST",
		RecordID:       "333",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "333", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response, err := g2engine.AddRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_AddRecordWithInfoWithReturnedRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.AddRecordWithInfoWithReturnedRecordIDRequest{
		DataSourceCode: "TEST",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_AddRecordWithReturnedRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.AddRecordWithReturnedRecordIDRequest{
		DataSourceCode: "TEST",
		JsonData:       `{"SOCIAL_HANDLE": "bobby", "DATE_OF_BIRTH": "1/2/1983", "ADDR_STATE": "WI", "ADDR_POSTAL_CODE": "54434", "SSN_NUMBER": "987-65-4321", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "Smith", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response, err := g2engine.AddRecordWithReturnedRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_CheckRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.CheckRecordRequest{
		Record:          `{"DATA_SOURCE": "TEST", "NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith", "NAME_MIDDLE": "M" }], "PASSPORT_NUMBER": "PP11111", "PASSPORT_COUNTRY": "US", "DRIVERS_LICENSE_NUMBER": "DL11111", "SSN_NUMBER": "111-11-1111"}`,
		RecordQueryList: `{"RECORDS": [{"DATA_SOURCE": "TEST","RECORD_ID": "111"},{"DATA_SOURCE": "TEST","RECORD_ID": "123456789"}]}`,
	}
	response, err := g2engine.CheckRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ExportJSONEntityReport(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ExportJSONEntityReportRequest{
		Flags: int64(0),
	}
	response, err := g2engine.ExportJSONEntityReport(ctx, request)
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
	request := &pb.FindInterestingEntitiesByEntityIDRequest{
		EntityID: 1,
		Flags:    0,
	}
	response, err := g2engine.FindInterestingEntitiesByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindInterestingEntitiesByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindInterestingEntitiesByRecordIDRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		Flags:          0,
	}
	response, err := g2engine.FindInterestingEntitiesByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindNetworkByEntityIDRequest{
		EntityList:     `{"ENTITIES": [{"ENTITY_ID": 1}, {"ENTITY_ID": 2}]}`,
		MaxDegree:      2,
		BuildOutDegree: 1,
		MaxEntities:    10,
	}
	response, err := g2engine.FindNetworkByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindNetworkByEntityID_V2Request{
		EntityList:     `{"ENTITIES": [{"ENTITY_ID": 1}, {"ENTITY_ID": 2}]}`,
		MaxDegree:      2,
		BuildOutDegree: 1,
		MaxEntities:    10,
		Flags:          0,
	}
	response, err := g2engine.FindNetworkByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindNetworkByRecordIDRequest{
		RecordList:     `{"RECORDS": [{"DATA_SOURCE": "TEST", "RECORD_ID": "111"}, {"DATA_SOURCE": "TEST", "RECORD_ID": "222"}, {"DATA_SOURCE": "TEST", "RECORD_ID": "333"}]}`,
		MaxDegree:      1,
		BuildOutDegree: 2,
		MaxEntities:    10,
	}
	response, err := g2engine.FindNetworkByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindNetworkByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindNetworkByRecordID_V2Request{
		RecordList:     `{"RECORDS": [{"DATA_SOURCE": "TEST", "RECORD_ID": "111"}, {"DATA_SOURCE": "TEST", "RECORD_ID": "222"}, {"DATA_SOURCE": "TEST", "RECORD_ID": "333"}]}`,
		MaxDegree:      1,
		BuildOutDegree: 2,
		MaxEntities:    10,
		Flags:          0,
	}
	response, err := g2engine.FindNetworkByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathByEntityIDRequest{
		EntityID1: 1,
		EntityID2: 2,
		MaxDegree: 1,
	}
	response, err := g2engine.FindPathByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathByEntityID_V2Request{
		EntityID1: 1,
		EntityID2: 2,
		MaxDegree: 1,
		Flags:     0,
	}
	response, err := g2engine.FindPathByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathByRecordIDRequest{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
	}
	response, err := g2engine.FindPathByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathByRecordID_V2Request{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		Flags:           0,
	}
	response, err := g2engine.FindPathByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathExcludingByEntityIDRequest{
		EntityID1:        1,
		EntityID2:        2,
		MaxDegree:        1,
		ExcludedEntities: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
	}
	response, err := g2engine.FindPathExcludingByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathExcludingByEntityID_V2Request{
		EntityID1:        1,
		EntityID2:        2,
		MaxDegree:        1,
		ExcludedEntities: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		Flags:            0,
	}
	response, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathExcludingByRecordIDRequest{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		ExcludedRecords: `{"RECORDS": [{ "DATA_SOURCE": "TEST", "RECORD_ID": "111"}]}`,
	}
	response, err := g2engine.FindPathExcludingByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathExcludingByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathExcludingByRecordID_V2Request{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		ExcludedRecords: `{"RECORDS": [{ "DATA_SOURCE": "TEST", "RECORD_ID": "111"}]}`,
		Flags:           0,
	}
	response, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathIncludingSourceByEntityIDRequest{
		EntityID1:        1,
		EntityID2:        2,
		MaxDegree:        1,
		ExcludedEntities: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		RequiredDsrcs:    `{"DATA_SOURCES": ["TEST"]}`,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathIncludingSourceByEntityID_V2Request{
		EntityID1:        1,
		EntityID2:        2,
		MaxDegree:        1,
		ExcludedEntities: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		RequiredDsrcs:    `{"DATA_SOURCES": ["TEST"]}`,
		Flags:            0,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathIncludingSourceByRecordIDRequest{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		ExcludedRecords: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		RequiredDsrcs:   `{"DATA_SOURCES": ["TEST"]}`,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_FindPathIncludingSourceByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathIncludingSourceByRecordID_V2Request{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		ExcludedRecords: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		RequiredDsrcs:   `{"DATA_SOURCES": ["TEST"]}`,
		Flags:           0,
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
	request := &pb.GetEntityByEntityIDRequest{
		EntityID: 1,
	}
	response, err := g2engine.GetEntityByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetEntityByEntityID_V2Request{
		EntityID: 1,
		Flags:    0,
	}
	response, err := g2engine.GetEntityByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetEntityByRecordIDRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
	}
	response, err := g2engine.GetEntityByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetEntityByRecordID_V2Request{
		DataSourceCode: "TEST",
		RecordID:       "111",
		Flags:          0,
	}
	response, err := g2engine.GetEntityByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
	}
	response, err := g2engine.GetRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetRecord_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetRecord_V2Request{
		DataSourceCode: "TEST",
		RecordID:       "111",
		Flags:          0,
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
	request := &pb.GetVirtualEntityByRecordIDRequest{
		RecordList: `{"RECORDS": [{"DATA_SOURCE": "TEST","RECORD_ID": "111"},{"DATA_SOURCE": "TEST","RECORD_ID": "222"}]}`,
	}
	response, err := g2engine.GetVirtualEntityByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_GetVirtualEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetVirtualEntityByRecordID_V2Request{
		RecordList: `{"RECORDS": [{"DATA_SOURCE": "TEST","RECORD_ID": "111"},{"DATA_SOURCE": "TEST","RECORD_ID": "222"}]}`,
		Flags:      0,
	}
	response, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_HowEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.HowEntityByEntityIDRequest{
		EntityID: 1,
	}
	response, err := g2engine.HowEntityByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_HowEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.HowEntityByEntityID_V2Request{
		EntityID: 1,
		Flags:    0,
	}
	response, err := g2engine.HowEntityByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_Init(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		assert.FailNow(test, jsonErr.Error())
	}
	request := &pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: 0,
	}
	response, err := g2engine.Init(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		assert.FailNow(test, jsonErr.Error())
	}
	request := &pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   1,
		VerboseLogging: 0,
	}
	response, err := g2engine.InitWithConfigID(ctx, request)
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
	request := &pb.ProcessRequest{
		Record: `{"DATA_SOURCE": "TEST", "SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "444", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
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
	request := &pb.ProcessRedoRecordWithInfoRequest{
		Flags: 0,
	}
	response, err := g2engine.ProcessRedoRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ProcessWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessWithInfoRequest{
		Record: `{"DATA_SOURCE": "TEST", "SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "555", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		Flags:  0,
	}
	response, err := g2engine.ProcessWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ProcessWithResponse(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessWithResponseRequest{
		Record: `{"DATA_SOURCE": "TEST", "SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "666", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
	}
	response, err := g2engine.ProcessWithResponse(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ProcessWithResponseResize(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessWithResponseResizeRequest{
		Record: `{"DATA_SOURCE": "TEST", "SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "777", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
	}
	response, err := g2engine.ProcessWithResponseResize(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReevaluateEntityRequest{
		EntityID: 1,
		Flags:    0,
	}
	response, err := g2engine.ReevaluateEntity(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateEntityWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReevaluateEntityWithInfoRequest{
		EntityID: 1,
		Flags:    0,
	}
	response, err := g2engine.ReevaluateEntityWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReevaluateRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		Flags:          0,
	}
	response, err := g2engine.ReevaluateRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReevaluateRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReevaluateRecordWithInfoRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		Flags:          0,
	}
	response, err := g2engine.ReevaluateRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
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

func TestG2engineServer_ReplaceRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReplaceRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1984", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "111", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response, err := g2engine.ReplaceRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_ReplaceRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReplaceRecordWithInfoRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1985", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "111", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
		Flags:          0,
	}
	response, err := g2engine.ReplaceRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_SearchByAttributes(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.SearchByAttributesRequest{
		JsonData: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`,
	}
	response, err := g2engine.SearchByAttributes(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_SearchByAttributes_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.SearchByAttributes_V2Request{
		JsonData: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`,
		Flags:    0,
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
	request := &pb.WhyEntitiesRequest{
		EntityID1: 1,
		EntityID2: 2,
	}
	response, err := g2engine.WhyEntities(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntities_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntities_V2Request{
		EntityID1: 1,
		EntityID2: 2,
		Flags:     0,
	}
	response, err := g2engine.WhyEntities_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntityByEntityIDRequest{
		EntityID: 1,
	}
	response, err := g2engine.WhyEntityByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntityByEntityID_V2Request{
		EntityID: 1,
		Flags:    0,
	}
	response, err := g2engine.WhyEntityByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntityByRecordIDRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
	}
	response, err := g2engine.WhyEntityByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntityByRecordID_V2Request{
		DataSourceCode: "TEST",
		RecordID:       "111",
		Flags:          0,
	}
	response, err := g2engine.WhyEntityByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyRecordsRequest{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
	}
	response, err := g2engine.WhyRecords(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_WhyRecords_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyRecords_V2Request{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		Flags:           0,
	}
	response, err := g2engine.WhyRecords_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_DeleteRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.DeleteRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "333",
		LoadID:         "TEST",
	}
	response, err := g2engine.DeleteRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

func TestG2engineServer_DeleteRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.DeleteRecordWithInfoRequest{
		DataSourceCode: "TEST",
		RecordID:       "333",
		LoadID:         "TEST",
		Flags:          0,
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
	testError(test, ctx, g2engine, err)
	printResponse(test, response)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

// PurgeRepository() is first to start with a clean database.
func ExampleG2EngineServer_PurgeRepository() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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

func ExampleG2EngineServer_AddRecord() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "111", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response, err := g2engine.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_AddRecord_secondRecord() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "222",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh2", "DATE_OF_BIRTH": "6/9/1983", "ADDR_STATE": "WI", "ADDR_POSTAL_CODE": "53543", "SSN_NUMBER": "153-33-5185", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "222", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "OCEANGUY", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response, err := g2engine.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_AddRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordWithInfoRequest{
		DataSourceCode: "TEST",
		RecordID:       "333",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "333", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response, err := g2engine.AddRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"333","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_AddRecordWithInfoWithReturnedRecordID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordWithInfoWithReturnedRecordIDRequest{
		DataSourceCode: "TEST",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
		Flags:          0,
	}
	response, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetWithInfo(), 37))
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_AddRecordWithReturnedRecordID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordWithReturnedRecordIDRequest{
		DataSourceCode: "TEST",
		JsonData:       `{"SOCIAL_HANDLE": "bobby", "DATE_OF_BIRTH": "1/2/1983", "ADDR_STATE": "WI", "ADDR_POSTAL_CODE": "54434", "SSN_NUMBER": "987-65-4321", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "Smith", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response, err := g2engine.AddRecordWithReturnedRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Length of record identifier is %d hexadecimal characters.\n", len(response.GetResult()))
	// Output: Length of record identifier is 40 hexadecimal characters.
}

func ExampleG2EngineServer_CheckRecord() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.CheckRecordRequest{
		Record:          `{"DATA_SOURCE": "TEST", "NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith", "NAME_MIDDLE": "M" }], "PASSPORT_NUMBER": "PP11111", "PASSPORT_COUNTRY": "US", "DRIVERS_LICENSE_NUMBER": "DL11111", "SSN_NUMBER": "111-11-1111"}`,
		RecordQueryList: `{"RECORDS": [{"DATA_SOURCE": "TEST","RECORD_ID": "111"},{"DATA_SOURCE": "TEST","RECORD_ID": "123456789"}]}`,
	}
	response, err := g2engine.CheckRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"CHECK_RECORD_RESPONSE":[{"DSRC_CODE":"TEST","RECORD_ID":"111","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","MATCH_KEY":"","ERRULE_CODE":"","ERRULE_ID":0,"CANDIDATE_MATCH":"N","NON_GENERIC_CANDIDATE_MATCH":"N"}]}
}

func ExampleG2EngineServer_CloseExport() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.CountRedoRecordsRequest{}
	response, err := g2engine.CountRedoRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: 0
}

func ExampleG2EngineServer_DeleteRecord() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.DeleteRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "333",
		LoadID:         "TEST",
	}
	response, err := g2engine.DeleteRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_DeleteRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.DeleteRecordWithInfoRequest{
		DataSourceCode: "TEST",
		RecordID:       "333",
		LoadID:         "TEST",
		Flags:          0,
	}
	response, err := g2engine.DeleteRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"333","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ExportCSVEntityReport() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)

	// Create a handle for the example.
	requestToExportJSONEntityReport := &pb.ExportJSONEntityReportRequest{
		Flags: 0,
	}
	responseFromExportJSONEntityReport, err := g2engine.ExportJSONEntityReport(ctx, requestToExportJSONEntityReport)

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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindInterestingEntitiesByEntityIDRequest{
		EntityID: 1,
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindInterestingEntitiesByRecordIDRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByEntityIDRequest{
		EntityList:     `{"ENTITIES": [{"ENTITY_ID": 1}, {"ENTITY_ID": 2}]}`,
		MaxDegree:      2,
		BuildOutDegree: 1,
		MaxEntities:    10,
	}
	response, err := g2engine.FindNetworkByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 124))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,...
}

func ExampleG2EngineServer_FindNetworkByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByEntityID_V2Request{
		EntityList:     `{"ENTITIES": [{"ENTITY_ID": 1}, {"ENTITY_ID": 2}]}`,
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}},{"RESOLVED_ENTITY":{"ENTITY_ID":3}}]}
}

func ExampleG2EngineServer_FindNetworkByRecordID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByRecordIDRequest{
		RecordList:     `{"RECORDS": [{"DATA_SOURCE": "TEST", "RECORD_ID": "111"}, {"DATA_SOURCE": "TEST", "RECORD_ID": "222"}]}`,
		MaxDegree:      1,
		BuildOutDegree: 2,
		MaxEntities:    10,
	}
	response, err := g2engine.FindNetworkByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 138))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_FindNetworkByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByRecordID_V2Request{
		RecordList:     `{"RECORDS": [{"DATA_SOURCE": "TEST", "RECORD_ID": "111"}, {"DATA_SOURCE": "TEST", "RECORD_ID": "222"}]}`,
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}},{"RESOLVED_ENTITY":{"ENTITY_ID":3}}]}
}

func ExampleG2EngineServer_FindPathByEntityID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByEntityIDRequest{
		EntityID1: 1,
		EntityID2: 2,
		MaxDegree: 1,
	}
	response, err := g2engine.FindPathByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 109))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByEntityID_V2Request{
		EntityID1: 1,
		EntityID2: 2,
		MaxDegree: 1,
		Flags:     0,
	}
	response, err := g2engine.FindPathByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathByRecordID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByRecordIDRequest{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
	}
	response, err := g2engine.FindPathByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 89))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":...
}

func ExampleG2EngineServer_FindPathByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByRecordID_V2Request{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		Flags:           0,
	}
	response, err := g2engine.FindPathByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathExcludingByEntityID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByEntityIDRequest{
		EntityID1:        1,
		EntityID2:        2,
		MaxDegree:        1,
		ExcludedEntities: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
	}
	response, err := g2engine.FindPathExcludingByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 109))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathExcludingByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByEntityID_V2Request{
		EntityID1:        1,
		EntityID2:        2,
		MaxDegree:        1,
		ExcludedEntities: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		Flags:            0,
	}
	response, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathExcludingByRecordID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByRecordIDRequest{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		ExcludedRecords: `{"RECORDS": [{ "DATA_SOURCE": "TEST", "RECORD_ID": "111"}]}`,
	}
	response, err := g2engine.FindPathExcludingByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 109))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathExcludingByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByRecordID_V2Request{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		ExcludedRecords: `{"RECORDS": [{ "DATA_SOURCE": "TEST", "RECORD_ID": "111"}]}`,
		Flags:           0,
	}
	response, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathIncludingSourceByEntityID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByEntityIDRequest{
		EntityID1:        1,
		EntityID2:        2,
		MaxDegree:        1,
		ExcludedEntities: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		RequiredDsrcs:    `{"DATA_SOURCES": ["TEST"]}`,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 106))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathIncludingSourceByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByEntityID_V2Request{
		EntityID1:        1,
		EntityID2:        2,
		MaxDegree:        1,
		ExcludedEntities: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		RequiredDsrcs:    `{"DATA_SOURCES": ["TEST"]}`,
		Flags:            0,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathIncludingSourceByRecordID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByRecordIDRequest{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		ExcludedRecords: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		RequiredDsrcs:   `{"DATA_SOURCES": ["TEST"]}`,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 119))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2EngineServer_FindPathIncludingSourceByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByRecordID_V2Request{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		MaxDegree:       1,
		ExcludedRecords: `{"ENTITIES": [{"ENTITY_ID": 1}]}`,
		RequiredDsrcs:   `{"DATA_SOURCES": ["TEST"]}`,
		Flags:           0,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_GetActiveConfigID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByEntityIDRequest{
		EntityID: 1,
	}
	response, err := g2engine.GetEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_GetEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByEntityID_V2Request{
		EntityID: 1,
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByRecordIDRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
	}
	response, err := g2engine.GetEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 35))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2EngineServer_GetEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByRecordID_V2Request{
		DataSourceCode: "TEST",
		RecordID:       "111",
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
	}
	response, err := g2engine.GetRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111","JSON_DATA":{"SOCIAL_HANDLE":"flavorh","DATE_OF_BIRTH":"4/8/1983","ADDR_STATE":"LA","ADDR_POSTAL_CODE":"71232","SSN_NUMBER":"053-39-3251","GENDER":"F","srccode":"MDMPER","CC_ACCOUNT_NUMBER":"5534202208773608","ADDR_CITY":"Delhi","DRIVERS_LICENSE_STATE":"DE","PHONE_NUMBER":"225-671-0796","NAME_LAST":"JOHNSON","entityid":"284430058","ADDR_LINE1":"772 Armstrong RD","DATA_SOURCE":"TEST","ENTITY_TYPE":"TEST","DSRC_ACTION":"A","RECORD_ID":"111"}}
}

func ExampleG2EngineServer_GetRecord_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRecord_V2Request{
		DataSourceCode: "TEST",
		RecordID:       "111",
		Flags:          0,
	}
	response, err := g2engine.GetRecord_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111"}
}

func ExampleG2EngineServer_GetRedoRecord() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRedoRecordRequest{}
	response, err := g2engine.GetRedoRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output:
}

func ExampleG2EngineServer_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetVirtualEntityByRecordIDRequest{
		RecordList: `{"RECORDS": [{"DATA_SOURCE": "TEST","RECORD_ID": "111"},{"DATA_SOURCE": "TEST","RECORD_ID": "222"}]}`,
	}
	response, err := g2engine.GetVirtualEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_GetVirtualEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetVirtualEntityByRecordID_V2Request{
		RecordList: `{"RECORDS": [{"DATA_SOURCE": "TEST","RECORD_ID": "111"},{"DATA_SOURCE": "TEST","RECORD_ID": "222"}]}`,
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.HowEntityByEntityIDRequest{
		EntityID: 1,
	}
	response, err := g2engine.HowEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919"}]}]}]}}}
}

func ExampleG2EngineServer_HowEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.HowEntityByEntityID_V2Request{
		EntityID: 1,
		Flags:    0,
	}
	response, err := g2engine.HowEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919"}]}]}]}}}
}

func ExampleG2EngineServer_Init() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_InitWithConfigID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_PrimeEngine() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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

func ExampleG2EngineServer_Process() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessRequest{
		Record: `{"DATA_SOURCE": "TEST", "SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "444", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
	}
	response, err := g2engine.Process(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ProcessRedoRecord() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessWithInfoRequest{
		Record: `{"DATA_SOURCE": "TEST", "SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "555", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		Flags:  0,
	}
	response, err := g2engine.ProcessWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"555","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ProcessWithResponse() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessWithResponseRequest{
		Record: `{"DATA_SOURCE": "TEST", "SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "666", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
	}
	response, err := g2engine.ProcessWithResponse(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
}

func ExampleG2EngineServer_ProcessWithResponseResize() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessWithResponseResizeRequest{
		Record: `{"DATA_SOURCE": "TEST", "SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "777", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
	}
	response, err := g2engine.ProcessWithResponseResize(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
}

func ExampleG2EngineServer_ReevaluateEntity() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateEntityRequest{
		EntityID: 1,
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateEntityWithInfoRequest{
		EntityID: 1,
		Flags:    0,
	}
	response, err := g2engine.ReevaluateEntityWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ReevaluateRecord() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
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
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateRecordWithInfoRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		Flags:          0,
	}
	response, err := g2engine.ReevaluateRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}

}

func ExampleG2EngineServer_Reinit() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)

	// Get a Senzing configuration ID for testing.
	requestToGetActiveConfigID := &pb.GetActiveConfigIDRequest{}
	responseFromGetActiveConfigID, err := g2engine.GetActiveConfigID(ctx, requestToGetActiveConfigID)

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

func ExampleG2EngineServer_ReplaceRecord() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReplaceRecordRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1985", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "111", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
	}
	response, err := g2engine.ReplaceRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ReplaceRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReplaceRecordWithInfoRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
		JsonData:       `{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1985", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "111", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "JOHNSON", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
		LoadID:         "TEST",
		Flags:          0,
	}
	response, err := g2engine.ReplaceRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_SearchByAttributes() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.SearchByAttributesRequest{
		JsonData: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`,
	}
	response, err := g2engine.SearchByAttributes(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 54))
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":...
}

func ExampleG2EngineServer_SearchByAttributes_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.SearchByAttributes_V2Request{
		JsonData: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`,
		Flags:    0,
	}
	response, err := g2engine.SearchByAttributes_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","MATCH_KEY":"+NAME+SSN","ERRULE_CODE":"SF1_PNAME_CSTAB"},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1}}}]}
}

func ExampleG2EngineServer_Stats() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.StatsRequest{}
	response, err := g2engine.Stats(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 125))
	// Output: { "workload": { "loadedRecords": 5,  "addedRecords": 2,  "deletedRecords": 0,  "reevaluations": 0,  "repairedEntities": 0,...
}

func ExampleG2EngineServer_WhyEntities() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntitiesRequest{
		EntityID1: 1,
		EntityID2: 2,
	}
	response, err := g2engine.WhyEntities(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 74))
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":2,"MATCH_INFO":{"WHY_KEY":...
}

func ExampleG2EngineServer_WhyEntities_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntities_V2Request{
		EntityID1: 1,
		EntityID2: 2,
		Flags:     0,
	}
	response, err := g2engine.WhyEntities_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":2,"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_WhyEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByEntityIDRequest{
		EntityID: 1,
	}
	response, err := g2engine.WhyEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 101))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_WhyEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByEntityID_V2Request{
		EntityID: 1,
		Flags:    0,
	}
	response, err := g2engine.WhyEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 101))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_WhyEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByRecordIDRequest{
		DataSourceCode: "TEST",
		RecordID:       "111",
	}
	response, err := g2engine.WhyEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 101))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_WhyEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByRecordID_V2Request{
		DataSourceCode: "TEST",
		RecordID:       "111",
		Flags:          0,
	}
	response, err := g2engine.WhyEntityByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 101))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_WhyRecords() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyRecordsRequest{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
	}
	response, err := g2engine.WhyRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 114))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"}],...
}

func ExampleG2EngineServer_WhyRecords_V2() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyRecords_V2Request{
		DataSourceCode1: "TEST",
		RecordID1:       "111",
		DataSourceCode2: "TEST",
		RecordID2:       "222",
		Flags:           0,
	}
	response, err := g2engine.WhyRecords_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":2,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"TEST","RECORD_ID":"222"}],"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-DOB-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_Destroy() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.DestroyRequest{}
	response, err := g2engine.Destroy(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

// PurgeRepository() is first to start with a clean database.
func ExampleG2EngineServer_PurgeRepository_finalCleanup() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2engineserver/g2engineserver_test.go
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

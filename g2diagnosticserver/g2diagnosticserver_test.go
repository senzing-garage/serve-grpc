package g2diagnosticserver

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
	pb "github.com/senzing/go-servegrpc/protobuf/g2diagnostic"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
)

var (
	g2diagnosticSingleton *G2DiagnosticServer
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) G2DiagnosticServer {
	if g2diagnosticSingleton == nil {
		g2diagnosticSingleton = &G2DiagnosticServer{}

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
		g2diagnosticSingleton.Init(ctx, request)
	}
	return *g2diagnosticSingleton
}

func getG2DiagnosticServer(ctx context.Context) G2DiagnosticServer {
	g2diagnosticServer := &G2DiagnosticServer{}
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
	g2diagnosticServer.Init(ctx, request)
	return *g2diagnosticServer
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if 1 == 0 {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, g2diagnostic G2DiagnosticServer, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, g2diagnostic G2DiagnosticServer, err error) {
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
	var err error = nil
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

func TestG2diagnosticserver_BuildSimpleSystemConfigurationJson(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2diagnosticserver_CheckDBPerf(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.CheckDBPerfRequest{
		SecondsToRun: int32(1),
	}
	actual, err := g2diagnostic.CheckDBPerf(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestEntityListBySize(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetEntityListBySizeRequest{
		EntitySize: int32(10),
	}
	actual, err := g2diagnostic.GetEntityListBySize(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)

	entityListBySizeHandle := actual.Result
	request2 := &pb.FetchNextEntityBySizeRequest{
		EntityListBySizeHandle: entityListBySizeHandle,
	}
	actual2, err2 := g2diagnostic.FetchNextEntityBySize(ctx, request2)
	testError(test, ctx, g2diagnostic, err2)
	test.Log("Actual:", actual2)

	request3 := &pb.CloseEntityListBySizeRequest{
		EntityListBySizeHandle: entityListBySizeHandle,
	}
	actual3, err3 := g2diagnostic.CloseEntityListBySize(ctx, request3)
	testError(test, ctx, g2diagnostic, err3)
	test.Log("Actual:", actual3)
}

func TestG2diagnosticserver_FindEntitiesByFeatureIDs(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.FindEntitiesByFeatureIDsRequest{
		Features: "{\"ENTITY_ID\":1,\"LIB_FEAT_IDS\":[1,3,4]}",
	}
	actual, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetAvailableMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetAvailableMemoryRequest{}
	actual, err := g2diagnostic.GetAvailableMemory(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetDataSourceCounts(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetDataSourceCountsRequest{}
	actual, err := g2diagnostic.GetDataSourceCounts(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetDBInfo(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetDBInfoRequest{}
	actual, err := g2diagnostic.GetDBInfo(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetEntityDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetEntityDetailsRequest{
		EntityID:                int64(1),
		IncludeInternalFeatures: 1,
	}
	actual, err := g2diagnostic.GetEntityDetails(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetEntityResume(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetEntityResumeRequest{
		EntityID: int64(1),
	}
	actual, err := g2diagnostic.GetEntityResume(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetEntitySizeBreakdown(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetEntitySizeBreakdownRequest{
		MinimumEntitySize:       int32(1),
		IncludeInternalFeatures: int32(1),
	}
	actual, err := g2diagnostic.GetEntitySizeBreakdown(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetFeature(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetFeatureRequest{
		LibFeatID: int64(1),
	}
	actual, err := g2diagnostic.GetFeature(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetGenericFeatures(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetGenericFeaturesRequest{
		FeatureType:           "PHONE",
		MaximumEstimatedCount: 10,
	}
	actual, err := g2diagnostic.GetGenericFeatures(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetLogicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetLogicalCoresRequest{}
	actual, err := g2diagnostic.GetLogicalCores(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetMappingStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: 1,
	}
	actual, err := g2diagnostic.GetMappingStatistics(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetPhysicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetPhysicalCoresRequest{}
	actual, err := g2diagnostic.GetPhysicalCores(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetRelationshipDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetRelationshipDetailsRequest{
		RelationshipID:          int64(1),
		IncludeInternalFeatures: 1,
	}
	actual, err := g2diagnostic.GetRelationshipDetails(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetResolutionStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetResolutionStatisticsRequest{}
	actual, err := g2diagnostic.GetResolutionStatistics(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_GetTotalSystemMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetTotalSystemMemoryRequest{}
	actual, err := g2diagnostic.GetTotalSystemMemory(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_Init(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		assert.FailNow(test, jsonErr.Error())
	}
	request := &pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	actual, err := g2diagnostic.Init(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		assert.FailNow(test, jsonErr.Error())
	}
	request := &pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   int64(1),
		VerboseLogging: int32(0),
	}
	actual, err := g2diagnostic.InitWithConfigID(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.ReinitRequest{
		InitConfigID: int64(testhelpers.TestConfigDataId),
	}
	actual, err := g2diagnostic.Reinit(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticserver_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.DestroyRequest{}
	actual, err := g2diagnostic.Destroy(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2DiagnosticServer_CheckDBPerf() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.CheckDBPerfRequest{
		SecondsToRun: int32(1),
	}
	result, err := g2diagnostic.CheckDBPerf(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result.Result, 25))
	// Output: {"numRecordsInserted":...
}

func ExampleG2DiagnosticServer_FindEntitiesByFeatureIDs() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.FindEntitiesByFeatureIDsRequest{
		Features: "{\"ENTITY_ID\":1,\"LIB_FEAT_IDS\":[1,3,4]}",
	}
	result, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: [{"LIB_FEAT_ID":4,"USAGE_TYPE":"","RES_ENT_ID":2}]
}

func ExampleG2DiagnosticServer_GetAvailableMemory() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetAvailableMemoryRequest{}
	result, err := g2diagnostic.GetAvailableMemory(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_GetDataSourceCounts() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetDataSourceCountsRequest{}
	result, err := g2diagnostic.GetDataSourceCounts(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: [{"DSRC_ID":1001,"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_ID":3,"ETYPE_CODE":"GENERIC","OBS_ENT_COUNT":2,"DSRC_RECORD_COUNT":3}]
}

func ExampleG2DiagnosticServer_GetDBInfo() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetDBInfoRequest{}
	result, err := g2diagnostic.GetDBInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result.Result, 52))
	// Output: {"Hybrid Mode":false,"Database Details":[{"Name":...
}

func ExampleG2DiagnosticServer_GetEntityDetails() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetEntityDetailsRequest{
		EntityID:                int64(1),
		IncludeInternalFeatures: 1,
	}
	result, err := g2diagnostic.GetEntityDetails(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: [{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"","FEAT_DESC":"JOHNSON"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","FEAT_DESC":"4/8/1983"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"No","FTYPE_CODE":"GENDER","USAGE_TYPE":"","FEAT_DESC":"F"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"","FEAT_DESC":"772 Armstrong RD Delhi LA 71232"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"","FEAT_DESC":"225-671-0796"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"No","FTYPE_CODE":"SSN","USAGE_TYPE":"","FEAT_DESC":"053-39-3251"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"No","FTYPE_CODE":"LOGIN_ID","USAGE_TYPE":"","FEAT_DESC":"flavorh"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"No","FTYPE_CODE":"ACCT_NUM","USAGE_TYPE":"CC","FEAT_DESC":"5534202208773608"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"JNSN"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"JNSN|DOB=80804"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"JNSN|POST=71232"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"JNSN|SSN=3251"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"772|ARMSTRNK||TL"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"772|ARMSTRNK||71232"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"ID_KEY","USAGE_TYPE":"","FEAT_DESC":"ACCT_NUM=5534202208773608"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"ID_KEY","USAGE_TYPE":"","FEAT_DESC":"SSN=053-39-3251"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"PHONE_KEY","USAGE_TYPE":"","FEAT_DESC":"2256710796"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"SEARCH_KEY","USAGE_TYPE":"","FEAT_DESC":"LOGIN_ID:FLAVORH|"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","DERIVED":"Yes","FTYPE_CODE":"SEARCH_KEY","USAGE_TYPE":"","FEAT_DESC":"SSN:3251|80804|"}]
}

func ExampleG2DiagnosticServer_GetEntityResume() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetEntityResumeRequest{
		EntityID: int64(1),
	}
	result, err := g2diagnostic.GetEntityResume(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: [{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9001","ENT_SRC_DESC":"JOHNSON","JSON_DATA":"{\"SOCIAL_HANDLE\":\"flavorh\",\"DATE_OF_BIRTH\":\"4/8/1983\",\"ADDR_STATE\":\"LA\",\"ADDR_POSTAL_CODE\":\"71232\",\"SSN_NUMBER\":\"053-39-3251\",\"GENDER\":\"F\",\"srccode\":\"MDMPER\",\"CC_ACCOUNT_NUMBER\":\"5534202208773608\",\"ADDR_CITY\":\"Delhi\",\"DRIVERS_LICENSE_STATE\":\"DE\",\"PHONE_NUMBER\":\"225-671-0796\",\"NAME_LAST\":\"JOHNSON\",\"entityid\":\"284430058\",\"ADDR_LINE1\":\"772 Armstrong RD\",\"DATA_SOURCE\":\"EXAMPLE_DATA_SOURCE\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"9001\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9002","ENT_SRC_DESC":"JOHNSON","JSON_DATA":"{\"SOCIAL_HANDLE\":\"flavorh\",\"DATE_OF_BIRTH\":\"4/8/1983\",\"ADDR_STATE\":\"LA\",\"ADDR_POSTAL_CODE\":\"71232\",\"SSN_NUMBER\":\"053-39-3251\",\"GENDER\":\"F\",\"srccode\":\"MDMPER\",\"CC_ACCOUNT_NUMBER\":\"5534202208773608\",\"ADDR_CITY\":\"Delhi\",\"DRIVERS_LICENSE_STATE\":\"DE\",\"PHONE_NUMBER\":\"225-671-0796\",\"NAME_LAST\":\"JOHNSON\",\"entityid\":\"284430058\",\"ADDR_LINE1\":\"772 Armstrong RD\",\"DATA_SOURCE\":\"EXAMPLE_DATA_SOURCE\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"9002\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","RECORD_ID":"9003","ENT_SRC_DESC":"Smith","JSON_DATA":"{\"ADDR_STATE\":\"LA\",\"ADDR_POSTAL_CODE\":\"71232\",\"GENDER\":\"M\",\"srccode\":\"MDMPER\",\"ADDR_CITY\":\"Delhi\",\"PHONE_NUMBER\":\"225-671-0796\",\"NAME_LAST\":\"Smith\",\"entityid\":\"284430058\",\"ADDR_LINE1\":\"772 Armstrong RD\",\"DATA_SOURCE\":\"EXAMPLE_DATA_SOURCE\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"9003\"}"}]
}

func ExampleG2DiagnosticServer_GetEntitySizeBreakdown() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetEntitySizeBreakdownRequest{
		MinimumEntitySize:       int32(1),
		IncludeInternalFeatures: int32(1),
	}
	result, err := g2diagnostic.GetEntitySizeBreakdown(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: [{"ENTITY_SIZE": 1,"ENTITY_COUNT": 2,"NAME": 1.00,"DOB": 0.50,"GENDER": 1.00,"ADDRESS": 1.00,"PHONE": 1.00,"SSN": 0.50,"LOGIN_ID": 0.50,"ACCT_NUM": 0.50,"NAME_KEY": 6.00,"ADDR_KEY": 2.00,"ID_KEY": 1.00,"PHONE_KEY": 1.00,"SEARCH_KEY": 1.00,"MIN_RES_ENT_ID": 1,"MAX_RES_ENT_ID": 2}]
}

func ExampleG2DiagnosticServer_GetFeature() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetFeatureRequest{
		LibFeatID: int64(1),
	}
	result, err := g2diagnostic.GetFeature(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: {"LIB_FEAT_ID":1,"FTYPE_CODE":"NAME","ELEMENTS":[{"FELEM_CODE":"TOKENIZED_NM","FELEM_VALUE":"JOHNSON"},{"FELEM_CODE":"CATEGORY","FELEM_VALUE":"PERSON"},{"FELEM_CODE":"CULTURE","FELEM_VALUE":"ANGLO"},{"FELEM_CODE":"SUR_NAME","FELEM_VALUE":"JOHNSON"},{"FELEM_CODE":"FULL_NAME","FELEM_VALUE":"JOHNSON"}]}
}

func ExampleG2DiagnosticServer_GetGenericFeatures() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetGenericFeaturesRequest{
		FeatureType:           "PHONE",
		MaximumEstimatedCount: 10,
	}
	result, err := g2diagnostic.GetGenericFeatures(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: []
}

func ExampleG2DiagnosticServer_GetLogicalCores() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetLogicalCoresRequest{}
	result, err := g2diagnostic.GetLogicalCores(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_GetMappingStatistics() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: 1,
	}
	result, err := g2diagnostic.GetMappingStatistics(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: [{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"","REC_COUNT":2,"REC_PCT":1.0,"UNIQ_COUNT":2,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"JOHNSON","MAX_FEAT_DESC":"Smith"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","REC_COUNT":1,"REC_PCT":0.5,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"4/8/1983","MAX_FEAT_DESC":"4/8/1983"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"GENDER","USAGE_TYPE":"","REC_COUNT":2,"REC_PCT":1.0,"UNIQ_COUNT":2,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"F","MAX_FEAT_DESC":"M"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"","REC_COUNT":2,"REC_PCT":1.0,"UNIQ_COUNT":1,"UNIQ_PCT":0.5,"MIN_FEAT_DESC":"772 Armstrong RD Delhi LA 71232","MAX_FEAT_DESC":"772 Armstrong RD Delhi LA 71232"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"","REC_COUNT":2,"REC_PCT":1.0,"UNIQ_COUNT":1,"UNIQ_PCT":0.5,"MIN_FEAT_DESC":"225-671-0796","MAX_FEAT_DESC":"225-671-0796"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"SSN","USAGE_TYPE":"","REC_COUNT":1,"REC_PCT":0.5,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"053-39-3251","MAX_FEAT_DESC":"053-39-3251"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"LOGIN_ID","USAGE_TYPE":"","REC_COUNT":1,"REC_PCT":0.5,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"flavorh","MAX_FEAT_DESC":"flavorh"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"ACCT_NUM","USAGE_TYPE":"CC","REC_COUNT":1,"REC_PCT":0.5,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"5534202208773608","MAX_FEAT_DESC":"5534202208773608"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","REC_COUNT":12,"REC_PCT":6.0,"UNIQ_COUNT":12,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"JNSN","MAX_FEAT_DESC":"SM0|POST=71232"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","REC_COUNT":4,"REC_PCT":2.0,"UNIQ_COUNT":2,"UNIQ_PCT":0.5,"MIN_FEAT_DESC":"772|ARMSTRNK||71232","MAX_FEAT_DESC":"772|ARMSTRNK||TL"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"ID_KEY","USAGE_TYPE":"","REC_COUNT":2,"REC_PCT":1.0,"UNIQ_COUNT":2,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"ACCT_NUM=5534202208773608","MAX_FEAT_DESC":"SSN=053-39-3251"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"PHONE_KEY","USAGE_TYPE":"","REC_COUNT":2,"REC_PCT":1.0,"UNIQ_COUNT":1,"UNIQ_PCT":0.5,"MIN_FEAT_DESC":"2256710796","MAX_FEAT_DESC":"2256710796"},{"DSRC_CODE":"EXAMPLE_DATA_SOURCE","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"SEARCH_KEY","USAGE_TYPE":"","REC_COUNT":2,"REC_PCT":1.0,"UNIQ_COUNT":2,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"LOGIN_ID:FLAVORH|","MAX_FEAT_DESC":"SSN:3251|80804|"}]
}

func ExampleG2DiagnosticServer_GetPhysicalCores() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetPhysicalCoresRequest{}
	result, err := g2diagnostic.GetPhysicalCores(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_GetRelationshipDetails() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetRelationshipDetailsRequest{
		RelationshipID:          int64(1),
		IncludeInternalFeatures: 1,
	}
	result, err := g2diagnostic.GetRelationshipDetails(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: [{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME","FEAT_DESC":"JOHNSON"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"DOB","FEAT_DESC":"4/8/1983"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"GENDER","FEAT_DESC":"F"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"ADDRESS","FEAT_DESC":"772 Armstrong RD Delhi LA 71232"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"PHONE","FEAT_DESC":"225-671-0796"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"SSN","FEAT_DESC":"053-39-3251"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"LOGIN_ID","FEAT_DESC":"flavorh"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"ACCT_NUM","FEAT_DESC":"5534202208773608"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JNSN"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JNSN|DOB=80804"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JNSN|POST=71232"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JNSN|SSN=3251"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"772|ARMSTRNK||TL"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"772|ARMSTRNK||71232"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"ID_KEY","FEAT_DESC":"ACCT_NUM=5534202208773608"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"ID_KEY","FEAT_DESC":"SSN=053-39-3251"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"PHONE_KEY","FEAT_DESC":"2256710796"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"SEARCH_KEY","FEAT_DESC":"LOGIN_ID:FLAVORH|"},{"RES_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"SEARCH_KEY","FEAT_DESC":"SSN:3251|80804|"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME","FEAT_DESC":"Smith"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"GENDER","FEAT_DESC":"M"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"ADDRESS","FEAT_DESC":"772 Armstrong RD Delhi LA 71232"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"PHONE","FEAT_DESC":"225-671-0796"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"SM0|POST=71232"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"SM0|ADDRESS.CITY_STD=TL"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"SM0|PHONE.PHONE_LAST_5=10796"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"SM0"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"772|ARMSTRNK||TL"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"772|ARMSTRNK||71232"},{"RES_ENT_ID":2,"ERRULE_CODE":"MFF","MATCH_KEY":"+ADDRESS+PHONE-GENDER","FTYPE_CODE":"PHONE_KEY","FEAT_DESC":"2256710796"}]
}

func ExampleG2DiagnosticServer_GetResolutionStatistics() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetResolutionStatisticsRequest{}
	result, err := g2diagnostic.GetResolutionStatistics(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result)
	// Output: [{"MATCH_LEVEL":3,"MATCH_KEY":"+ADDRESS+PHONE-GENDER","RAW_MATCH_KEYS":[{"MATCH_KEY":"+ADDRESS+PHONE-GENDER"}],"ERRULE_ID":200,"ERRULE_CODE":"MFF","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":1,"MAX_RES_ENT_ID":2,"MIN_RES_REL_ID":1,"MAX_RES_REL_ID":1}]
}

func ExampleG2DiagnosticServer_GetTotalSystemMemory() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.GetTotalSystemMemoryRequest{}
	result, err := g2diagnostic.GetTotalSystemMemory(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.Result > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_Init() {
	ctx := context.TODO()
	g2diagnostic := &G2DiagnosticServer{}
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	request := &pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	result, err := g2diagnostic.Init(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2DiagnosticServer_InitWithConfigID() {
	ctx := context.TODO()
	g2diagnostic := &G2DiagnosticServer{}
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	request := &pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   int64(1),
		VerboseLogging: int32(0),
	}
	result, err := g2diagnostic.InitWithConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2DiagnosticServer_Reinit() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.ReinitRequest{
		InitConfigID: int64(testhelpers.TestConfigDataId),
	}
	result, err := g2diagnostic.Reinit(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2DiagnosticServer_Destroy() {
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &pb.DestroyRequest{}
	result, err := g2diagnostic.Destroy(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

package g2diagnosticserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go-base/g2config"
	"github.com/senzing/g2-sdk-go-base/g2configmgr"
	"github.com/senzing/g2-sdk-go-base/g2engine"
	"github.com/senzing/g2-sdk-go/g2error"
	g2configmgrpb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-common/truthset"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/serve-grpc/g2configmgrserver"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2configmgrServerSingleton *g2configmgrserver.G2ConfigmgrServer
	g2diagnosticTestSingleton  *G2DiagnosticServer
	localLogger                logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return g2error.Cast(localLogger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) G2DiagnosticServer {
	if g2diagnosticTestSingleton == nil {
		g2diagnosticTestSingleton = &G2DiagnosticServer{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkG2diagnostic().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *g2diagnosticTestSingleton
}

func getG2DiagnosticServer(ctx context.Context) G2DiagnosticServer {
	if g2diagnosticTestSingleton == nil {
		g2diagnosticTestSingleton = &G2DiagnosticServer{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkG2diagnostic().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *g2diagnosticTestSingleton
}

func getG2ConfigmgrServer(ctx context.Context) g2configmgrserver.G2ConfigmgrServer {
	if g2configmgrServerSingleton == nil {
		g2configmgrServerSingleton = &g2configmgrserver.G2ConfigmgrServer{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = g2configmgrserver.GetSdkG2configmgr().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *g2configmgrServerSingleton
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
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

func expectError(test *testing.T, ctx context.Context, g2diagnostic G2DiagnosticServer, err error, messageId string) {
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
		if g2error.Is(err, g2error.G2Unrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2Retryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2BadUserInput) {
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

func setupSenzingConfig(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	now := time.Now()

	aG2config := &g2config.G2config{}
	err := aG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
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

	configStr, err := aG2config.Save(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	err = aG2config.Close(ctx, configHandle)
	if err != nil {
		return createError(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return createError(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &g2configmgr.G2configmgr{}
	err = aG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return createError(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return createError(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func setupAddRecords(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {

	aG2engine := &g2engine.G2engine{}
	err := aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5916, err)
	}

	testRecordIds := []string{"1001", "1002", "1003", "1004", "1005", "1039", "1040"}
	for _, testRecordId := range testRecordIds {
		testRecord := truthset.CustomerRecords[testRecordId]
		err := aG2engine.AddRecord(ctx, testRecord.DataSource, testRecord.Id, testRecord.Json, "G2Diagnostic_test")
		if err != nil {
			return createError(5917, err)
		}
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return createError(5918, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	aG2engine := &g2engine.G2engine{}
	err := aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5903, err)
	}

	err = aG2engine.PurgeRepository(ctx)
	if err != nil {
		return createError(5904, err)
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return createError(5905, err)
	}
	return err
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	moduleName := "Test module name"
	verboseLogging := 0
	localLogger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages)
	if err != nil {
		panic(err)
	}

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
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

	// Add records.

	err = setupAddRecords(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5922, err)
	}

	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestBuildSimpleSystemConfigurationJsonUsingEnvVars(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
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
	request := &g2pb.CheckDBPerfRequest{
		SecondsToRun: int32(1),
	}
	response, err := g2diagnostic.CheckDBPerf(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_EntityListBySize(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)

	// GetEntityListBySize()
	requestForGetEntityListBySizeRequest := &g2pb.GetEntityListBySizeRequest{
		EntitySize: int32(10),
	}
	responseFromGetEntityListBySize, err := g2diagnostic.GetEntityListBySize(ctx, requestForGetEntityListBySizeRequest)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, responseFromGetEntityListBySize)

	// FetchNextEntityBySize()
	requestFetchNextEntityBySize := &g2pb.FetchNextEntityBySizeRequest{
		EntityListBySizeHandle: responseFromGetEntityListBySize.GetResult(),
	}
	responseFromFetchNextEntityBySize, err2 := g2diagnostic.FetchNextEntityBySize(ctx, requestFetchNextEntityBySize)
	testError(test, ctx, g2diagnostic, err2)
	printActual(test, responseFromFetchNextEntityBySize)

	// CloseEntityListBySize()
	requestCloseEntityListBySize := &g2pb.CloseEntityListBySizeRequest{
		EntityListBySizeHandle: responseFromGetEntityListBySize.GetResult(),
	}
	responseFromCloseEntityListBySize, err3 := g2diagnostic.CloseEntityListBySize(ctx, requestCloseEntityListBySize)
	testError(test, ctx, g2diagnostic, err3)
	printActual(test, responseFromCloseEntityListBySize)

}

func TestG2diagnosticserver_FindEntitiesByFeatureIDs(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.FindEntitiesByFeatureIDsRequest{
		Features: "{\"ENTITY_ID\":1,\"LIB_FEAT_IDS\":[1,3,4]}",
	}
	response, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetAvailableMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetAvailableMemoryRequest{}
	response, err := g2diagnostic.GetAvailableMemory(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetDataSourceCounts(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetDataSourceCountsRequest{}
	response, err := g2diagnostic.GetDataSourceCounts(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetDBInfo(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetDBInfoRequest{}
	response, err := g2diagnostic.GetDBInfo(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetEntityDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetEntityDetailsRequest{
		EntityID:                int64(1),
		IncludeInternalFeatures: 1,
	}
	response, err := g2diagnostic.GetEntityDetails(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetEntityResume(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetEntityResumeRequest{
		EntityID: int64(1),
	}
	response, err := g2diagnostic.GetEntityResume(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetEntitySizeBreakdown(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetEntitySizeBreakdownRequest{
		MinimumEntitySize:       int32(1),
		IncludeInternalFeatures: int32(1),
	}
	response, err := g2diagnostic.GetEntitySizeBreakdown(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetFeature(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetFeatureRequest{
		LibFeatID: int64(1),
	}
	response, err := g2diagnostic.GetFeature(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetGenericFeatures(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetGenericFeaturesRequest{
		FeatureType:           "PHONE",
		MaximumEstimatedCount: 10,
	}
	response, err := g2diagnostic.GetGenericFeatures(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetLogicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetLogicalCoresRequest{}
	response, err := g2diagnostic.GetLogicalCores(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetMappingStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: 1,
	}
	response, err := g2diagnostic.GetMappingStatistics(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetPhysicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetPhysicalCoresRequest{}
	response, err := g2diagnostic.GetPhysicalCores(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetRelationshipDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetRelationshipDetailsRequest{
		RelationshipID:          int64(1),
		IncludeInternalFeatures: 1,
	}
	response, err := g2diagnostic.GetRelationshipDetails(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetResolutionStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetResolutionStatisticsRequest{}
	response, err := g2diagnostic.GetResolutionStatistics(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_GetTotalSystemMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.GetTotalSystemMemoryRequest{}
	response, err := g2diagnostic.GetTotalSystemMemory(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_Init(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	response, err := g2diagnostic.Init(ctx, request)
	expectError(test, ctx, g2diagnostic, err, "senzing-60134002")
	printActual(test, response)
}

func TestG2diagnosticserver_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   int64(1),
		VerboseLogging: int32(0),
	}
	response, err := g2diagnostic.InitWithConfigID(ctx, request)
	expectError(test, ctx, g2diagnostic, err, "senzing-60134003")
	printActual(test, response)
}

func TestG2diagnosticserver_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	g2configmgr := getG2ConfigmgrServer(ctx)
	getDefaultConfigIDRequest := &g2configmgrpb.GetDefaultConfigIDRequest{}
	getDefaultConfigIDResponse, err := g2configmgr.GetDefaultConfigID(ctx, getDefaultConfigIDRequest)
	testError(test, ctx, g2diagnostic, err)
	request := &g2pb.ReinitRequest{
		InitConfigID: getDefaultConfigIDResponse.ConfigID,
	}
	response, err := g2diagnostic.Reinit(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &g2pb.DestroyRequest{}
	response, err := g2diagnostic.Destroy(ctx, request)
	expectError(test, ctx, g2diagnostic, err, "senzing-60134001")
	printActual(test, response)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2DiagnosticServer_CheckDBPerf() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.CheckDBPerfRequest{
		SecondsToRun: int32(1),
	}
	response, err := g2diagnostic.CheckDBPerf(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 25))
	// Output: {"numRecordsInserted":...
}

// func ExampleG2diagnosticImpl_CloseEntityListBySize() {

// func ExampleG2diagnosticImpl_FetchNextEntityBySize() {

func ExampleG2DiagnosticServer_FindEntitiesByFeatureIDs() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.FindEntitiesByFeatureIDsRequest{
		Features: "{\"ENTITY_ID\":1,\"LIB_FEAT_IDS\":[1,3,4]}",
	}
	response, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: []
}

func ExampleG2DiagnosticServer_GetAvailableMemory() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetAvailableMemoryRequest{}
	response, err := g2diagnostic.GetAvailableMemory(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_GetDataSourceCounts() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetDataSourceCountsRequest{}
	response, err := g2diagnostic.GetDataSourceCounts(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: [{"DSRC_ID":1001,"DSRC_CODE":"CUSTOMERS","ETYPE_ID":3,"ETYPE_CODE":"GENERIC","OBS_ENT_COUNT":7,"DSRC_RECORD_COUNT":7}]
}

func ExampleG2DiagnosticServer_GetDBInfo() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetDBInfoRequest{}
	response, err := g2diagnostic.GetDBInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 52))
	// Output: {"Hybrid Mode":false,"Database Details":[{"Name":...
}

func ExampleG2DiagnosticServer_GetEntityDetails() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetEntityDetailsRequest{
		EntityID:                int64(1),
		IncludeInternalFeatures: 1,
	}
	response, err := g2diagnostic.GetEntityDetails(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: [{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"Robert Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","FEAT_DESC":"12/11/1978"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"MAILING","FEAT_DESC":"123 Main Street, Las Vegas NV 89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"HOME","FEAT_DESC":"702-919-1300"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"EMAIL","USAGE_TYPE":"","FEAT_DESC":"bsmith@work.com"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMYY_HASH=1278"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|ADDRESS.CITY_STD=LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|POST=89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|PHONE.PHONE_LAST_5=91300"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"123|MN||LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"123|MN||89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"PHONE_KEY","USAGE_TYPE":"","FEAT_DESC":"7029191300"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"EMAIL_KEY","USAGE_TYPE":"","FEAT_DESC":"bsmith@WORK.COM"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"Bob Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","FEAT_DESC":"11/12/1978"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"HOME","FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"MOBILE","FEAT_DESC":"702-919-1300"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|ADDRESS.CITY_STD=LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|PHONE.PHONE_LAST_5=91300"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMYY_HASH=1178"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|POST=89111"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|ADDRESS.CITY_STD=LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|POST=89111"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB.MMYY_HASH=1178"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|PHONE.PHONE_LAST_5=91300"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"1515|ATL||89111"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"1515|ATL||LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"PHONE_KEY","USAGE_TYPE":"","FEAT_DESC":"7029191300"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"Bob J Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","FEAT_DESC":"12/11/1978"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"No","FTYPE_CODE":"EMAIL","USAGE_TYPE":"","FEAT_DESC":"bsmith@work.com"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMYY_HASH=1278"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|J|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|J|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|J|SM0|DOB.MMYY_HASH=1278"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB.MMYY_HASH=1278"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|J|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"EMAIL_KEY","USAGE_TYPE":"","FEAT_DESC":"bsmith@WORK.COM"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"B Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","FEAT_DESC":"11/12/1979"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"HOME","FEAT_DESC":"1515 Adela Ln Las Vegas NV 89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"EMAIL","USAGE_TYPE":"","FEAT_DESC":"bsmith@work.com"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|DOB.MMYY_HASH=1179"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|ADDRESS.CITY_STD=LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|POST=89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"1515|ATL||LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"1515|ATL||89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"EMAIL_KEY","USAGE_TYPE":"","FEAT_DESC":"bsmith@WORK.COM"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"Robbie Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"MAILING","FEAT_DESC":"123 E Main St Henderson NV 89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"No","FTYPE_CODE":"DRLIC","USAGE_TYPE":"","FEAT_DESC":"112233 NV"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|POST=89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|ADDRESS.CITY_STD=HNTRSN"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RB|SM0|POST=89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RB|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RB|SM0|ADDRESS.CITY_STD=HNTRSN"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"123|MN||89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"123|MN||HNTRSN"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"ID_KEY","USAGE_TYPE":"","FEAT_DESC":"DRLIC=112233"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"}]
}

// func ExampleG2diagnosticImpl_GetEntityListBySize() {

func ExampleG2DiagnosticServer_GetEntityResume() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetEntityResumeRequest{
		EntityID: int64(1),
	}
	response, err := g2diagnostic.GetEntityResume(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: [{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","ENT_SRC_DESC":"Robert Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"Robert\",\"DATE_OF_BIRTH\":\"12/11/1978\",\"ADDR_TYPE\":\"MAILING\",\"ADDR_LINE1\":\"123 Main Street, Las Vegas NV 89132\",\"PHONE_TYPE\":\"HOME\",\"PHONE_NUMBER\":\"702-919-1300\",\"EMAIL_ADDRESS\":\"bsmith@work.com\",\"DATE\":\"1/2/18\",\"STATUS\":\"Active\",\"AMOUNT\":\"100\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1001\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","ENT_SRC_DESC":"Bob Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"Bob\",\"DATE_OF_BIRTH\":\"11/12/1978\",\"ADDR_TYPE\":\"HOME\",\"ADDR_LINE1\":\"1515 Adela Lane\",\"ADDR_CITY\":\"Las Vegas\",\"ADDR_STATE\":\"NV\",\"ADDR_POSTAL_CODE\":\"89111\",\"PHONE_TYPE\":\"MOBILE\",\"PHONE_NUMBER\":\"702-919-1300\",\"DATE\":\"3/10/17\",\"STATUS\":\"Inactive\",\"AMOUNT\":\"200\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1002\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","ENT_SRC_DESC":"Bob J Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"Bob\",\"PRIMARY_NAME_MIDDLE\":\"J\",\"DATE_OF_BIRTH\":\"12/11/1978\",\"EMAIL_ADDRESS\":\"bsmith@work.com\",\"DATE\":\"4/9/16\",\"STATUS\":\"Inactive\",\"AMOUNT\":\"300\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1003\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","ENT_SRC_DESC":"B Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"B\",\"DATE_OF_BIRTH\":\"11/12/1979\",\"ADDR_TYPE\":\"HOME\",\"ADDR_LINE1\":\"1515 Adela Ln\",\"ADDR_CITY\":\"Las Vegas\",\"ADDR_STATE\":\"NV\",\"ADDR_POSTAL_CODE\":\"89132\",\"EMAIL_ADDRESS\":\"bsmith@work.com\",\"DATE\":\"1/5/15\",\"STATUS\":\"Inactive\",\"AMOUNT\":\"400\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1004\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","ENT_SRC_DESC":"Robbie Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"Robbie\",\"DRIVERS_LICENSE_NUMBER\":\"112233\",\"DRIVERS_LICENSE_STATE\":\"NV\",\"ADDR_TYPE\":\"MAILING\",\"ADDR_LINE1\":\"123 E Main St\",\"ADDR_CITY\":\"Henderson\",\"ADDR_STATE\":\"NV\",\"ADDR_POSTAL_CODE\":\"89132\",\"DATE\":\"7/16/19\",\"STATUS\":\"Active\",\"AMOUNT\":\"500\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1005\"}"}]
}

func ExampleG2DiagnosticServer_GetEntitySizeBreakdown() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetEntitySizeBreakdownRequest{
		MinimumEntitySize:       int32(1),
		IncludeInternalFeatures: int32(1),
	}
	response, err := g2diagnostic.GetEntitySizeBreakdown(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: [{"ENTITY_SIZE": 5,"ENTITY_COUNT": 1,"NAME": 5.00,"DOB": 3.00,"ADDRESS": 4.00,"PHONE": 2.00,"DRLIC": 1.00,"EMAIL": 1.00,"NAME_KEY": 31.00,"ADDR_KEY": 6.00,"ID_KEY": 1.00,"PHONE_KEY": 1.00,"RECORD_TYPE": 1.00,"EMAIL_KEY": 1.00,"MIN_RES_ENT_ID": 1,"MAX_RES_ENT_ID": 1},{"ENTITY_SIZE": 1,"ENTITY_COUNT": 2,"NAME": 1.00,"DOB": 1.00,"GENDER": 0.50,"ADDRESS": 1.00,"NAME_KEY": 6.00,"ADDR_KEY": 2.00,"RECORD_TYPE": 1.00,"MIN_RES_ENT_ID": 6,"MAX_RES_ENT_ID": 7}]
}

func ExampleG2DiagnosticServer_GetFeature() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetFeatureRequest{
		LibFeatID: int64(1),
	}
	response, err := g2diagnostic.GetFeature(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"LIB_FEAT_ID":1,"FTYPE_CODE":"NAME","ELEMENTS":[{"FELEM_CODE":"TOKENIZED_NM","FELEM_VALUE":"ROBERT|SMITH"},{"FELEM_CODE":"CATEGORY","FELEM_VALUE":"PERSON"},{"FELEM_CODE":"CULTURE","FELEM_VALUE":"ANGLO"},{"FELEM_CODE":"GIVEN_NAME","FELEM_VALUE":"Robert"},{"FELEM_CODE":"SUR_NAME","FELEM_VALUE":"Smith"},{"FELEM_CODE":"FULL_NAME","FELEM_VALUE":"Robert Smith"}]}
}

func ExampleG2DiagnosticServer_GetGenericFeatures() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetGenericFeaturesRequest{
		FeatureType:           "PHONE",
		MaximumEstimatedCount: 10,
	}
	response, err := g2diagnostic.GetGenericFeatures(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: []
}

func ExampleG2DiagnosticServer_GetLogicalCores() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetLogicalCoresRequest{}
	response, err := g2diagnostic.GetLogicalCores(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_GetMappingStatistics() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: 1,
	}
	response, err := g2diagnostic.GetMappingStatistics(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: [{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","REC_COUNT":7,"REC_PCT":1.0,"UNIQ_COUNT":6,"UNIQ_PCT":0.8571428571428571,"MIN_FEAT_DESC":"B Smith","MAX_FEAT_DESC":"Robert Smith"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","REC_COUNT":6,"REC_PCT":0.8571428571428571,"UNIQ_COUNT":5,"UNIQ_PCT":0.8333333333333334,"MIN_FEAT_DESC":"10/10/70","MAX_FEAT_DESC":"3/15/90"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"GENDER","USAGE_TYPE":"","REC_COUNT":1,"REC_PCT":0.14285714285714285,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"M","MAX_FEAT_DESC":"M"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"HOME","REC_COUNT":4,"REC_PCT":0.5714285714285714,"UNIQ_COUNT":3,"UNIQ_PCT":0.75,"MIN_FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","MAX_FEAT_DESC":"3212 W. 32nd St Palm Harbor, FL 60527"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"MAILING","REC_COUNT":2,"REC_PCT":0.2857142857142857,"UNIQ_COUNT":2,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"123 E Main St Henderson NV 89132","MAX_FEAT_DESC":"123 Main Street, Las Vegas NV 89132"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"HOME","REC_COUNT":1,"REC_PCT":0.14285714285714285,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"702-919-1300","MAX_FEAT_DESC":"702-919-1300"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"MOBILE","REC_COUNT":1,"REC_PCT":0.14285714285714285,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"702-919-1300","MAX_FEAT_DESC":"702-919-1300"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"DRLIC","USAGE_TYPE":"","REC_COUNT":1,"REC_PCT":0.14285714285714285,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"112233 NV","MAX_FEAT_DESC":"112233 NV"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"EMAIL","USAGE_TYPE":"","REC_COUNT":3,"REC_PCT":0.42857142857142855,"UNIQ_COUNT":1,"UNIQ_PCT":0.3333333333333333,"MIN_FEAT_DESC":"bsmith@work.com","MAX_FEAT_DESC":"bsmith@work.com"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","REC_COUNT":57,"REC_PCT":8.142857142857142,"UNIQ_COUNT":40,"UNIQ_PCT":0.7017543859649122,"MIN_FEAT_DESC":"BB|J|SM0","MAX_FEAT_DESC":"RB|SM0|POST=89132"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","REC_COUNT":12,"REC_PCT":1.7142857142857142,"UNIQ_COUNT":8,"UNIQ_PCT":0.6666666666666666,"MIN_FEAT_DESC":"123|MN||89132","MAX_FEAT_DESC":"3212|NT||PLM HRBR"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"ID_KEY","USAGE_TYPE":"","REC_COUNT":1,"REC_PCT":0.14285714285714285,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"DRLIC=112233","MAX_FEAT_DESC":"DRLIC=112233"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"PHONE_KEY","USAGE_TYPE":"","REC_COUNT":2,"REC_PCT":0.2857142857142857,"UNIQ_COUNT":1,"UNIQ_PCT":0.5,"MIN_FEAT_DESC":"7029191300","MAX_FEAT_DESC":"7029191300"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","REC_COUNT":7,"REC_PCT":1.0,"UNIQ_COUNT":1,"UNIQ_PCT":0.14285714285714285,"MIN_FEAT_DESC":"PERSON","MAX_FEAT_DESC":"PERSON"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"EMAIL_KEY","USAGE_TYPE":"","REC_COUNT":3,"REC_PCT":0.42857142857142855,"UNIQ_COUNT":1,"UNIQ_PCT":0.3333333333333333,"MIN_FEAT_DESC":"bsmith@WORK.COM","MAX_FEAT_DESC":"bsmith@WORK.COM"}]
}

func ExampleG2DiagnosticServer_GetPhysicalCores() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetPhysicalCoresRequest{}
	response, err := g2diagnostic.GetPhysicalCores(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_GetRelationshipDetails() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetRelationshipDetailsRequest{
		RelationshipID:          int64(1),
		IncludeInternalFeatures: 1,
	}
	response, err := g2diagnostic.GetRelationshipDetails(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: [{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME","FEAT_DESC":"John Smith"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"DOB","FEAT_DESC":"10/10/70"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"GENDER","FEAT_DESC":"M"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDRESS","FEAT_DESC":"3212 W. 32nd St Palm Harbor, FL 60527"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|ADDRESS.CITY_STD=PLM HRBR"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB.MMDD_HASH=1010"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB=71010"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|POST=60527"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB.MMYY_HASH=1070"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"3212|NT||60527"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"3212|NT||PLM HRBR"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"RECORD_TYPE","FEAT_DESC":"PERSON"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME","FEAT_DESC":"John Smith"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"DOB","FEAT_DESC":"3/15/90"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDRESS","FEAT_DESC":"3212 W. 32nd St Palm Harbor, FL 60527"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|ADDRESS.CITY_STD=PLM HRBR"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|POST=60527"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB.MMYY_HASH=0390"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB.MMDD_HASH=1503"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB=91503"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"3212|NT||60527"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"3212|NT||PLM HRBR"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"RECORD_TYPE","FEAT_DESC":"PERSON"}]
}

func ExampleG2DiagnosticServer_GetResolutionStatistics() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetResolutionStatisticsRequest{}
	response, err := g2diagnostic.GetResolutionStatistics(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: [{"MATCH_LEVEL":1,"MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","RAW_MATCH_KEYS":[{"MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME"}],"ERRULE_ID":110,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":1,"MAX_RES_ENT_ID":1,"MIN_RES_REL_ID":0,"MAX_RES_REL_ID":0},{"MATCH_LEVEL":1,"MATCH_KEY":"+NAME+DOB+EMAIL","RAW_MATCH_KEYS":[{"MATCH_KEY":"+NAME+DOB+EMAIL"}],"ERRULE_ID":120,"ERRULE_CODE":"SF1_PNAME_CSTAB","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":1,"MAX_RES_ENT_ID":1,"MIN_RES_REL_ID":0,"MAX_RES_REL_ID":0},{"MATCH_LEVEL":1,"MATCH_KEY":"+NAME+DOB+PHONE","RAW_MATCH_KEYS":[{"MATCH_KEY":"+NAME+DOB+PHONE"}],"ERRULE_ID":160,"ERRULE_CODE":"CNAME_CFF_CEXCL","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":1,"MAX_RES_ENT_ID":1,"MIN_RES_REL_ID":0,"MAX_RES_REL_ID":0},{"MATCH_LEVEL":1,"MATCH_KEY":"+NAME+ADDRESS","RAW_MATCH_KEYS":[{"MATCH_KEY":"+NAME+ADDRESS"}],"ERRULE_ID":162,"ERRULE_CODE":"CNAME_CFF","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":1,"MAX_RES_ENT_ID":1,"MIN_RES_REL_ID":0,"MAX_RES_REL_ID":0},{"MATCH_LEVEL":2,"MATCH_KEY":"+NAME+ADDRESS-DOB","RAW_MATCH_KEYS":[{"MATCH_KEY":"+NAME+ADDRESS-DOB"}],"ERRULE_ID":164,"ERRULE_CODE":"CNAME_CFF_DEXCL","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":6,"MAX_RES_ENT_ID":7,"MIN_RES_REL_ID":1,"MAX_RES_REL_ID":1}]
}

func ExampleG2DiagnosticServer_GetTotalSystemMemory() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetTotalSystemMemoryRequest{}
	response, err := g2diagnostic.GetTotalSystemMemory(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_Init() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := &G2DiagnosticServer{}
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	response, err := g2diagnostic.Init(ctx, request)
	if err != nil {
		// This should produce a "senzing-60134002" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2DiagnosticServer_InitWithConfigID() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := &G2DiagnosticServer{}
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	request := &g2pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   int64(1),
		VerboseLogging: int32(0),
	}
	response, err := g2diagnostic.InitWithConfigID(ctx, request)
	if err != nil {
		// This should produce a "senzing-60134003" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2DiagnosticServer_Reinit() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	g2configmgr := getG2ConfigmgrServer(ctx)
	getDefaultConfigIDRequest := &g2configmgrpb.GetDefaultConfigIDRequest{}
	getDefaultConfigIDResponse, err := g2configmgr.GetDefaultConfigID(ctx, getDefaultConfigIDRequest)
	if err != nil {
		fmt.Println(err)
	}
	request := &g2pb.ReinitRequest{
		InitConfigID: getDefaultConfigIDResponse.ConfigID,
	}
	_, err = g2diagnostic.Reinit(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2DiagnosticServer_Destroy() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.DestroyRequest{}
	response, err := g2diagnostic.Destroy(ctx, request)
	if err != nil {
		// This should produce a "senzing-60134001" error.
	}
	fmt.Println(response)
	// Output:
}

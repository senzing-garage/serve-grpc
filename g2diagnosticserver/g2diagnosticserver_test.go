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

// func getG2DiagnosticServer(ctx context.Context) G2DiagnosticServer {
// 	g2diagnostic := &G2diagnosticImpl{}
// 	moduleName := "Test module name"
// 	verboseLogging := 0 // 0 for no Senzing logging; 1 for logging
// 	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	g2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
// 	return g2diagnostic
// }

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

func TestG2diagnosticImpl_BuildSimpleSystemConfigurationJson(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Test interface functions - names begin with "Test"
// ----------------------------------------------------------------------------

func TestG2diagnosticImpl_CheckDBPerf(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.CheckDBPerfRequest{
		SecondsToRun: int32(1),
	}
	actual, err := g2diagnostic.CheckDBPerf(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

//func TestEntityListBySize(test *testing.T) {
//	ctx := context.TODO()
//	g2diagnostic, _ := getTestObject(ctx)
//	request := &pb.GetEntityListBySizeRequest{
//		EntitySize: int32(10),
//	}
//	actual, err := g2diagnostic.GetEntityListBySize(ctx, request)
//	testError(test, ctx, err)
//	test.Log("Actual:", actual)
//
//	entityListBySizeHandle := actual.Result
//	request2 := pb.FetchNextEntityBySizeRequest{
//		EntityListBySizeHandle: entityListBySizeHandle,
//	}
//	actual2, err2 := g2diagnostic.FetchNextEntityBySize(ctx, request2)
//	testError(test, ctx, err2)
//	test.Log("Actual:", actual2)
//
//	request3 := pb.CloseEntityListBySizeRequest{
//		EntityListBySizeHandle: entityListBySizeHandle,
//	}
//	actual3, err3 := g2diagnostic.CloseEntityListBySize(ctx, request3)
//	testError(test, ctx, err3)
//	test.Log("Actual:", actual3)
//}

func TestG2diagnosticImpl_FindEntitiesByFeatureIDs(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.FindEntitiesByFeatureIDsRequest{
		Features: "{\"ENTITY_ID\":1,\"LIB_FEAT_IDS\":[1,3,4]}",
	}
	actual, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetAvailableMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetAvailableMemoryRequest{}
	actual, err := g2diagnostic.GetAvailableMemory(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetDataSourceCounts(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetDataSourceCountsRequest{}
	actual, err := g2diagnostic.GetDataSourceCounts(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetDBInfo(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetDBInfoRequest{}
	actual, err := g2diagnostic.GetDBInfo(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetEntityDetails(test *testing.T) {
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

func TestG2diagnosticImpl_GetEntityResume(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetEntityResumeRequest{
		EntityID: int64(1),
	}
	actual, err := g2diagnostic.GetEntityResume(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetEntitySizeBreakdown(test *testing.T) {
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

func TestG2diagnosticImpl_GetFeature(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetFeatureRequest{
		LibFeatID: int64(1),
	}
	actual, err := g2diagnostic.GetFeature(ctx, request)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetGenericFeatures(test *testing.T) {
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

func TestG2diagnosticImpl_GetLogicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetLogicalCoresRequest{}
	actual, err := g2diagnostic.GetLogicalCores(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetMappingStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: 1,
	}
	actual, err := g2diagnostic.GetMappingStatistics(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetPhysicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetPhysicalCoresRequest{}
	actual, err := g2diagnostic.GetPhysicalCores(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetRelationshipDetails(test *testing.T) {
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

func TestG2diagnosticImpl_GetResolutionStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetResolutionStatisticsRequest{}
	actual, err := g2diagnostic.GetResolutionStatistics(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_GetTotalSystemMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.GetTotalSystemMemoryRequest{}
	actual, err := g2diagnostic.GetTotalSystemMemory(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_Init(test *testing.T) {
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

func TestG2diagnosticImpl_InitWithConfigID(test *testing.T) {
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

func TestG2diagnosticImpl_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.ReinitRequest{
		InitConfigID: int64(testhelpers.TestConfigDataId),
	}
	actual, err := g2diagnostic.Reinit(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticImpl_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	request := &pb.DestroyRequest{}
	actual, err := g2diagnostic.Destroy(ctx, request)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

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
	if 1 == 0 {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
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

func TestG2engineserver_BuildSimpleSystemConfigurationJson(test *testing.T) {
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

// PurgeRepository() is first to start with a clean database.
func TestG2engineServer_PurgeRepository(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.PurgeRepositoryRequest{}
	actual, err := g2engine.PurgeRepository(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_AddRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.AddRecordRequest{}
	actual, err := g2engine.AddRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_AddRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.AddRecordWithInfoRequest{}
	actual, err := g2engine.AddRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_AddRecordWithInfoWithReturnedRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.AddRecordWithInfoWithReturnedRecordIDRequest{}
	actual, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_AddRecordWithReturnedRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.AddRecordWithReturnedRecordIDRequest{}
	actual, err := g2engine.AddRecordWithReturnedRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_CheckRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.CheckRecordRequest{}
	actual, err := g2engine.CheckRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

// FAIL:
func TestG2engineServer_ExportJSONEntityReport(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ExportJSONEntityReportRequest{}
	actual, err := g2engine.ExportJSONEntityReport(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_CountRedoRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.CountRedoRecordsRequest{}
	actual, err := g2engine.CountRedoRecords(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ExportConfigAndConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ExportConfigAndConfigIDRequest{}
	actual, err := g2engine.ExportConfigAndConfigID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.PurgeRepositoryRequest{}
	actual, err := g2engine.PurgeRepository(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

//func TestG2engineServer_ExportCSVEntityReport(test *testing.T) {
// ctx := context.TODO()
// g2engine := getTestObject(ctx, test)
// request := &pb.PurgeRepositoryRequest{}
// actual, err := g2engine.PurgeRepository(ctx, request)
// testError(test, ctx, g2engine, err)
// printActual(test, actual)
//}

func TestG2engineServer_FindInterestingEntitiesByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindInterestingEntitiesByEntityIDRequest{}
	actual, err := g2engine.FindInterestingEntitiesByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindInterestingEntitiesByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindInterestingEntitiesByRecordIDRequest{}
	actual, err := g2engine.FindInterestingEntitiesByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindNetworkByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindNetworkByEntityIDRequest{}
	actual, err := g2engine.FindNetworkByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindNetworkByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindNetworkByEntityID_V2Request{}
	actual, err := g2engine.FindNetworkByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindNetworkByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindNetworkByRecordIDRequest{}
	actual, err := g2engine.FindNetworkByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindNetworkByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindNetworkByRecordID_V2Request{}
	actual, err := g2engine.FindNetworkByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathByEntityIDRequest{}
	actual, err := g2engine.FindPathByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathByEntityID_V2Request{}
	actual, err := g2engine.FindPathByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathByRecordIDRequest{}
	actual, err := g2engine.FindPathByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathByRecordID_V2Request{}
	actual, err := g2engine.FindPathByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathExcludingByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathExcludingByEntityIDRequest{}
	actual, err := g2engine.FindPathExcludingByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathExcludingByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathExcludingByEntityID_V2Request{}
	actual, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathExcludingByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathExcludingByRecordIDRequest{}
	actual, err := g2engine.FindPathExcludingByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathExcludingByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathExcludingByRecordID_V2Request{}
	actual, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathIncludingSourceByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathIncludingSourceByEntityIDRequest{}
	actual, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathIncludingSourceByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathIncludingSourceByEntityID_V2Request{}
	actual, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathIncludingSourceByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathIncludingSourceByRecordIDRequest{}
	actual, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_FindPathIncludingSourceByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.FindPathIncludingSourceByRecordID_V2Request{}
	actual, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetActiveConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetActiveConfigIDRequest{}
	actual, err := g2engine.GetActiveConfigID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetEntityByEntityIDRequest{}
	actual, err := g2engine.GetEntityByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetEntityByEntityID_V2Request{}
	actual, err := g2engine.GetEntityByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetEntityByRecordIDRequest{}
	actual, err := g2engine.GetEntityByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetEntityByRecordID_V2Request{}
	actual, err := g2engine.GetEntityByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetRecordRequest{}
	actual, err := g2engine.GetRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetRecord_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetRecord_V2Request{}
	actual, err := g2engine.GetRecord_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetRedoRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetRedoRecordRequest{}
	actual, err := g2engine.GetRedoRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetRepositoryLastModifiedTime(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetRepositoryLastModifiedTimeRequest{}
	actual, err := g2engine.GetRepositoryLastModifiedTime(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetVirtualEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetVirtualEntityByRecordIDRequest{}
	actual, err := g2engine.GetVirtualEntityByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_GetVirtualEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.GetVirtualEntityByRecordID_V2Request{}
	actual, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_HowEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.HowEntityByEntityIDRequest{}
	actual, err := g2engine.HowEntityByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_HowEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.HowEntityByEntityID_V2Request{}
	actual, err := g2engine.HowEntityByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_Init(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.InitRequest{}
	actual, err := g2engine.Init(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.InitWithConfigIDRequest{}
	actual, err := g2engine.InitWithConfigID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_PrimeEngine(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.PrimeEngineRequest{}
	actual, err := g2engine.PrimeEngine(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_Process(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessRequest{}
	actual, err := g2engine.Process(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ProcessRedoRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessRedoRecordRequest{}
	actual, err := g2engine.ProcessRedoRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ProcessRedoRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessRedoRecordWithInfoRequest{}
	actual, err := g2engine.ProcessRedoRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ProcessWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessWithInfoRequest{}
	actual, err := g2engine.ProcessWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ProcessWithResponse(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessWithResponseRequest{}
	actual, err := g2engine.ProcessWithResponse(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ProcessWithResponseResize(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ProcessWithResponseResizeRequest{}
	actual, err := g2engine.ProcessWithResponseResize(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReevaluateEntityRequest{}
	actual, err := g2engine.ReevaluateEntity(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ReevaluateEntityWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReevaluateEntityWithInfoRequest{}
	actual, err := g2engine.ReevaluateEntityWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ReevaluateRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReevaluateRecordRequest{}
	actual, err := g2engine.ReevaluateRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ReevaluateRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReevaluateRecordWithInfoRequest{}
	actual, err := g2engine.ReevaluateRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReinitRequest{}
	actual, err := g2engine.Reinit(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ReplaceRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReplaceRecordRequest{}
	actual, err := g2engine.ReplaceRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_ReplaceRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.ReplaceRecordWithInfoRequest{}
	actual, err := g2engine.ReplaceRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_SearchByAttributes(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.SearchByAttributesRequest{}
	actual, err := g2engine.SearchByAttributes(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_SearchByAttributes_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.SearchByAttributes_V2Request{}
	actual, err := g2engine.SearchByAttributes_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_Stats(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.StatsRequest{}
	actual, err := g2engine.Stats(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_WhyEntities(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntitiesRequest{}
	actual, err := g2engine.WhyEntities(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_WhyEntities_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntities_V2Request{}
	actual, err := g2engine.WhyEntities_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_WhyEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntityByEntityIDRequest{}
	actual, err := g2engine.WhyEntityByEntityID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_WhyEntityByEntityID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntityByEntityID_V2Request{}
	actual, err := g2engine.WhyEntityByEntityID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_WhyEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntityByRecordIDRequest{}
	actual, err := g2engine.WhyEntityByRecordID(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_WhyEntityByRecordID_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyEntityByRecordID_V2Request{}
	actual, err := g2engine.WhyEntityByRecordID_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyRecordsRequest{}
	actual, err := g2engine.WhyRecords(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_WhyRecords_V2(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.WhyRecords_V2Request{}
	actual, err := g2engine.WhyRecords_V2(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_DeleteRecord(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.DeleteRecordRequest{}
	actual, err := g2engine.DeleteRecord(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_DeleteRecordWithInfo(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.DeleteRecordWithInfoRequest{}
	actual, err := g2engine.DeleteRecordWithInfo(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

func TestG2engineServer_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2engine := getTestObject(ctx, test)
	request := &pb.DestroyRequest{}
	actual, err := g2engine.Destroy(ctx, request)
	testError(test, ctx, g2engine, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

// PurgeRepository() is first to start with a clean database.
func ExampleG2EngineServer_PurgeRepository() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.PurgeRepositoryRequest{}
	result, err := g2engine.PurgeRepository(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_AddRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordRequest{}
	result, err := g2engine.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_AddRecord_secondRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordRequest{}
	result, err := g2engine.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_AddRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordWithInfoRequest{}
	result, err := g2engine.AddRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"333","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_AddRecordWithInfoWithReturnedRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordWithInfoWithReturnedRecordIDRequest{}
	result, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_AddRecordWithReturnedRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.AddRecordWithReturnedRecordIDRequest{}
	result, err := g2engine.AddRecordWithReturnedRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: Length of record identifier is 40 hexadecimal characters.
}

func ExampleG2EngineServer_CheckRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.CheckRecordRequest{}
	result, err := g2engine.CheckRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"CHECK_RECORD_RESPONSE":[{"DSRC_CODE":"TEST","RECORD_ID":"111","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","MATCH_KEY":"","ERRULE_CODE":"","ERRULE_ID":0,"CANDIDATE_MATCH":"N","NON_GENERIC_CANDIDATE_MATCH":"N"}]}
}

func ExampleG2EngineServer_CloseExport() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.CloseExportRequest{}
	result, err := g2engine.CloseExport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_CountRedoRecords() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.CountRedoRecordsRequest{}
	result, err := g2engine.CountRedoRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: 0
}

func ExampleG2EngineServer_DeleteRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.DeleteRecordRequest{}
	result, err := g2engine.DeleteRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_DeleteRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.DeleteRecordWithInfoRequest{}
	result, err := g2engine.DeleteRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"333","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ExportCSVEntityReport() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ExportCSVEntityReportRequest{}
	result, err := g2engine.ExportCSVEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: true
}

func ExampleG2EngineServer_ExportConfig() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ExportConfigRequest{}
	result, err := g2engine.ExportConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"G2_CONFIG":{"CFG_ETYPE":[{"ETYPE_ID":...
}

func ExampleG2EngineServer_ExportConfigAndConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ExportConfigAndConfigIDRequest{}
	result, err := g2engine.ExportConfigAndConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: true
}

func ExampleG2EngineServer_ExportJSONEntityReport() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ExportJSONEntityReportRequest{}
	result, err := g2engine.ExportJSONEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: true
}

func ExampleG2EngineServer_FetchNext() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FetchNextRequest{}
	result, err := g2engine.FetchNext(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: true
}

func ExampleG2EngineServer_FindInterestingEntitiesByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindInterestingEntitiesByEntityIDRequest{}
	result, err := g2engine.FindInterestingEntitiesByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_FindInterestingEntitiesByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindInterestingEntitiesByRecordIDRequest{}
	result, err := g2engine.FindInterestingEntitiesByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_FindNetworkByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByEntityIDRequest{}
	result, err := g2engine.FindNetworkByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,...
}

func ExampleG2EngineServer_FindNetworkByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByEntityID_V2Request{}
	result, err := g2engine.FindNetworkByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}},{"RESOLVED_ENTITY":{"ENTITY_ID":3}}]}
}

func ExampleG2EngineServer_FindNetworkByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByRecordIDRequest{}
	result, err := g2engine.FindNetworkByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_FindNetworkByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindNetworkByRecordID_V2Request{}
	result, err := g2engine.FindNetworkByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}},{"RESOLVED_ENTITY":{"ENTITY_ID":3}}]}
}

func ExampleG2EngineServer_FindPathByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByEntityIDRequest{}
	result, err := g2engine.FindPathByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByEntityID_V2Request{}
	result, err := g2engine.FindPathByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByRecordIDRequest{}
	result, err := g2engine.FindPathByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":...
}

func ExampleG2EngineServer_FindPathByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathByRecordID_V2Request{}
	result, err := g2engine.FindPathByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathExcludingByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByEntityIDRequest{}
	result, err := g2engine.FindPathExcludingByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathExcludingByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByEntityID_V2Request{}
	result, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathExcludingByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByRecordIDRequest{}
	result, err := g2engine.FindPathExcludingByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathExcludingByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathExcludingByRecordID_V2Request{}
	result, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathIncludingSourceByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByEntityIDRequest{}
	result, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathIncludingSourceByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByEntityID_V2Request{}
	result, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_FindPathIncludingSourceByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByRecordIDRequest{}
	result, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2EngineServer_FindPathIncludingSourceByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.FindPathIncludingSourceByRecordID_V2Request{}
	result, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_GetActiveConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetActiveConfigIDRequest{}
	result, err := g2engine.GetActiveConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: true
}

func ExampleG2EngineServer_GetEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByEntityIDRequest{}
	result, err := g2engine.GetEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_GetEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByEntityID_V2Request{}
	result, err := g2engine.GetEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2EngineServer_GetEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByRecordIDRequest{}
	result, err := g2engine.GetEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2EngineServer_GetEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetEntityByRecordID_V2Request{}
	result, err := g2engine.GetEntityByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2EngineServer_GetRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRecordRequest{}
	result, err := g2engine.GetRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111","JSON_DATA":{"SOCIAL_HANDLE":"flavorh","DATE_OF_BIRTH":"4/8/1983","ADDR_STATE":"LA","ADDR_POSTAL_CODE":"71232","SSN_NUMBER":"053-39-3251","GENDER":"F","srccode":"MDMPER","CC_ACCOUNT_NUMBER":"5534202208773608","ADDR_CITY":"Delhi","DRIVERS_LICENSE_STATE":"DE","PHONE_NUMBER":"225-671-0796","NAME_LAST":"JOHNSON","entityid":"284430058","ADDR_LINE1":"772 Armstrong RD","DATA_SOURCE":"TEST","ENTITY_TYPE":"TEST","DSRC_ACTION":"A","RECORD_ID":"111"}}
}

func ExampleG2EngineServer_GetRecord_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRecord_V2Request{}
	result, err := g2engine.GetRecord_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111"}
}

func ExampleG2EngineServer_GetRedoRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRedoRecordRequest{}
	result, err := g2engine.GetRedoRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetRepositoryLastModifiedTimeRequest{}
	result, err := g2engine.GetRepositoryLastModifiedTime(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: true
}

func ExampleG2EngineServer_GetVirtualEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetVirtualEntityByRecordIDRequest{}
	result, err := g2engine.GetVirtualEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_GetVirtualEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.GetVirtualEntityByRecordID_V2Request{}
	result, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2EngineServer_HowEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.HowEntityByEntityIDRequest{}
	result, err := g2engine.HowEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919"}]}]}]}}}
}

func ExampleG2EngineServer_HowEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.HowEntityByEntityID_V2Request{}
	result, err := g2engine.HowEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919"}]}]}]}}}
}

func ExampleG2EngineServer_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.InitRequest{}
	result, err := g2engine.Init(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_InitWithConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.InitWithConfigIDRequest{}
	result, err := g2engine.InitWithConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_PrimeEngine() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.PrimeEngineRequest{}
	result, err := g2engine.PrimeEngine(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_Process() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessRequest{}
	result, err := g2engine.Process(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_ProcessRedoRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessRedoRecordRequest{}
	result, err := g2engine.ProcessRedoRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_ProcessRedoRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessRedoRecordWithInfoRequest{}
	result, err := g2engine.ProcessRedoRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_ProcessWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessWithInfoRequest{}
	result, err := g2engine.ProcessWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"555","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ProcessWithResponse() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessWithResponseRequest{}
	result, err := g2engine.ProcessWithResponse(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
}

func ExampleG2EngineServer_ProcessWithResponseResize() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ProcessWithResponseResizeRequest{}
	result, err := g2engine.ProcessWithResponseResize(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
}

func ExampleG2EngineServer_ReevaluateEntity() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateEntityRequest{}
	result, err := g2engine.ReevaluateEntity(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}
func ExampleG2EngineServer_ReevaluateEntityWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateEntityWithInfoRequest{}
	result, err := g2engine.ReevaluateEntityWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ReevaluateRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateRecordRequest{}
	result, err := g2engine.ReevaluateRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_ReevaluateRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReevaluateRecordWithInfoRequest{}
	result, err := g2engine.ReevaluateRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}

}

func ExampleG2EngineServer_Reinit() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReinitRequest{}
	result, err := g2engine.Reinit(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_ReplaceRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReplaceRecordRequest{}
	result, err := g2engine.ReplaceRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_ReplaceRecordWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.ReplaceRecordWithInfoRequest{}
	result, err := g2engine.ReplaceRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"TEST","RECORD_ID":"111","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.DestroyRequest{}
	result, err := g2engine.Destroy(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2EngineServer_SearchByAttributes() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.SearchByAttributesRequest{}
	result, err := g2engine.SearchByAttributes(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":...
}

func ExampleG2EngineServer_SearchByAttributes_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.SearchByAttributes_V2Request{}
	result, err := g2engine.SearchByAttributes_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","MATCH_KEY":"+NAME+SSN","ERRULE_CODE":"SF1_PNAME_CSTAB"},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1}}}]}
}

func ExampleG2EngineServer_Stats() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.StatsRequest{}
	result, err := g2engine.Stats(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: { "workload": { "loadedRecords": 5,  "addedRecords": 2,  "deletedRecords": 0,  "reevaluations": 0,  "repairedEntities": 0,...
}

func ExampleG2EngineServer_WhyEntities() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntitiesRequest{}
	result, err := g2engine.WhyEntities(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":2,"MATCH_INFO":{"WHY_KEY":...
}

func ExampleG2EngineServer_WhyEntities_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntities_V2Request{}
	result, err := g2engine.WhyEntities_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":2,"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

func ExampleG2EngineServer_WhyEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByEntityIDRequest{}
	result, err := g2engine.WhyEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_WhyEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByEntityID_V2Request{}
	result, err := g2engine.WhyEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_WhyEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByRecordIDRequest{}
	result, err := g2engine.WhyEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_WhyEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyEntityByRecordID_V2Request{}
	result, err := g2engine.WhyEntityByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":...
}

func ExampleG2EngineServer_WhyRecords() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyRecordsRequest{}
	result, err := g2engine.WhyRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"}],...
}

func ExampleG2EngineServer_WhyRecords_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.WhyRecords_V2Request{}
	result, err := g2engine.WhyRecords_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":2,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"TEST","RECORD_ID":"222"}],"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-DOB-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":2}}]}
}

// PurgeRepository() is first to start with a clean database.
func ExampleG2EngineServer_PurgeRepository_finalCleanup() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &pb.PurgeRepositoryRequest{}
	result, err := g2engine.PurgeRepository(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

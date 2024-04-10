package szdiagnosticserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/serve-grpc/szconfigmanagerserver"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	g2configmgrpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	szConfigManagerServerSingleton *szconfigmanagerserver.SzConfigManagerServer
	szDiagnoisticTestSingleton     *SzDiagnosticServer
	localLogger                    logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(localLogger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) SzDiagnosticServer {
	if szDiagnoisticTestSingleton == nil {
		szDiagnoisticTestSingleton = &SzDiagnosticServer{}
		instanceName := "Test module name"
		verboseLogging := int64(0)
		settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkSzDiagnostic().Init(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *szDiagnoisticTestSingleton
}

func getSzDiagnosticServer(ctx context.Context) SzDiagnosticServer {
	if szDiagnoisticTestSingleton == nil {
		szDiagnoisticTestSingleton = &SzDiagnosticServer{}
		moduleName := "Test module name"
		verboseLogging := int64(0)
		iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkSzDiagnostic().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *szDiagnoisticTestSingleton
}

func getSzConfigManagerServer(ctx context.Context) szconfigmanagerserver.SzConfigManagerServer {
	if szConfigManagerServerSingleton == nil {
		szConfigManagerServerSingleton = &szconfigmanagerserver.SzConfigManagerServer{}
		instanceName := "Test module name"
		verboseLogging := int64(0)
		settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = szconfigmanagerserver.GetSdkSzConfigManager().Init(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *szConfigManagerServerSingleton
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

func testError(test *testing.T, ctx context.Context, szDiagnosticServer SzDiagnosticServer, err error) {
	_ = ctx
	_ = szDiagnosticServer
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, szDiagnosticServer SzDiagnosticServer, err error, messageId string) {
	_ = ctx
	_ = szDiagnosticServer
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

func testErrorNoFail(test *testing.T, ctx context.Context, szDiagnosticServer SzDiagnosticServer, err error) {
	_ = ctx
	_ = szDiagnosticServer
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

func setupAddRecords(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {

	aG2engine := &szengine.Szengine{}
	err := aG2engine.Initialize(ctx, moduleName, iniParams, verboseLogging)
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

func setupPurgeRepository(ctx context.Context, instancename string, settings string, verboseLogging int64) error {
	aG2diagnostic := &szdiagnostic.Szdiagnostic{}
	err := aG2diagnostic.Initialize(ctx, instancename, settings, verboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
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
	actual, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
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
	szDiagnostic := getTestObject(ctx, test)
	request := &g2pb.CheckDBPerfRequest{
		SecondsToRun: int32(1),
	}
	response, err := szDiagnostic.CheckDBPerf(ctx, request)
	testError(test, ctx, szDiagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_Init(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int64(0),
	}
	response, err := szDiagnostic.Init(ctx, request)
	expectError(test, ctx, szDiagnostic, err, "senzing-60134002")
	printActual(test, response)
}

func TestG2diagnosticserver_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   int64(1),
		VerboseLogging: int64(0),
	}
	response, err := szDiagnostic.InitWithConfigID(ctx, request)
	expectError(test, ctx, szDiagnostic, err, "senzing-60134003")
	printActual(test, response)
}

func TestG2diagnosticserver_Reinit(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	getDefaultConfigIDRequest := &g2configmgrpb.GetDefaultConfigIdRequest{}
	getDefaultConfigIDResponse, err := szConfigManagerServer.GetDefaultConfigID(ctx, getDefaultConfigIDRequest)
	testError(test, ctx, szDiagnostic, err)
	request := &g2pb.ReinitRequest{
		InitConfigID: getDefaultConfigIDResponse.ConfigID,
	}
	response, err := szDiagnostic.Reinit(ctx, request)
	testError(test, ctx, szDiagnostic, err)
	printActual(test, response)
}

func TestG2diagnosticserver_Destroy(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	request := &g2pb.DestroyRequest{}
	response, err := szDiagnostic.Destroy(ctx, request)
	expectError(test, ctx, szDiagnostic, err, "senzing-60134001")
	printActual(test, response)
}

package szdiagnosticserver

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	szDiagnosticServerSingleton *SzDiagnosticServer
	localLogger                 logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(localLogger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) SzDiagnosticServer {
	if szDiagnosticServerSingleton == nil {
		szDiagnosticServerSingleton = &SzDiagnosticServer{}
		instanceName := "Test name"
		verboseLogging := sz.SZ_NO_LOGGING
		configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
		settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkSzDiagnostic().Initialize(ctx, instanceName, settings, verboseLogging, configId)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *szDiagnosticServerSingleton
}

func getSzDiagnosticServer(ctx context.Context) SzDiagnosticServer {
	if szDiagnosticServerSingleton == nil {
		szDiagnosticServerSingleton = &SzDiagnosticServer{}
		instanceName := "Test name"
		verboseLogging := sz.SZ_NO_LOGGING
		configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
		settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkSzDiagnostic().Initialize(ctx, instanceName, settings, verboseLogging, configId)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *szDiagnosticServerSingleton
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

// func expectError(test *testing.T, ctx context.Context, szDiagnosticServer SzDiagnosticServer, err error, messageId string) {
// 	_ = ctx
// 	_ = szDiagnosticServer
// 	if err != nil {
// 		var dictionary map[string]interface{}
// 		unmarshalErr := json.Unmarshal([]byte(err.Error()), &dictionary)
// 		if unmarshalErr != nil {
// 			test.Log("Unmarshal Error:", unmarshalErr.Error())
// 		}
// 		assert.Equal(test, messageId, dictionary["id"].(string))
// 	} else {
// 		assert.FailNow(test, "Should have failed with", messageId)
// 	}
// }

// func testErrorNoFail(test *testing.T, ctx context.Context, szDiagnosticServer SzDiagnosticServer, err error) {
// 	_ = ctx
// 	_ = szDiagnosticServer
// 	if err != nil {
// 		test.Log("Error:", err.Error())
// 	}
// }

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

	configComments := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	configID, err := szConfigManager.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configID)
	if err != nil {
		return createError(5914, err)
	}

	err = szConfigManager.Destroy(ctx)
	if err != nil {
		return createError(5915, err)
	}
	return err
}

func setupAddRecords(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	szEngine := &szengine.Szengine{}
	configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
	err := szEngine.Initialize(ctx, instanceName, settings, verboseLogging, configId)
	if err != nil {
		return createError(5916, err)
	}

	flags := sz.SZ_NO_FLAGS
	testRecordIds := []string{"1001", "1002", "1003", "1004", "1005", "1039", "1040"}
	for _, testRecordId := range testRecordIds {
		testRecord := truthset.CustomerRecords[testRecordId]
		_, err := szEngine.AddRecord(ctx, testRecord.DataSource, testRecord.Id, testRecord.Json, flags)
		if err != nil {
			return createError(5917, err)
		}
	}

	err = szEngine.Destroy(ctx)
	if err != nil {
		return createError(5918, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, instancename string, settings string, verboseLogging int64) error {
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	err := szDiagnostic.Initialize(ctx, instancename, settings, verboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
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
	instanceName := "Test module name"
	verboseLogging := sz.SZ_NO_LOGGING
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

	err = setupPurgeRepository(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return createError(5921, err)
	}

	// Add records.

	err = setupAddRecords(ctx, instanceName, settings, verboseLogging)
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

func TestSzDiagnosticServer_CheckDatabasePerformance(test *testing.T) {
	ctx := context.TODO()
	szDiagnosticServer := getTestObject(ctx, test)
	request := &szpb.CheckDatabasePerformanceRequest{
		SecondsToRun: int32(1),
	}
	response, err := szDiagnosticServer.CheckDatabasePerformance(ctx, request)
	testError(test, ctx, szDiagnosticServer, err)
	printActual(test, response)
}

func TestSzDiagnosticServer_PurgeRepository(test *testing.T) {
	ctx := context.TODO()
	szDiagnosticServer := getTestObject(ctx, test)
	request := &szpb.PurgeRepositoryRequest{}
	response, err := szDiagnosticServer.PurgeRepository(ctx, request)
	testError(test, ctx, szDiagnosticServer, err)
	printActual(test, response)
}

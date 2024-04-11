package szconfigmanagerserver

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
	"github.com/senzing-garage/serve-grpc/szconfigserver"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	szConfigManagerServerSingleton *SzConfigManagerServer
	localLogger                    logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(localLogger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) SzConfigManagerServer {
	if szConfigManagerServerSingleton == nil {
		szConfigManagerServerSingleton = &SzConfigManagerServer{}
		instanceName := "Test name"
		verboseLogging := int64(0)
		settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkSzConfigManager().Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *szConfigManagerServerSingleton
}

func getSzConfigManagerServer(ctx context.Context) SzConfigManagerServer {
	if szConfigManagerServerSingleton == nil {
		szConfigManagerServerSingleton = &SzConfigManagerServer{}
		moduleName := "Test module name"
		verboseLogging := int64(0)
		iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkSzConfigManager().Initialize(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *szConfigManagerServerSingleton
}

func getSzConfigServer(ctx context.Context) szconfigserver.SzConfigServer {
	szConfigServer := &szconfigserver.SzConfigServer{}
	instanceName := "Test name"
	verboseLogging := int64(0)
	settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	err = szconfigserver.GetSdkSzConfig().Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	return *szConfigServer
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

func testError(test *testing.T, ctx context.Context, szConfigManager SzConfigManagerServer, err error) {
	_ = ctx
	_ = szConfigManager
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, szConfigManager SzConfigManagerServer, err error, messageId string) {
	_ = ctx
	_ = szConfigManager
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

func setupSenzingConfig(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	now := time.Now()

	szConfig := &szconfig.Szconfig{}
	err := szConfig.Initialize(ctx, moduleName, iniParams, verboseLogging)
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

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
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
	err = szConfigManager.Initialize(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return createError(5912, err)
	}

	configComment := fmt.Sprintf("Created by szconfigmanagerserver_test at %s", now.UTC())
	configId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
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

func setupPurgeRepository(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	err := szDiagnostic.Initialize(ctx, instanceName, settings, verboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
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

func TestSzConfigManagerServer_AddConfig(test *testing.T) {
	ctx := context.TODO()
	now := time.Now()
	szConfigServer := getSzConfigServer(ctx)

	// SzConfig Create() to create a Senzing configuration.
	requestToCreate := &szconfigpb.CreateConfigRequest{}
	responseFromCreate, err := szConfigServer.CreateConfig(ctx, requestToCreate)
	if err != nil {
		test.Log(err)
	}

	// SzConfig ExportConfig() to create a JSON string.
	requestToSave := &szconfigpb.ExportConfigRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := szConfigServer.ExportConfig(ctx, requestToSave)
	if err != nil {
		test.Log(err)
	}

	// Test.
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.AddConfigRequest{
		ConfigDefinition: responseFromSave.GetResult(),
		ConfigComment:    fmt.Sprintf("g2configmgrserver_test at %s", now.UTC()),
	}
	response, err := szConfigManagerServer.AddConfig(ctx, request)
	testError(test, ctx, szConfigManagerServer, err)
	printActual(test, response)
}

func TestSzConfigManagerServerImpl_GetConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)

	// GetDefaultConfigId().
	requestToGetDefaultConfigId := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigId, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigId)
	testError(test, ctx, szConfigManagerServer, err)
	printActual(test, responseFromGetDefaultConfigId)

	// Test.
	request := &szpb.GetConfigRequest{
		ConfigId: responseFromGetDefaultConfigId.GetResult(),
	}
	response, err := szConfigManagerServer.GetConfig(ctx, request)
	testError(test, ctx, szConfigManagerServer, err)
	printActual(test, response)
}

func TestSzConfigManagerServerImpl_GetConfigList(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.GetConfigListRequest{}
	response, err := szConfigManagerServer.GetConfigList(ctx, request)
	testError(test, ctx, szConfigManagerServer, err)
	printActual(test, response)
}

func TestSzConfigManagerServerImpl_GetDefaultConfigId(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.GetDefaultConfigIdRequest{}
	response, err := szConfigManagerServer.GetDefaultConfigId(ctx, request)
	testError(test, ctx, szConfigManagerServer, err)
	printActual(test, response)
}

func TestSzConfigManagerServerImpl_ReplaceDefaultConfigId(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)

	// GetDefaultConfigId()
	requestToGetDefaultConfigId := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigId, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigId)
	testError(test, ctx, szConfigManagerServer, err)

	// Test. Note: Cheating a little with replacing with same configId.
	request := &szpb.ReplaceDefaultConfigIdRequest{
		CurrentDefaultConfigId: responseFromGetDefaultConfigId.GetResult(),
		NewDefaultConfigId:     responseFromGetDefaultConfigId.GetResult(),
	}
	response, err := szConfigManagerServer.ReplaceDefaultConfigId(ctx, request)
	testError(test, ctx, szConfigManagerServer, err)
	printActual(test, response)
}

func TestSzConfigManagerServerImpl_SetDefaultConfigId(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)

	// GetDefaultConfigId()
	requestToGetDefaultConfigId := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigId, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigId)
	testError(test, ctx, szConfigManagerServer, err)

	// Test.
	request := &szpb.SetDefaultConfigIdRequest{
		ConfigId: responseFromGetDefaultConfigId.GetResult(),
	}
	response, err := szConfigManagerServer.SetDefaultConfigId(ctx, request)
	testError(test, ctx, szConfigManagerServer, err)
	printActual(test, response)
}

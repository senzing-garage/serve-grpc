package szconfigmanagerserver

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-grpc/szconfigserver"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultTruncation = 76
	observerID        = "Observer 1"
	observerOrigin    = "Observer 1 origin"
	printResults      = false
)

var (
	logger            logging.Logging
	logLevelName      = "INFO"
	observerSingleton = &observer.NullObserver{
		ID:       observerID,
		IsSilent: true,
	}
	szConfigManagerServerSingleton *SzConfigManagerServer
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzConfigManagerServer_GetConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	require.NoError(test, err)

	// Test GetConfig.

	request := &szpb.GetConfigRequest{
		ConfigId: responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.GetConfig(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzConfigManagerServer_GetConfigs(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.GetConfigsRequest{}
	response, err := szConfigManagerServer.GetConfigs(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzConfigManagerServer_GetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.GetDefaultConfigIdRequest{}
	response, err := szConfigManagerServer.GetDefaultConfigId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzConfigManagerServer_GetTemplateConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.GetTemplateConfigRequest{}
	response, err := szConfigManagerServer.GetTemplateConfig(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzConfigManagerServer_RegisterConfig(test *testing.T) {
	now := time.Now()
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	require.NoError(test, err)

	// Add DataSource to the Senzing configuration.

	requestToAddDataSource := &szconfigpb.AddDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}
	responseFromAddDataSource, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	require.NoError(test, err)

	// Test RegisterConfig.

	request := &szpb.RegisterConfigRequest{
		ConfigDefinition: responseFromAddDataSource.GetResult(),
		ConfigComment:    fmt.Sprintf("szconfigmanagerserver_test at %s", now.UTC()),
	}
	response, err := szConfigManagerServer.RegisterConfig(ctx, request)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzConfigManagerServer_ReplaceDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	require.NoError(test, err)

	// Test. Note: Cheating a little with replacing with same configId.

	request := &szpb.ReplaceDefaultConfigIdRequest{
		CurrentDefaultConfigId: responseFromGetDefaultConfigID.GetResult(),
		NewDefaultConfigId:     responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.ReplaceDefaultConfigId(ctx, request)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzConfigManagerServer_SetDefaultConfig(test *testing.T) {}

func TestSzConfigManagerServer_SetDefaultConfigId(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	require.NoError(test, err)

	// Test GetDefaultConfigID.

	request := &szpb.SetDefaultConfigIdRequest{
		ConfigId: responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.SetDefaultConfigId(ctx, request)
	require.NoError(test, err)
	printActual(test, response)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzConfigManagerServer_RegisterObserver(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

func TestSzConfigManagerServer_SetLogLevel(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	require.NoError(test, err)
}

func TestSzConfigManagerServer__SetLogLevel_badLevelName(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	require.Error(test, err)
}

func TestSzConfigManagerServer_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	testObject.SetObserverOrigin(ctx, observerOrigin)
}

func TestSzConfigManagerServer_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	actual := testObject.GetObserverOrigin(ctx)
	assert.Equal(test, observerOrigin, actual)
}

func TestSzConfigManagerServer_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestBuildSimpleSystemConfigurationJsonUsingEnvVars(test *testing.T) {
	actual, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getSzConfigManagerServer(ctx context.Context) SzConfigManagerServer {
	if szConfigManagerServerSingleton == nil {
		szConfigManagerServerSingleton = &SzConfigManagerServer{}
		instanceName := "Test instance name"
		verboseLogging := senzing.SzNoLogging
		settings, err := settings.BuildSimpleSettingsUsingEnvVars()
		panicOnError(err)
		osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
		if len(osenvLogLevel) > 0 {
			logLevelName = osenvLogLevel
		}
		err = szConfigManagerServerSingleton.SetLogLevel(ctx, logLevelName)
		panicOnError(err)
		err = GetSdkSzConfigManager().Initialize(ctx, instanceName, settings, verboseLogging)
		panicOnError(err)
	}
	return *szConfigManagerServerSingleton
}

func getSzConfigServer(ctx context.Context) szconfigserver.SzConfigServer {
	szConfigServer := &szconfigserver.SzConfigServer{}
	instanceName := "Test instance name"
	verboseLogging := senzing.SzNoLogging
	settings, err := settings.BuildSimpleSettingsUsingEnvVars()
	panicOnError(err)
	osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
	if len(osenvLogLevel) > 0 {
		logLevelName = osenvLogLevel
	}
	err = szConfigServer.SetLogLevel(ctx, logLevelName)
	panicOnError(err)
	err = szConfigServer.GetSdkSzConfig().Initialize(ctx, instanceName, settings, verboseLogging)
	panicOnError(err)
	return *szConfigServer
}

func getTestObject(ctx context.Context, test *testing.T) SzConfigManagerServer {
	_ = test
	return getSzConfigManagerServer(ctx)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setupSenzingConfig(ctx context.Context, instanceName string, settings string, verboseLogging int64) {
	now := time.Now()

	szAbstractFactory := &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	panicOnError(err)

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	panicOnError(err)

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range datasourceNames {
		_, err := szConfig.AddDataSource(ctx, dataSourceCode)
		panicOnError(err)
	}

	configComment := fmt.Sprintf("Created by szconfigmanagerserver_test at %s", now.UTC())
	configDefinition, err := szConfig.Export(ctx)
	panicOnError(err)

	configID, err := szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	panicOnError(err)

	szAbstractFactory.Reinitialize(ctx, configID)
}

func setupPurgeRepository(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	err := szDiagnostic.Initialize(
		ctx,
		instanceName,
		settings,
		senzing.SzInitializeWithDefaultConfiguration,
		verboseLogging,
	)
	panicOnError(err)

	err = szDiagnostic.PurgeRepository(ctx)
	panicOnError(err)

	err = szDiagnostic.Destroy(ctx)
	panicOnError(err)

	return err
}

func setup() {
	ctx := context.TODO()
	instanceName := "Test name"
	verboseLogging := senzing.SzNoLogging
	settings, err := settings.BuildSimpleSettingsUsingEnvVars()
	panicOnError(err)
	setupSenzingConfig(ctx, instanceName, settings, verboseLogging)
	err = setupPurgeRepository(ctx, instanceName, settings, verboseLogging)
	panicOnError(err)
}

package szconfigmanagerserver_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-grpc/szconfigmanagerserver"
	"github.com/senzing-garage/serve-grpc/szconfigserver"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultTruncation = 76
	jsonIndentation   = "    "
	observerID        = "Observer 1"
	observerOrigin    = "Observer 1 origin"
	printErrors       = false
	printResults      = false
)

// Bad parameters

const (
	badConfigDefinition       = "\n\t"
	badConfigID               = int64(0)
	badCurrentDefaultConfigID = int64(0)
	badLogLevelName           = "BadLogLevelName"
	badNewDefaultConfigID     = int64(0)
	baseTen                   = 10
)

// Nil/empty parameters

var nilConfigDefinition string

var (
	logLevelName      = "INFO"
	observerSingleton = &observer.NullObserver{
		ID:       observerID,
		IsSilent: true,
	}
	szConfigManagerServerSingleton *szconfigmanagerserver.SzConfigManagerServer
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzConfigManagerServer_GetConfig(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	printError(test, err)
	require.NoError(test, err)

	// Test GetConfig.

	request := &szpb.GetConfigRequest{
		ConfigId: responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.GetConfig(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzConfigManagerServer_GetConfig_badConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.GetConfigRequest{
		ConfigId: badConfigID,
	}
	response, err := szConfigManagerServer.GetConfig(ctx, request)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanagerserver.(*SzConfigManagerServer).GetConfig","text":"CreateConfigFromConfigID","error":{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromConfigID","error":{"function":"szconfigmanager.(*Szconfigmanager).createConfigFromConfigIDChoreography","text":"getConfig(0)","error":{"id":"SZSDK60024003","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, response)
}

func TestSzConfigManagerServer_GetConfigs(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.GetConfigsRequest{}
	response, err := szConfigManagerServer.GetConfigs(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzConfigManagerServer_GetDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.GetDefaultConfigIdRequest{}
	response, err := szConfigManagerServer.GetDefaultConfigId(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzConfigManagerServer_GetTemplateConfig(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.GetTemplateConfigRequest{}
	response, err := szConfigManagerServer.GetTemplateConfig(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzConfigManagerServer_RegisterConfig(test *testing.T) {
	now := time.Now()
	ctx := test.Context()
	szConfigServer := getSzConfigServer(ctx)
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	printError(test, err)
	require.NoError(test, err)

	// Add DataSource to the Senzing configuration.

	requestToAddDataSource := &szconfigpb.AddDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}
	responseFromAddDataSource, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	printError(test, err)
	require.NoError(test, err)

	// Test RegisterConfig.

	request := &szpb.RegisterConfigRequest{
		ConfigDefinition: responseFromAddDataSource.GetResult(),
		ConfigComment:    fmt.Sprintf("szconfigmanagerserver_test at %s", now.UTC()),
	}
	response, err := szConfigManagerServer.RegisterConfig(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzConfigManagerServer_RegisterConfig_badConfigDefinition(test *testing.T) {
	now := time.Now()
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.RegisterConfigRequest{
		ConfigDefinition: badConfigDefinition,
		ConfigComment:    fmt.Sprintf("szconfigmanagerserver_test at %s", now.UTC()),
	}
	response, err := szConfigManagerServer.RegisterConfig(ctx, request)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanagerserver.(*SzConfigManagerServer).RegisterConfig","error":{"function":"szconfigmanager.(*Szconfigmanager).RegisterConfig","error":{"id":"SZSDK60024001","reason":"SENZ0028|Invalid JSON config document"}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, response)
}

func TestSzConfigManagerServer_RegisterConfig_nilConfigDefinition(test *testing.T) {
	now := time.Now()
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)
	request := &szpb.RegisterConfigRequest{
		ConfigDefinition: nilConfigDefinition,
		ConfigComment:    fmt.Sprintf("szconfigmanagerserver_test at %s", now.UTC()),
	}
	response, err := szConfigManagerServer.RegisterConfig(ctx, request)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanagerserver.(*SzConfigManagerServer).RegisterConfig","error":{"function":"szconfigmanager.(*Szconfigmanager).RegisterConfig","error":{"id":"SZSDK60024001","reason":"SENZ0028|Invalid JSON config document"}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, response)
}

func TestSzConfigManagerServer_ReplaceDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	printError(test, err)
	require.NoError(test, err)

	// Test. Note: Cheating a little with replacing with same configId.

	request := &szpb.ReplaceDefaultConfigIdRequest{
		CurrentDefaultConfigId: responseFromGetDefaultConfigID.GetResult(),
		NewDefaultConfigId:     responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.ReplaceDefaultConfigId(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzConfigManagerServer_ReplaceDefaultConfigID_badCurrentDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	printError(test, err)
	require.NoError(test, err)

	request := &szpb.ReplaceDefaultConfigIdRequest{
		CurrentDefaultConfigId: badCurrentDefaultConfigID,
		NewDefaultConfigId:     responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.ReplaceDefaultConfigId(ctx, request)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzReplaceConflict)

	expectedErr := `{"function":"szconfigmanagerserver.(*SzConfigManagerServer).ReplaceDefaultConfigId","error":{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7245|Current configuration ID does not match specified data ID [0]."}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, response)
}

func TestSzConfigManagerServer_ReplaceDefaultConfigID_badNewDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	printError(test, err)
	require.NoError(test, err)

	request := &szpb.ReplaceDefaultConfigIdRequest{
		CurrentDefaultConfigId: responseFromGetDefaultConfigID.GetResult(),
		NewDefaultConfigId:     badNewDefaultConfigID,
	}
	response, err := szConfigManagerServer.ReplaceDefaultConfigId(ctx, request)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanagerserver.(*SzConfigManagerServer).ReplaceDefaultConfigId","error":{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, response)
}

func TestSzConfigManagerServer_SetDefaultConfig(test *testing.T) {
	ctx := test.Context()

	szConfigManagerServer := getTestObject(ctx, test)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	printError(test, err)
	require.NoError(test, err)

	defaultConfigID := responseFromGetDefaultConfigID.GetResult()

	// Get the definition of the default Senzing configuration.

	requestToGetConfig := &szpb.GetConfigRequest{
		ConfigId: defaultConfigID,
	}
	responseFromGetConfig, err := szConfigManagerServer.GetConfig(ctx, requestToGetConfig)
	printError(test, err)
	require.NoError(test, err)

	defaultConfigDefinition := responseFromGetConfig.GetResult()

	// Set the new Senzing configuration.
	// Note: This cheats a little.  Normally a non-default configuration is used in SetDefaultConfig.

	requestToSetDefaultConfig := &szpb.SetDefaultConfigRequest{
		ConfigDefinition: defaultConfigDefinition,
		ConfigComment:    "Just a test",
	}
	_, err = szConfigManagerServer.SetDefaultConfig(ctx, requestToSetDefaultConfig)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzConfigManagerServer_SetDefaultConfig_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)
	requestToSetDefaultConfig := &szpb.SetDefaultConfigRequest{
		ConfigDefinition: badConfigDefinition,
		ConfigComment:    "Just a test",
	}
	_, err := szConfigManagerServer.SetDefaultConfig(ctx, requestToSetDefaultConfig)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanagerserver.(*SzConfigManagerServer).SetDefaultConfig","error":{"function":"szconfigmanager.(*Szconfigmanager).SetDefaultConfig","error":{"function":"szconfigmanager.(*Szconfigmanager).setDefaultConfigChoreography","text":"registerConfig","error":{"id":"SZSDK60024001","reason":"SENZ0028|Invalid JSON config document"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigManagerServer_SetDefaultConfigId(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	printError(test, err)
	require.NoError(test, err)

	// Test GetDefaultConfigID.

	request := &szpb.SetDefaultConfigIdRequest{
		ConfigId: responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.SetDefaultConfigId(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzConfigManagerServer_SetDefaultConfigId_badConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getTestObject(ctx, test)

	// Test GetDefaultConfigID.

	request := &szpb.SetDefaultConfigIdRequest{
		ConfigId: badConfigID,
	}
	response, err := szConfigManagerServer.SetDefaultConfigId(ctx, request)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanagerserver.(*SzConfigManagerServer).SetDefaultConfigId","error":{"function":"szconfigmanager.(*Szconfigmanager).SetDefaultConfigID","error":{"id":"SZSDK60024008","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, response)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzConfigManagerServer_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzConfigManagerServer_SetLogLevel(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	printError(test, err)
	require.NoError(test, err)
}

func TestSzConfigManagerServer__SetLogLevel_badLevelName(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	printError(test, err)
	require.Error(test, err)
}

func TestSzConfigManagerServer_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	testObject.SetObserverOrigin(ctx, observerOrigin)
}

func TestSzConfigManagerServer_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	actual := testObject.GetObserverOrigin(ctx)
	assert.Equal(test, observerOrigin, actual)
}

func TestSzConfigManagerServer_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.UnregisterObserver(ctx, observerSingleton)
	printError(test, err)
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

func getSzConfigManagerServer(ctx context.Context) *szconfigmanagerserver.SzConfigManagerServer {
	if szConfigManagerServerSingleton == nil {
		szConfigManagerServerSingleton = &szconfigmanagerserver.SzConfigManagerServer{}
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
		err = szconfigmanagerserver.GetSdkSzConfigManager().Initialize(ctx, instanceName, settings, verboseLogging)
		panicOnError(err)
	}

	return szConfigManagerServerSingleton
}

func getSzConfigServer(ctx context.Context) *szconfigserver.SzConfigServer {
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
	err = szconfigserver.GetSdkSzConfigManager().Initialize(ctx, instanceName, settings, verboseLogging)
	panicOnError(err)

	return szConfigServer
}

func getTestObject(ctx context.Context, t *testing.T) *szconfigmanagerserver.SzConfigManagerServer {
	t.Helper()

	return getSzConfigManagerServer(ctx)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func printActual(t *testing.T, actual interface{}) {
	t.Helper()
	printResult(t, "Actual", actual)
}

func printError(t *testing.T, err error) {
	t.Helper()

	if printErrors {
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}
}

func printResult(t *testing.T, title string, result interface{}) {
	t.Helper()

	if printResults {
		t.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
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

	err = szAbstractFactory.Reinitialize(ctx, configID)
	panicOnError(err)
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

	return wraperror.Errorf(err, wraperror.NoMessage)
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

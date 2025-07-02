package szconfigserver_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-grpc/szconfigmanagerserver"
	"github.com/senzing-garage/serve-grpc/szconfigserver"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
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
	badConfigDefinition = "}{"
	badConfigHandle     = uintptr(0)
	badDataSourceCode   = "\n\tGO_TEST"
	badLogLevelName     = "BadLogLevelName"
	badSettings         = "{]"
)

var (
	logLevelName      = "INFO"
	observerSingleton = &observer.NullObserver{
		ID:       observerID,
		IsSilent: true,
	}
	szConfigManagerServerSingleton *szconfigmanagerserver.SzConfigManagerServer
)

// Nil/empty parameters

var nilDataSourceCode string

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzConfigServer_RegisterDataSource(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Add DataSource to the Senzing configuration.

	requestToRegisterDataSource := &szpb.RegisterDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}
	responseFromRegisterDataSource, err := szConfigServer.RegisterDataSource(ctx, requestToRegisterDataSource)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromRegisterDataSource.GetResult())
}

func TestSzConfigServer_RegisterDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Add DataSource to the Senzing configuration.

	requestToRegisterDataSource := &szpb.RegisterDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   badDataSourceCode,
	}
	_, err = szConfigServer.RegisterDataSource(ctx, requestToRegisterDataSource)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).RegisterDataSource","text":"RegisterDataSource: \n\tGO_TEST","error":{"function":"szconfig.(*Szconfig).RegisterDataSource","error":{"function":"szconfig.(*Szconfig).registerDataSourceChoreography","text":"registerDataSource: \n\tGO_TEST","error":{"id":"SZSDK60014001","reason":"SENZ3121|JSON Parsing Failure [code=12,offset=15]"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigServer_RegisterDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Add DataSource to the Senzing configuration.

	requestToRegisterDataSource := &szpb.RegisterDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   nilDataSourceCode,
	}
	_, err = szConfigServer.RegisterDataSource(ctx, requestToRegisterDataSource)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).RegisterDataSource","text":"RegisterDataSource: ","error":{"function":"szconfig.(*Szconfig).RegisterDataSource","error":{"function":"szconfig.(*Szconfig).registerDataSourceChoreography","text":"registerDataSource: ","error":{"id":"SZSDK60014001","reason":"SENZ7313|A non-empty value for [DSRC_CODE] must be specified."}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigServer_UnregisterDataSource(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Delete DataSource to the Senzing configuration.

	requestToUnregisterDataSource := &szpb.UnregisterDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}
	_, err = szConfigServer.UnregisterDataSource(ctx, requestToUnregisterDataSource)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzConfigServer_UnregisterDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Delete DataSource to the Senzing configuration.

	requestToUnregisterDataSource := &szpb.UnregisterDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   badDataSourceCode,
	}
	_, err = szConfigServer.UnregisterDataSource(ctx, requestToUnregisterDataSource)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).UnregisterDataSource","text":"UnregisterDataSource: \n\tGO_TEST","error":{"function":"szconfig.(*Szconfig).UnregisterDataSource","error":{"function":"szconfig.(*Szconfig).unregisterDataSourceChoreography","text":"unregisterDataSource(\n\tGO_TEST)","error":{"id":"SZSDK60014004","reason":"SENZ3121|JSON Parsing Failure [code=12,offset=15]"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigServer_UnregisterDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Delete DataSource to the Senzing configuration.

	requestToUnregisterDataSource := &szpb.UnregisterDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   nilDataSourceCode,
	}
	_, err = szConfigServer.UnregisterDataSource(ctx, requestToUnregisterDataSource)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).UnregisterDataSource","text":"UnregisterDataSource: ","error":{"function":"szconfig.(*Szconfig).UnregisterDataSource","error":{"function":"szconfig.(*Szconfig).unregisterDataSourceChoreography","text":"unregisterDataSource()","error":{"id":"SZSDK60014004","reason":"SENZ7313|A non-empty value for [DSRC_CODE] must be specified."}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigServer_GetDataSourceRegistry(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Delete DataSource to the Senzing configuration.

	requestToGetDataSourceRegistry := &szpb.GetDataSourceRegistryRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
	}
	responseFromGetDataSourceRegistry, err := szConfigServer.GetDataSourceRegistry(ctx, requestToGetDataSourceRegistry)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetDataSourceRegistry.GetResult())
}

func TestSzConfigServer_VerifyConfig(test *testing.T) {
	ctx := test.Context()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Delete DataSource to the Senzing configuration.

	requestToVerifyConfig := &szpb.VerifyConfigRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
	}
	responseFromGetDataSourceRegistry, err := szConfigServer.VerifyConfig(ctx, requestToVerifyConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetDataSourceRegistry.GetResult())
}

func TestSzConfigServer_VerifyConfig_bad_config(test *testing.T) {
	ctx := test.Context()
	badConfigDefinition := "}{"

	szConfigServer := getTestObject(ctx, test)

	// Delete DataSource to the Senzing configuration.

	requestToVerifyConfig := &szpb.VerifyConfigRequest{
		ConfigDefinition: badConfigDefinition,
	}
	responseFromGetDataSourceRegistry, err := szConfigServer.VerifyConfig(ctx, requestToVerifyConfig)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).createSzConfig","error":{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromStringChoreography","text":"VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).verifyConfigDefinitionChoreography","text":"load","error":{"id":"SZSDK60014009","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, responseFromGetDataSourceRegistry.GetResult())
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzConfigServer_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzConfigServer_SetLogLevel(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	printError(test, err)
	require.NoError(test, err)
}

func TestSzConfigServer__SetLogLevel_badLevelName(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	printError(test, err)
	require.Error(test, err)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).SetLogLevel","text":"invalid error level: BADLEVELNAME","error":"szconfigserver"}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigServer_UnregisterObserver(test *testing.T) {
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
	panicOnError(err)
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

func getTestObject(ctx context.Context, t *testing.T) *szconfigserver.SzConfigServer {
	t.Helper()

	return getSzConfigServer(ctx)
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

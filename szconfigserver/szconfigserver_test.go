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
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzConfigServer_AddDataSource(test *testing.T) {
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

	requestToAddDataSource := &szpb.AddDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}
	responseFromAddDataSource, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromAddDataSource.GetResult())
}

func TestSzConfigServer_AddDataSource_badDataSourceCode(test *testing.T) {
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

	requestToAddDataSource := &szpb.AddDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   badDataSourceCode,
	}
	_, err = szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).AddDataSource","text":"AddDataSource: \n\tGO_TEST","error":{"function":"szconfig.(*Szconfig).AddDataSource","error":{"function":"szconfig.(*Szconfig).addDataSourceChoreography","text":"addDataSource: \n\tGO_TEST","error":{"id":"SZSDK60014001","reason":"SENZ3121|JSON Parsing Failure [code=12,offset=15]"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigServer_AddDataSource_nilDataSourceCode(test *testing.T) {
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

	requestToAddDataSource := &szpb.AddDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   nilDataSourceCode,
	}
	_, err = szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).AddDataSource","text":"AddDataSource: ","error":{"function":"szconfig.(*Szconfig).AddDataSource","error":{"function":"szconfig.(*Szconfig).addDataSourceChoreography","text":"addDataSource: ","error":{"id":"SZSDK60014001","reason":"SENZ7313|A non-empty value for [DSRC_CODE] must be specified."}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigServer_DeleteDataSource(test *testing.T) {
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

	requestToDeleteDataSource := &szpb.DeleteDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}
	_, err = szConfigServer.DeleteDataSource(ctx, requestToDeleteDataSource)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzConfigServer_DeleteDataSource_badDataSourceCode(test *testing.T) {
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

	requestToDeleteDataSource := &szpb.DeleteDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   badDataSourceCode,
	}
	_, err = szConfigServer.DeleteDataSource(ctx, requestToDeleteDataSource)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).DeleteDataSource","text":"DeleteDataSource: \n\tGO_TEST","error":{"function":"szconfig.(*Szconfig).DeleteDataSource","error":{"function":"szconfig.(*Szconfig).deleteDataSourceChoreography","text":"deleteDataSource(\n\tGO_TEST)","error":{"id":"SZSDK60014004","reason":"SENZ3121|JSON Parsing Failure [code=12,offset=15]"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigServer_DeleteDataSource_nilDataSourceCode(test *testing.T) {
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

	requestToDeleteDataSource := &szpb.DeleteDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   nilDataSourceCode,
	}
	_, err = szConfigServer.DeleteDataSource(ctx, requestToDeleteDataSource)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).DeleteDataSource","text":"DeleteDataSource: ","error":{"function":"szconfig.(*Szconfig).DeleteDataSource","error":{"function":"szconfig.(*Szconfig).deleteDataSourceChoreography","text":"deleteDataSource()","error":{"id":"SZSDK60014004","reason":"SENZ7313|A non-empty value for [DSRC_CODE] must be specified."}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzConfigServer_GetDataSources(test *testing.T) {
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

	requestToGetDataSources := &szpb.GetDataSourcesRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
	}
	responseFromGetDataSources, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetDataSources.GetResult())
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
	responseFromGetDataSources, err := szConfigServer.VerifyConfig(ctx, requestToVerifyConfig)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, responseFromGetDataSources.GetResult())
}

func TestSzConfigServer_VerifyConfig_bad_config(test *testing.T) {
	ctx := test.Context()
	badConfigDefinition := "}{"

	szConfigServer := getTestObject(ctx, test)

	// Delete DataSource to the Senzing configuration.

	requestToVerifyConfig := &szpb.VerifyConfigRequest{
		ConfigDefinition: badConfigDefinition,
	}
	responseFromGetDataSources, err := szConfigServer.VerifyConfig(ctx, requestToVerifyConfig)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigserver.(*SzConfigServer).createSzConfig","error":{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromStringChoreography","text":"VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).verifyConfigDefinitionChoreography","text":"load","error":{"id":"SZSDK60014009","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, responseFromGetDataSources.GetResult())
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

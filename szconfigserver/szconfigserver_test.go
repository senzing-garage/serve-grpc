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
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
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
	logLevelName      = "INFO"
	observerSingleton = &observer.NullObserver{
		ID:       observerID,
		IsSilent: true,
	}
	szConfigManagerTestSingleton *szconfigmanagerserver.SzConfigManagerServer
	szConfigTestSingleton        *szconfigserver.SzConfigServer
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzConfigServer_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Add DataSource to the Senzing configuration.

	requestToAddDataSource := &szpb.AddDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}
	responseFromAddDataSource, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	require.NoError(test, err)
	printActual(test, responseFromAddDataSource.GetResult())
}

func TestSzConfigServer_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Delete DataSource to the Senzing configuration.

	requestToDeleteDataSource := &szpb.DeleteDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}
	_, err = szConfigServer.DeleteDataSource(ctx, requestToDeleteDataSource)
	require.NoError(test, err)
}

func TestSzConfigServer_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getTestObject(ctx, test)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	require.NoError(test, err)
	printActual(test, responseFromGetTemplateConfig.GetResult())

	// Delete DataSource to the Senzing configuration.

	requestToGetDataSources := &szpb.GetDataSourcesRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
	}
	responseFromGetDataSources, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	require.NoError(test, err)
	printActual(test, responseFromGetDataSources.GetResult())
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzConfigServer_RegisterObserver(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

func TestSzConfigServer_SetLogLevel(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	require.NoError(test, err)
}

func TestSzConfigServer__SetLogLevel_badLevelName(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	require.Error(test, err)
}

func TestSzConfigServer_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	testObject.SetObserverOrigin(ctx, observerOrigin)
}

func TestSzConfigServer_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	actual := testObject.GetObserverOrigin(ctx)
	assert.Equal(test, observerOrigin, actual)
}

func TestSzConfigServer_UnregisterObserver(test *testing.T) {
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
	panicOnError(err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getSzConfigManagerServer(ctx context.Context) *szconfigmanagerserver.SzConfigManagerServer {
	if szConfigManagerTestSingleton == nil {
		szConfigManagerTestSingleton = &szconfigmanagerserver.SzConfigManagerServer{}
		osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
		if len(osenvLogLevel) > 0 {
			logLevelName = osenvLogLevel
		}
		err := szConfigTestSingleton.SetLogLevel(ctx, logLevelName)
		panicOnError(err)
	}
	return szConfigManagerTestSingleton
}

func getSzConfigServer(ctx context.Context) *szconfigserver.SzConfigServer {
	if szConfigTestSingleton == nil {
		szConfigTestSingleton = &szconfigserver.SzConfigServer{}
		osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
		if len(osenvLogLevel) > 0 {
			logLevelName = osenvLogLevel
		}
		err := szConfigTestSingleton.SetLogLevel(ctx, logLevelName)
		panicOnError(err)
	}
	return szConfigTestSingleton
}

func getTestObject(ctx context.Context, test *testing.T) *szconfigserver.SzConfigServer {
	_ = test
	return getSzConfigServer(ctx)
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

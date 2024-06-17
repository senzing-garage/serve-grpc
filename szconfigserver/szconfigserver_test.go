package szconfigserver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
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
	szConfigTestSingleton *SzConfigServer
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzConfigServer_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &szpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	require.NoError(test, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// AddDataSource.
	requestToAddDataSource := &szpb.AddDataSourceRequest{
		ConfigHandle:   responseFromCreateConfig.GetResult(),
		DataSourceCode: "GO_TEST",
	}
	responseFromAddDataSource, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	require.NoError(test, err)
	printActual(test, responseFromAddDataSource.GetResult())

	// Close.
	requestToCloseConfig := &szpb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	require.NoError(test, err)
}

func TestSzConfigServer_CloseConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &szpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	require.NoError(test, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// Close.
	requestToCloseConfig := &szpb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	require.NoError(test, err)
}

func TestSzConfigServer_CreateConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)
	requestToCreate := &szpb.CreateConfigRequest{}
	response, err := szConfigServer.CreateConfig(ctx, requestToCreate)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzConfigServer_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &szpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	require.NoError(test, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// GetDataSources #1.
	requestToGetDataSources := &szpb.GetDataSourcesRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromGetDataSources, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	require.NoError(test, err)
	initialDataSources := responseFromGetDataSources.GetResult()
	printActual(test, initialDataSources)

	// AddDataSource.
	requestToAddDataSource := &szpb.AddDataSourceRequest{
		ConfigHandle:   responseFromCreateConfig.GetResult(),
		DataSourceCode: "GO_TEST",
	}
	responseFromAddDataSource, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	require.NoError(test, err)
	printActual(test, responseFromAddDataSource.GetResult())

	// GetDataSources #2.
	responseFromListDataSources2, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	require.NoError(test, err)
	printActual(test, responseFromListDataSources2.GetResult())

	// DeleteDataSource.
	requestToDeleteDataSource := &szpb.DeleteDataSourceRequest{
		ConfigHandle:   responseFromCreateConfig.GetResult(),
		DataSourceCode: "GO_TEST",
	}
	_, err = szConfigServer.DeleteDataSource(ctx, requestToDeleteDataSource)
	require.NoError(test, err)

	// ListDataSources #3.
	responseFromGetDataSources3, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	require.NoError(test, err)
	printActual(test, responseFromGetDataSources3.GetResult())

	// Close.
	requestToCloseConfig := &szpb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	require.NoError(test, err)

	assert.Equal(test, initialDataSources, responseFromGetDataSources3.GetResult())
}

func TestSzConfigServer_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &szpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	require.NoError(test, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// ListDataSources.
	requestToGetDataSources := &szpb.GetDataSourcesRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromGetDataSources, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	require.NoError(test, err)
	printActual(test, responseFromGetDataSources.GetResult())

	// Close.
	requestToCloseConfig := &szpb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	require.NoError(test, err)
}

func TestSzConfigServer_ImportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreate := &szpb.CreateConfigRequest{}
	responseFromCreate, err := szConfigServer.CreateConfig(ctx, requestToCreate)
	require.NoError(test, err)
	printActual(test, responseFromCreate.GetResult())

	// Export Config to string.
	requestToExportConfig := &szpb.ExportConfigRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromExportConfig, err := szConfigServer.ExportConfig(ctx, requestToExportConfig)
	require.NoError(test, err)
	printActual(test, responseFromExportConfig.GetResult())

	// Close.
	requestToCloseConfig := &szpb.CloseConfigRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	require.NoError(test, err)

	// Load.
	requestToImportConfig := &szpb.ImportConfigRequest{
		ConfigDefinition: responseFromExportConfig.GetResult(),
	}
	responseFromLoad, err := szConfigServer.ImportConfig(ctx, requestToImportConfig)
	require.NoError(test, err)
	printActual(test, responseFromLoad.GetResult())

	// Close.
	requestToCloseConfig = &szpb.CloseConfigRequest{
		ConfigHandle: responseFromLoad.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	require.NoError(test, err)
}

func TestSzConfigServer_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &szpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	require.NoError(test, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// Save.
	requestToExportConfig := &szpb.ExportConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromExportConfig, err := szConfigServer.ExportConfig(ctx, requestToExportConfig)
	require.NoError(test, err)
	printActual(test, responseFromExportConfig.GetResult())

	// Close.
	requestToCloseConfig := &szpb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	require.NoError(test, err)
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
// Internal functions
// ----------------------------------------------------------------------------

func getSzConfigServer(ctx context.Context) SzConfigServer {
	if szConfigTestSingleton == nil {
		szConfigTestSingleton = &SzConfigServer{}
		instanceName := "Test instance name"
		verboseLogging := senzing.SzNoLogging
		settings, err := settings.BuildSimpleSettingsUsingEnvVars()
		panicOnError(err)
		osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
		if len(osenvLogLevel) > 0 {
			logLevelName = osenvLogLevel
		}
		err = szConfigTestSingleton.SetLogLevel(ctx, logLevelName)
		panicOnError(err)
		err = GetSdkSzConfig().Initialize(ctx, instanceName, settings, verboseLogging)
		panicOnError(err)
	}
	return *szConfigTestSingleton
}

func getTestObject(ctx context.Context, test *testing.T) SzConfigServer {
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

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		if errors.Is(err, szerror.ErrSzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzBadInput) {
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

func setup() error {
	var err error
	return err
}

func teardown() error {
	var err error
	return err
}

func TestBuildSimpleSystemConfigurationJsonUsingEnvVars(test *testing.T) {
	actual, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printActual(test, actual)
}

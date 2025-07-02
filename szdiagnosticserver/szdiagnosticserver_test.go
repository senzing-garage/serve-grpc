package szdiagnosticserver_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-grpc/szconfigmanagerserver"
	"github.com/senzing-garage/serve-grpc/szdiagnosticserver"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
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
	badFeatureID    = int64(-1)
	badLogLevelName = "BadLogLevelName"
	badSecondsToRun = -1
)

// Nil/empty parameters

var (
	nilSecondsToRun int32
	nilFeatureID    int64
)

var (
	logLevelName      = "INFO"
	observerSingleton = &observer.NullObserver{
		ID:       observerID,
		IsSilent: true,
	}
	szConfigManagerServerSingleton *szconfigmanagerserver.SzConfigManagerServer
	szDiagnosticServerSingleton    *szdiagnosticserver.SzDiagnosticServer
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzDiagnosticServer_CheckRepositoryPerformance(test *testing.T) {
	ctx := test.Context()
	szDiagnosticServer := getTestObject(ctx, test)
	request := &szpb.CheckRepositoryPerformanceRequest{
		SecondsToRun: int32(1),
	}
	response, err := szDiagnosticServer.CheckRepositoryPerformance(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzDiagnosticServer_CheckRepositoryPerformance_badSecondsToRun(test *testing.T) {
	ctx := test.Context()
	szDiagnosticServer := getTestObject(ctx, test)
	request := &szpb.CheckRepositoryPerformanceRequest{
		SecondsToRun: badSecondsToRun,
	}
	response, err := szDiagnosticServer.CheckRepositoryPerformance(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzDiagnosticServer_CheckRepositoryPerformance_nilSecondsToRun(test *testing.T) {
	ctx := test.Context()
	szDiagnosticServer := getTestObject(ctx, test)
	request := &szpb.CheckRepositoryPerformanceRequest{
		SecondsToRun: nilSecondsToRun,
	}
	response, err := szDiagnosticServer.CheckRepositoryPerformance(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzDiagnosticServer_GetRepositoryInfo(test *testing.T) {
	ctx := test.Context()
	szDiagnosticServer := getTestObject(ctx, test)
	request := &szpb.GetRepositoryInfoRequest{}
	response, err := szDiagnosticServer.GetRepositoryInfo(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzDiagnosticServer_GetFeature(test *testing.T) {
	ctx := test.Context()
	szDiagnosticServer := getTestObject(ctx, test)
	request := &szpb.GetFeatureRequest{
		FeatureId: int64(1),
	}
	response, err := szDiagnosticServer.GetFeature(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzDiagnosticServer_GetFeature_badFeatureID(test *testing.T) {
	ctx := test.Context()
	szDiagnosticServer := getTestObject(ctx, test)
	request := &szpb.GetFeatureRequest{
		FeatureId: badFeatureID,
	}
	response, err := szDiagnosticServer.GetFeature(ctx, request)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szdiagnosticserver.(*SzDiagnosticServer).GetFeature","error":{"function":"szdiagnostic.(*Szdiagnostic).GetFeature","error":{"id":"SZSDK60034004","reason":"SENZ0057|Unknown feature ID value '-1'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, response)
}

func TestSzDiagnosticServer_GetFeature_nilFeatureID(test *testing.T) {
	ctx := test.Context()
	szDiagnosticServer := getTestObject(ctx, test)
	request := &szpb.GetFeatureRequest{
		FeatureId: nilFeatureID,
	}
	response, err := szDiagnosticServer.GetFeature(ctx, request)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szdiagnosticserver.(*SzDiagnosticServer).GetFeature","error":{"function":"szdiagnostic.(*Szdiagnostic).GetFeature","error":{"id":"SZSDK60034004","reason":"SENZ0057|Unknown feature ID value '0'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
	printActual(test, response)
}

// func TestSzDiagnosticServer_PurgeRepository(test *testing.T) {
// 	ctx := test.Context()
// 	szDiagnosticServer := getTestObject(ctx, test)
// 	request := &szpb.PurgeRepositoryRequest{}
// 	response, err := szDiagnosticServer.PurgeRepository(ctx, request)
// 	require.NoError(test, err)
// 	printActual(test, response)
// }

func TestSzDiagnosticServer_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(ctx, test)
	szConfigManager := getSzConfigManagerServer(ctx)
	getDefaultConfigIDRequest := &szconfigmanagerpb.GetDefaultConfigIdRequest{}
	getDefaultConfigIDResponse, err := szConfigManager.GetDefaultConfigId(ctx, getDefaultConfigIDRequest)
	printError(test, err)
	require.NoError(test, err)

	request := &szpb.ReinitializeRequest{
		ConfigId: getDefaultConfigIDResponse.GetResult(),
	}
	response, err := szDiagnostic.Reinitialize(ctx, request)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, response)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzDiagnosticServer_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzDiagnosticServer_SetLogLevel(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	printError(test, err)
	require.NoError(test, err)
}

func TestSzDiagnosticServer__SetLogLevel_badLevelName(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	printError(test, err)
	require.Error(test, err)
}

func TestSzDiagnosticServer_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	testObject.SetObserverOrigin(ctx, observerOrigin)
}

func TestSzDiagnosticServer_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	actual := testObject.GetObserverOrigin(ctx)
	assert.Equal(test, observerOrigin, actual)
}

func TestSzDiagnosticServer_UnregisterObserver(test *testing.T) {
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

func getSzConfigManagerServer(ctx context.Context) szconfigmanagerserver.SzConfigManagerServer {
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

	return *szConfigManagerServerSingleton
}

func getSzDiagnosticServer(ctx context.Context) *szdiagnosticserver.SzDiagnosticServer {
	if szDiagnosticServerSingleton == nil {
		szDiagnosticServerSingleton = &szdiagnosticserver.SzDiagnosticServer{}
		instanceName := "Test instance name"
		verboseLogging := senzing.SzNoLogging
		configID := senzing.SzInitializeWithDefaultConfiguration
		settings, err := settings.BuildSimpleSettingsUsingEnvVars()
		panicOnError(err)

		osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
		if len(osenvLogLevel) > 0 {
			logLevelName = osenvLogLevel
		}

		err = szDiagnosticServerSingleton.SetLogLevel(ctx, logLevelName)
		panicOnError(err)
		err = szdiagnosticserver.GetSdkSzDiagnostic().Initialize(ctx, instanceName, settings, configID, verboseLogging)
		panicOnError(err)
	}

	return szDiagnosticServerSingleton
}

func getTestObject(ctx context.Context, t *testing.T) *szdiagnosticserver.SzDiagnosticServer {
	t.Helper()

	return getSzDiagnosticServer(ctx)
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
		_, err := szConfig.RegisterDataSource(ctx, dataSourceCode)
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

func setupAddRecords(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	szEngine := &szengine.Szengine{}
	configID := senzing.SzInitializeWithDefaultConfiguration
	err := szEngine.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	panicOnError(err)

	flags := senzing.SzNoLogging

	testRecordIDs := []string{"1001", "1002", "1003", "1004", "1005", "1039", "1040"}
	for _, testRecordID := range testRecordIDs {
		testRecord := truthset.CustomerRecords[testRecordID]
		_, err := szEngine.AddRecord(ctx, testRecord.DataSource, testRecord.ID, testRecord.JSON, flags)
		panicOnError(err)
	}

	err = szEngine.Destroy(ctx)
	panicOnError(err)

	return wraperror.Errorf(err, wraperror.NoMessage)
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
	err = setupAddRecords(ctx, instanceName, settings, verboseLogging)
	panicOnError(err)
}

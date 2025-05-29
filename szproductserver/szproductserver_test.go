package szproductserver_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-grpc/szproductserver"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultTruncation = 76
	jsonIndentation   = "    "
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
	szProductTestSingleton *szproductserver.SzProductServer
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzProductServer_GetLicense(test *testing.T) {
	ctx := test.Context()
	szProductServer := getTestObject(ctx, test)
	request := &szpb.GetLicenseRequest{}
	actual, err := szProductServer.GetLicense(ctx, request)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzProductServer_GetVersion(test *testing.T) {
	ctx := test.Context()
	szProductServer := getTestObject(ctx, test)
	request := &szpb.GetVersionRequest{}
	actual, err := szProductServer.GetVersion(ctx, request)
	require.NoError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzProductServer_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

func TestSzProductServer_SetLogLevel(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	require.NoError(test, err)
}

func TestSzProductServer__SetLogLevel_badLevelName(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	require.Error(test, err)
}

func TestSzProductServer_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	testObject.SetObserverOrigin(ctx, observerOrigin)
}

func TestSzProductServer_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	actual := testObject.GetObserverOrigin(ctx)
	assert.Equal(test, observerOrigin, actual)
}

func TestSzProductServer_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
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

func getSzProductServer(ctx context.Context) *szproductserver.SzProductServer {
	if szProductTestSingleton == nil {
		szProductTestSingleton = &szproductserver.SzProductServer{}
		instanceName := "Test instance name"
		verboseLogging := senzing.SzNoLogging
		settings, err := settings.BuildSimpleSettingsUsingEnvVars()
		panicOnError(err)

		osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
		if len(osenvLogLevel) > 0 {
			logLevelName = osenvLogLevel
		}

		err = szProductTestSingleton.SetLogLevel(ctx, logLevelName)
		panicOnError(err)
		err = szproductserver.GetSdkSzProduct().Initialize(ctx, instanceName, settings, verboseLogging)
		panicOnError(err)
	}

	return szProductTestSingleton
}

func getTestObject(ctx context.Context, t *testing.T) *szproductserver.SzProductServer {
	t.Helper()

	return getSzProductServer(ctx)
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

func printResult(t *testing.T, title string, result interface{}) {
	t.Helper()

	if printResults {
		t.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

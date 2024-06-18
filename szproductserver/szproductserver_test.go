package szproductserver

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
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
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
	szProductTestSingleton *SzProductServer
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzProductServer_GetLicense(test *testing.T) {
	ctx := context.TODO()
	szProductServer := getTestObject(ctx, test)
	request := &szpb.GetLicenseRequest{}
	actual, err := szProductServer.GetLicense(ctx, request)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzProductServer_GetVersion(test *testing.T) {
	ctx := context.TODO()
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
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

func TestSzProductServer_SetLogLevel(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	require.NoError(test, err)
}

func TestSzProductServer__SetLogLevel_badLevelName(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	require.Error(test, err)
}

func TestSzProductServer_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	testObject.SetObserverOrigin(ctx, observerOrigin)
}

func TestSzProductServer_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	actual := testObject.GetObserverOrigin(ctx)
	assert.Equal(test, observerOrigin, actual)
}

func TestSzProductServer_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getSzProductServer(ctx context.Context) SzProductServer {
	if szProductTestSingleton == nil {
		szProductTestSingleton = &SzProductServer{}
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
		err = GetSdkSzProduct().Initialize(ctx, instanceName, settings, verboseLogging)
		panicOnError(err)

	}
	return *szProductTestSingleton
}

func getTestObject(ctx context.Context, test *testing.T) SzProductServer {
	_ = test
	return getSzProductServer(ctx)
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

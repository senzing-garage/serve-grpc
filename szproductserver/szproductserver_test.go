package szproductserver

import (
	"context"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/settings"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
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
// Internal functions
// ----------------------------------------------------------------------------

func getSzProductServer(ctx context.Context) SzProductServer {
	if szProductTestSingleton == nil {
		szProductTestSingleton = &SzProductServer{}
		instanceName := "Test module name"
		verboseLogging := int64(0)
		settings, err := settings.BuildSimpleSettingsUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkSzProduct().Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *szProductTestSingleton
}

func getTestObject(ctx context.Context, test *testing.T) SzProductServer {
	if szProductTestSingleton == nil {
		szProductTestSingleton = &SzProductServer{}
		instanceName := "Test module name"
		verboseLogging := int64(0)
		settings, err := settings.BuildSimpleSettingsUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkSzProduct().Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *szProductTestSingleton
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

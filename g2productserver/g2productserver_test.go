package g2productserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	g2pb "github.com/senzing-garage/g2-sdk-proto/go/g2product"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2productTestSingleton *G2ProductServer
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) G2ProductServer {
	if g2productTestSingleton == nil {
		g2productTestSingleton = &G2ProductServer{}
		moduleName := "Test module name"
		verboseLogging := int64(0)
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkG2product().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *g2productTestSingleton
}

func getG2ProductServer(ctx context.Context) G2ProductServer {
	if g2productTestSingleton == nil {
		g2productTestSingleton = &G2ProductServer{}
		moduleName := "Test module name"
		verboseLogging := int64(0)
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkG2product().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *g2productTestSingleton
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, g2product G2ProductServer, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2product G2ProductServer, err error, messageId string) {
	if err != nil {
		var dictionary map[string]interface{}
		unmarshalErr := json.Unmarshal([]byte(err.Error()), &dictionary)
		if unmarshalErr != nil {
			test.Log("Unmarshal Error:", unmarshalErr.Error())
		}
		assert.Equal(test, messageId, dictionary["id"].(string))
	} else {
		assert.FailNow(test, "Should have failed with", messageId)
	}
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
	var err error = nil
	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestBuildSimpleSystemConfigurationJsonUsingEnvVars(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2productServer_Init(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int64(0),
	}
	response, err := g2product.Init(ctx, request)
	expectError(test, ctx, g2product, err, "senzing-60164002")
	printActual(test, response)
}

func TestG2productServer_License(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	request := &g2pb.LicenseRequest{}
	actual, err := g2product.License(ctx, request)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productServer_Version(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	request := &g2pb.VersionRequest{}
	actual, err := g2product.Version(ctx, request)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productServer_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	request := &g2pb.DestroyRequest{}
	actual, err := g2product.Destroy(ctx, request)
	expectError(test, ctx, g2product, err, "senzing-60164001")
	printActual(test, actual)
}

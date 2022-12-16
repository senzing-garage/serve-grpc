package g2productserver

import (
	"context"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	pb "github.com/senzing/go-servegrpc/protobuf/g2product"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
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
		verboseLogging := 0
		iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if jsonErr != nil {
			test.Logf("Cannot construct system configuration. Error: %v", jsonErr)
		}

		request := &pb.InitRequest{
			ModuleName:     moduleName,
			IniParams:      iniParams,
			VerboseLogging: int32(verboseLogging),
		}
		g2productTestSingleton.Init(ctx, request)
	}
	return *g2productTestSingleton
}

func getG2ProductServer(ctx context.Context) G2ProductServer {
	g2productserver := &G2ProductServer{}
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	request := &pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	g2productserver.Init(ctx, request)
	return *g2productserver
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if 1 == 0 {
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

func testErrorNoFail(test *testing.T, ctx context.Context, g2product G2ProductServer, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
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

func TestG2productserver_BuildSimpleSystemConfigurationJson(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
	request := &pb.InitRequest{}
	actual, err := g2product.Init(ctx, request)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productServer_License(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	request := &pb.LicenseRequest{}
	actual, err := g2product.License(ctx, request)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productServer_ValidateLicenseFile(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	request := &pb.ValidateLicenseFileRequest{}
	actual, err := g2product.ValidateLicenseFile(ctx, request)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productServer_ValidateLicenseStringBase64(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	request := &pb.ValidateLicenseStringBase64Request{}
	actual, err := g2product.ValidateLicenseStringBase64(ctx, request)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productServer_Version(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	request := &pb.VersionRequest{}
	actual, err := g2product.Version(ctx, request)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

func TestG2productServer_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2product := getTestObject(ctx, test)
	request := &pb.DestroyRequest{}
	actual, err := g2product.Destroy(ctx, request)
	testError(test, ctx, g2product, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2ProductServer_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2product/g2product_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &pb.DestroyRequest{}
	result, err := g2product.Destroy(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2ProductServer_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2product/g2product_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &pb.InitRequest{}
	result, err := g2product.Init(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2ProductServer_License() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2product/g2product_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &pb.LicenseRequest{}
	result, err := g2product.License(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"customer":"Senzing Public Test License","contract":"EVALUATION - support@senzing.com","issueDate":"2022-11-29","licenseType":"EVAL (Solely for non-productive use)","licenseLevel":"STANDARD","billing":"MONTHLY","expireDate":"2023-11-29","recordLimit":50000}
}

func ExampleG2ProductServer_ValidateLicenseFile() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2product/g2product_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &pb.ValidateLicenseFileRequest{}
	result, err := g2product.ValidateLicenseFile(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: Success
}

func ExampleG2ProductServer_ValidateLicenseStringBase64() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2product/g2product_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &pb.ValidateLicenseStringBase64Request{}
	result, err := g2product.ValidateLicenseStringBase64(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: Success
}

func ExampleG2ProductServer_Version() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2product/g2product_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &pb.VersionRequest{}
	result, err := g2product.Version(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"PRODUCT_NAME":"Senzing API","VERSION":"3.3.2","BUILD_VERSION":"3.3.2.22299","BUILD_DATE":"2022-10-26","BUILD_NUMBER":"2022_10_26__19_38","COMPATIBILITY_VERSION":{"CONFIG_VERSION":"10"},"SCHEMA_VERSION":{"ENGINE_SCHEMA_VERSION":"3.3","MINIMUM_REQUIRED_SCHEMA_VERSION":"3.0","MAXIMUM_REQUIRED_SCHEMA_VERSION":"3.99"}}
}

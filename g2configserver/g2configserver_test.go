package g2configserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go/g2error"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2configTestSingleton *G2ConfigServer
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) G2ConfigServer {
	if g2configTestSingleton == nil {
		g2configTestSingleton = &G2ConfigServer{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkG2config().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *g2configTestSingleton
}

func getG2ConfigServer(ctx context.Context) G2ConfigServer {
	if g2configTestSingleton == nil {
		g2configTestSingleton = &G2ConfigServer{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkG2config().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *g2configTestSingleton
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

func testError(test *testing.T, ctx context.Context, g2config G2ConfigServer, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2config G2ConfigServer, err error, messageId string) {
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
		if g2error.Is(err, g2error.G2Unrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2Retryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2BadUserInput) {
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

func TestG2configserver_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// AddDataSource.
	requestToAddDataSource := &g2pb.AddDataSourceRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		InputJson:    `{"DSRC_CODE": "GO_TEST"}`,
	}
	responseFromAddDataSource, err := g2config.AddDataSource(ctx, requestToAddDataSource)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromAddDataSource.GetResult())

	// Close.
	requestToClose := &g2pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Close(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// Close.
	requestToClose := &g2pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Create(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	requestToCreate := &g2pb.CreateRequest{}
	response, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, response.GetResult())
}

func TestG2configserver_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// ListDataSources #1.
	requestToListDataSources := &g2pb.ListDataSourcesRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromListDataSources, err := g2config.ListDataSources(ctx, requestToListDataSources)
	testError(test, ctx, g2config, err)
	listBefore := responseFromListDataSources.GetResult()
	printActual(test, listBefore)

	// AddDataSource.
	requestToAddDataSource := &g2pb.AddDataSourceRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		InputJson:    `{"DSRC_CODE": "GO_TEST"}`,
	}
	responseFromAddDataSource, err := g2config.AddDataSource(ctx, requestToAddDataSource)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromAddDataSource.GetResult())

	// ListDataSources #2.
	responseFromListDataSources2, err := g2config.ListDataSources(ctx, requestToListDataSources)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromListDataSources2.GetResult())

	// DeleteDataSource.
	requestToDeleteDataSource := &g2pb.DeleteDataSourceRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		InputJson:    `{"DSRC_CODE": "GO_TEST"}`,
	}
	_, err = g2config.DeleteDataSource(ctx, requestToDeleteDataSource)
	testError(test, ctx, g2config, err)

	// ListDataSources #3.
	responseFromListDataSources3, err := g2config.ListDataSources(ctx, requestToListDataSources)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromListDataSources3.GetResult())

	// Close.
	requestToClose := &g2pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)

	assert.Equal(test, listBefore, responseFromListDataSources3.GetResult())
}

func TestG2configserver_ListDataSources(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// ListDataSources.
	requestToList := &g2pb.ListDataSourcesRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromListDataSources, err := g2config.ListDataSources(ctx, requestToList)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromListDataSources.GetResult())

	// Close.
	requestToClose := &g2pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Load(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// Save.
	requestToSave := &g2pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromSave.GetResult())

	// Close.
	requestToClose := &g2pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)

	// Load.
	requestToLoad := &g2pb.LoadRequest{
		JsonConfig: responseFromSave.GetResult(),
	}
	responseFromLoad, err := g2config.Load(ctx, requestToLoad)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromLoad.GetResult())

	// Close.
	requestToClose = &g2pb.CloseRequest{
		ConfigHandle: responseFromLoad.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Save(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// Save.
	requestToSave := &g2pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromSave.GetResult())

	// Close.
	requestToClose := &g2pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Init(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	response, err := g2config.Init(ctx, request)
	expectError(test, ctx, g2config, err, "senzing-60114002")
	printActual(test, response)
}

func TestG2configserver_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	request := &g2pb.DestroyRequest{}
	response, err := g2config.Destroy(ctx, request)
	expectError(test, ctx, g2config, err, "senzing-60114001")
	printActual(test, response)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2ConfigServer_AddDataSource() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configserver/g2configserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.AddDataSourceRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		InputJson:    `{"DSRC_CODE": "GO_TEST"}`,
	}
	response, err := g2config.AddDataSource(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DSRC_ID":1001}
}

func ExampleG2ConfigServer_Close() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configserver/g2configserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2ConfigServer_Create() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configserver/g2configserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)
	request := &g2pb.CreateRequest{}
	response, err := g2config.Create(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2ConfigServer_DeleteDataSource() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configserver/g2configserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.DeleteDataSourceRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		InputJson:    `{"DSRC_CODE": "TEST"}`,
	}
	_, err = g2config.DeleteDataSource(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2ConfigServer_ListDataSources() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configserver/g2configserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.ListDataSourcesRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	response, err := g2config.ListDataSources(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleG2ConfigServer_Load() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configserver/g2configserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Save() to create a JSON string.
	requestToSave := &g2pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.LoadRequest{
		JsonConfig: responseFromSave.GetResult(),
	}
	responseFromLoad, err := g2config.Load(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseFromLoad.GetResult() > 0)
	// Output: true
}

func ExampleG2ConfigServer_Save() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configserver/g2configserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	response, err := g2config.Save(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"ADVANCED":"Yes","INTERNAL":"No"},...
}

func ExampleG2ConfigServer_Init() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configserver/g2configserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	response, err := g2config.Init(ctx, request)
	if err != nil {
		// This should produce a "senzing-60114002" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2ConfigServer_Destroy() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configserver/g2configserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)
	request := &g2pb.DestroyRequest{}
	response, err := g2config.Destroy(ctx, request)
	if err != nil {
		// This should produce a "senzing-60114001" error.
	}
	fmt.Println(response)
	// Output:
}

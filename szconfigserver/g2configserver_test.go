package szconfigserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2configTestSingleton *SzConfigServer
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) SzConfigServer {
	if g2configTestSingleton == nil {
		g2configTestSingleton = &SzConfigServer{}
		moduleName := "Test module name"
		verboseLogging := int64(0)
		iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
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

func getSzConfigServer(ctx context.Context) SzConfigServer {
	if g2configTestSingleton == nil {
		g2configTestSingleton = &SzConfigServer{}
		moduleName := "Test module name"
		verboseLogging := int64(0)
		iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
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

func testError(test *testing.T, ctx context.Context, g2config SzConfigServer, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2config SzConfigServer, err error, messageId string) {
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
		if szerror.Is(err, szerror.SzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzBadInput) {
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
	actual, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
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
	iniParams, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int64(0),
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

package g2configserver

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go/g2config"
	"github.com/senzing/g2-sdk-go/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2engine"
	"github.com/senzing/g2-sdk-go/testhelpers"
	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/servegrpc/protobuf/g2config"
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
		iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if jsonErr != nil {
			test.Logf("Cannot construct system configuration. Error: %v", jsonErr)
		}

		request := &pb.InitRequest{
			ModuleName:     moduleName,
			IniParams:      iniParams,
			VerboseLogging: int32(verboseLogging),
		}
		g2configTestSingleton.Init(ctx, request)
	}
	return *g2configTestSingleton
}

func getG2ConfigServer(ctx context.Context) G2ConfigServer {
	G2ConfigServer := &G2ConfigServer{}
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
	G2ConfigServer.Init(ctx, request)
	return *G2ConfigServer
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

func testErrorNoFail(test *testing.T, ctx context.Context, g2config G2ConfigServer, err error) {
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
	ctx := context.TODO()
	now := time.Now()
	moduleName := "Test module name"
	verboseLogging := 0
	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	// if err != nil {
	// 	return logger.Error(5901, err)
	// }

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		return logger.Error(5902, err)
	}

	// Add Data Sources to in-memory Senzing configuration.

	aG2config := &g2config.G2configImpl{}
	err = aG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return logger.Error(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
	if err != nil {
		return logger.Error(5907, err)
	}

	for _, testDataSource := range testhelpers.TestDataSources {
		_, err := aG2config.AddDataSource(ctx, configHandle, testDataSource.Data)
		if err != nil {
			return logger.Error(5908, err)
		}
	}

	configStr, err := aG2config.Save(ctx, configHandle)
	if err != nil {
		return logger.Error(5909, err)
	}

	err = aG2config.Close(ctx, configHandle)
	if err != nil {
		return logger.Error(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return logger.Error(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &g2configmgr.G2configmgrImpl{}
	err = aG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return logger.Error(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return logger.Error(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return logger.Error(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return logger.Error(5915, err)
	}

	// Purge repository.

	aG2engine := &g2engine.G2engineImpl{}
	err = aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return logger.Error(5903, err)
	}

	err = aG2engine.PurgeRepository(ctx)
	if err != nil {
		return logger.Error(5904, err)
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return logger.Error(5905, err)
	}

	// Add records.

	aG2engine = &g2engine.G2engineImpl{}
	err = aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return logger.Error(5916, err)
	}

	for _, testRecord := range testhelpers.TestRecords {
		err := aG2engine.AddRecord(ctx, testRecord.DataSource, testRecord.Id, testRecord.Data, testRecord.LoadId)
		if err != nil {
			return logger.Error(5917, err)
		}
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return logger.Error(5918, err)
	}

	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestG2configserver_BuildSimpleSystemConfigurationJson(test *testing.T) {
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

func TestG2configserver_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// AddDataSource.
	requestToAddDataSource := &pb.AddDataSourceRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		InputJson:    `{"DSRC_CODE": "GO_TEST"}`,
	}
	responseFromAddDataSource, err := g2config.AddDataSource(ctx, requestToAddDataSource)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromAddDataSource.GetResult())

	// Close.
	requestToClose := &pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Close(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// Close.
	requestToClose := &pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Create(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	requestToCreate := &pb.CreateRequest{}
	response, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, response.GetResult())
}

func TestG2configserver_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// ListDataSources #1.
	requestToListDataSources := &pb.ListDataSourcesRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromListDataSources, err := g2config.ListDataSources(ctx, requestToListDataSources)
	testError(test, ctx, g2config, err)
	listBefore := responseFromListDataSources.GetResult()
	printActual(test, listBefore)

	// AddDataSource.
	requestToAddDataSource := &pb.AddDataSourceRequest{
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
	requestToDeleteDataSource := &pb.DeleteDataSourceRequest{
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
	requestToClose := &pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)

	assert.Equal(test, listBefore, responseFromListDataSources3.GetResult())
}

func TestG2configserver_Init(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		assert.FailNow(test, jsonErr.Error())
	}
	request := &pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	response, err := g2config.Init(ctx, request)
	testError(test, ctx, g2config, err)
	printActual(test, response)
}

func TestG2configserver_ListDataSources(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// ListDataSources.
	requestToList := &pb.ListDataSourcesRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromListDataSources, err := g2config.ListDataSources(ctx, requestToList)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromListDataSources.GetResult())

	// Close.
	requestToClose := &pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Load(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// Save.
	requestToSave := &pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromSave.GetResult())

	// Load.
	requestToLoad := &pb.LoadRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		JsonConfig:   responseFromSave.GetResult(),
	}
	_, err = g2config.Load(ctx, requestToLoad)
	testError(test, ctx, g2config, err)

	// Close.
	requestToClose := &pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Save(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)

	// Create.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromCreate.GetResult())

	// Save.
	requestToSave := &pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	testError(test, ctx, g2config, err)
	printActual(test, responseFromSave.GetResult())

	// Close.
	requestToClose := &pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, requestToClose)
	testError(test, ctx, g2config, err)
}

func TestG2configserver_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	request := &pb.DestroyRequest{}
	response, err := g2config.Destroy(ctx, request)
	testError(test, ctx, g2config, err)
	printActual(test, response)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2ConfigServer_AddDataSource() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &pb.AddDataSourceRequest{
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
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2ConfigServer_Create() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)
	request := &pb.CreateRequest{}
	response, err := g2config.Create(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2ConfigServer_DeleteDataSource() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &pb.DeleteDataSourceRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		InputJson:    `{"DSRC_CODE": "TEST"}`,
	}
	_, err = g2config.DeleteDataSource(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2ConfigServer_Destroy() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)
	request := &pb.DestroyRequest{}
	response, err := g2config.Destroy(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2ConfigServer_Init() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	request := &pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	response, err := g2config.Init(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2ConfigServer_ListDataSources() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &pb.ListDataSourcesRequest{
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
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Save() to create a JSON string.
	requestToSave := &pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &pb.LoadRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		JsonConfig:   responseFromSave.GetResult(),
	}
	response, err := g2config.Load(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2ConfigServer_Save() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	response, err := g2config.Save(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"ADVANCED":"Yes","INTERNAL":"No"},...
}

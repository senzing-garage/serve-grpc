package g2configmgrserver

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
	"github.com/senzing/go-servegrpc/g2configserver"
	pbconfig "github.com/senzing/go-servegrpc/protobuf/g2config"
	pb "github.com/senzing/go-servegrpc/protobuf/g2configmgr"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2configmgrTestSingleton *G2ConfigmgrServer
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) G2ConfigmgrServer {
	if g2configmgrTestSingleton == nil {
		g2configmgrTestSingleton = &G2ConfigmgrServer{}

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
		g2configmgrTestSingleton.Init(ctx, request)
	}
	return *g2configmgrTestSingleton
}

func getG2ConfigServer(ctx context.Context) g2configserver.G2ConfigServer {
	G2ConfigServer := &g2configserver.G2ConfigServer{}
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	request := &pbconfig.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	G2ConfigServer.Init(ctx, request)
	return *G2ConfigServer
}

func getG2ConfigmgrServer(ctx context.Context) G2ConfigmgrServer {
	G2ConfigmgrServer := &G2ConfigmgrServer{}
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
	G2ConfigmgrServer.Init(ctx, request)
	return *G2ConfigmgrServer
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

func testError(test *testing.T, ctx context.Context, g2configmgr G2ConfigmgrServer, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, g2configmgr G2ConfigmgrServer, err error) {
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

	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestG2ConfigmgrServer_BuildSimpleSystemConfigurationJson(test *testing.T) {
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

func TestG2configmgrserver_AddConfig(test *testing.T) {
	ctx := context.TODO()
	now := time.Now()
	g2config := getG2ConfigServer(ctx)

	// G2config Create() to create a Senzing configuration.
	requestToCreate := &pbconfig.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		test.Log(err)
	}

	// G2config Save() to create a JSON string.
	requestToSave := &pbconfig.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	if err != nil {
		test.Log(err)
	}

	// Test.
	g2configmgr := getTestObject(ctx, test)
	request := &pb.AddConfigRequest{
		ConfigStr:      responseFromSave.GetResult(),
		ConfigComments: fmt.Sprintf("g2configmgrserver_test at %s", now.UTC()),
	}
	response, err := g2configmgr.AddConfig(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrImpl_GetConfig(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)

	// GetDefaultConfigID().
	requestToGetDefaultConfigID := &pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestToGetDefaultConfigID)
	testError(test, ctx, g2configmgr, err)
	printActual(test, responseFromGetDefaultConfigID)

	// Test.
	request := &pb.GetConfigRequest{
		ConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.GetConfig(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrImpl_GetConfigList(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	request := &pb.GetConfigListRequest{}
	response, err := g2configmgr.GetConfigList(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrImpl_GetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	request := &pb.GetDefaultConfigIDRequest{}
	response, err := g2configmgr.GetDefaultConfigID(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrImpl_Init(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		assert.FailNow(test, jsonErr.Error())
	}
	request := &pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	response, err := g2configmgr.Init(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrImpl_ReplaceDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)

	// GetDefaultConfigID()
	requestToGetDefaultConfigID := &pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestToGetDefaultConfigID)
	testError(test, ctx, g2configmgr, err)

	// Test.
	request := &pb.ReplaceDefaultConfigIDRequest{
		OldConfigID: responseFromGetDefaultConfigID.GetConfigID(),
		NewConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.ReplaceDefaultConfigID(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrImpl_SetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)

	// GetDefaultConfigID()
	requestToGetDefaultConfigID := &pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestToGetDefaultConfigID)
	testError(test, ctx, g2configmgr, err)

	// Test.
	request := &pb.SetDefaultConfigIDRequest{
		ConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.SetDefaultConfigID(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrImpl_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	request := &pb.DestroyRequest{}
	response, err := g2configmgr.Destroy(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2configmgrImpl_AddConfig() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	now := time.Now()
	g2config := getG2ConfigServer(ctx)

	// G2config Create() to create a Senzing configuration.
	requestToCreate := &pbconfig.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// G2config Save() to create a JSON string.
	requestToSave := &pbconfig.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	if err != nil {
		fmt.Println(err)
	}

	// Test.
	g2configmgr := getG2ConfigmgrServer(ctx)
	request := &pb.AddConfigRequest{
		ConfigStr:      responseFromSave.GetResult(),
		ConfigComments: fmt.Sprintf("g2configmgrserver_test at %s", now.UTC()),
	}
	response, err := g2configmgr.AddConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true

}

func ExampleG2configmgrImpl_Destroy() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)
	request := &pb.DestroyRequest{}
	result, err := g2configmgr.Destroy(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2configmgrImpl_GetConfig() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)
	requestToGetDefaultConfigID := &pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestToGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}
	request := &pb.GetConfigRequest{
		ConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.GetConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), defaultTruncation))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR...
}

func ExampleG2configmgrImpl_GetConfigList() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)
	request := &pb.GetConfigListRequest{}
	response, err := g2configmgr.GetConfigList(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleG2configmgrImpl_GetDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)
	request := &pb.GetDefaultConfigIDRequest{}
	response, err := g2configmgr.GetDefaultConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetConfigID() > 0) // Dummy output.
	// Output: true
}

func ExampleG2configmgrImpl_Init() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2config := getG2ConfigmgrServer(ctx)
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

func ExampleG2configmgrImpl_ReplaceDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	now := time.Now()
	g2config := getG2ConfigServer(ctx)
	g2configmgr := getG2ConfigmgrServer(ctx)

	// GetDefaultConfigID()
	requestForGetDefaultConfigID := &pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestForGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}

	// G2config Create() to create a Senzing configuration.
	requestToCreate := &pbconfig.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// G2config Save() to create a JSON string.
	requestToSave := &pbconfig.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	if err != nil {
		fmt.Println(err)
	}

	// AddConfig()
	requestForAddConfig := &pb.AddConfigRequest{
		ConfigStr:      responseFromSave.GetResult(),
		ConfigComments: fmt.Sprintf("g2configmgrserver_test at %s", now.UTC()),
	}
	responseFromAddConfig, err := g2configmgr.AddConfig(ctx, requestForAddConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Test.
	request := &pb.ReplaceDefaultConfigIDRequest{
		OldConfigID: responseFromGetDefaultConfigID.GetConfigID(),
		NewConfigID: responseFromAddConfig.GetResult(),
	}
	response, err := g2configmgr.ReplaceDefaultConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2configmgrImpl_SetDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/go-servegrpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)

	// GetDefaultConfigID()
	requestForGetDefaultConfigID := &pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestForGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}

	// Test.
	request := &pb.SetDefaultConfigIDRequest{
		ConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.SetDefaultConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

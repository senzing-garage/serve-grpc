package g2configmgrserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go-base/g2config"
	"github.com/senzing/g2-sdk-go-base/g2configmgr"
	"github.com/senzing/g2-sdk-go-base/g2engine"
	g2configpb "github.com/senzing/g2-sdk-proto/go/g2config"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-common/truthset"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/servegrpc/g2configserver"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2configmgrServerSingleton *G2ConfigmgrServer
	localLogger                messagelogger.MessageLoggerInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) G2ConfigmgrServer {
	if g2configmgrServerSingleton == nil {
		g2configmgrServerSingleton = &G2ConfigmgrServer{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkG2configmgr().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *g2configmgrServerSingleton
}

func getG2ConfigmgrServer(ctx context.Context) G2ConfigmgrServer {
	if g2configmgrServerSingleton == nil {
		g2configmgrServerSingleton = &G2ConfigmgrServer{}
		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkG2configmgr().Init(ctx, moduleName, iniParams, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *g2configmgrServerSingleton
}

func getG2ConfigServer(ctx context.Context) g2configserver.G2ConfigServer {
	G2ConfigServer := &g2configserver.G2ConfigServer{}
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	err = g2configserver.GetSdkG2config().Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
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

func testError(test *testing.T, ctx context.Context, g2configmgr G2ConfigmgrServer, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2configmgr G2ConfigmgrServer, err error, messageId string) {
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

func setupSenzingConfig(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	now := time.Now()

	aG2config := &g2config.G2config{}
	err := aG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
	if err != nil {
		return localLogger.Error(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := aG2config.AddDataSource(ctx, configHandle, datasource.Json)
		if err != nil {
			return localLogger.Error(5908, err)
		}
	}

	configStr, err := aG2config.Save(ctx, configHandle)
	if err != nil {
		return localLogger.Error(5909, err)
	}

	err = aG2config.Close(ctx, configHandle)
	if err != nil {
		return localLogger.Error(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return localLogger.Error(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &g2configmgr.G2configmgr{}
	err = aG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return localLogger.Error(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return localLogger.Error(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return localLogger.Error(5915, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	aG2engine := &g2engine.G2engine{}
	err := aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5903, err)
	}

	err = aG2engine.PurgeRepository(ctx)
	if err != nil {
		return localLogger.Error(5904, err)
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return localLogger.Error(5905, err)
	}
	return err
}

func setup() error {
	ctx := context.TODO()
	var err error = nil
	moduleName := "Test module name"
	verboseLogging := 0
	localLogger, err = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	if err != nil {
		return localLogger.Error(5901, err)
	}

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		return localLogger.Error(5902, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5920, err)
	}

	// Purge repository.

	err = setupPurgeRepository(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5921, err)
	}

	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestBuildSimpleSystemConfigurationJson(test *testing.T) {
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
	requestToCreate := &g2configpb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		test.Log(err)
	}

	// G2config Save() to create a JSON string.
	requestToSave := &g2configpb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	if err != nil {
		test.Log(err)
	}

	// Test.
	g2configmgr := getTestObject(ctx, test)
	request := &g2pb.AddConfigRequest{
		ConfigStr:      responseFromSave.GetResult(),
		ConfigComments: fmt.Sprintf("g2configmgrserver_test at %s", now.UTC()),
	}
	response, err := g2configmgr.AddConfig(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrserverImpl_GetConfig(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)

	// GetDefaultConfigID().
	requestToGetDefaultConfigID := &g2pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestToGetDefaultConfigID)
	testError(test, ctx, g2configmgr, err)
	printActual(test, responseFromGetDefaultConfigID)

	// Test.
	request := &g2pb.GetConfigRequest{
		ConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.GetConfig(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrserverImpl_GetConfigList(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	request := &g2pb.GetConfigListRequest{}
	response, err := g2configmgr.GetConfigList(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrserverImpl_GetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	request := &g2pb.GetDefaultConfigIDRequest{}
	response, err := g2configmgr.GetDefaultConfigID(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrserverImpl_ReplaceDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)

	// GetDefaultConfigID()
	requestToGetDefaultConfigID := &g2pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestToGetDefaultConfigID)
	testError(test, ctx, g2configmgr, err)

	// Test.
	request := &g2pb.ReplaceDefaultConfigIDRequest{
		OldConfigID: responseFromGetDefaultConfigID.GetConfigID(),
		NewConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.ReplaceDefaultConfigID(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrserverImpl_SetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)

	// GetDefaultConfigID()
	requestToGetDefaultConfigID := &g2pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestToGetDefaultConfigID)
	testError(test, ctx, g2configmgr, err)

	// Test.
	request := &g2pb.SetDefaultConfigIDRequest{
		ConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.SetDefaultConfigID(ctx, request)
	testError(test, ctx, g2configmgr, err)
	printActual(test, response)
}

func TestG2configmgrserverImpl_Init(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		assert.FailNow(test, err.Error())
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	response, err := g2configmgr.Init(ctx, request)
	expectError(test, ctx, g2configmgr, err, "senzing-60124002")
	printActual(test, response)
}

func TestG2configmgrserverImpl_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	request := &g2pb.DestroyRequest{}
	response, err := g2configmgr.Destroy(ctx, request)
	expectError(test, ctx, g2configmgr, err, "senzing-60124001")
	printActual(test, response)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2ConfigmgrServer_AddConfig() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2configmgrserver/g2configmgrserver_test.go
	ctx := context.TODO()
	now := time.Now()
	g2config := getG2ConfigServer(ctx)

	// G2config Create() to create a Senzing configuration.
	requestToCreate := &g2configpb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// G2config Save() to create a JSON string.
	requestToSave := &g2configpb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	g2configmgr := getG2ConfigmgrServer(ctx)
	request := &g2pb.AddConfigRequest{
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

func ExampleG2ConfigmgrServer_GetConfig() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2configmgrserver/g2configmgrserver_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)

	// GetDefaultConfigID() to get an example configuration ID.
	requestToGetDefaultConfigID := &g2pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestToGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.GetConfigRequest{
		ConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.GetConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), defaultTruncation))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR...
}

func ExampleG2ConfigmgrServer_GetConfigList() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2configmgrserver/g2configmgrserver_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)
	request := &g2pb.GetConfigListRequest{}
	response, err := g2configmgr.GetConfigList(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleG2ConfigmgrServer_GetDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2configmgrserver/g2configmgrserver_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)
	request := &g2pb.GetDefaultConfigIDRequest{}
	response, err := g2configmgr.GetDefaultConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetConfigID() > 0) // Dummy output.
	// Output: true
}

func ExampleG2ConfigmgrServer_ReplaceDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2configmgrserver/g2configmgrserver_test.go
	ctx := context.TODO()
	now := time.Now()
	g2config := getG2ConfigServer(ctx)
	g2configmgr := getG2ConfigmgrServer(ctx)

	// GetDefaultConfigID() to get the current configuration ID.
	requestForGetDefaultConfigID := &g2pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestForGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}

	// G2config Create() to create a Senzing configuration.
	requestToCreate := &g2configpb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// G2config Save() to create a JSON string.
	requestToSave := &g2configpb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	if err != nil {
		fmt.Println(err)
	}

	// AddConfig() to modify the configuration.
	requestForAddConfig := &g2pb.AddConfigRequest{
		ConfigStr:      responseFromSave.GetResult(),
		ConfigComments: fmt.Sprintf("g2configmgrserver_test at %s", now.UTC()),
	}
	responseFromAddConfig, err := g2configmgr.AddConfig(ctx, requestForAddConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.ReplaceDefaultConfigIDRequest{
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

func ExampleG2ConfigmgrServer_SetDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2configmgrserver/g2configmgrserver_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)

	// GetDefaultConfigID() to get an example configuration ID.
	requestForGetDefaultConfigID := &g2pb.GetDefaultConfigIDRequest{}
	responseFromGetDefaultConfigID, err := g2configmgr.GetDefaultConfigID(ctx, requestForGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.SetDefaultConfigIDRequest{
		ConfigID: responseFromGetDefaultConfigID.GetConfigID(),
	}
	response, err := g2configmgr.SetDefaultConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2ConfigmgrServer_Init() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2configmgrserver/g2configmgrserver_test.go
	ctx := context.TODO()
	g2config := getG2ConfigmgrServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
		// This should produce a "senzing-60124002" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2ConfigmgrServer_Destroy() {
	// For more information, visit https://github.com/Senzing/servegrpc/blob/main/g2configmgrserver/g2configmgrserver_test.go
	ctx := context.TODO()
	g2configmgr := getG2ConfigmgrServer(ctx)
	request := &g2pb.DestroyRequest{}
	response, err := g2configmgr.Destroy(ctx, request)
	if err != nil {
		// This should produce a "senzing-60124001" error.
	}
	fmt.Println(response)
	// Output:
}

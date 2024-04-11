package szconfigserver

import (
	"context"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	szConfigTestSingleton *SzConfigServer
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) SzConfigServer {
	if szConfigTestSingleton == nil {
		szConfigTestSingleton = &SzConfigServer{}
		instanceName := "Test name"
		verboseLogging := sz.SZ_NO_LOGGING
		settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Logf("Cannot construct system configuration. Error: %v", err)
		}
		err = GetSdkSzConfig().Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			test.Logf("Cannot Init. Error: %v", err)
		}
	}
	return *szConfigTestSingleton
}

func getSzConfigServer(ctx context.Context) SzConfigServer {
	if szConfigTestSingleton == nil {
		szConfigTestSingleton = &SzConfigServer{}
		instanceName := "Test name"
		verboseLogging := sz.SZ_NO_LOGGING
		settings, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			fmt.Println(err)
		}
		err = GetSdkSzConfig().Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return *szConfigTestSingleton
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

func testError(test *testing.T, ctx context.Context, err error) {
	_ = ctx
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

// func expectError(test *testing.T, ctx context.Context, err error, messageId string) {
// 	_ = ctx
// 	if err != nil {
// 		var dictionary map[string]interface{}
// 		unmarshalErr := json.Unmarshal([]byte(err.Error()), &dictionary)
// 		if unmarshalErr != nil {
// 			test.Log("Unmarshal Error:", unmarshalErr.Error())
// 		}
// 		assert.Equal(test, messageId, dictionary["id"].(string))
// 	} else {
// 		assert.FailNow(test, "Should have failed with", messageId)
// 	}
// }

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

func TestSzConfigServer_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &g2pb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	testError(test, ctx, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// AddDataSource.
	requestToAddDataSource := &g2pb.AddDataSourceRequest{
		ConfigHandle:   responseFromCreateConfig.GetResult(),
		DataSourceCode: "GO_TEST",
	}
	responseFromAddDataSource, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	testError(test, ctx, err)
	printActual(test, responseFromAddDataSource.GetResult())

	// Close.
	requestToCloseConfig := &g2pb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	testError(test, ctx, err)
}

func TestSzConfigServer_CloseConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &g2pb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	testError(test, ctx, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// Close.
	requestToCloseConfig := &g2pb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	testError(test, ctx, err)
}

func TestSzConfigServer_CreateConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)
	requestToCreate := &g2pb.CreateConfigRequest{}
	response, err := szConfigServer.CreateConfig(ctx, requestToCreate)
	testError(test, ctx, err)
	printActual(test, response.GetResult())
}

func TestSzConfigServer_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &g2pb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	testError(test, ctx, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// GetDataSources #1.
	requestToGetDataSources := &g2pb.GetDataSourcesRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromGetDataSources, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	testError(test, ctx, err)
	initialDataSources := responseFromGetDataSources.GetResult()
	printActual(test, initialDataSources)

	// AddDataSource.
	requestToAddDataSource := &g2pb.AddDataSourceRequest{
		ConfigHandle:   responseFromCreateConfig.GetResult(),
		DataSourceCode: "GO_TEST",
	}
	responseFromAddDataSource, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	testError(test, ctx, err)
	printActual(test, responseFromAddDataSource.GetResult())

	// GetDataSources #2.
	responseFromListDataSources2, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	testError(test, ctx, err)
	printActual(test, responseFromListDataSources2.GetResult())

	// DeleteDataSource.
	requestToDeleteDataSource := &g2pb.DeleteDataSourceRequest{
		ConfigHandle:   responseFromCreateConfig.GetResult(),
		DataSourceCode: "GO_TEST",
	}
	_, err = szConfigServer.DeleteDataSource(ctx, requestToDeleteDataSource)
	testError(test, ctx, err)

	// ListDataSources #3.
	responseFromGetDataSources3, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	testError(test, ctx, err)
	printActual(test, responseFromGetDataSources3.GetResult())

	// Close.
	requestToCloseConfig := &g2pb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	testError(test, ctx, err)

	assert.Equal(test, initialDataSources, responseFromGetDataSources3.GetResult())
}

func TestSzConfigServer_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &g2pb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	testError(test, ctx, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// ListDataSources.
	requestToGetDataSources := &g2pb.GetDataSourcesRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromGetDataSources, err := szConfigServer.GetDataSources(ctx, requestToGetDataSources)
	testError(test, ctx, err)
	printActual(test, responseFromGetDataSources.GetResult())

	// Close.
	requestToCloseConfig := &g2pb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	testError(test, ctx, err)
}

func TestSzConfigServer_ImportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreate := &g2pb.CreateConfigRequest{}
	responseFromCreate, err := szConfigServer.CreateConfig(ctx, requestToCreate)
	testError(test, ctx, err)
	printActual(test, responseFromCreate.GetResult())

	// Export Config to string.
	requestToExportConfig := &g2pb.ExportConfigRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromExportConfig, err := szConfigServer.ExportConfig(ctx, requestToExportConfig)
	testError(test, ctx, err)
	printActual(test, responseFromExportConfig.GetResult())

	// Close.
	requestToCloseConfig := &g2pb.CloseConfigRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	testError(test, ctx, err)

	// Load.
	requestToImportConfig := &g2pb.ImportConfigRequest{
		ConfigDefinition: responseFromExportConfig.GetResult(),
	}
	responseFromLoad, err := szConfigServer.ImportConfig(ctx, requestToImportConfig)
	testError(test, ctx, err)
	printActual(test, responseFromLoad.GetResult())

	// Close.
	requestToCloseConfig = &g2pb.CloseConfigRequest{
		ConfigHandle: responseFromLoad.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	testError(test, ctx, err)
}

func TestSzConfigServer_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigServer := getTestObject(ctx, test)

	// Create.
	requestToCreateConfig := &g2pb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	testError(test, ctx, err)
	printActual(test, responseFromCreateConfig.GetResult())

	// Save.
	requestToExportConfig := &g2pb.ExportConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromExportConfig, err := szConfigServer.ExportConfig(ctx, requestToExportConfig)
	testError(test, ctx, err)
	printActual(test, responseFromExportConfig.GetResult())

	// Close.
	requestToCloseConfig := &g2pb.CloseConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, requestToCloseConfig)
	testError(test, ctx, err)
}

//go:build linux

package szconfigmanagerserver

import (
	"context"
	"fmt"
	"time"

	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
)

// ----------------------------------------------------------------------------
// Interface functions - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzConfigManagerServer_GetConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigmanagerserver/szconfigmanagerserver_examples_test.go
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)

	// GetDefaultConfigId() to get an example configuration ID.
	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &szpb.GetConfigRequest{
		ConfigId: responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.GetConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), defaultTruncation))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR...
}

func ExampleSzConfigManagerServer_GetConfigs() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigmanagerserver/szconfigmanagerserver_examples_test.go
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	request := &szpb.GetConfigsRequest{}
	response, err := szConfigManagerServer.GetConfigs(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleSzConfigManagerServer_GetDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigmanagerserver/szconfigmanagerserver_examples_test.go
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	request := &szpb.GetDefaultConfigIdRequest{}
	response, err := szConfigManagerServer.GetDefaultConfigId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzConfigManagerServer_GetTemplateConfig() {
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	request := &szpb.GetTemplateConfigRequest{}
	response, err := szConfigManagerServer.GetTemplateConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: true
}

func ExampleSzConfigManagerServer_RegisterConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigmanagerserver/szconfigmanagerserver_examples_test.go

	now := time.Now()
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)
	szConfigManagerServer := getSzConfigManagerServer(ctx)

	// Get the template configuration.

	requestToGetTemplateConfig := &szpb.GetTemplateConfigRequest{}
	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Add DataSource to the Senzing configuration.

	requestToAddDataSource := &szconfigpb.AddDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}
	responseFromAddDataSource, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	if err != nil {
		fmt.Println(err)
	}

	// Test RegisterConfig.

	request := &szpb.RegisterConfigRequest{
		ConfigDefinition: responseFromAddDataSource.GetResult(),
		ConfigComment:    fmt.Sprintf("szconfigmanagerserver_test at %s", now.UTC()),
	}
	response, err := szConfigManagerServer.RegisterConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzConfigManagerServer_ReplaceDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigmanagerserver/szconfigmanagerserver_examples_test.go
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)

	// Get the ConfigID of the default Senzing configuration.

	requestToGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestToGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}

	// Test. Note: Cheating a little with replacing with same configId.

	request := &szpb.ReplaceDefaultConfigIdRequest{
		CurrentDefaultConfigId: responseFromGetDefaultConfigID.GetResult(),
		NewDefaultConfigId:     responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.ReplaceDefaultConfigId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleSzConfigManagerServer_SetDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigmanagerserver/szconfigmanagerserver_examples_test.go
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)

	// GetDefaultConfigId() to get an example configuration ID.
	requestForGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestForGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &szpb.SetDefaultConfigIdRequest{
		ConfigId: responseFromGetDefaultConfigID.GetResult(),
	}
	response, err := szConfigManagerServer.SetDefaultConfigId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

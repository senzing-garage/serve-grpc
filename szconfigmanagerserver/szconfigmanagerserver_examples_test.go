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

func ExampleSzConfigManagerServer_AddConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigmanagerserver/szconfigmanagerserver_examples_test.go
	ctx := context.TODO()
	now := time.Now()
	szConfigServer := getSzConfigServer(ctx)

	// SzConfig CreateConfig() to create a Senzing configuration.
	requestToCreateConfig := &szconfigpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// SzConfig ExportConfig() to create a JSON string.
	requestToExportConfig := &szconfigpb.ExportConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromExportConfig, err := szConfigServer.ExportConfig(ctx, requestToExportConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	request := &szpb.AddConfigRequest{
		ConfigDefinition: responseFromExportConfig.GetResult(),
		ConfigComment:    fmt.Sprintf("szconfigmanagerserver_test at %s", now.UTC()),
	}
	response, err := szConfigManagerServer.AddConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true

}

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

func ExampleSzConfigManagerServer_ReplaceDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigmanagerserver/szconfigmanagerserver_examples_test.go
	ctx := context.TODO()
	now := time.Now()
	szConfigServer := getSzConfigServer(ctx)
	szConfigManagerServer := getSzConfigManagerServer(ctx)

	// GetDefaultConfigId() to get the current configuration ID.
	requestForGetDefaultConfigID := &szpb.GetDefaultConfigIdRequest{}
	responseFromGetDefaultConfigID, err := szConfigManagerServer.GetDefaultConfigId(ctx, requestForGetDefaultConfigID)
	if err != nil {
		fmt.Println(err)
	}

	// SzConfig CreateConfig() to create a Senzing configuration.
	requestToCreateConfig := &szconfigpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// SzConfig ExportConfig() to create a JSON string.
	requestToExportConfig := &szconfigpb.ExportConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromExportConfig, err := szConfigServer.ExportConfig(ctx, requestToExportConfig)
	if err != nil {
		fmt.Println(err)
	}

	// AddConfig() to modify the configuration.
	requestForAddConfig := &szpb.AddConfigRequest{
		ConfigDefinition: responseFromExportConfig.GetResult(),
		ConfigComment:    fmt.Sprintf("szconfigmanagerserver_test at %s", now.UTC()),
	}
	responseFromAddConfig, err := szConfigManagerServer.AddConfig(ctx, requestForAddConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &szpb.ReplaceDefaultConfigIdRequest{
		CurrentDefaultConfigId: responseFromGetDefaultConfigID.GetResult(),
		NewDefaultConfigId:     responseFromAddConfig.GetResult(),
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

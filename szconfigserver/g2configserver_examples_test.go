//go:build linux

package szconfigserver

import (
	"context"
	"fmt"

	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzConfigServer_AddDataSource() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreateConfig := &szpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &szpb.AddDataSourceRequest{
		ConfigHandle:   responseFromCreateConfig.GetResult(),
		DataSourceCode: "GO_TEST",
	}
	response, err := szConfigServer.AddDataSource(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DSRC_ID":1001}
}

func ExampleSzConfigServer_CloseConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &szpb.CreateConfigRequest{}
	responseFromCreate, err := szConfigServer.CreateConfig(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &szpb.CloseConfigRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfigServer_CreateConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)
	request := &szpb.CreateConfigRequest{}
	response, err := szConfigServer.CreateConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzConfigServer_DeleteDataSource() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &szpb.CreateConfigRequest{}
	responseFromCreate, err := szConfigServer.CreateConfig(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &szpb.DeleteDataSourceRequest{
		ConfigHandle:   responseFromCreate.GetResult(),
		DataSourceCode: "TEST",
	}
	_, err = szConfigServer.DeleteDataSource(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfigServer_GetDataSources() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &szpb.CreateConfigRequest{}
	responseFromCreate, err := szConfig.CreateConfig(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &szpb.GetDataSourcesRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	response, err := szConfig.GetDataSources(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleSzConfigServer_ImportConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreateConfig := &szpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfig.CreateConfig(ctx, requestToCreateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Save() to create a JSON string.
	requestToExportConfig := &szpb.ExportConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromExportConfig, err := szConfig.ExportConfig(ctx, requestToExportConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &szpb.ImportConfigRequest{
		ConfigDefinition: responseFromExportConfig.GetResult(),
	}
	responseFromImportConfig, err := szConfig.ImportConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseFromImportConfig.GetResult() > 0)
	// Output: true
}

func ExampleSzConfigServer_ExportConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreateConfig := &szpb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfig.CreateConfig(ctx, requestToCreateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &szpb.ExportConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	response, err := szConfig.ExportConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"INTERNAL":"No"},{"ATTR_ID":1003,"...
}

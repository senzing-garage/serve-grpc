//go:build linux

package szconfigserver

import (
	"context"
	"fmt"

	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzConfigServer_AddDataSource() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreateConfig := &g2pb.CreateConfigRequest{}
	responseFromCreateConfig, err := szConfigServer.CreateConfig(ctx, requestToCreateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.AddDataSourceRequest{
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
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateConfigRequest{}
	responseFromCreate, err := szConfigServer.CreateConfig(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.CloseConfigRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = szConfigServer.CloseConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfigServer_CreateConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)
	request := &g2pb.CreateConfigRequest{}
	response, err := szConfigServer.CreateConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzConfigServer_DeleteDataSource() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	szConfigServer := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateConfigRequest{}
	responseFromCreate, err := szConfigServer.CreateConfig(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.DeleteDataSourceRequest{
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
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateConfigRequest{}
	responseFromCreate, err := g2config.CreateConfig(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.GetDataSourcesRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	response, err := g2config.GetDataSources(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleSzConfigServer_ImportConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreateConfig := &g2pb.CreateConfigRequest{}
	responseFromCreateConfig, err := g2config.CreateConfig(ctx, requestToCreateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Save() to create a JSON string.
	requestToExportConfig := &g2pb.ExportConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	responseFromExportConfig, err := g2config.ExportConfig(ctx, requestToExportConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.ImportConfigRequest{
		ConfigDefinition: responseFromExportConfig.GetResult(),
	}
	responseFromImportConfig, err := g2config.ImportConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseFromImportConfig.GetResult() > 0)
	// Output: true
}

func ExampleSzConfigServer_Save() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getSzConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreateConfig := &g2pb.CreateConfigRequest{}
	responseFromCreateConfig, err := g2config.CreateConfig(ctx, requestToCreateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.ExportConfigRequest{
		ConfigHandle: responseFromCreateConfig.GetResult(),
	}
	response, err := g2config.ExportConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"ADVANCED":"Yes","INTERNAL":"No"},...
}

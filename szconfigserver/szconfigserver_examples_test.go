//go:build linux

package szconfigserver_test

import (
	"context"
	"fmt"

	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
)

// ----------------------------------------------------------------------------
// Interface functions - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzConfigServer_AddDataSource() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getSzConfigServer(ctx)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}

	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Add DataSource to the Senzing configuration.

	requestToAddDataSource := &szpb.AddDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}

	response, err := szConfigServer.AddDataSource(ctx, requestToAddDataSource)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output: {"DSRC_ID":1001}
}

func ExampleSzConfigServer_DeleteDataSource() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getSzConfigServer(ctx)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}

	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Add DataSource to the Senzing configuration.

	requestToDeleteDataSource := &szpb.DeleteDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}

	response, err := szConfigServer.DeleteDataSource(ctx, requestToDeleteDataSource)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzConfigServer_GetDataSources() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
	ctx := context.TODO()
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	szConfigServer := getSzConfigServer(ctx)

	// Get the template configuration.

	requestToGetTemplateConfig := &szconfigmanagerpb.GetTemplateConfigRequest{}

	responseFromGetTemplateConfig, err := szConfigManagerServer.GetTemplateConfig(ctx, requestToGetTemplateConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Get Datasources in the Senzing configuration.

	requestToDeleteDataSource := &szpb.GetDataSourcesRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
	}

	response, err := szConfigServer.GetDataSources(ctx, requestToDeleteDataSource)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

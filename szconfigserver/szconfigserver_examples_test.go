//go:build linux

package szconfigserver_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzConfigServer_RegisterDataSource() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
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

	requestToRegisterDataSource := &szpb.RegisterDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}

	response, err := szConfigServer.RegisterDataSource(ctx, requestToRegisterDataSource)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "DSRC_ID": 1001
	// }
}

func ExampleSzConfigServer_UnregisterDataSource() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
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

	requestToUnregisterDataSource := &szpb.UnregisterDataSourceRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
		DataSourceCode:   "GO_TEST",
	}

	response, err := szConfigServer.UnregisterDataSource(ctx, requestToUnregisterDataSource)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzConfigServer_GetDataSourceRegistry() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szconfigserver/szconfigserver_examples_test.go
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

	requestToUnregisterDataSource := &szpb.GetDataSourceRegistryRequest{
		ConfigDefinition: responseFromGetTemplateConfig.GetResult(),
	}

	response, err := szConfigServer.GetDataSourceRegistry(ctx, requestToUnregisterDataSource)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCES": [
	//         {
	//             "DSRC_ID": 1,
	//             "DSRC_CODE": "TEST"
	//         },
	//         {
	//             "DSRC_ID": 2,
	//             "DSRC_CODE": "SEARCH"
	//         }
	//     ]
	// }
}

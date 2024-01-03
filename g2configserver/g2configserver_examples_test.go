//go:build linux

package g2configserver

import (
	"context"
	"fmt"

	g2pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2ConfigServer_AddDataSource() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.AddDataSourceRequest{
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
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.CloseRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	_, err = g2config.Close(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2ConfigServer_Create() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)
	request := &g2pb.CreateRequest{}
	response, err := g2config.Create(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2ConfigServer_DeleteDataSource() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.DeleteDataSourceRequest{
		ConfigHandle: responseFromCreate.GetResult(),
		InputJson:    `{"DSRC_CODE": "TEST"}`,
	}
	_, err = g2config.DeleteDataSource(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2ConfigServer_ListDataSources() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.ListDataSourcesRequest{
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
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Save() to create a JSON string.
	requestToSave := &g2pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	responseFromSave, err := g2config.Save(ctx, requestToSave)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.LoadRequest{
		JsonConfig: responseFromSave.GetResult(),
	}
	responseFromLoad, err := g2config.Load(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseFromLoad.GetResult() > 0)
	// Output: true
}

func ExampleG2ConfigServer_Save() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)

	// Create() to create an in-memory Senzing configuration.
	requestToCreate := &g2pb.CreateRequest{}
	responseFromCreate, err := g2config.Create(ctx, requestToCreate)
	if err != nil {
		fmt.Println(err)
	}

	// Example
	request := &g2pb.SaveRequest{
		ConfigHandle: responseFromCreate.GetResult(),
	}
	response, err := g2config.Save(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"ADVANCED":"Yes","INTERNAL":"No"},...
}

func ExampleG2ConfigServer_Init() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int64(0),
	}
	response, err := g2config.Init(ctx, request)
	if err != nil {
		// This should produce a "senzing-60114002" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2ConfigServer_Destroy() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2configserver/g2configserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigServer(ctx)
	request := &g2pb.DestroyRequest{}
	response, err := g2config.Destroy(ctx, request)
	if err != nil {
		// This should produce a "senzing-60114001" error.
	}
	fmt.Println(response)
	// Output:
}

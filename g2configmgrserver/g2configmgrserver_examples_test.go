//go:build linux

package g2configmgrserver

import (
	"context"
	"fmt"
	"time"

	g2configpb "github.com/senzing/g2-sdk-proto/go/g2config"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/go-common/g2engineconfigurationjson"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2ConfigmgrServer_AddConfig() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configmgrserver/g2configmgrserver_examples_test.go
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
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configmgrserver/g2configmgrserver_examples_test.go
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
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configmgrserver/g2configmgrserver_examples_test.go
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
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configmgrserver/g2configmgrserver_examples_test.go
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
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configmgrserver/g2configmgrserver_examples_test.go
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
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configmgrserver/g2configmgrserver_examples_test.go
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
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configmgrserver/g2configmgrserver_examples_test.go
	ctx := context.TODO()
	g2config := getG2ConfigmgrServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
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
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2configmgrserver/g2configmgrserver_examples_test.go
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

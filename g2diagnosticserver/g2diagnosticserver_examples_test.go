//go:build linux

package g2diagnosticserver

import (
	"context"
	"fmt"

	g2configmgrpb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2DiagnosticServer_CheckDBPerf() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.CheckDBPerfRequest{
		SecondsToRun: int32(1),
	}
	response, err := g2diagnostic.CheckDBPerf(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 25))
	// Output: {"numRecordsInserted":...
}

// func ExampleG2diagnosticImpl_CloseEntityListBySize() {

func ExampleG2DiagnosticServer_GetAvailableMemory() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetAvailableMemoryRequest{}
	response, err := g2diagnostic.GetAvailableMemory(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_GetDBInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetDBInfoRequest{}
	response, err := g2diagnostic.GetDBInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 52))
	// Output: {"Hybrid Mode":false,"Database Details":[{"Name":...
}

func ExampleG2DiagnosticServer_GetLogicalCores() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetLogicalCoresRequest{}
	response, err := g2diagnostic.GetLogicalCores(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_GetPhysicalCores() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetPhysicalCoresRequest{}
	response, err := g2diagnostic.GetPhysicalCores(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_GetTotalSystemMemory() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.GetTotalSystemMemoryRequest{}
	response, err := g2diagnostic.GetTotalSystemMemory(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2DiagnosticServer_Init() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := &G2DiagnosticServer{}
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int64(0),
	}
	response, err := g2diagnostic.Init(ctx, request)
	if err != nil {
		// This should produce a "senzing-60134002" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2DiagnosticServer_InitWithConfigID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := &G2DiagnosticServer{}
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	request := &g2pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   int64(1),
		VerboseLogging: int64(0),
	}
	response, err := g2diagnostic.InitWithConfigID(ctx, request)
	if err != nil {
		// This should produce a "senzing-60134003" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2DiagnosticServer_Reinit() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	g2configmgr := getG2ConfigmgrServer(ctx)
	getDefaultConfigIDRequest := &g2configmgrpb.GetDefaultConfigIDRequest{}
	getDefaultConfigIDResponse, err := g2configmgr.GetDefaultConfigID(ctx, getDefaultConfigIDRequest)
	if err != nil {
		fmt.Println(err)
	}
	request := &g2pb.ReinitRequest{
		InitConfigID: getDefaultConfigIDResponse.ConfigID,
	}
	_, err = g2diagnostic.Reinit(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2DiagnosticServer_Destroy() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2DiagnosticServer(ctx)
	request := &g2pb.DestroyRequest{}
	response, err := g2diagnostic.Destroy(ctx, request)
	if err != nil {
		// This should produce a "senzing-60134001" error.
	}
	fmt.Println(response)
	// Output:
}

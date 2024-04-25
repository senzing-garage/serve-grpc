//go:build linux

package szdiagnosticserver

import (
	"context"
	"fmt"

	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
)

// ----------------------------------------------------------------------------
// Interface functions - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzDiagnosticServer_CheckDatabasePerformance() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szdiagnosticserver/szdiagnosticserver_examples_test.go
	ctx := context.TODO()
	szDiagnosticServer := getSzDiagnosticServer(ctx)
	request := &szpb.CheckDatabasePerformanceRequest{
		SecondsToRun: int32(1),
	}
	response, err := szDiagnosticServer.CheckDatabasePerformance(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 25))
	// Output: {"numRecordsInserted":...
}

func ExampleSzDiagnosticServer_PurgeRepository() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szdiagnosticserver/szdiagnosticserver_test.go
	ctx := context.TODO()
	szDiagnosticServer := getSzDiagnosticServer(ctx)
	request := &szpb.PurgeRepositoryRequest{}
	response, err := szDiagnosticServer.PurgeRepository(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleSzDiagnosticServer_Reinitialize() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szdiagnosticserver/szdiagnosticserver_test.go
	ctx := context.TODO()
	szDiagnosticServer := getSzDiagnosticServer(ctx)
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	getDefaultConfigIdRequest := &szconfigmanagerpb.GetDefaultConfigIdRequest{}
	getDefaultConfigIdResponse, err := szConfigManagerServer.GetDefaultConfigId(ctx, getDefaultConfigIdRequest)
	if err != nil {
		fmt.Println(err)
	}
	request := &szpb.ReinitializeRequest{
		ConfigId: getDefaultConfigIdResponse.GetResult(),
	}
	_, err = szDiagnosticServer.Reinitialize(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

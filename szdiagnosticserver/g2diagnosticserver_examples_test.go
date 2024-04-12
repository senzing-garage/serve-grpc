//go:build linux

package szdiagnosticserver

import (
	"context"
	"fmt"

	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
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

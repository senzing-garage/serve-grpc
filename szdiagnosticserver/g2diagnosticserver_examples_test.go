//go:build linux

package szdiagnosticserver

import (
	"context"
	"fmt"

	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzDiagnosticServer_CheckDatabasePerformance() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2diagnosticserver/g2diagnosticserver_examples_test.go
	ctx := context.TODO()
	szDiagnosticServer := getSzDiagnosticServer(ctx)
	request := &g2pb.CheckDatabasePerformanceRequest{
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
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szDiagnosticServer := getSzDiagnosticServer(ctx)
	request := &g2pb.PurgeRepositoryRequest{}
	response, err := szDiagnosticServer.PurgeRepository(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

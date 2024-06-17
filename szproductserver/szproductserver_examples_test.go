//go:build linux

package szproductserver

import (
	"context"
	"fmt"

	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
)

// ----------------------------------------------------------------------------
// Interface functions - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzProductServer_GetLicense() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szproductserver/szproductserver_test.go
	ctx := context.TODO()
	szProductServer := getSzProductServer(ctx)
	request := &szpb.GetLicenseRequest{}
	response, err := szProductServer.GetLicense(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"customer":"Senzing Public Test License","contract":"Senzing Public Test License","issueDate":"2024-05-02","licenseType":"EVAL (Solely for non-productive use)","licenseLevel":"STANDARD","billing":"YEARLY","expireDate":"2025-05-02","recordLimit":50000}
}

func ExampleSzProductServer_GetVersion() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/szproductserver/szproductserver_test.go
	ctx := context.TODO()
	szProductServer := getSzProductServer(ctx)
	request := &szpb.GetVersionRequest{}
	response, err := szProductServer.GetVersion(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 43))
	// Output: {"PRODUCT_NAME":"Senzing API","VERSION":...
}

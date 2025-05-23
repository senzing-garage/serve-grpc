//go:build linux

package szproductserver_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
)

// ----------------------------------------------------------------------------
// Interface functions - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzProductServer_GetLicense() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szproductserver/szproductserver_test.go
	ctx := context.TODO()
	szProductServer := getSzProductServer(ctx)
	request := &szpb.GetLicenseRequest{}

	response, err := szProductServer.GetLicense(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output: {"customer":"Senzing Public Test License","contract":"Senzing Public Test License","issueDate":"2025-04-10","licenseType":"EVAL (Solely for non-productive use)","licenseLevel":"STANDARD","billing":"YEARLY","expireDate":"2026-04-10","recordLimit":50000}
}

func ExampleSzProductServer_GetVersion() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szproductserver/szproductserver_test.go
	ctx := context.TODO()
	szProductServer := getSzProductServer(ctx)
	request := &szpb.GetVersionRequest{}

	response, err := szProductServer.GetVersion(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output: {"PRODUCT_NAME":"Senzing SDK","VERSION":...
}

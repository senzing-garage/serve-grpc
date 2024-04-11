//go:build linux

package szproductserver

import (
	"context"
	"fmt"

	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzProductServer_GetLicense() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2productserver/g2productserver_test.go
	ctx := context.TODO()
	szProductServer := getSzProductServer(ctx)
	request := &g2pb.GetLicenseRequest{}
	response, err := szProductServer.GetLicense(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"customer":"Senzing Public Test License","contract":"Senzing Public Test - 50K records test","issueDate":"2023-11-02","licenseType":"EVAL (Solely for non-productive use)","licenseLevel":"STANDARD","billing":"YEARLY","expireDate":"2024-11-02","recordLimit":50000}
}

func ExampleSzProductServer_GetVersion() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2productserver/g2productserver_test.go
	ctx := context.TODO()
	szProductServer := getSzProductServer(ctx)
	request := &g2pb.GetVersionRequest{}
	response, err := szProductServer.GetVersion(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 43))
	// Output: {"PRODUCT_NAME":"Senzing API","VERSION":...
}

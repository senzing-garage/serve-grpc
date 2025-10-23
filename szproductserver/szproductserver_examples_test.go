//go:build linux

package szproductserver_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
)

const AllLines = -1

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
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

	redactKeys := []string{"issueDate", "expireDate", "BUILD_VERSION"}
	fmt.Println(jsonutil.PrettyPrint(jsonutil.Truncate(response.GetResult(), AllLines, redactKeys...), jsonIndentation))
	// Output:
	// {
	//     "advSearch": 0,
	//     "billing": "",
	//     "contract": "",
	//     "customer": "",
	//     "licenseLevel": "",
	//     "licenseType": "EVAL (Solely for non-productive use)",
	//     "recordLimit": 500
	// }
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

	// fmt.Println(truncate(response.GetResult(), 43))

	redactKeys := []string{"BUILD_DATE", "BUILD_NUMBER", "BUILD_VERSION", "ENGINE_SCHEMA_VERSION", "VERSION"}
	fmt.Println(jsonutil.PrettyPrint(jsonutil.Truncate(response.GetResult(), AllLines, redactKeys...), jsonIndentation))
	// Output:
	// {
	//     "COMPATIBILITY_VERSION": {
	//         "CONFIG_VERSION": "11"
	//     },
	//     "PRODUCT_NAME": "Senzing SDK",
	//     "SCHEMA_VERSION": {
	//         "MAXIMUM_REQUIRED_SCHEMA_VERSION": "4.99",
	//         "MINIMUM_REQUIRED_SCHEMA_VERSION": "4.0"
	//     }
	// }
}

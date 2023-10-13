//go:build linux

package g2productserver

import (
	"context"
	"fmt"

	g2pb "github.com/senzing/g2-sdk-proto/go/g2product"
	"github.com/senzing/go-common/g2engineconfigurationjson"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2ProductServer_Init() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2productserver/g2productserver_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int64(0),
	}
	response, err := g2product.Init(ctx, request)
	if err != nil {
		// This should produce a "senzing-60164002" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2ProductServer_License() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2productserver/g2productserver_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &g2pb.LicenseRequest{}
	response, err := g2product.License(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"customer":"Senzing Public Test License","contract":"EVALUATION - support@senzing.com","issueDate":"2022-11-29","licenseType":"EVAL (Solely for non-productive use)","licenseLevel":"STANDARD","billing":"MONTHLY","expireDate":"2023-11-29","recordLimit":50000}
}

func ExampleG2ProductServer_Version() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2productserver/g2productserver_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &g2pb.VersionRequest{}
	response, err := g2product.Version(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 43))
	// Output: {"PRODUCT_NAME":"Senzing API","VERSION":...
}

func ExampleG2ProductServer_Destroy() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2productserver/g2productserver_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &g2pb.DestroyRequest{}
	response, err := g2product.Destroy(ctx, request)
	if err != nil {
		// This should produce a "senzing-60164001" error.
	}
	fmt.Println(response)
	// Output:
}

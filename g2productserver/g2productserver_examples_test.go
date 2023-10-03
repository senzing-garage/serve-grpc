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
		VerboseLogging: int32(0),
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

func ExampleG2ProductServer_ValidateLicenseFile() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2productserver/g2productserver_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &g2pb.ValidateLicenseFileRequest{
		LicenseFilePath: licenseFilePath,
	}
	response, err := g2product.ValidateLicenseFile(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: Success
}

func ExampleG2ProductServer_ValidateLicenseStringBase64() {
	// For more information, visit https://github.com/Senzing/serve-grpc/blob/main/g2productserver/g2productserver_test.go
	ctx := context.TODO()
	g2product := getG2ProductServer(ctx)
	request := &g2pb.ValidateLicenseStringBase64Request{
		LicenseString: "AQAAADgCAAAAAAAAU2VuemluZyBQdWJsaWMgVGVzdCBMaWNlbnNlAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARVZBTFVBVElPTiAtIHN1cHBvcnRAc2VuemluZy5jb20AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADIwMjItMTEtMjkAAAAAAAAAAAAARVZBTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFNUQU5EQVJEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFDDAAAAAAAAMjAyMy0xMS0yOQAAAAAAAAAAAABNT05USExZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACQfw5e19QAHetkvd+vk0cYHtLaQCLmgx2WUfLorDfLQq15UXmOawNIXc1XguPd8zJtnOaeI6CB2smxVaj10mJE2ndGPZ1JjGk9likrdAj3rw+h6+C/Lyzx/52U8AuaN1kWgErDKdNE9qL6AnnN5LLi7Xs87opP7wbVMOdzsfXx2Xi3H7dSDIam7FitF6brSFoBFtIJac/V/Zc3b8jL/a1o5b1eImQldaYcT4jFrRZkdiVO/SiuLslEb8or3alzT0XsoUJnfQWmh0BjehBK9W74jGw859v/L1SGn1zBYKQ4m8JBiUOytmc9ekLbUKjIg/sCdmGMIYLywKqxb9mZo2TLZBNOpYWVwfaD/6O57jSixfJEHcLx30RPd9PKRO0Nm+4nPdOMMLmd4aAcGPtGMpI6ldTiK9hQyUfrvc9z4gYE3dWhz2Qu3mZFpaAEuZLlKtxaqEtVLWIfKGxwxPargPEfcLsv+30fdjSy8QaHeU638tj67I0uCEgnn5aB8pqZYxLxJx67hvVKOVsnbXQRTSZ00QGX1yTA+fNygqZ5W65wZShhICq5Fz8wPUeSbF7oCcE5VhFfDnSyi5v0YTNlYbF8LOAqXPTi+0KP11Wo24PjLsqYCBVvmOg9ohZ89iOoINwUB32G8VucRfgKKhpXhom47jObq4kSnihxRbTwJRx4o",
	}
	response, err := g2product.ValidateLicenseStringBase64(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: Success
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

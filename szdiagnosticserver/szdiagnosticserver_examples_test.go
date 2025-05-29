//go:build linux

package szdiagnosticserver_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
)

const AllLines = -1

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzDiagnosticServer_CheckDatastorePerformance() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szdiagnosticserver/szdiagnosticserver_examples_test.go
	ctx := context.TODO()
	szDiagnosticServer := getSzDiagnosticServer(ctx)
	request := &szpb.CheckDatastorePerformanceRequest{
		SecondsToRun: int32(1),
	}

	response, err := szDiagnosticServer.CheckDatastorePerformance(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	redactKeys := []string{"numRecordsInserted"}
	fmt.Println(jsonutil.PrettyPrint(jsonutil.Truncate(response.GetResult(), AllLines, redactKeys...), jsonIndentation))
	// Output:
	// {
	//     "insertTime": 1000
	// }
}

func ExampleSzDiagnosticServer_GetDatastoreInfo() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szdiagnosticserver/szdiagnosticserver_examples_test.go
	ctx := context.TODO()
	szDiagnosticServer := getSzDiagnosticServer(ctx)
	request := &szpb.GetDatastoreInfoRequest{}

	response, err := szDiagnosticServer.GetDatastoreInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "dataStores": [
	//         {
	//             "id": "CORE",
	//             "type": "sqlite3",
	//             "location": "nowhere"
	//         }
	//     ]
	// }
}

func ExampleSzDiagnosticServer_GetFeature() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szdiagnosticserver/szdiagnosticserver_examples_test.go
	ctx := context.TODO()
	szDiagnosticServer := getSzDiagnosticServer(ctx)
	request := &szpb.GetFeatureRequest{
		FeatureId: int64(1),
	}

	response, err := szDiagnosticServer.GetFeature(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "LIB_FEAT_ID": 1,
	//     "FTYPE_CODE": "NAME",
	//     "ELEMENTS": [
	//         {
	//             "FELEM_CODE": "FULL_NAME",
	//             "FELEM_VALUE": "Robert Smith"
	//         },
	//         {
	//             "FELEM_CODE": "SUR_NAME",
	//             "FELEM_VALUE": "Smith"
	//         },
	//         {
	//             "FELEM_CODE": "GIVEN_NAME",
	//             "FELEM_VALUE": "Robert"
	//         },
	//         {
	//             "FELEM_CODE": "CULTURE",
	//             "FELEM_VALUE": "ANGLO"
	//         },
	//         {
	//             "FELEM_CODE": "CATEGORY",
	//             "FELEM_VALUE": "PERSON"
	//         },
	//         {
	//             "FELEM_CODE": "TOKENIZED_NM",
	//             "FELEM_VALUE": "ROBERT|SMITH"
	//         }
	//     ]
	// }
}

func ExampleSzDiagnosticServer_PurgeRepository() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szdiagnosticserver/szdiagnosticserver_test.go
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
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szdiagnosticserver/szdiagnosticserver_test.go
	ctx := context.TODO()
	szDiagnosticServer := getSzDiagnosticServer(ctx)
	szConfigManagerServer := getSzConfigManagerServer(ctx)
	getDefaultConfigIDRequest := &szconfigmanagerpb.GetDefaultConfigIdRequest{}

	getDefaultConfigIDResponse, err := szConfigManagerServer.GetDefaultConfigId(ctx, getDefaultConfigIDRequest)
	if err != nil {
		fmt.Println(err)
	}

	request := &szpb.ReinitializeRequest{
		ConfigId: getDefaultConfigIDResponse.GetResult(),
	}

	_, err = szDiagnosticServer.Reinitialize(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

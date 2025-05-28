//go:build linux

package szengineserver_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func ExampleSzEngineServer_AddRecord() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.AddRecordRequest{
		DataSourceCode:   "CUSTOMERS",
		RecordId:         "1001",
		RecordDefinition: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		Flags:            senzing.SzWithoutInfo,
	}

	response, err := szEngineServer.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzEngineServer_AddRecord_secondRecord() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.AddRecordRequest{
		DataSourceCode:   "CUSTOMERS",
		RecordId:         "1002",
		RecordDefinition: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`,
		Flags:            senzing.SzWithoutInfo,
	}

	response, err := szEngineServer.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzEngineServer_AddRecord_withInfo() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.AddRecordRequest{
		DataSourceCode:   "CUSTOMERS",
		RecordId:         "1003",
		RecordDefinition: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`,
		Flags:            senzing.SzWithInfo,
	}

	response, err := szEngineServer.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCE": "CUSTOMERS",
	//     "RECORD_ID": "1003",
	//     "AFFECTED_ENTITIES": [
	//         {
	//             "ENTITY_ID": 100012
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_CloseExport() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)

	// Create a handle for the example.
	requestToExportJSONEntityReport := &szpb.ExportJsonEntityReportRequest{
		Flags: senzing.SzNoFlags,
	}

	responseFromExportJSONEntityReport, err := szEngineServer.ExportJsonEntityReport(
		ctx,
		requestToExportJSONEntityReport,
	)
	if err != nil {
		fmt.Println(err)
	}
	// Example
	request := &szpb.CloseExportRequest{
		ExportHandle: responseFromExportJSONEntityReport.GetResult(),
	}

	response, err := szEngineServer.CloseExport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response)
	// Output:
}

func ExampleSzEngineServer_CountRedoRecords() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.CountRedoRecordsRequest{}

	response, err := szEngineServer.CountRedoRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output: 3
}

func ExampleSzEngineServer_ExportCsvEntityReport() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.ExportCsvEntityReportRequest{
		CsvColumnList: "",
		Flags:         senzing.SzNoFlags,
	}

	response, err := szEngineServer.ExportCsvEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngineServer_ExportJsonEntityReport() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.ExportJsonEntityReportRequest{
		Flags: senzing.SzNoFlags,
	}

	response, err := szEngineServer.ExportJsonEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngineServer_FetchNext() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)

	// Create a handle for the example.
	requestToExportJSONEntityReport := &szpb.ExportJsonEntityReportRequest{
		Flags: senzing.SzNoFlags,
	}

	responseFromExportJSONEntityReport, err := szEngineServer.ExportJsonEntityReport(
		ctx,
		requestToExportJSONEntityReport,
	)
	if err != nil {
		fmt.Println(err)
	}
	// Example
	request := &szpb.FetchNextRequest{
		ExportHandle: responseFromExportJSONEntityReport.GetResult(),
	}

	response, err := szEngineServer.FetchNext(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(response.GetResult()) > 0) // Dummy output.
	// Output: false
}

func ExampleSzEngineServer_FindInterestingEntitiesByEntityId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngine := getSzEngineServer(ctx)
	request := &szpb.FindInterestingEntitiesByEntityIdRequest{
		EntityId: getEntityIDForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}

	response, err := szEngine.FindInterestingEntitiesByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "INTERESTING_ENTITIES": {
	//         "ENTITIES": []
	//     }
	// }
}

func ExampleSzEngineServer_FindInterestingEntitiesByRecordId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngine := getSzEngineServer(ctx)
	request := &szpb.FindInterestingEntitiesByRecordIdRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1001",
		Flags:          0,
	}

	response, err := szEngine.FindInterestingEntitiesByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "INTERESTING_ENTITIES": {
	//         "ENTITIES": []
	//     }
	// }
}

func ExampleSzEngineServer_FindNetworkByEntityId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord(
		"CUSTOMERS",
		"1001",
	) + `}, {"ENTITY_ID": ` + getEntityIDStringForRecord(
		"CUSTOMERS",
		"1002",
	) + `}]}`
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     1,
		BuildOutMaxEntities: 10,
		EntityIds:           entityIDs,
		Flags:               senzing.SzNoFlags,
		MaxDegrees:          2,
	}

	response, err := szEngineServer.FindNetworkByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100012
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_FindNetworkByRecordId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.FindNetworkByRecordIdRequest{
		BuildOutDegrees:     2,
		BuildOutMaxEntities: 10,
		Flags:               senzing.SzNoFlags,
		MaxDegrees:          1,
		RecordKeys:          `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`,
	}

	response, err := szEngineServer.FindNetworkByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100012
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_FindPathByEntityId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.FindPathByEntityIdRequest{
		StartEntityId: getEntityIDForRecord("CUSTOMERS", "1001"),
		EndEntityId:   getEntityIDForRecord("CUSTOMERS", "1002"),
		MaxDegrees:    1,
		Flags:         senzing.SzNoFlags,
	}

	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100012,
	//             "END_ENTITY_ID": 100012,
	//             "ENTITIES": [
	//                 100012
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100012
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_FindPathByEntityId_avoiding() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &szpb.FindPathByEntityIdRequest{
		StartEntityId:  getEntityIDForRecord("CUSTOMERS", "1001"),
		EndEntityId:    getEntityIDForRecord("CUSTOMERS", "1002"),
		MaxDegrees:     1,
		AvoidEntityIds: avoidEntityIDs,
		Flags:          senzing.SzNoFlags,
	}

	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100012,
	//             "END_ENTITY_ID": 100012,
	//             "ENTITIES": [
	//                 100012
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100012
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_FindPathByEntityId_avoidingAndIncluding() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &szpb.FindPathByEntityIdRequest{
		StartEntityId:       getEntityIDForRecord("CUSTOMERS", "1001"),
		EndEntityId:         getEntityIDForRecord("CUSTOMERS", "1002"),
		MaxDegrees:          1,
		AvoidEntityIds:      avoidEntityIDs,
		RequiredDataSources: `{"DATA_SOURCES": ["CUSTOMERS"]}`,
		Flags:               senzing.SzNoFlags,
	}

	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100012,
	//             "END_ENTITY_ID": 100012,
	//             "ENTITIES": []
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100012
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_FindPathByRecordId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.FindPathByRecordIdRequest{
		StartDataSourceCode: "CUSTOMERS",
		StartRecordId:       "1001",
		EndDataSourceCode:   "CUSTOMERS",
		EndRecordId:         "1002",
		MaxDegrees:          1,
		Flags:               senzing.SzNoFlags,
	}

	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100012,
	//             "END_ENTITY_ID": 100012,
	//             "ENTITIES": [
	//                 100012
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100012
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_FindPathByRecordId_avoiding() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.FindPathByRecordIdRequest{
		StartDataSourceCode: "CUSTOMERS",
		StartRecordId:       "1001",
		EndDataSourceCode:   "CUSTOMERS",
		EndRecordId:         "1002",
		MaxDegrees:          1,
		AvoidRecordKeys:     `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`,
		Flags:               senzing.SzNoFlags,
	}

	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100012,
	//             "END_ENTITY_ID": 100012,
	//             "ENTITIES": [
	//                 100012
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100012
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_FindPathByRecordId_avoidingAndIncluding() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.FindPathByRecordIdRequest{
		StartDataSourceCode: "CUSTOMERS",
		StartRecordId:       "1001",
		EndDataSourceCode:   "CUSTOMERS",
		EndRecordId:         "1002",
		MaxDegrees:          1,
		AvoidRecordKeys:     `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord("CUSTOMERS", "1003") + `}]}`,
		RequiredDataSources: `{"DATA_SOURCES": ["CUSTOMERS"]}`,
		Flags:               senzing.SzNoFlags,
	}

	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100012,
	//             "END_ENTITY_ID": 100012,
	//             "ENTITIES": []
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100012
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_GetActiveConfigId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.GetActiveConfigIdRequest{}

	response, err := szEngineServer.GetActiveConfigId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngineServer_GetEntityByEntityId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.GetEntityByEntityIdRequest{
		EntityId: getEntityIDForRecord("CUSTOMERS", "1001"),
		Flags:    senzing.SzNoFlags,
	}

	response, err := szEngineServer.GetEntityByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITY": {
	//         "ENTITY_ID": 100012
	//     }
	// }
}

func ExampleSzEngineServer_GetEntityByRecordId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.GetEntityByRecordIdRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1001",
		Flags:          senzing.SzNoFlags,
	}

	response, err := szEngineServer.GetEntityByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITY": {
	//         "ENTITY_ID": 100012
	//     }
	// }
}

func ExampleSzEngineServer_GetRecord() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.GetRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1001",
		Flags:          senzing.SzNoFlags,
	}

	response, err := szEngineServer.GetRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCE": "CUSTOMERS",
	//     "RECORD_ID": "1001"
	// }
}

func ExampleSzEngineServer_GetRedoRecord() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.GetRedoRecordRequest{}

	response, err := szEngineServer.GetRedoRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "REASON": "deferred delete",
	//     "DATA_SOURCE": "CUSTOMERS",
	//     "RECORD_ID": "1003",
	//     "REEVAL_ITERATION": 1,
	//     "_ENTITY_CORRUPTION_TRANSIENT": true,
	//     "DSRC_ACTION": "X"
	// }
}

func ExampleSzEngineServer_GetStats() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.GetStatsRequest{}

	response, err := szEngineServer.GetStats(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.Truncate(response.GetResult(), 5))
	// Output: {"workload":{"abortedUnresolve":0,"actualAmbiguousTest":0,"addedRecords":14,...
}

func ExampleSzEngineServer_GetVirtualEntityByRecordId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.GetVirtualEntityByRecordIdRequest{
		RecordKeys: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`,
		Flags:      senzing.SzNoFlags,
	}

	response, err := szEngineServer.GetVirtualEntityByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITY": {
	//         "ENTITY_ID": 100012
	//     }
	// }
}

func ExampleSzEngineServer_HowEntityByEntityId() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.HowEntityByEntityIdRequest{
		EntityId: getEntityIDForRecord("CUSTOMERS", "1001"),
		Flags:    senzing.SzNoFlags,
	}

	response, err := szEngineServer.HowEntityByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.Truncate(response.GetResult(), 5))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[...
}

func ExampleSzEngineServer_PreprocessRecord() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.PreprocessRecordRequest{
		RecordDefinition: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		Flags:            senzing.SzNoFlags,
	}

	response, err := szEngineServer.PreprocessRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output: {}
}

func ExampleSzEngineServer_PrimeEngine() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.PrimeEngineRequest{}

	response, err := szEngineServer.PrimeEngine(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response)
	// Output:
}

func ExampleSzEngineServer_ProcessRedoRecord() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	// IMPROVE: Document ExampleSzEngineServer_ProcessRedoRecord

	// Output:
}

func ExampleSzEngineServer_SearchByAttributes() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.SearchByAttributesRequest{
		Attributes: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`,
		Flags:      senzing.SzNoFlags,
	}

	response, err := szEngineServer.SearchByAttributes(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(
		jsonutil.PrettyPrint(
			jsonutil.Flatten(
				jsonutil.Redact(
					jsonutil.Flatten(jsonutil.NormalizeAndSort(response.GetResult())),
					"FIRST_SEEN_DT",
					"LAST_SEEN_DT",
				),
			), jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITIES": [
	//         {
	//             "ENTITY": {
	//                 "RESOLVED_ENTITY": {
	//                     "ENTITY_ID": 100012
	//                 }
	//             },
	//             "MATCH_INFO": {
	//                 "ERRULE_CODE": "SF1",
	//                 "MATCH_KEY": "+PNAME+EMAIL",
	//                 "MATCH_LEVEL_CODE": "POSSIBLY_RELATED"
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_SearchByAttributes_searchProfile() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	// IMPROVE: Fix SearchProfile value
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.SearchByAttributesRequest{
		Attributes:    `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`,
		SearchProfile: "SEARCH",
		Flags:         senzing.SzNoFlags,
	}

	response, err := szEngineServer.SearchByAttributes(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(
		jsonutil.PrettyPrint(
			jsonutil.Flatten(
				jsonutil.Redact(
					jsonutil.Flatten(jsonutil.NormalizeAndSort(response.GetResult())),
					"FIRST_SEEN_DT",
					"LAST_SEEN_DT",
				),
			), jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITIES": [
	//         {
	//             "ENTITY": {
	//                 "RESOLVED_ENTITY": {
	//                     "ENTITY_ID": 100012
	//                 }
	//             },
	//             "MATCH_INFO": {
	//                 "ERRULE_CODE": "SF1",
	//                 "MATCH_KEY": "+PNAME+EMAIL",
	//                 "MATCH_LEVEL_CODE": "POSSIBLY_RELATED"
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_WhyEntities() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.WhyEntitiesRequest{
		EntityId_1: getEntityIDForRecord("CUSTOMERS", "1001"),
		EntityId_2: getEntityIDForRecord("CUSTOMERS", "1002"),
		Flags:      senzing.SzNoFlags,
	}

	response, err := szEngineServer.WhyEntities(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "WHY_RESULTS": [
	//         {
	//             "ENTITY_ID": 100012,
	//             "ENTITY_ID_2": 100012,
	//             "MATCH_INFO": {
	//                 "WHY_KEY": "+NAME+DOB+ADDRESS+PHONE+EMAIL",
	//                 "WHY_ERRULE_CODE": "SF1_SNAME_CFF_CSTAB",
	//                 "MATCH_LEVEL_CODE": "RESOLVED"
	//             }
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100012
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_WhyRecords() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.WhyRecordsRequest{
		DataSourceCode_1: "CUSTOMERS",
		RecordId_1:       "1001",
		DataSourceCode_2: "CUSTOMERS",
		RecordId_2:       "1002",
		Flags:            senzing.SzNoFlags,
	}

	response, err := szEngineServer.WhyRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.Truncate(response.GetResult(), 7))
	// Output: {"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100012}}...
}

// ----------------------------------------------------------------------------
// Interface methods - Destructive, so order needs to be maintained
// ----------------------------------------------------------------------------

func ExampleSzEngineServer_ReevaluateEntity() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.ReevaluateEntityRequest{
		EntityId: getEntityIDForRecord("CUSTOMERS", "1001"),
		Flags:    senzing.SzWithoutInfo,
	}

	response, err := szEngineServer.ReevaluateEntity(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzEngineServer_ReevaluateEntity_withInfo() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	// IMPROVE: Fix Output
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.ReevaluateEntityRequest{
		EntityId: getEntityIDForRecord("CUSTOMERS", "1001"),
		Flags:    senzing.SzWithInfo,
	}

	response, err := szEngineServer.ReevaluateEntity(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "AFFECTED_ENTITIES": [
	//         {
	//             "ENTITY_ID": 100012
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_ReevaluateRecord() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.ReevaluateRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1001",
		Flags:          senzing.SzWithoutInfo,
	}

	response, err := szEngineServer.ReevaluateRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzEngineServer_ReevaluateRecord_withInfo() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	// IMPROVE: Fix Output
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.ReevaluateRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1001",
		Flags:          senzing.SzWithInfo,
	}

	response, err := szEngineServer.ReevaluateRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCE": "CUSTOMERS",
	//     "RECORD_ID": "1001",
	//     "AFFECTED_ENTITIES": [
	//         {
	//             "ENTITY_ID": 100012
	//         }
	//     ]
	// }
}

func ExampleSzEngineServer_Reinitialize() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)

	// Get a Senzing configuration ID for testing.
	requestToGetActiveConfigID := &szpb.GetActiveConfigIdRequest{}

	responseFromGetActiveConfigID, err := szEngineServer.GetActiveConfigId(ctx, requestToGetActiveConfigID)
	if err != nil {
		fmt.Println(err)
	}
	// Example
	request := &szpb.ReinitializeRequest{
		ConfigId: responseFromGetActiveConfigID.GetResult(),
	}

	response, err := szEngineServer.Reinitialize(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response)
	// Output:
}

func ExampleSzEngineServer_DeleteRecord() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.DeleteRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1003",
		Flags:          senzing.SzWithoutInfo,
	}

	response, err := szEngineServer.DeleteRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzEngineServer_DeleteRecord_withInfo() {
	// For more information,
	// visit https://github.com/senzing-garage/serve-grpc/blob/main/szengineserver/szengineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &szpb.DeleteRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1003",
		Flags:          senzing.SzWithInfo,
	}

	response, err := szEngineServer.DeleteRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonutil.PrettyPrint(response.GetResult(), jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCE": "CUSTOMERS",
	//     "RECORD_ID": "1003",
	//     "AFFECTED_ENTITIES": []
	// }
}

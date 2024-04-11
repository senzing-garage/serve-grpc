//go:build linux

package szengineserver

import (
	"context"
	"fmt"

	"github.com/senzing-garage/sz-sdk-go/sz"
	g2pb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzEngineServer_AddRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.AddRecordRequest{
		DataSourceCode:   "CUSTOMERS",
		RecordId:         "1001",
		RecordDefinition: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		Flags:            sz.SZ_WITHOUT_INFO,
	}
	response, err := szEngineServer.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {}
}

func ExampleSzEngineServer_AddRecord_secondRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.AddRecordRequest{
		DataSourceCode:   "CUSTOMERS",
		RecordId:         "1002",
		RecordDefinition: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`,
		Flags:            sz.SZ_WITHOUT_INFO,
	}
	response, err := szEngineServer.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {}
}

func ExampleSzEngineServer_AddRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.AddRecordRequest{
		DataSourceCode:   "CUSTOMERS",
		RecordId:         "1003",
		RecordDefinition: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`,
		Flags:            sz.SZ_WITH_INFO,
	}
	response, err := szEngineServer.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzEngineServer_CloseExport() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)

	// Create a handle for the example.
	requestToExportJsonEntityReport := &g2pb.ExportJsonEntityReportRequest{
		Flags: sz.SZ_NO_FLAGS,
	}
	responseFromExportJsonEntityReport, err := szEngineServer.ExportJsonEntityReport(ctx, requestToExportJsonEntityReport)
	if err != nil {
		fmt.Println(err)
	}
	// Example
	request := &g2pb.CloseExportRequest{
		ResponseHandle: responseFromExportJsonEntityReport.GetResult(),
	}
	response, err := szEngineServer.CloseExport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleSzEngineServer_CountRedoRecords() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.CountRedoRecordsRequest{}
	response, err := szEngineServer.CountRedoRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: 1
}

func ExampleSzEngineServer_ExportCsvEntityReport() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.ExportCsvEntityReportRequest{
		CsvColumnList: "",
		Flags:         sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.ExportCsvEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngineServer_ExportJsonEntityReport() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.ExportJsonEntityReportRequest{
		Flags: sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.ExportJsonEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngineServer_FetchNext() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)

	// Create a handle for the example.
	requestToExportJsonEntityReport := &g2pb.ExportJsonEntityReportRequest{
		Flags: sz.SZ_NO_FLAGS,
	}
	responseFromExportJsonEntityReport, err := szEngineServer.ExportJsonEntityReport(ctx, requestToExportJsonEntityReport)
	if err != nil {
		fmt.Println(err)
	}
	// Example
	request := &g2pb.FetchNextRequest{
		ResponseHandle: responseFromExportJsonEntityReport.GetResult(),
	}
	response, err := szEngineServer.FetchNext(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(response.GetResult()) >= 0) // Dummy output.
	// Output: true
}

func ExampleSzEngineServer_FindNetworkByEntityId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	request := &g2pb.FindNetworkByEntityIdRequest{
		EntityList:     entityList,
		MaxDegrees:     2,
		BuildOutDegree: 1,
		MaxEntities:    10,
		Flags:          sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.FindNetworkByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleSzEngineServer_FindNetworkByRecordId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.FindNetworkByRecordIdRequest{
		RecordList:     `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`,
		MaxDegrees:     1,
		BuildOutDegree: 2,
		MaxEntities:    10,
		Flags:          sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.FindNetworkByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleSzEngineServer_FindPathByEntityId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.FindPathByEntityIdRequest{
		StartEntityId: getEntityIdForRecord("CUSTOMERS", "1001"),
		EndEntityId:   getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegrees:    1,
		Flags:         sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngineServer_FindPathByEntityId_exclusions() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &g2pb.FindPathByEntityIdRequest{
		StartEntityId: getEntityIdForRecord("CUSTOMERS", "1001"),
		EndEntityId:   getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegrees:    1,
		Exclusions:    exclusions,
		Flags:         sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngineServer_FindPathByEntityId_inclusions() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &g2pb.FindPathByEntityIdRequest{
		StartEntityId:       getEntityIdForRecord("CUSTOMERS", "1001"),
		EndEntityId:         getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegrees:          1,
		Exclusions:          excludedEntities,
		RequiredDataSources: `{"DATA_SOURCES": ["CUSTOMERS"]}`,
		Flags:               sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 106))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngineServer_FindPathByRecordId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.FindPathByRecordIdRequest{
		StartDataSourceCode: "CUSTOMERS",
		StartRecordId:       "1001",
		EndDataSourceCode:   "CUSTOMERS",
		EndRecordId:         "1002",
		MaxDegrees:          1,
		Flags:               sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 87))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":...
}

func ExampleSzEngineServer_FindPathByRecordId_exclusions() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.FindPathByRecordIdRequest{
		StartDataSourceCode: "CUSTOMERS",
		StartRecordId:       "1001",
		EndDataSourceCode:   "CUSTOMERS",
		EndRecordId:         "1002",
		MaxDegrees:          1,
		Exclusions:          `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`,
		Flags:               sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngineServer_FindPathByRecordId_inclusions() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.FindPathByRecordIdRequest{
		StartDataSourceCode: "CUSTOMERS",
		StartRecordId:       "1001",
		EndDataSourceCode:   "CUSTOMERS",
		EndRecordId:         "1002",
		MaxDegrees:          1,
		Exclusions:          `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`,
		RequiredDataSources: `{"DATA_SOURCES": ["CUSTOMERS"]}`,
		Flags:               sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 119))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleSzEngineServer_GetActiveConfigId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.GetActiveConfigIdRequest{}
	response, err := szEngineServer.GetActiveConfigId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngineServer_GetEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.GetEntityByEntityIdRequest{
		EntityId: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.GetEntityByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleSzEngineServer_GetEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.GetEntityByRecordIdRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1001",
		Flags:          sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.GetEntityByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleSzEngineServer_GetRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.GetRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1001",
		Flags:          sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.GetRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}
}

func ExampleSzEngineServer_GetRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.GetRedoRecordRequest{}
	response, err := szEngineServer.GetRedoRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","DSRC_ACTION":"X"}
}

func ExampleSzEngineServer_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.GetRepositoryLastModifiedTimeRequest{}
	response, err := szEngineServer.GetRepositoryLastModifiedTime(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngineServer_GetStats() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.GetStatsRequest{}
	response, err := szEngineServer.GetStats(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 16))
	// Output: { "workload":...
}

func ExampleSzEngineServer_GetVirtualEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.GetVirtualEntityByRecordIdRequest{
		RecordList: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`,
		Flags:      sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.GetVirtualEntityByRecordId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleSzEngineServer_HowEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.HowEntityByEntityIdRequest{
		EntityId: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.HowEntityByEntityId(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL"}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1-S1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V4","MEMBER_RECORDS":[{"INTERNAL_ID":4,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","RESULT_VIRTUAL_ENTITY_ID":"V1-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB"}}],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1-S2","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":4,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleSzEngineServer_PrimeEngine() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.PrimeEngineRequest{}
	response, err := szEngineServer.PrimeEngine(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleSzEngineServer_SearchByAttributes() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.SearchByAttributesRequest{
		Attributes: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`,
		Flags:      sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.SearchByAttributes(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1"},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1}}}]}
}

func ExampleSzEngineServer_SearchByAttributes_searchProfile() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	// TODO: Fix SearchProfile value
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.SearchByAttributesRequest{
		Attributes:    `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`,
		SearchProfile: "SEARCH",
		Flags:         sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.SearchByAttributes(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 1962))
	// Output:
}

func ExampleSzEngineServer_WhyEntities() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.WhyEntitiesRequest{
		EntityId1: getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityId2: getEntityIdForRecord("CUSTOMERS", "1002"),
		Flags:     sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.WhyEntities(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 74))
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":1,"MATCH_INFO":{"WHY_KEY":...
}

func ExampleSzEngineServer_WhyRecords() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.WhyRecordsRequest{
		DataSourceCode1: "CUSTOMERS",
		RecordId1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordId2:       "1002",
		Flags:           sz.SZ_NO_FLAGS,
	}
	response, err := szEngineServer.WhyRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 115))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],...
}

func ExampleSzEngineServer_ReevaluateEntity() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.ReevaluateEntityRequest{
		EntityId: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    sz.SZ_WITHOUT_INFO,
	}
	response, err := szEngineServer.ReevaluateEntity(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzEngineServer_ReevaluateEntity_withInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.ReevaluateEntityRequest{
		EntityId: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    sz.SZ_WITH_INFO,
	}
	response, err := szEngineServer.ReevaluateEntity(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzEngineServer_ReevaluateRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.ReevaluateRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1001",
		Flags:          sz.SZ_WITHOUT_INFO,
	}
	response, err := szEngineServer.ReevaluateRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzEngineServer_ReevaluateRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.ReevaluateRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1001",
		Flags:          sz.SZ_WITH_INFO,
	}
	response, err := szEngineServer.ReevaluateRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzEngineServer_ReplaceRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.ReplaceRecordRequest{
		DataSourceCode:   "CUSTOMERS",
		RecordId:         "1001",
		RecordDefinition: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		Flags:            sz.SZ_WITHOUT_INFO,
	}
	response, err := szEngineServer.ReplaceRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output:
}

func ExampleSzEngineServer_ReplaceRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.ReplaceRecordRequest{
		DataSourceCode:   "CUSTOMERS",
		RecordId:         "1001",
		RecordDefinition: `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		Flags:            sz.SZ_WITH_INFO,
	}
	response, err := szEngineServer.ReplaceRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzEngineServer_DeleteRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.DeleteRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1003",
		Flags:          sz.SZ_WITHOUT_INFO,
	}
	response, err := szEngineServer.DeleteRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {}
}

func ExampleSzEngineServer_DeleteRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	szEngineServer := getSzEngineServer(ctx)
	request := &g2pb.DeleteRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordId:       "1003",
		Flags:          sz.SZ_WITH_INFO,
	}
	response, err := szEngineServer.DeleteRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

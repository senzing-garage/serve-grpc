//go:build linux

package g2engineserver

import (
	"context"
	"fmt"

	g2pb "github.com/senzing-garage/g2-sdk-proto/go/g2engine"
	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2EngineServer_AddRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.AddRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_AddRecord_secondRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.AddRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1002",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`,
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.AddRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_AddRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.AddRecordWithInfoRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1003",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`,
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.AddRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_CloseExport() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)

	// Create a handle for the example.
	requestToExportJSONEntityReport := &g2pb.ExportJSONEntityReportRequest{
		Flags: 0,
	}
	responseFromExportJSONEntityReport, err := g2engine.ExportJSONEntityReport(ctx, requestToExportJSONEntityReport)
	if err != nil {
		fmt.Println(err)
	}
	// Example
	request := &g2pb.CloseExportRequest{
		ResponseHandle: responseFromExportJSONEntityReport.GetResult(),
	}
	response, err := g2engine.CloseExport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_CountRedoRecords() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.CountRedoRecordsRequest{}
	response, err := g2engine.CountRedoRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: 1
}

func ExampleG2EngineServer_ExportCSVEntityReport() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ExportCSVEntityReportRequest{
		CsvColumnList: "",
		Flags:         0,
	}
	response, err := g2engine.ExportCSVEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_ExportConfig() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ExportConfigRequest{}
	response, err := g2engine.ExportConfig(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 42))
	// Output: {"G2_CONFIG":{"CFG_ETYPE":[{"ETYPE_ID":...
}

func ExampleG2EngineServer_ExportConfigAndConfigID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ExportConfigAndConfigIDRequest{}
	response, err := g2engine.ExportConfigAndConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetConfigID() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_ExportJSONEntityReport() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ExportJSONEntityReportRequest{
		Flags: 0,
	}
	response, err := g2engine.ExportJSONEntityReport(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_FetchNext() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)

	// Create a handle for the example.
	requestToExportJSONEntityReport := &g2pb.ExportJSONEntityReportRequest{
		Flags: 0,
	}
	responseFromExportJSONEntityReport, err := g2engine.ExportJSONEntityReport(ctx, requestToExportJSONEntityReport)
	if err != nil {
		fmt.Println(err)
	}
	// Example
	request := &g2pb.FetchNextRequest{
		ResponseHandle: responseFromExportJSONEntityReport.GetResult(),
	}
	response, err := g2engine.FetchNext(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(response.GetResult()) >= 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_FindInterestingEntitiesByEntityID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindInterestingEntitiesByEntityIDRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.FindInterestingEntitiesByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_FindInterestingEntitiesByRecordID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindInterestingEntitiesByRecordIDRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.FindInterestingEntitiesByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_FindNetworkByEntityID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	request := &g2pb.FindNetworkByEntityIDRequest{
		EntityList:     entityList,
		MaxDegree:      2,
		BuildOutDegree: 1,
		MaxEntities:    10,
	}
	response, err := g2engine.FindNetworkByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleG2EngineServer_FindNetworkByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	request := &g2pb.FindNetworkByEntityID_V2Request{
		EntityList:     entityList,
		MaxDegree:      2,
		BuildOutDegree: 1,
		MaxEntities:    10,
		Flags:          0,
	}
	response, err := g2engine.FindNetworkByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindNetworkByRecordID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindNetworkByRecordIDRequest{
		RecordList:     `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`,
		MaxDegree:      1,
		BuildOutDegree: 2,
		MaxEntities:    10,
	}
	response, err := g2engine.FindNetworkByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleG2EngineServer_FindNetworkByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindNetworkByRecordID_V2Request{
		RecordList:     `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`,
		MaxDegree:      1,
		BuildOutDegree: 2,
		MaxEntities:    10,
		Flags:          0,
	}
	response, err := g2engine.FindNetworkByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathByEntityID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindPathByEntityIDRequest{
		EntityID1: getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2: getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree: 1,
	}
	response, err := g2engine.FindPathByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindPathByEntityID_V2Request{
		EntityID1: getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2: getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree: 1,
		Flags:     0,
	}
	response, err := g2engine.FindPathByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathByRecordID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindPathByRecordIDRequest{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
	}
	response, err := g2engine.FindPathByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 87))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":...
}

func ExampleG2EngineServer_FindPathByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindPathByRecordID_V2Request{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		Flags:           0,
	}
	response, err := g2engine.FindPathByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathExcludingByEntityID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &g2pb.FindPathExcludingByEntityIDRequest{
		EntityID1:        getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2:        getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree:        1,
		ExcludedEntities: excludedEntities,
	}
	response, err := g2engine.FindPathExcludingByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathExcludingByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &g2pb.FindPathExcludingByEntityID_V2Request{
		EntityID1:        getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2:        getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree:        1,
		ExcludedEntities: excludedEntities,
		Flags:            0,
	}
	response, err := g2engine.FindPathExcludingByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathExcludingByRecordID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindPathExcludingByRecordIDRequest{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		ExcludedRecords: `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`,
	}
	response, err := g2engine.FindPathExcludingByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathExcludingByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindPathExcludingByRecordID_V2Request{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		ExcludedRecords: `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`,
		Flags:           0,
	}
	response, err := g2engine.FindPathExcludingByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathIncludingSourceByEntityID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &g2pb.FindPathIncludingSourceByEntityIDRequest{
		EntityID1:        getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2:        getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree:        1,
		ExcludedEntities: excludedEntities,
		RequiredDsrcs:    `{"DATA_SOURCES": ["CUSTOMERS"]}`,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 106))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2EngineServer_FindPathIncludingSourceByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	request := &g2pb.FindPathIncludingSourceByEntityID_V2Request{
		EntityID1:        getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2:        getEntityIdForRecord("CUSTOMERS", "1002"),
		MaxDegree:        1,
		ExcludedEntities: excludedEntities,
		RequiredDsrcs:    `{"DATA_SOURCES": ["CUSTOMERS"]}`,
		Flags:            0,
	}
	response, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_FindPathIncludingSourceByRecordID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindPathIncludingSourceByRecordIDRequest{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		ExcludedRecords: `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`,
		RequiredDsrcs:   `{"DATA_SOURCES": ["CUSTOMERS"]}`,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 119))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2EngineServer_FindPathIncludingSourceByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.FindPathIncludingSourceByRecordID_V2Request{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		MaxDegree:       1,
		ExcludedRecords: `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`,
		RequiredDsrcs:   `{"DATA_SOURCES": ["CUSTOMERS"]}`,
		Flags:           0,
	}
	response, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_GetActiveConfigID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetActiveConfigIDRequest{}
	response, err := g2engine.GetActiveConfigID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_GetEntityByEntityID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetEntityByEntityIDRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
	}
	response, err := g2engine.GetEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_GetEntityByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetEntityByEntityID_V2Request{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.GetEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2EngineServer_GetEntityByRecordID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetEntityByRecordIDRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
	}
	response, err := g2engine.GetEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 35))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2EngineServer_GetEntityByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetEntityByRecordID_V2Request{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.GetEntityByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2EngineServer_GetRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
	}
	response, err := g2engine.GetRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","JSON_DATA":{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","RECORD_TYPE":"PERSON","PRIMARY_NAME_LAST":"Smith","PRIMARY_NAME_FIRST":"Robert","DATE_OF_BIRTH":"12/11/1978","ADDR_TYPE":"MAILING","ADDR_LINE1":"123 Main Street, Las Vegas NV 89132","PHONE_TYPE":"HOME","PHONE_NUMBER":"702-919-1300","EMAIL_ADDRESS":"bsmith@work.com","DATE":"1/2/18","STATUS":"Active","AMOUNT":"100"}}
}

func ExampleG2EngineServer_GetRecord_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetRecord_V2Request{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.GetRecord_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}
}

func ExampleG2EngineServer_GetRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetRedoRecordRequest{}
	response, err := g2engine.GetRedoRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","ENTITY_TYPE":"GENERIC","DSRC_ACTION":"X"}
}

func ExampleG2EngineServer_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetRepositoryLastModifiedTimeRequest{}
	response, err := g2engine.GetRepositoryLastModifiedTime(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult() > 0) // Dummy output.
	// Output: true
}

func ExampleG2EngineServer_GetVirtualEntityByRecordID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetVirtualEntityByRecordIDRequest{
		RecordList: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`,
	}
	response, err := g2engine.GetVirtualEntityByRecordID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2EngineServer_GetVirtualEntityByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.GetVirtualEntityByRecordID_V2Request{
		RecordList: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`,
		Flags:      0,
	}
	response, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2EngineServer_HowEntityByEntityID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.HowEntityByEntityIDRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
	}
	response, err := g2engine.HowEntityByEntityID(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL","FEATURE_SCORES":{"ADDRESS":[{"INBOUND_FEAT_ID":20,"INBOUND_FEAT":"1515 Adela Lane Las Vegas NV 89111","INBOUND_FEAT_USAGE_TYPE":"HOME","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":42,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":19,"INBOUND_FEAT":"11/12/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":95,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"FMES"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":1,"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":97,"GNR_SN":100,"GNR_GN":95,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"702-919-1300","INBOUND_FEAT_USAGE_TYPE":"MOBILE","CANDIDATE_FEAT_ID":4,"CANDIDATE_FEAT":"702-919-1300","CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1-S1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V100001","MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","RESULT_VIRTUAL_ENTITY_ID":"V1-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"DOB":[{"INBOUND_FEAT_ID":2,"INBOUND_FEAT":"12/11/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FMES"}],"EMAIL":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":93,"GNR_SN":100,"GNR_GN":93,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"},{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"Robert Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":90,"GNR_SN":100,"GNR_GN":88,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}}],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1-S2","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleG2EngineServer_HowEntityByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.HowEntityByEntityID_V2Request{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.HowEntityByEntityID_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL"}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V1-S1","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V100001","MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","RESULT_VIRTUAL_ENTITY_ID":"V1-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB"}}],"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V1-S2","MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleG2EngineServer_PrimeEngine() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.PrimeEngineRequest{}
	response, err := g2engine.PrimeEngine(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_SearchByAttributes() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.SearchByAttributesRequest{
		JsonData: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`,
	}
	response, err := g2engine.SearchByAttributes(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 1962))
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1","FEATURE_SCORES":{"EMAIL":[{"INBOUND_FEAT":"bsmith@work.com","CANDIDATE_FEAT":"bsmith@work.com","FULL_SCORE":100}],"NAME":[{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Bob J Smith","GNR_FN":83,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1},{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Robert Smith","GNR_FN":88,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1}]}},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","FEATURES":{"ADDRESS":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20}]},{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3,"USAGE_TYPE":"MAILING","FEAT_DESC_VALUES":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3}]}],"DOB":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2},{"FEAT_DESC":"11/12/1978","LIB_FEAT_ID":19}]}],"EMAIL":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5}]}],"NAME":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1,"USAGE_TYPE":"PRIMARY","FEAT_DESC_VALUES":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1},{"FEAT_DESC":"Bob J Smith","LIB_FEAT_ID":32},{"FEAT_DESC":"Bob Smith","LIB_FEAT_ID":18}]}],"PHONE":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]},{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"MOBILE","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]}],"RECORD_TYPE":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleG2EngineServer_SearchByAttributes_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.SearchByAttributes_V2Request{
		JsonData: `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`,
		Flags:    0,
	}
	response, err := g2engine.SearchByAttributes_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1"},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1}}}]}
}

// func ExampleG2engineImpl_SetLogLevel() {
// }

func ExampleG2EngineServer_Stats() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.StatsRequest{}
	response, err := g2engine.Stats(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 16))
	// Output: { "workload":...
}

func ExampleG2EngineServer_WhyEntities() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.WhyEntitiesRequest{
		EntityID1: getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2: getEntityIdForRecord("CUSTOMERS", "1002"),
	}
	response, err := g2engine.WhyEntities(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 74))
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":1,"MATCH_INFO":{"WHY_KEY":...
}

func ExampleG2EngineServer_WhyEntities_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.WhyEntities_V2Request{
		EntityID1: getEntityIdForRecord("CUSTOMERS", "1001"),
		EntityID2: getEntityIdForRecord("CUSTOMERS", "1002"),
		Flags:     0,
	}
	response, err := g2engine.WhyEntities_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":1,"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+ADDRESS+PHONE+EMAIL","WHY_ERRULE_CODE":"SF1_SNAME_CFF_CSTAB","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_WhyRecords() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.WhyRecordsRequest{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
	}
	response, err := g2engine.WhyRecords(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(response.GetResult(), 115))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],...
}

func ExampleG2EngineServer_WhyRecords_V2() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.WhyRecords_V2Request{
		DataSourceCode1: "CUSTOMERS",
		RecordID1:       "1001",
		DataSourceCode2: "CUSTOMERS",
		RecordID2:       "1002",
		Flags:           0,
	}
	response, err := g2engine.WhyRecords_V2(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":1,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}],"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+PHONE","WHY_ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2EngineServer_ReevaluateEntity() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ReevaluateEntityRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.ReevaluateEntity(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ReevaluateEntityWithInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ReevaluateEntityWithInfoRequest{
		EntityID: getEntityIdForRecord("CUSTOMERS", "1001"),
		Flags:    0,
	}
	response, err := g2engine.ReevaluateEntityWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ReevaluateRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ReevaluateRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.ReevaluateRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ReevaluateRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ReevaluateRecordWithInfoRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		Flags:          0,
	}
	response, err := g2engine.ReevaluateRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_ReplaceRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ReplaceRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.ReplaceRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_ReplaceRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.ReplaceRecordWithInfoRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1001",
		JsonData:       `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`,
		LoadID:         "G2Engine_test",
		Flags:          0,
	}
	response, err := g2engine.ReplaceRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_DeleteRecord() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.DeleteRecordRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1003",
		LoadID:         "G2Engine_test",
	}
	response, err := g2engine.DeleteRecord(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_DeleteRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.DeleteRecordWithInfoRequest{
		DataSourceCode: "CUSTOMERS",
		RecordID:       "1003",
		LoadID:         "G2Engine_test",
		Flags:          0,
	}
	response, err := g2engine.DeleteRecordWithInfo(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.GetResult())
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2EngineServer_Init() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)

	}
	request := &g2pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: 0,
	}
	response, err := g2engine.Init(ctx, request)
	if err != nil {
		// This should produce a "senzing-60144002" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_InitWithConfigID() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Println(err)

	}
	request := &g2pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   1,
		VerboseLogging: 0,
	}
	response, err := g2engine.InitWithConfigID(ctx, request)
	if err != nil {
		// This should produce a "senzing-60144003" error.
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_Reinit() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)

	// Get a Senzing configuration ID for testing.
	requestToGetActiveConfigID := &g2pb.GetActiveConfigIDRequest{}
	responseFromGetActiveConfigID, err := g2engine.GetActiveConfigID(ctx, requestToGetActiveConfigID)
	if err != nil {
		fmt.Println(err)
	}
	// Example
	request := &g2pb.ReinitRequest{
		InitConfigID: responseFromGetActiveConfigID.GetResult(),
	}
	response, err := g2engine.Reinit(ctx, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output:
}

func ExampleG2EngineServer_Destroy() {
	// For more information, visit https://github.com/senzing-garage/serve-grpc/blob/main/g2engineserver/g2engineserver_test.go
	ctx := context.TODO()
	g2engine := getG2EngineServer(ctx)
	request := &g2pb.DestroyRequest{}
	response, err := g2engine.Destroy(ctx, request)
	if err != nil {
		// This should produce a "senzing-60164001" error.
	}
	fmt.Println(response)
	// Output:
}

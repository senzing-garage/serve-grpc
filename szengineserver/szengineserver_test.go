package szengineserver_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-grpc/szengineserver"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	attributes        = `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	defaultTruncation = 76
	observerID        = "Observer 1"
	observerOrigin    = "Observer 1 origin"
	printResults      = false
)

type GetEntityByRecordIDResponse struct {
	ResolvedEntity struct {
		EntityID int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

var (
	logLevelName      = "INFO"
	observerSingleton = &observer.NullObserver{
		ID:       observerID,
		IsSilent: true,
	}
	szEngineTestSingleton *szengineserver.SzEngineServer
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzEngineServer_AddRecord(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	request1 := &szpb.AddRecordRequest{
		DataSourceCode:   record1.DataSource,
		Flags:            senzing.SzWithInfo,
		RecordDefinition: record1.JSON,
		RecordId:         record1.ID,
	}
	response1, err := szEngineServer.AddRecord(ctx, request1)
	require.NoError(test, err)
	require.NotEmpty(test, response1.GetResult())
	printActual(test, response1.GetResult())

	request2 := &szpb.AddRecordRequest{
		DataSourceCode:   record2.DataSource,
		Flags:            senzing.SzWithInfo,
		RecordDefinition: record2.JSON,
		RecordId:         record2.ID,
	}
	response2, err := szEngineServer.AddRecord(ctx, request2)
	require.NoError(test, err)
	require.NotEmpty(test, response1.GetResult())
	printActual(test, response2.GetResult())
}

func TestSzEngineServer_AddRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	request := &szpb.AddRecordRequest{
		DataSourceCode:   record.DataSource,
		Flags:            senzing.SzWithInfo,
		RecordDefinition: record.JSON,
		RecordId:         record.ID,
	}
	response, err := szEngineServer.AddRecord(ctx, request)
	require.NoError(test, err)
	require.NotEmpty(test, response.GetResult())
	printActual(test, response.GetResult())
}

func TestSzEngineServer_CountRedoRecords(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.CountRedoRecordsRequest{}
	response, err := szEngineServer.CountRedoRecords(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_ExportJsonEntityReport(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	flags := senzing.SzNoFlags
	request := &szpb.ExportJsonEntityReportRequest{
		Flags: flags,
	}
	response, err := szEngineServer.ExportJsonEntityReport(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_ExportCsvEntityReport(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.ExportCsvEntityReportRequest{
		CsvColumnList: "",
		Flags:         senzing.SzNoFlags,
	}
	response, err := szEngineServer.ExportCsvEntityReport(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_FindInterestingEntitiesByEntityId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	flags := int64(0)
	request := &szpb.FindInterestingEntitiesByEntityIdRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	response, err := szEngineServer.FindInterestingEntitiesByEntityId(ctx, request)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzEngineServer_FindInterestingEntitiesByRecordId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := int64(0)
	request := &szpb.FindInterestingEntitiesByRecordIdRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	response, err := szEngineServer.FindInterestingEntitiesByRecordId(ctx, request)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzEngineServer_FindNetworkByEntityId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDString(
		record1,
	) + `}, {"ENTITY_ID": ` + getEntityIDString(
		record2,
	) + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	buildOutMaxEntities := int64(10)
	flags := senzing.SzNoFlags
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     buildOutDegree,
		BuildOutMaxEntities: buildOutMaxEntities,
		EntityIds:           entityIDs,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
	}
	response, err := szEngineServer.FindNetworkByEntityId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_FindNetworkByEntityId_Not_Using_CUSTOMERS(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getSzEngineServer(ctx)
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord(
		"REFERENCE",
		"2012",
	) + `}, {"ENTITY_ID": ` + getEntityIDStringForRecord(
		"REFERENCE",
		"2013",
	) + `}]}`
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     1,
		BuildOutMaxEntities: 10,
		EntityIds:           entityIDs,
		Flags:               senzing.SzNoFlags,
		MaxDegrees:          2,
	}

	_, err := szEngineServer.FindNetworkByEntityId(ctx, request)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
}

func TestSzEngineServer_FindNetworkByRecordId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.ID + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.ID + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.ID + `"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	buildOutMaxEntities := int64(10)
	flags := senzing.SzNoFlags
	request := &szpb.FindNetworkByRecordIdRequest{
		BuildOutDegrees:     buildOutDegree,
		BuildOutMaxEntities: buildOutMaxEntities,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		RecordKeys:          recordKeys,
	}
	response, err := szEngineServer.FindNetworkByRecordId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_FindPathByEntityId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	startEntityID := getEntityID(truthset.CustomerRecords["1001"])
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	flags := senzing.SzNoFlags
	request := &szpb.FindPathByEntityIdRequest{
		EndEntityId:   endEntityID,
		Flags:         flags,
		MaxDegrees:    maxDegrees,
		StartEntityId: startEntityID,
	}
	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_FindPathByEntityId_exclusions(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	startEntityID := getEntityID(record1)
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDString(record1) + `}]}`
	flags := senzing.SzNoFlags
	request := &szpb.FindPathByEntityIdRequest{
		AvoidEntityIds: avoidEntityIDs,
		EndEntityId:    endEntityID,
		Flags:          flags,
		MaxDegrees:     maxDegrees,
		StartEntityId:  startEntityID,
	}
	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_FindPathByEntityId_inclusions(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	startEntityID := getEntityID(record1)
	endEntityID := getEntityID(truthset.CustomerRecords["1002"])
	maxDegrees := int64(1)
	avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDString(record1) + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	request := &szpb.FindPathByEntityIdRequest{
		AvoidEntityIds:      avoidEntityIDs,
		EndEntityId:         endEntityID,
		MaxDegrees:          maxDegrees,
		RequiredDataSources: requiredDataSources,
		StartEntityId:       startEntityID,
	}
	response, err := szEngineServer.FindPathByEntityId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_FindPathByRecordId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegrees := int64(1)
	flags := senzing.SzNoFlags
	request := &szpb.FindPathByRecordIdRequest{
		EndDataSourceCode:   record2.DataSource,
		EndRecordId:         record2.ID,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		StartDataSourceCode: record1.DataSource,
		StartRecordId:       record1.ID,
	}
	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_FindPathByRecordId_exclusions(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegrees := int64(1)
	avoidRecordKeys := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.ID + `"}]}`
	flags := senzing.SzNoFlags
	request := &szpb.FindPathByRecordIdRequest{
		AvoidRecordKeys:     avoidRecordKeys,
		EndDataSourceCode:   record2.DataSource,
		EndRecordId:         record2.ID,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		StartDataSourceCode: record1.DataSource,
		StartRecordId:       record1.ID,
	}
	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_FindPathByRecordId_inclusions(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegrees := int64(1)
	avoidRecordKeys := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDString(record1) + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := senzing.SzNoFlags
	request := &szpb.FindPathByRecordIdRequest{
		AvoidRecordKeys:     avoidRecordKeys,
		EndDataSourceCode:   record2.DataSource,
		EndRecordId:         record1.ID,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		RequiredDataSources: requiredDataSources,
		StartDataSourceCode: record1.DataSource,
		StartRecordId:       record1.ID,
	}
	response, err := szEngineServer.FindPathByRecordId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_GetActiveConfigId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.GetActiveConfigIdRequest{}
	response, err := szEngineServer.GetActiveConfigId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_GetEntityByEntityId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	flags := senzing.SzNoFlags
	request := &szpb.GetEntityByEntityIdRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	response, err := szEngineServer.GetEntityByEntityId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_GetEntityByRecordId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.GetEntityByRecordIdRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	response, err := szEngineServer.GetEntityByRecordId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_GetRecord(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.GetRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	response, err := szEngineServer.GetRecord(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_GetRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.GetRedoRecordRequest{}
	response, err := szEngineServer.GetRedoRecord(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_GetVirtualEntityByRecordId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.ID + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.ID + `"}]}`
	flags := senzing.SzNoFlags
	request := &szpb.GetVirtualEntityByRecordIdRequest{
		Flags:      flags,
		RecordKeys: recordKeys,
	}
	response, err := szEngineServer.GetVirtualEntityByRecordId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_HowEntityByEntityId(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	flags := senzing.SzNoFlags
	request := &szpb.HowEntityByEntityIdRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	response, err := szEngineServer.HowEntityByEntityId(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_PreprocessRecord(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	request := &szpb.PreprocessRecordRequest{
		Flags:            senzing.SzNoFlags,
		RecordDefinition: record.JSON,
	}
	response, err := szEngineServer.PreprocessRecord(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_PrimeEngine(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.PrimeEngineRequest{}
	response, err := szEngineServer.PrimeEngine(ctx, request)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzEngineServer_ProcessRedoRecord(test *testing.T) {
	// IMPROVE: Implement TestSzEngineServer_ProcessRedoRecord
	_ = test
}

func TestSzEngineServer_ReevaluateEntity(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	flags := senzing.SzWithoutInfo
	request := &szpb.ReevaluateEntityRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	response, err := szEngineServer.ReevaluateEntity(ctx, request)
	require.NoError(test, err)
	require.Empty(test, response.GetResult())
	printActual(test, response.GetResult())
}

func TestSzEngineServer_ReevaluateEntity_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	flags := senzing.SzWithInfo
	request := &szpb.ReevaluateEntityRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	response, err := szEngineServer.ReevaluateEntity(ctx, request)
	require.NoError(test, err)
	require.NotEmpty(test, response.GetResult())
	printActual(test, response.GetResult())
}

func TestSzEngineServer_ReevaluateRecord(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	request := &szpb.ReevaluateRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	response, err := szEngineServer.ReevaluateRecord(ctx, request)
	require.NoError(test, err)
	require.Empty(test, response.GetResult())
	printActual(test, response.GetResult())
}

func TestSzEngineServer_ReevaluateRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithInfo
	request := &szpb.ReevaluateRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	response, err := szEngineServer.ReevaluateRecord(ctx, request)
	require.NoError(test, err)
	require.NotEmpty(test, response.GetResult())
	printActual(test, response.GetResult())
}

func TestSzEngineServer_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)

	requestToGetActiveConfigID := &szpb.GetActiveConfigIdRequest{}
	responseFromGetActiveConfigID, err := szEngineServer.GetActiveConfigId(ctx, requestToGetActiveConfigID)
	require.NoError(test, err)

	request := &szpb.ReinitializeRequest{
		ConfigId: responseFromGetActiveConfigID.GetResult(),
	}
	response, err := szEngineServer.Reinitialize(ctx, request)
	require.NoError(test, err)
	printActual(test, response)
}

func TestSzEngineServer_SearchByAttributes(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	flags := senzing.SzNoFlags
	request := &szpb.SearchByAttributesRequest{
		Attributes: attributes,
		Flags:      flags,
	}
	response, err := szEngineServer.SearchByAttributes(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_SearchByAttributes_searchProfile(test *testing.T) {
	// IMPROVE:  Use actual searchProfile
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	searchProfile := "SEARCH"
	flags := senzing.SzNoFlags
	request := &szpb.SearchByAttributesRequest{
		Attributes:    attributes,
		Flags:         flags,
		SearchProfile: searchProfile,
	}
	response, err := szEngineServer.SearchByAttributes(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_Stats(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	request := &szpb.GetStatsRequest{}
	response, err := szEngineServer.GetStats(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_StreamExportCsvEntityReport(test *testing.T) {
	// IMPROVE: Implement TestSzEngineServer_StreamExportCsvEntityReport
	_ = test
}

func TestSzEngineServer_StreamExportJsonEntityReport(test *testing.T) {
	// IMPROVE: Implement TestSzEngineServer_StreamExportJsonEntityReport
	_ = test
}

func TestSzEngineServer_WhyEntities(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	entityID1 := getEntityID(truthset.CustomerRecords["1001"])
	entityID2 := getEntityID(truthset.CustomerRecords["1002"])
	flags := senzing.SzNoFlags
	request := &szpb.WhyEntitiesRequest{
		EntityId_1: entityID1,
		EntityId_2: entityID2,
		Flags:      flags,
	}
	response, err := szEngineServer.WhyEntities(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_WhyRecordInEntity(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.WhyRecordInEntityRequest{
		DataSourceCode: record1.DataSource,
		Flags:          flags,
		RecordId:       record1.ID,
	}
	response, err := szEngineServer.WhyRecordInEntity(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_WhyRecords(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	request := &szpb.WhyRecordsRequest{
		DataSourceCode_1: record1.DataSource,
		DataSourceCode_2: record2.DataSource,
		Flags:            flags,
		RecordId_1:       record1.ID,
		RecordId_2:       record2.ID,
	}
	response, err := szEngineServer.WhyRecords(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_WhySearch(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	searchProfile := "SEARCH"

	flags := senzing.SzNoFlags
	request := &szpb.WhySearchRequest{
		Attributes:    attributes,
		EntityId:      entityID,
		SearchProfile: searchProfile,
		Flags:         flags,
	}
	response, err := szEngineServer.WhySearch(ctx, request)
	require.NoError(test, err)
	printActual(test, response.GetResult())
}

func TestSzEngineServer_DeleteRecord(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	flags := senzing.SzWithoutInfo
	request := &szpb.DeleteRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	response, err := szEngineServer.DeleteRecord(ctx, request)
	require.NoError(test, err)
	require.Empty(test, response.GetResult())
	printActual(test, response.GetResult())
}

func TestSzEngineServer_DeleteRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngineServer := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1003"]
	flags := senzing.SzWithInfo
	request := &szpb.DeleteRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	response, err := szEngineServer.DeleteRecord(ctx, request)
	require.NoError(test, err)
	require.NotEmpty(test, response.GetResult())
	printActual(test, response.GetResult())
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzEngineServer_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

func TestSzEngineServer_SetLogLevel(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	require.NoError(test, err)
}

func TestSzEngineServer__SetLogLevel_badLevelName(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	require.Error(test, err)
}

func TestSzEngineServer_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	testObject.SetObserverOrigin(ctx, observerOrigin)
}

func TestSzEngineServer_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	actual := testObject.GetObserverOrigin(ctx)
	assert.Equal(test, observerOrigin, actual)
}

func TestSzEngineServer_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestBuildSimpleSystemConfigurationJsonUsingEnvVars(test *testing.T) {
	actual, err := settings.BuildSimpleSettingsUsingEnvVars()
	panicOnError(err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getEntityID(record record.Record) int64 {
	return getEntityIDForRecord(record.DataSource, record.ID)
}

func getEntityIDForRecord(datasource string, recordID string) int64 {
	ctx := context.TODO()

	var result int64

	szEngine := getSzEngineServer(ctx)
	request := &szpb.GetEntityByRecordIdRequest{
		DataSourceCode: datasource,
		RecordId:       recordID,
	}

	response, err := szEngine.GetEntityByRecordId(ctx, request)
	if err != nil {
		return result
	}

	getEntityByRecordIDResponse := &GetEntityByRecordIDResponse{}

	err = json.Unmarshal([]byte(response.GetResult()), &getEntityByRecordIDResponse)
	if err != nil {
		return result
	}

	return getEntityByRecordIDResponse.ResolvedEntity.EntityID
}

func getEntityIDString(record record.Record) string {
	entityID := getEntityID(record)

	return strconv.FormatInt(entityID, 10)
}

func getEntityIDStringForRecord(datasource string, id string) string {
	entityID := getEntityIDForRecord(datasource, id)

	return strconv.FormatInt(entityID, 10)
}

func getSzEngineServer(ctx context.Context) *szengineserver.SzEngineServer {
	if szEngineTestSingleton == nil {
		szEngineTestSingleton = &szengineserver.SzEngineServer{}
		instanceName := "Test instance name"
		verboseLogging := senzing.SzNoLogging
		configID := senzing.SzInitializeWithDefaultConfiguration
		setting, err := settings.BuildSimpleSettingsUsingEnvVars()
		panicOnError(err)

		osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
		if len(osenvLogLevel) > 0 {
			logLevelName = osenvLogLevel
		}

		err = szEngineTestSingleton.SetLogLevel(ctx, logLevelName)
		panicOnError(err)
		err = szengineserver.GetSdkSzEngine().Initialize(ctx, instanceName, setting, configID, verboseLogging)
		panicOnError(err)
	}

	return szEngineTestSingleton
}

func getTestObject(ctx context.Context, t *testing.T) *szengineserver.SzEngineServer {
	t.Helper()

	return getSzEngineServer(ctx)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func printActual(t *testing.T, actual interface{}) {
	t.Helper()
	printResult(t, "Actual", actual)
}

func printResult(t *testing.T, title string, result interface{}) {
	t.Helper()

	if printResults {
		t.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	setup()

	code := m.Run()
	os.Exit(code)
}

func setupSenzingConfig(ctx context.Context, instanceName string, settings string, verboseLogging int64) {
	now := time.Now()

	szAbstractFactory := &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	panicOnError(err)

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	panicOnError(err)

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range datasourceNames {
		_, err := szConfig.AddDataSource(ctx, dataSourceCode)
		panicOnError(err)
	}

	configComment := fmt.Sprintf("Created by szconfigmanagerserver_test at %s", now.UTC())
	configDefinition, err := szConfig.Export(ctx)
	panicOnError(err)

	configID, err := szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	panicOnError(err)

	err = szAbstractFactory.Reinitialize(ctx, configID)
	panicOnError(err)
}

func setupPurgeRepository(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	err := szDiagnostic.Initialize(
		ctx,
		instanceName,
		settings,
		senzing.SzInitializeWithDefaultConfiguration,
		verboseLogging,
	)
	panicOnError(err)

	err = szDiagnostic.PurgeRepository(ctx)
	panicOnError(err)

	err = szDiagnostic.Destroy(ctx)
	panicOnError(err)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

func setup() {
	ctx := context.TODO()
	instanceName := "Test name"
	verboseLogging := senzing.SzNoLogging
	settings, err := settings.BuildSimpleSettingsUsingEnvVars()
	panicOnError(err)
	setupSenzingConfig(ctx, instanceName, settings, verboseLogging)
	err = setupPurgeRepository(ctx, instanceName, settings, verboseLogging)
	panicOnError(err)
}

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
	attributes          = `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	avoidEntityIDs      = senzing.SzNoAvoidance
	avoidRecordKeys     = senzing.SzNoAvoidance
	baseTen             = 10
	buildOutDegrees     = int64(2)
	buildOutMaxEntities = int64(10)
	defaultTruncation   = 76
	instanceName        = "SzEngine Test"
	jsonIndentation     = "    "
	maxDegrees          = int64(2)
	observerID          = "Observer 1"
	observerOrigin      = "SzEngine observer"
	originMessage       = "Machine: nn; Task: UnitTest"
	printErrors         = false
	printResults        = false
	searchAttributes    = `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile       = senzing.SzNoSearchProfile
	verboseLogging      = senzing.SzNoLogging
)

// Bad parameters

const (
	badAttributes          = "}{"
	badAvoidRecordKeys     = "}{"
	badBuildOutDegrees     = int64(-1)
	badBuildOutMaxEntities = int64(-1)
	badCsvColumnList       = "BAD, CSV, COLUMN, LIST"
	badDataSourceCode      = "BadDataSourceCode"
	badEntityID            = int64(-1)
	badExportHandle        = uintptr(0)
	badLogLevelName        = "BadLogLevelName"
	badMaxDegrees          = int64(-1)
	badRecordDefinition    = "}{"
	badRecordID            = "BadRecordID"
	badRedoRecord          = "{}"
	badSearchProfile       = "}{"
)

// Nil/empty parameters

var (
	nilAttributes          string
	nilBuildOutDegrees     int64
	nilBuildOutMaxEntities int64
	nilCsvColumnList       string
	nilDataSourceCode      string
	nilEntityID            int64
	nilExportHandle        uintptr
	nilMaxDegrees          int64
	nilRecordDefinition    string
	nilRecordID            string
	nilRedoRecord          string
	nilRequiredDataSources string
	nilSearchProfile       string
)

var (
	logLevelName      = "INFO"
	observerSingleton = &observer.NullObserver{
		ID:       observerID,
		IsSilent: true,
	}
	szEngineTestSingleton *szengineserver.SzEngineServer
)

type GetEntityByRecordIDResponse struct {
	ResolvedEntity struct {
		EntityID int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func TestSzEngine_AddRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	for _, record := range records {
		request := &szpb.AddRecordRequest{
			DataSourceCode:   record.DataSource,
			Flags:            flags,
			RecordDefinition: record.JSON,
			RecordId:         record.ID,
		}
		actual, err := szEngine.AddRecord(ctx, request)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.Empty(test, actual.GetResult())
	}

	for _, record := range records {
		request := &szpb.DeleteRecordRequest{
			DataSourceCode: record.DataSource,
			Flags:          flags,
			RecordId:       record.ID,
		}
		actual, err := szEngine.DeleteRecord(ctx, request)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.Empty(test, actual.GetResult())
	}
}

func TestSzEngine_AddRecord_badDataSourceCodeInJSON(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithoutInfo
	record := truthset.CustomerRecords["1002"]
	badRecordJSON := `{"DATA_SOURCE": "BOB", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	request := &szpb.AddRecordRequest{
		DataSourceCode:   record.DataSource,
		Flags:            flags,
		RecordDefinition: badRecordJSON,
		RecordId:         record.ID,
	}
	actual, err := szEngine.AddRecord(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function": "szengineserver.(*SzEngineServer).AddRecord", "error": {"function": "szengine.(*Szengine).AddRecord", "error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'CUSTOMERS' and 'BOB'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo

	request := &szpb.AddRecordRequest{
		DataSourceCode:   badDataSourceCode,
		Flags:            flags,
		RecordDefinition: record.JSON,
		RecordId:         record.ID,
	}
	actual, err := szEngine.AddRecord(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_badRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	request := &szpb.AddRecordRequest{
		DataSourceCode:   record.DataSource,
		Flags:            flags,
		RecordDefinition: record.JSON,
		RecordId:         badRecordID,
	}
	actual, err := szEngine.AddRecord(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function": "szengineserver.(*SzEngineServer).AddRecord", "error": {"function": "szengine.(*Szengine).AddRecord", "error": {"id":"SZSDK60044001","reason":"SENZ0024|Conflicting RECORD_ID values 'BadRecordID' and '1001'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_badRecordDefinition(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	request := &szpb.AddRecordRequest{
		DataSourceCode:   record.DataSource,
		Flags:            flags,
		RecordDefinition: badRecordDefinition,
		RecordId:         record.ID,
	}
	actual, err := szEngine.AddRecord(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function": "szengineserver.(*SzEngineServer).AddRecord", "error": {"function": "szengine.(*Szengine).AddRecord", "error": {"id":"SZSDK60044001","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	request := &szpb.AddRecordRequest{
		DataSourceCode:   nilDataSourceCode,
		Flags:            flags,
		RecordDefinition: record.JSON,
		RecordId:         record.ID,
	}
	actual, err := szEngine.AddRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual.GetResult())
}

func TestSzEngine_AddRecord_nilRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	request := &szpb.AddRecordRequest{
		DataSourceCode:   record.DataSource,
		Flags:            flags,
		RecordDefinition: record.JSON,
		RecordId:         nilRecordID,
	}
	actual, err := szEngine.AddRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual.GetResult())
}

func TestSzEngine_AddRecord_nilRecordDefinition(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	request := &szpb.AddRecordRequest{
		DataSourceCode:   record.DataSource,
		Flags:            flags,
		RecordDefinition: nilRecordDefinition,
		RecordId:         record.ID,
	}
	actual, err := szEngine.AddRecord(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	expectedErr := `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ3121|JSON Parsing Failure [code=1,offset=0]"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_AddRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	records := []record.Record{
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
	}

	for _, record := range records {
		request := &szpb.AddRecordRequest{
			DataSourceCode:   record.DataSource,
			Flags:            flags,
			RecordDefinition: record.JSON,
			RecordId:         record.ID,
		}
		actual, err := szEngine.AddRecord(ctx, request)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.NotEmpty(test, actual.GetResult())
	}

	for _, record := range records {
		request := &szpb.DeleteRecordRequest{
			DataSourceCode: record.DataSource,
			Flags:          flags,
			RecordId:       record.ID,
		}
		actual, err := szEngine.DeleteRecord(ctx, request)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.NotEmpty(test, actual.GetResult())
	}
}

func TestSzEngine_AddRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzWithInfo
	records := []record.Record{
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
	}

	for _, record := range records {
		request := &szpb.AddRecordRequest{
			DataSourceCode:   badDataSourceCode,
			Flags:            flags,
			RecordDefinition: record.JSON,
			RecordId:         record.ID,
		}
		actual, err := szEngine.AddRecord(ctx, request)
		printDebug(test, err, actual)
		require.ErrorIs(test, err, szerror.ErrSzBadInput)

		expectedErr := `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}}`
		require.JSONEq(test, expectedErr, err.Error())
	}

	for _, record := range records {
		request := &szpb.DeleteRecordRequest{
			DataSourceCode: record.DataSource,
			Flags:          flags,
			RecordId:       record.ID,
		}
		actual, err := szEngine.DeleteRecord(ctx, request)
		printDebug(test, err, actual)
		require.NoError(test, err)
		require.NotEmpty(test, actual.GetResult())
	}
}

func TestSzEngine_CountRedoRecords(test *testing.T) {
	ctx := test.Context()
	expected := int64(2)
	szEngine := getTestObject(test)
	request := &szpb.CountRedoRecordsRequest{}
	actual, err := szEngine.CountRedoRecords(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Equal(test, expected, actual.GetResult())
}

func TestSzEngine_DeleteRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1005"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo

	request := &szpb.DeleteRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.DeleteRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual.GetResult())
}

func TestSzEngine_DeleteRecord_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1005"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo

	request := &szpb.DeleteRecordRequest{
		DataSourceCode: badDataSourceCode,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.DeleteRecord(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	expectedErr := `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_DeleteRecord_badRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1005"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo

	request := &szpb.DeleteRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       badRecordID,
	}
	actual, err := szEngine.DeleteRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual.GetResult())
}

func TestSzEngine_DeleteRecord_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1005"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo

	request := &szpb.DeleteRecordRequest{
		DataSourceCode: nilDataSourceCode,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.DeleteRecord(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_DeleteRecord_nilRecordID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	records := []record.Record{
		truthset.CustomerRecords["1005"],
	}
	addRecords(ctx, records)

	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo

	request := &szpb.DeleteRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       nilRecordID,
	}
	actual, err := szEngine.DeleteRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual.GetResult())
}

func TestSzEngine_DeleteRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1003"]
	flags := senzing.SzWithInfo
	request := &szpb.DeleteRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.DeleteRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.NotEmpty(test, actual.GetResult())
}

func TestSzEngine_DeleteRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1003"]
	flags := senzing.SzWithInfo
	request := &szpb.DeleteRecordRequest{
		DataSourceCode: badDataSourceCode,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.DeleteRecord(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	expectedErr := `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044005","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ExportCsvEntityReport(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	csvColumnList := ""
	flags := senzing.SzExportIncludeAllEntities

	request := &szpb.ExportCsvEntityReportRequest{
		CsvColumnList: csvColumnList,
		Flags:         flags,
	}
	actual, err := szEngine.ExportCsvEntityReport(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ExportCsvEntityReport_badCsvColumnList(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzExportIncludeAllEntities

	request := &szpb.ExportCsvEntityReportRequest{
		CsvColumnList: badCsvColumnList,
		Flags:         flags,
	}
	actual, err := szEngine.ExportCsvEntityReport(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	expectedErr := `{"function":"szengineserver.(*SzEngineServer).ExportCsvEntityReport","error":{"function":"szengine.(*Szengine).ExportCsvEntityReport","error":{"id":"SZSDK60044007","reason":"SENZ3131|Invalid column [BAD] requested for CSV export."}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_ExportCsvEntityReport_nilCsvColumnList(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzExportIncludeAllEntities

	request := &szpb.ExportCsvEntityReportRequest{
		CsvColumnList: nilCsvColumnList,
		Flags:         flags,
	}
	actual, err := szEngine.ExportCsvEntityReport(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ExportJSONEntityReport(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	request := &szpb.ExportJsonEntityReportRequest{
		Flags: flags,
	}
	actual, err := szEngine.ExportJsonEntityReport(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ExportJSONEntityReport_65536(test *testing.T) {
	// IMPROVE:
}

func TestSzEngine_FindInterestingEntitiesByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])

	flags := senzing.SzNoFlags
	request := &szpb.FindInterestingEntitiesByEntityIdRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	actual, err := szEngine.FindInterestingEntitiesByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindInterestingEntitiesByEntityID_badEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)

	flags := senzing.SzNoFlags
	request := &szpb.FindInterestingEntitiesByEntityIdRequest{
		EntityId: badEntityID,
		Flags:    flags,
	}
	actual, err := szEngine.FindInterestingEntitiesByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByEntityId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByEntityID","error":{"id":"SZSDK60044010","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByEntityID_nilEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)

	flags := senzing.SzNoFlags
	request := &szpb.FindInterestingEntitiesByEntityIdRequest{
		EntityId: nilEntityID,
		Flags:    flags,
	}
	actual, err := szEngine.FindInterestingEntitiesByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByEntityId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByEntityID","error":{"id":"SZSDK60044010","reason":"SENZ0037|Unknown resolved entity value '0'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.FindInterestingEntitiesByRecordIdRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.FindInterestingEntitiesByRecordId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindInterestingEntitiesByRecordID_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.FindInterestingEntitiesByRecordIdRequest{
		DataSourceCode: badDataSourceCode,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.FindInterestingEntitiesByRecordId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByRecordId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByRecordID_badRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.FindInterestingEntitiesByRecordIdRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       badRecordID,
	}
	actual, err := szEngine.FindInterestingEntitiesByRecordId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByRecordId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByRecordID_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.FindInterestingEntitiesByRecordIdRequest{
		DataSourceCode: nilDataSourceCode,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.FindInterestingEntitiesByRecordId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByRecordId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ2207|Data source code [] does not exist."}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindInterestingEntitiesByRecordID_nilRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.FindInterestingEntitiesByRecordIdRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       nilRecordID,
	}
	actual, err := szEngine.FindInterestingEntitiesByRecordId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByRecordId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     buildOutDegrees,
		BuildOutMaxEntities: buildOutMaxEntities,
		EntityIds:           entityIDs,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
	}
	actual, err := szEngine.FindNetworkByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByEntityID_badEntityIDs(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	badEntityID1 := 0
	badEntityID2 := 1
	badEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + strconv.Itoa(
		badEntityID1,
	) + `}, {"ENTITY_ID": ` + strconv.Itoa(
		badEntityID2,
	) + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     buildOutDegrees,
		BuildOutMaxEntities: buildOutMaxEntities,
		EntityIds:           badEntityIDs,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
	}
	actual, err := szEngine.FindNetworkByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)

	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindNetworkByEntityId","error":{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0037|Unknown resolved entity value '0'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID_badMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags

	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     buildOutDegrees,
		BuildOutMaxEntities: buildOutMaxEntities,
		EntityIds:           entityIDs,
		Flags:               flags,
		MaxDegrees:          badMaxDegrees,
	}
	actual, err := szEngine.FindNetworkByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindNetworkByEntityId","error":{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0031|Invalid value of max degree '-1'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID_badBuildOutDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     badBuildOutDegrees,
		BuildOutMaxEntities: buildOutMaxEntities,
		EntityIds:           entityIDs,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
	}
	actual, err := szEngine.FindNetworkByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindNetworkByEntityId","error":{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0032|Invalid value of build out degree '-1'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID_badBuildOutMaxEntities(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     buildOutDegrees,
		BuildOutMaxEntities: badBuildOutMaxEntities,
		EntityIds:           entityIDs,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
	}
	actual, err := szEngine.FindNetworkByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindNetworkByEntityId","error":{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0029|Invalid value of max entities '-1'"}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByEntityID_nilMaxDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags

	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     buildOutDegrees,
		BuildOutMaxEntities: buildOutMaxEntities,
		EntityIds:           entityIDs,
		Flags:               flags,
		MaxDegrees:          nilMaxDegrees,
	}
	actual, err := szEngine.FindNetworkByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByEntityID_nilBuildOutDegrees(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     nilBuildOutDegrees,
		BuildOutMaxEntities: buildOutMaxEntities,
		EntityIds:           entityIDs,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
	}
	actual, err := szEngine.FindNetworkByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByEntityID_nilBuildOutMaxEntities(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1 := getEntityIDString(record1)
	entityID2 := getEntityIDString(record2)

	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     buildOutDegrees,
		BuildOutMaxEntities: nilBuildOutMaxEntities,
		EntityIds:           entityIDs,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
	}
	actual, err := szEngine.FindNetworkByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

// func TestSzEngine_FindNetworkByEntityID_Not_Using_CUSTOMERS(test *testing.T) {
// 	ctx := test.Context()
// 	szEngine := getSzEngineServer(ctx)
// 	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` +
// 		getEntityIDStringForRecord("REFERENCE", "2012") +
// 		`}, {"ENTITY_ID": ` +
// 		getEntityIDStringForRecord("REFERENCE", "2013") +
// 		`}]}`
// 	request := &szpb.FindNetworkByEntityIdRequest{
// 		BuildOutDegrees:     1,
// 		BuildOutMaxEntities: 10,
// 		EntityIds:           entityIDs,
// 		Flags:               senzing.SzNoFlags,
// 		MaxDegrees:          2,
// 	}

// 	actual, err := szEngine.FindNetworkByEntityId(ctx, request)
// 	printDebug(test, err, actual)
// 	require.ErrorIs(test, err, szerror.ErrSzNotFound)
// }

func TestSzEngine_FindNetworkByRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}, {"DATA_SOURCE": "` +
		record3.DataSource +
		`", "RECORD_ID": "` +
		record3.ID +
		`"}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	request := &szpb.FindNetworkByRecordIdRequest{
		BuildOutDegrees:     buildOutDegrees,
		BuildOutMaxEntities: buildOutMaxEntities,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		RecordKeys:          recordKeys,
	}
	actual, err := szEngine.FindNetworkByRecordId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_FindNetworkByRecordID_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	badRecordKeys := `{"RECORDS": [{"DATA_SOURCE": "` +
		badDataSourceCode +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}, {"DATA_SOURCE": "` +
		record3.DataSource +
		`", "RECORD_ID": "` +
		record3.ID +
		`"}]}`
	flags := senzing.SzFindNetworkDefaultFlags
	request := &szpb.FindNetworkByRecordIdRequest{
		BuildOutDegrees:     buildOutDegrees,
		BuildOutMaxEntities: buildOutMaxEntities,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		RecordKeys:          badRecordKeys,
	}
	actual, err := szEngine.FindNetworkByRecordId(ctx, request)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)

	expectedErr := `{"function":"szengineserver.(*SzEngineServer).FindNetworkByRecordId","error":{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzEngine_FindNetworkByRecordID_badRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_FindNetworkByRecordID_nilMaxDegrees(test *testing.T) {
	// FIXME:
}

func TestSzEngine_FindNetworkByRecordID_nilBuildOutDegrees(test *testing.T) {
	// FIXME:
}

func TestSzEngine_FindNetworkByRecordID_nilBuildOutMaxEntities(test *testing.T) {
	// FIXME:
}

func TestSzEngine_FindPathByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindPathByEntityID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {

			// Insert test data.

			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(test)
			startEntityID := getEntityID(truthset.CustomerRecords["1001"])
			endEntityID := getEntityID(truthset.CustomerRecords["1002"])

			avoidEntityIDs := senzing.SzNoAvoidance
			if testCase.avoidEntityIDs != nil {
				avoidEntityIDs = testCase.avoidEntityIDs()
			}

			requiredDataSources := senzing.SzNoRequiredDatasources
			if testCase.requiredDataSources != nil {
				requiredDataSources = testCase.requiredDataSources()
			}

			// Test.

			request := &szpb.FindPathByEntityIdRequest{
				AvoidEntityIds:      avoidEntityIDs,
				EndEntityId:         xInt64(testCase.endEntityID, endEntityID),
				Flags:               testCase.flags,
				MaxDegrees:          testCase.maxDegrees,
				RequiredDataSources: requiredDataSources,
				StartEntityId:       xInt64(testCase.startEntityID, startEntityID),
			}
			actual, err := szEngine.FindPathByEntityId(ctx, request)
			printDebug(test, err, actual)
			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_FindPathByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindPathByRecordID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {

			// Insert test data.

			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(test)

			record1 := truthset.CustomerRecords["1001"]
			startDataSourceCode := record1.DataSource
			startRecordId := record1.ID

			record2 := truthset.CustomerRecords["1002"]
			endDataSourceCode := record2.DataSource
			endRecordId := record2.ID

			avoidRecordKeys := senzing.SzNoAvoidance
			if testCase.avoidRecordKeys != nil {
				avoidRecordKeys = testCase.avoidRecordKeys()
			}

			requiredDataSources := senzing.SzNoRequiredDatasources
			if testCase.requiredDataSources != nil {
				requiredDataSources = testCase.requiredDataSources()
			}

			// Test.

			request := &szpb.FindPathByRecordIdRequest{
				AvoidRecordKeys:     avoidRecordKeys,
				EndDataSourceCode:   xString(testCase.endDataSourceCode, endDataSourceCode),
				EndRecordId:         xString(testCase.endRecordId, endRecordId),
				Flags:               testCase.flags,
				MaxDegrees:          testCase.maxDegrees,
				RequiredDataSources: requiredDataSources,
				StartDataSourceCode: xString(testCase.startDataSourceCode, startDataSourceCode),
				StartRecordId:       xString(testCase.startRecordId, startRecordId),
			}
			actual, err := szEngine.FindPathByRecordId(ctx, request)
			printDebug(test, err, actual)
			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)
			}
		})
	}
}

func TestSzEngine_GetActiveConfigID(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	request := &szpb.GetActiveConfigIdRequest{}
	actual, err := szEngine.GetActiveConfigId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	flags := senzing.SzNoFlags
	request := &szpb.GetEntityByEntityIdRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	actual, err := szEngine.GetEntityByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetEntityByEntityID_badEntityID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetEntityByEntityID_nilEntityID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetEntityByRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.GetEntityByRecordIdRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.GetEntityByRecordId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetEntityByRecordID_badDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetEntityByRecordID_badRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetEntityByRecordID_nilDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetEntityByRecordID_nilRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetRecord(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.GetRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.GetRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetRecord_badDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetRecord_badRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetRecord_nilDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetRecord_nilRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	request := &szpb.GetRedoRecordRequest{}
	actual, err := szEngine.GetRedoRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetStats(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	request := &szpb.GetStatsRequest{}
	actual, err := szEngine.GetStats(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetVirtualEntityByRecordID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.ID + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.ID + `"}]}`
	flags := senzing.SzNoFlags
	request := &szpb.GetVirtualEntityByRecordIdRequest{
		Flags:      flags,
		RecordKeys: recordKeys,
	}
	actual, err := szEngine.GetVirtualEntityByRecordId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetVirtualEntityByRecordID_badDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_GetVirtualEntityByRecordID_badRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_HowEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	flags := senzing.SzNoFlags
	request := &szpb.HowEntityByEntityIdRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	actual, err := szEngine.HowEntityByEntityId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_HowEntityByEntityID_badEntityID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_HowEntityByEntityID_nilEntityID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_PreprocessRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	request := &szpb.PreprocessRecordRequest{
		Flags:            senzing.SzNoFlags,
		RecordDefinition: record.JSON,
	}
	actual, err := szEngine.PreprocessRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_PreprocessRecord_badRecordDefinition(test *testing.T) {
	// FIXME:
}

func TestSzEngine_PrimeEngine(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)
	request := &szpb.PrimeEngineRequest{}
	actual, err := szEngine.PrimeEngine(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ProcessRedoRecord(test *testing.T) {
	// FIXME: Implement TestSzEngine_ProcessRedoRecord
	_ = test
}

func TestSzEngine_ProcessRedoRecord_badRedoRecord(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ProcessRedoRecord_nilRedoRecord(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ProcessRedoRecord_withInfo(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ProcessRedoRecord_withInfo_badRedoRecord(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ProcessRedoRecord_withInfo_nilRedoRecord(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateEntity(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	flags := senzing.SzWithoutInfo
	request := &szpb.ReevaluateEntityRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	actual, err := szEngine.ReevaluateEntity(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual.GetResult())
}

func TestSzEngine_ReevaluateEntity_badEntityID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateEntity_nilEntityID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateEntity_withInfo(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	flags := senzing.SzWithInfo
	request := &szpb.ReevaluateEntityRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	actual, err := szEngine.ReevaluateEntity(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.NotEmpty(test, actual.GetResult())
}

func TestSzEngine_ReevaluateEntity_withInfo_badEntityID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateEntity_withInfo_nilEntityID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateRecord(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	request := &szpb.ReevaluateRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.ReevaluateRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Empty(test, actual.GetResult())
}

func TestSzEngine_ReevaluateRecord_badDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateRecord_badRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateRecord_nilDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateRecord_nilRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateRecord_withInfo(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithInfo
	request := &szpb.ReevaluateRecordRequest{
		DataSourceCode: record.DataSource,
		Flags:          flags,
		RecordId:       record.ID,
	}
	actual, err := szEngine.ReevaluateRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.NotEmpty(test, actual.GetResult())
}

func TestSzEngine_ReevaluateRecord_withInfo_badDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_ReevaluateRecord_withInfo_nilDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(test)

	requestToGetActiveConfigID := &szpb.GetActiveConfigIdRequest{}
	responseFromGetActiveConfigID, err := szEngine.GetActiveConfigId(ctx, requestToGetActiveConfigID)
	printDebug(test, err, responseFromGetActiveConfigID)
	require.NoError(test, err)

	request := &szpb.ReinitializeRequest{
		ConfigId: responseFromGetActiveConfigID.GetResult(),
	}
	actual, err := szEngine.Reinitialize(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_SearchByAttributes(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	request := &szpb.SearchByAttributesRequest{
		Attributes: attributes,
		Flags:      flags,
	}
	actual, err := szEngine.SearchByAttributes(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_SearchByAttributes_badAttributes(test *testing.T) {
	// FIXME:
}

func TestSzEngine_SearchByAttributes_badSearchProfile(test *testing.T) {
	// FIXME:
}

func TestSzEngine_SearchByAttributes_nilAttributes(test *testing.T) {
	// FIXME:
}

func TestSzEngine_SearchByAttributes_nilSearchProfile(test *testing.T) {
	// FIXME:
}

func TestSzEngine_SearchByAttributes_withSearchProfile(test *testing.T) {
	// FIXME:
}

func TestSzEngine_SearchByAttributes_searchProfile(test *testing.T) {
	// IMPROVE:  Use actual searchProfile
	ctx := test.Context()
	szEngine := getTestObject(test)
	flags := senzing.SzNoFlags
	request := &szpb.SearchByAttributesRequest{
		Attributes:    attributes,
		Flags:         flags,
		SearchProfile: searchProfile,
	}
	actual, err := szEngine.SearchByAttributes(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_StreamExportCsvEntityReport(test *testing.T) {
	// IMPROVE: Implement TestSzEngine_StreamExportCsvEntityReport
	_ = test
}

func TestSzEngine_StreamExportJsonEntityReport(test *testing.T) {
	// IMPROVE: Implement TestSzEngine_StreamExportJsonEntityReport
	_ = test
}

func TestSzEngine_WhyEntities(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID1 := getEntityID(truthset.CustomerRecords["1001"])
	entityID2 := getEntityID(truthset.CustomerRecords["1002"])
	flags := senzing.SzNoFlags
	request := &szpb.WhyEntitiesRequest{
		EntityId_1: entityID1,
		EntityId_2: entityID2,
		Flags:      flags,
	}
	actual, err := szEngine.WhyEntities(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_WhyEntities_badEnitity1(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyEntities_badEnitity2(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyEntities_nilEnitity1(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyEntities_nilEnitity2(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyRecordInEntity(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	record1 := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	request := &szpb.WhyRecordInEntityRequest{
		DataSourceCode: record1.DataSource,
		Flags:          flags,
		RecordId:       record1.ID,
	}
	actual, err := szEngine.WhyRecordInEntity(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_WhyRecordInEntity_badDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyRecordInEntity_badRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyRecordInEntity_nilDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyRecordInEntity_nilRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyRecords(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
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
	actual, err := szEngine.WhyRecords(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_WhyRecords_badDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyRecords_badRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyRecords_nilDataSourceCode(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhyRecords_nilRecordID(test *testing.T) {
	// FIXME:
}

func TestSzEngine_WhySearch(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(test)
	entityID := getEntityID(truthset.CustomerRecords["1001"])
	searchProfile := "SEARCH"

	flags := senzing.SzNoFlags
	request := &szpb.WhySearchRequest{
		Attributes:    attributes,
		EntityId:      entityID,
		SearchProfile: searchProfile,
		Flags:         flags,
	}
	actual, err := szEngine.WhySearch(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzEngine_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzEngine_SetLogLevel(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzEngine__SetLogLevel_badLevelName(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	printDebug(test, err)
	require.Error(test, err)
}

func TestSzEngine_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(test)
	testObject.SetObserverOrigin(ctx, observerOrigin)
}

func TestSzEngine_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(test)
	actual := testObject.GetObserverOrigin(ctx)
	assert.Equal(test, observerOrigin, actual)
}

func TestSzEngine_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(test)
	err := testObject.UnregisterObserver(ctx, observerSingleton)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Test
// ----------------------------------------------------------------------------

func TestBuildSimpleSystemConfigurationJsonUsingEnvVars(test *testing.T) {
	actual, err := settings.BuildSimpleSettingsUsingEnvVars()
	printDebug(test, err, actual)
	panicOnError(err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func addRecords(ctx context.Context, records []record.Record) {
	szEngine := getSzEngineServer(ctx)

	flags := senzing.SzWithoutInfo
	for _, record := range records {
		request := &szpb.AddRecordRequest{
			DataSourceCode:   record.DataSource,
			Flags:            flags,
			RecordDefinition: record.JSON,
			RecordId:         record.ID,
		}
		_, err := szEngine.AddRecord(ctx, request)
		panicOnError(err)
	}
}

func deleteRecords(ctx context.Context, records []record.Record) {
	szEngine := getSzEngineServer(ctx)

	flags := senzing.SzWithoutInfo
	for _, record := range records {
		request := &szpb.DeleteRecordRequest{
			DataSourceCode: record.DataSource,
			Flags:          flags,
			RecordId:       record.ID,
		}
		_, err := szEngine.DeleteRecord(ctx, request)
		panicOnError(err)
	}
}

func getEntityID(record record.Record) int64 {
	return getEntityIDForRecord(record.DataSource, record.ID)
}

func getEntityIDForRecord(datasource string, recordID string) int64 {
	ctx := context.TODO()
	szEngine := getSzEngineServer(ctx)
	request := &szpb.GetEntityByRecordIdRequest{
		DataSourceCode: datasource,
		RecordId:       recordID,
	}
	response, err := szEngine.GetEntityByRecordId(ctx, request)
	panicOnError(err)

	getEntityByRecordIDResponse := &GetEntityByRecordIDResponse{}
	err = json.Unmarshal([]byte(response.GetResult()), &getEntityByRecordIDResponse)
	panicOnError(err)

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

func getTestObject(t *testing.T) *szengineserver.SzEngineServer {
	t.Helper()

	return getSzEngineServer(t.Context())
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

func printDebug(t *testing.T, err error, items ...any) {
	t.Helper()
	printError(t, err)

	for item := range items {
		printActual(t, item)
	}
}

func printError(t *testing.T, err error) {
	t.Helper()

	if printErrors {
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}
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

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

type TestMetadataForFindPathByEntityID struct {
	avoidEntityIDs      func() string
	endEntityID         int64
	expectedErr         error
	expectedErrMessage  string
	flags               int64
	maxDegrees          int64
	name                string
	requiredDataSources func() string
	startEntityID       int64
}

type TestMetadataForFindPathByRecordID struct {
	avoidRecordKeys     func() string
	endDataSourceCode   string
	endRecordId         string
	expectedErr         error
	expectedErrMessage  string
	flags               int64
	maxDegrees          int64
	name                string
	requiredDataSources func() string
	startDataSourceCode string
	startRecordId       string
}

func getTestCasesForFindPathByEntityID() []TestMetadataForFindPathByEntityID {
	result := []TestMetadataForFindPathByEntityID{
		{
			name: "default",
		},
		{
			name:               "badStartEntityID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
			startEntityID:      badEntityID,
		},
		{
			name:               "badEndEntityID",
			endEntityID:        badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
		},
		{
			name:               "badMaxDegrees",
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0031|Invalid value of max degree '-1'"}}}`,
			maxDegrees:         badMaxDegrees,
		},
		{
			name:               "badAvoidEntityIDs",
			avoidEntityIDs:     badAvoidEntityIDsFunc,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044021","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
		},
		{
			name:                "badRequiredDataSource",
			expectedErr:         szerror.ErrSzBadInput,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044025","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
			requiredDataSources: badRequiredDataSourcesFunc,
		},
		{
			name:       "nilMaxDegrees",
			maxDegrees: nilMaxDegrees,
		},
		{
			name:           "nilAvoidEntityIDs",
			avoidEntityIDs: nilAvoidEntityIDsFunc,
		},
		{
			name:                "nilRequiredDataSource",
			requiredDataSources: nilRequiredDataSourcesFunc,
		},
		{
			name:           "avoiding",
			avoidEntityIDs: avoidEntityIDsFunc,
		},
		{
			name:               "avoiding_badStartEntityID",
			avoidEntityIDs:     avoidEntityIDsFunc,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044021","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
			startEntityID:      badEntityID,
		},
		{
			name:                "avoiding_and_including",
			avoidEntityIDs:      avoidEntityIDsFunc,
			maxDegrees:          1,
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:                "including",
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:                "including_badStartEntityID",
			expectedErr:         szerror.ErrSzNotFound,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044025","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
			requiredDataSources: requiredDataSourcesFunc,
			startEntityID:       badEntityID,
		},
	}
	return result
}

func getTestCasesForFindPathByRecordID() []TestMetadataForFindPathByRecordID {
	result := []TestMetadataForFindPathByRecordID{
		{
			name: "default",
		},
		{
			name:                "badDataSourceCode",
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
			startDataSourceCode: badDataSourceCode,
		},
		{
			name:               "badDataRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
			startRecordId:      badRecordID,
		},
		{
			name:               "badMaxDegrees",
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ0031|Invalid value of max degree '-1'"}}}`,
			maxDegrees:         badMaxDegrees,
		},
		{
			name:               "badAvoidRecordKeys",
			avoidRecordKeys:    badAvoidRecordIDsFunc,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
		},
		{
			name:               "badRequiredDataSources",
			avoidRecordKeys:    badRequiredDataSourcesFunc,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
		},
		{
			name:            "avoiding",
			avoidRecordKeys: avoidRecordIDsFunc,
			maxDegrees:      1,
		},
		{
			name:                "avoiding_badStartDataSourceCode",
			avoidRecordKeys:     avoidRecordIDsFunc,
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
			startDataSourceCode: badDataSourceCode,
		},
		{
			name:                "avoiding_and_including",
			avoidRecordKeys:     avoidRecordIDsFunc,
			maxDegrees:          1,
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:                "including",
			avoidRecordKeys:     avoidRecordIDsFunc,
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:                "including_badDataSourceCode",
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044027","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
			requiredDataSources: requiredDataSourcesFunc,
			startDataSourceCode: badDataSourceCode,
		},
	}
	return result
}

// Return first non-zero length candidate.  Last candidate is default.
func xString(candidates ...string) string {
	var result string
	for _, result = range candidates {
		if len(result) > 0 {
			return result
		}
	}
	return result
}

// Return first non-zero candidate.  Last candidate is default.
func xInt64(candidates ...int64) int64 {
	var result int64
	for _, result = range candidates {
		if result != 0 {
			return result
		}
	}
	return result
}

func avoidEntityIDsFunc() string {
	result := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord("CUSTOMERS", "1001") + `}]}`

	return result
}

func avoidRecordIDsFunc() string {
	result := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord("CUSTOMERS", "1001") + `}]}`

	return result
}

func badAvoidEntityIDsFunc() string {
	return "}{"
}

func badAvoidRecordIDsFunc() string {
	return "}{"
}

func badRequiredDataSourcesFunc() string {
	return "}{"
}

func nilAvoidEntityIDsFunc() string {
	var result string

	return result
}

func nilRequiredDataSourcesFunc() string {
	var result string

	return result
}

func requiredDataSourcesFunc() string {
	record := truthset.CustomerRecords["1001"]

	return `{"DATA_SOURCES": ["` + record.DataSource + `"]}`
}

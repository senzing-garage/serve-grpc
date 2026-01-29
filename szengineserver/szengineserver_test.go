package szengineserver_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-grpc/getversion"
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
	baseTen                    = 10
	defaultAttributes          = `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	defaultBuildOutDegrees     = int64(2)
	defaultBuildOutMaxEntities = int64(10)
	defaultMaxDegrees          = int64(2)
	defaultSearchProfile       = senzing.SzNoSearchProfile
	defaultTruncation          = 76
	jsonIndentation            = "    "
	observerID                 = "Observer 1"
	observerOrigin             = "SzEngine observer"
	printErrors                = false
	printResults               = false
)

// Bad parameters

const (
	badAttributes          = "}{"
	badBuildOutDegrees     = int64(-1)
	badBuildOutMaxEntities = int64(-1)
	badCsvColumnList       = "BAD, CSV, COLUMN, LIST"
	badDataSourceCode      = "BadDataSourceCode"
	badEntityID            = int64(-1)
	badMaxDegrees          = int64(-1)
	badRecordDefinition    = "}{"
	badRecordID            = "BadRecordID"
	badRedoRecord          = "{}"
	badSearchProfile       = "}{"
	nilSemaphoreString     = "xyzzy"
	nilSemaphoreInt64      = int64(-9999)
)

// Nil/empty parameters

var (
	nilAttributes          = nilSemaphoreString
	nilBuildOutDegrees     = nilSemaphoreInt64
	nilBuildOutMaxEntities = nilSemaphoreInt64
	nilCsvColumnList       string
	nilDataSourceCode      = nilSemaphoreString
	nilEntityID            = nilSemaphoreInt64
	nilMaxDegrees          = nilSemaphoreInt64
	nilRecordDefinition    = nilSemaphoreString
	nilRecordID            = nilSemaphoreString
	nilRedoRecord          = nilSemaphoreString
	nilSearchProfile       = nilSemaphoreString
)

var (
	logLevelName      = "INFO"
	observerSingleton = &observer.NullObserver{
		ID:       observerID,
		IsSilent: true,
	}
	senzingVersion        = 0
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

func TestSzEngine_PrintSenzingVersion(test *testing.T) {
	test.Logf("Detected Senzing version: %d", senzingVersion)
}

func TestSzEngine_AddRecord(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForAddRecord()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Defaults.
			szEngine := getTestObject(ctx, test)
			record := truthset.CustomerRecords["1001"]

			// Test.

			dataSourceCode := xString(testCase.dataSourceCode, record.DataSource)
			recordID := xString(testCase.recordID, record.ID)

			request := &szpb.AddRecordRequest{
				DataSourceCode:   dataSourceCode,
				Flags:            xInt64(testCase.flags, senzing.SzNoFlags),
				RecordDefinition: xString(testCase.recordDefinition, record.JSON),
				RecordId:         recordID,
			}

			actual, err := szEngine.AddRecord(ctx, request)
			printDebug(test, err, actual)

			if testCase.expectedErr != nil {
				require.ErrorIs(test, err, testCase.expectedErr)
				require.JSONEq(test, testCase.expectedErrMessage, err.Error())
			} else {
				require.NoError(test, err)

				deleteRequest := &szpb.DeleteRecordRequest{
					DataSourceCode: dataSourceCode,
					Flags:          senzing.SzNoFlags,
					RecordId:       recordID,
				}
				actual, err := szEngine.DeleteRecord(ctx, deleteRequest)
				printDebug(test, err, actual)
			}
		})
	}
}

func TestSzEngine_CountRedoRecords(test *testing.T) {
	ctx := test.Context()
	expected := int64(0)
	szEngine := getTestObject(ctx, test)
	request := &szpb.CountRedoRecordsRequest{}
	actual, err := szEngine.CountRedoRecords(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
	require.Equal(test, expected, actual.GetResult())
}

func TestSzEngine_DeleteRecord(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForDeleteRecord()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1005 := truthset.CustomerRecords["1005"]

			records := []record.Record{
				record1005,
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			// Test.

			dataSourceCode := xString(testCase.dataSourceCode, record1005.DataSource)
			recordID := xString(testCase.recordID, record1005.ID)

			request := &szpb.DeleteRecordRequest{
				DataSourceCode: dataSourceCode,
				Flags:          xInt64(testCase.flags, senzing.SzNoFlags),
				RecordId:       recordID,
			}

			actual, err := szEngine.DeleteRecord(ctx, request)
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

func TestSzEngine_ExportCsvEntityReport(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szEngine := getTestObject(ctx, test)
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

	szEngine := getTestObject(ctx, test)
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

	szEngine := getTestObject(ctx, test)
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

	szEngine := getTestObject(ctx, test)
	flags := senzing.SzNoFlags
	request := &szpb.ExportJsonEntityReportRequest{
		Flags: flags,
	}
	actual, err := szEngine.ExportJsonEntityReport(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ExportJSONEntityReport_65536(test *testing.T) {
	// IMPROVE: Implement TestSzEngine_ExportJSONEntityReport_65536
	_ = test
}

func TestSzEngine_FindInterestingEntitiesByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindInterestingEntitiesByEntityID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)
			entityID := getEntityID(ctx, record1001)

			// Test.

			request := &szpb.FindInterestingEntitiesByEntityIdRequest{
				EntityId: xInt64(testCase.entityID, entityID),
				Flags:    xInt64(testCase.flags, senzing.SzNoFlags),
			}
			actual, err := szEngine.FindInterestingEntitiesByEntityId(ctx, request)
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

func TestSzEngine_FindInterestingEntitiesByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindInterestingEntitiesByRecordID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			// Test.

			request := &szpb.FindInterestingEntitiesByRecordIdRequest{
				DataSourceCode: xString(testCase.dataSourceCode, record1001.DataSource),
				Flags:          xInt64(testCase.flags, senzing.SzNoFlags),
				RecordId:       xString(testCase.recordID, record1001.ID),
			}

			actual, err := szEngine.FindInterestingEntitiesByRecordId(ctx, request)
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

func TestSzEngine_FindNetworkByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindNetworkByEntityID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]
			record1002 := truthset.CustomerRecords["1002"]

			records := []record.Record{
				record1001,
				record1002,
				truthset.CustomerRecords["1003"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			entityIDs := entityIDsJSON(
				getEntityID(ctx, record1001),
				getEntityID(ctx, record1002))

			// Test.

			request := &szpb.FindNetworkByEntityIdRequest{
				EntityIds:           entityIDs,
				Flags:               xInt64(testCase.flags, senzing.SzFindNetworkDefaultFlags),
				MaxDegrees:          xInt64(testCase.maxDegrees, defaultMaxDegrees),
				BuildOutDegrees:     xInt64(testCase.buildOutDegrees, defaultBuildOutDegrees),
				BuildOutMaxEntities: xInt64(testCase.buildOutMaxEntities, defaultBuildOutMaxEntities),
			}
			actual, err := szEngine.FindNetworkByEntityId(ctx, request)
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

// func TestSzEngine_FindNetworkByEntityID_Not_Using_CUSTOMERS(test *testing.T) {
// 	ctx := test.Context()

// 	record2012 := truthset.ReferenceRecords["2012"]
// 	record2013 := truthset.ReferenceRecords["2013"]

// 	records := []record.Record{
// 		record2012,
// 		record2013,
// 		truthset.CustomerRecords["1003"],
// 	}

// 	defer func() { deleteRecords(ctx, records) }()

// 	addRecords(ctx, records)

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
// 	// require.ErrorIs(test, err, szerror.ErrSzNotFound)
// 	expectedErrMessage := ""
// 	require.JSONEq(test, expectedErrMessage, err.Error())
// }

func TestSzEngine_FindNetworkByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForFindNetworkByRecordID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			recordKeys := recordKeysFunc()
			if testCase.recordKeys != nil {
				recordKeys = testCase.recordKeys()
			}

			// Test.

			request := &szpb.FindNetworkByRecordIdRequest{
				RecordKeys:          recordKeys,
				Flags:               xInt64(testCase.flags, senzing.SzFindNetworkDefaultFlags),
				MaxDegrees:          xInt64(testCase.maxDegrees, defaultMaxDegrees),
				BuildOutDegrees:     xInt64(testCase.buildOutDegrees, defaultBuildOutDegrees),
				BuildOutMaxEntities: xInt64(testCase.buildOutMaxEntities, defaultBuildOutMaxEntities),
			}
			actual, err := szEngine.FindNetworkByRecordId(ctx, request)
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

			szEngine := getTestObject(ctx, test)
			startEntityID := getEntityID(ctx, truthset.CustomerRecords["1001"])
			endEntityID := getEntityID(ctx, truthset.CustomerRecords["1002"])

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
				Flags:               xInt64(testCase.flags),
				MaxDegrees:          xInt64(testCase.maxDegrees),
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

			szEngine := getTestObject(ctx, test)

			record1 := truthset.CustomerRecords["1001"]
			startDataSourceCode := record1.DataSource
			startRecordID := record1.ID

			record2 := truthset.CustomerRecords["1002"]
			endDataSourceCode := record2.DataSource
			endRecordID := record2.ID

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
				EndRecordId:         xString(testCase.endRecordID, endRecordID),
				Flags:               xInt64(testCase.flags, senzing.SzFindNetworkDefaultFlags),
				MaxDegrees:          xInt64(testCase.maxDegrees, defaultMaxDegrees),
				RequiredDataSources: requiredDataSources,
				StartDataSourceCode: xString(testCase.startDataSourceCode, startDataSourceCode),
				StartRecordId:       xString(testCase.startRecordID, startRecordID),
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
	szEngine := getTestObject(ctx, test)
	request := &szpb.GetActiveConfigIdRequest{}
	actual, err := szEngine.GetActiveConfigId(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetEntityByEntityID()

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

			szEngine := getTestObject(ctx, test)
			entityID := getEntityID(ctx, truthset.CustomerRecords["1001"])

			// Test.

			request := &szpb.GetEntityByEntityIdRequest{
				EntityId: xInt64(testCase.entityID, entityID),
				Flags:    xInt64(testCase.flags, senzing.SzEntityDefaultFlags),
			}
			actual, err := szEngine.GetEntityByEntityId(ctx, request)
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

func TestSzEngine_GetEntityByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetEntityByRecordID()

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

			szEngine := getTestObject(ctx, test)
			record := truthset.CustomerRecords["1001"]

			// Test.

			request := &szpb.GetEntityByRecordIdRequest{
				DataSourceCode: xString(testCase.dataSourceCode, record.DataSource),
				Flags:          xInt64(testCase.flags, senzing.SzEntityDefaultFlags),
				RecordId:       xString(testCase.recordID, record.ID),
			}
			actual, err := szEngine.GetEntityByRecordId(ctx, request)
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

func TestSzEngine_GetRecord(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetRecord()

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

			szEngine := getTestObject(ctx, test)
			record := truthset.CustomerRecords["1001"]

			// Test.

			request := &szpb.GetRecordRequest{
				DataSourceCode: xString(testCase.dataSourceCode, record.DataSource),
				Flags:          xInt64(testCase.flags, senzing.SzRecordDefaultFlags),
				RecordId:       xString(testCase.recordID, record.ID),
			}
			actual, err := szEngine.GetRecord(ctx, request)
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

func TestSzEngine_GetRedoRecord(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	request := &szpb.GetRedoRecordRequest{}
	actual, err := szEngine.GetRedoRecord(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetStats(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	request := &szpb.GetStatsRequest{}
	actual, err := szEngine.GetStats(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_GetVirtualEntityByRecordID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetVirtualEntityByRecordID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			recordKeys := recordKeysFunc()
			if testCase.recordKeys != nil {
				recordKeys = testCase.recordKeys()
			}

			// Test.

			request := &szpb.GetVirtualEntityByRecordIdRequest{
				Flags:      xInt64(testCase.flags, senzing.SzRecordDefaultFlags),
				RecordKeys: recordKeys,
			}
			actual, err := szEngine.GetVirtualEntityByRecordId(ctx, request)
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

func TestSzEngine_HowEntityByEntityID(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForHowEntityByEntityID()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
				truthset.CustomerRecords["1002"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)
			entityID := getEntityID(ctx, record1001)

			// Test.

			request := &szpb.HowEntityByEntityIdRequest{
				Flags:    xInt64(testCase.flags, senzing.SzRecordDefaultFlags),
				EntityId: xInt64(testCase.entityID, entityID),
			}
			actual, err := szEngine.HowEntityByEntityId(ctx, request)
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

func TestSzEngine_GetRecordPreview(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForGetRecordPreview()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			// Test.

			request := &szpb.GetRecordPreviewRequest{
				Flags:            xInt64(testCase.flags, senzing.SzRecordDefaultFlags),
				RecordDefinition: xString(testCase.recordDefinition, record1001.JSON),
			}
			actual, err := szEngine.GetRecordPreview(ctx, request)
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

func TestSzEngine_PrimeEngine(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)
	request := &szpb.PrimeEngineRequest{}
	actual, err := szEngine.PrimeEngine(ctx, request)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzEngine_ProcessRedoRecord(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForProcessRedoRecord()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
				truthset.CustomerRecords["1004"],
				truthset.CustomerRecords["1005"],
				truthset.CustomerRecords["1009"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			redoRecordRequest := &szpb.GetRedoRecordRequest{}
			redoRecordResponse, err := szEngine.GetRedoRecord(ctx, redoRecordRequest)
			require.NoError(test, err)

			redoRecord := xString(testCase.redoRecord, redoRecordResponse.GetResult())

			// Test.

			request := &szpb.ProcessRedoRecordRequest{
				Flags:      xInt64(testCase.flags, senzing.SzRecordDefaultFlags),
				RedoRecord: redoRecord,
			}
			actual, err := szEngine.ProcessRedoRecord(ctx, request)
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

func TestSzEngine_ReevaluateEntity(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForReevaluateEntity()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)
			entityID := getEntityID(ctx, record1001)

			// Test.

			request := &szpb.ReevaluateEntityRequest{
				Flags:    xInt64(testCase.flags, senzing.SzNoFlags),
				EntityId: xInt64(testCase.entityID, entityID),
			}
			actual, err := szEngine.ReevaluateEntity(ctx, request)
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

func TestSzEngine_ReevaluateRecord(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForReevaluateRecord()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			// Test.

			request := &szpb.ReevaluateRecordRequest{
				DataSourceCode: xString(testCase.dataSourceCode, record1001.DataSource),
				Flags:          xInt64(testCase.flags, senzing.SzNoFlags),
				RecordId:       xString(testCase.recordID, record1001.ID),
			}
			actual, err := szEngine.ReevaluateRecord(ctx, request)
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

func TestSzEngine_Reinitialize(test *testing.T) {
	ctx := test.Context()
	szEngine := getTestObject(ctx, test)

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
	testCases := getTestCasesForSearchByAttributes()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			// Test.

			request := &szpb.SearchByAttributesRequest{
				Attributes:    xString(testCase.attributes, defaultAttributes),
				Flags:         xInt64(testCase.flags, senzing.SzNoFlags),
				SearchProfile: xString(testCase.searchProfile, defaultSearchProfile),
			}
			actual, err := szEngine.SearchByAttributes(ctx, request)
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
	testCases := getTestCasesForWhyEntities()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]
			record1002 := truthset.CustomerRecords["1002"]

			records := []record.Record{
				record1001,
				record1002,
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)
			entityID1 := getEntityID(ctx, record1001)
			entityID2 := getEntityID(ctx, record1002)

			// Test.

			request := &szpb.WhyEntitiesRequest{
				EntityId_1: xInt64(testCase.entityID1, entityID1),
				EntityId_2: xInt64(testCase.entityID2, entityID2),
				Flags:      xInt64(testCase.flags, senzing.SzNoFlags),
			}
			actual, err := szEngine.WhyEntities(ctx, request)
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

func TestSzEngine_WhyRecordInEntity(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForWhyRecordInEntity()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]

			records := []record.Record{
				record1001,
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			// Test.

			request := &szpb.WhyRecordInEntityRequest{
				DataSourceCode: xString(testCase.dataSourceCode, record1001.DataSource),
				Flags:          xInt64(testCase.flags, senzing.SzNoFlags),
				RecordId:       xString(testCase.recordID, record1001.ID),
			}
			actual, err := szEngine.WhyRecordInEntity(ctx, request)
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

func TestSzEngine_WhyRecords(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForWhyRecords()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			record1001 := truthset.CustomerRecords["1001"]
			record1002 := truthset.CustomerRecords["1002"]

			records := []record.Record{
				record1001,
				record1002,
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)

			// Test.

			request := &szpb.WhyRecordsRequest{
				DataSourceCode_1: xString(testCase.dataSourceCode1, record1001.DataSource),
				DataSourceCode_2: xString(testCase.dataSourceCode2, record1002.DataSource),
				Flags:            xInt64(testCase.flags, senzing.SzNoFlags),
				RecordId_1:       xString(testCase.recordID1, record1001.ID),
				RecordId_2:       xString(testCase.recordID2, record1002.ID),
			}
			actual, err := szEngine.WhyRecords(ctx, request)
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

func TestSzEngine_WhySearch(test *testing.T) {
	ctx := test.Context()
	testCases := getTestCasesForWhySearch()

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			// Insert test data.
			records := []record.Record{
				truthset.CustomerRecords["1001"],
				truthset.CustomerRecords["1002"],
				truthset.CustomerRecords["1003"],
			}

			defer func() { deleteRecords(ctx, records) }()

			addRecords(ctx, records)

			// Defaults.

			szEngine := getTestObject(ctx, test)
			entityID := getEntityID(ctx, truthset.CustomerRecords["1001"])

			// Test.

			request := &szpb.WhySearchRequest{
				Attributes:    xString(testCase.attributes, defaultAttributes),
				EntityId:      xInt64(testCase.entityID, entityID),
				Flags:         xInt64(testCase.flags, senzing.SzNoFlags),
				SearchProfile: xString(testCase.searchProfile, defaultSearchProfile),
			}
			actual, err := szEngine.WhySearch(ctx, request)
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

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzEngine_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.RegisterObserver(ctx, observerSingleton)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzEngine_SetLogLevel(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "DEBUG")
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzEngine_SetLogLevel_badLevelName(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, "BADLEVELNAME")
	printDebug(test, err)
	require.Error(test, err)
}

func TestSzEngine_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	testObject.SetObserverOrigin(ctx, observerOrigin)
}

func TestSzEngine_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	actual := testObject.GetObserverOrigin(ctx)
	assert.Equal(test, observerOrigin, actual)
}

func TestSzEngine_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
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

func getEntityID(ctx context.Context, record record.Record) int64 {
	return getEntityIDForRecord(ctx, record.DataSource, record.ID)
}

func getEntityIDForRecord(ctx context.Context, datasource string, recordID string) int64 {
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

// func getEntityIDString(record record.Record) string {
// 	entityID := getEntityID(record)

// 	return strconv.FormatInt(entityID, 10)
// }

func getEntityIDStringForRecord(datasource string, id string) string { //nolint
	ctx := context.Background()
	entityID := getEntityIDForRecord(ctx, datasource, id)

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
	ctx := context.Background()
	senzingVersion = getversion.GetSenzingVersion(ctx)

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
		_, err := szConfig.RegisterDataSource(ctx, dataSourceCode)
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
// Test structs
// ----------------------------------------------------------------------------

type TestMetadataForAddRecord struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordDefinition   string
	recordID           string
}

type TestMetadataForDeleteRecord struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForFindInterestingEntitiesByEntityID struct {
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}
type TestMetadataForFindInterestingEntitiesByRecordID struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForFindNetworkByEntityID struct {
	buildOutDegrees     int64
	buildOutMaxEntities int64
	entityIDs           func() string
	expectedErr         error
	expectedErrMessage  string
	flags               int64
	maxDegrees          int64
	name                string
}

type TestMetadataForFindNetworkByRecordID struct {
	buildOutDegrees     int64
	buildOutMaxEntities int64
	expectedErr         error
	expectedErrMessage  string
	flags               int64
	maxDegrees          int64
	name                string
	recordKeys          func() string
}

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
	endRecordID         string
	expectedErr         error
	expectedErrMessage  string
	flags               int64
	maxDegrees          int64
	name                string
	requiredDataSources func() string
	startDataSourceCode string
	startRecordID       string
}

type TestMetadataForGetEntityByEntityID struct {
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}

type TestMetadataForGetEntityByRecordID struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForGetRecord struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForGetVirtualEntityByRecordID struct {
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordKeys         func() string
}

type TestMetadataForHowEntityByEntityID struct {
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}

type TestMetadataForGetRecordPreview struct {
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordDefinition   string
}

type TestMetadataForProcessRedoRecord struct {
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	redoRecord         string
}

type TestMetadataForReevaluateEntity struct {
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}

type TestMetadataForReevaluateRecord struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForSearchByAttributes struct {
	attributes         string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	searchProfile      string
}

type TestMetadataForWhyEntities struct {
	entityID1          int64
	entityID2          int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
}

type TestMetadataForWhyRecordInEntity struct {
	dataSourceCode     string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID           string
}

type TestMetadataForWhyRecords struct {
	dataSourceCode1    string
	dataSourceCode2    string
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	recordID1          string
	recordID2          string
}

type TestMetadataForWhySearch struct {
	attributes         string
	entityID           int64
	expectedErr        error
	expectedErrMessage string
	flags              int64
	name               string
	searchProfile      string
}

// ----------------------------------------------------------------------------
// Test data
// ----------------------------------------------------------------------------

func getTestCasesForAddRecord() []TestMetadataForAddRecord {
	record1002 := truthset.CustomerRecords["1002"]
	result := []TestMetadataForAddRecord{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}}`,
		},
		{
			name:               "badDataSourceCode_asErrSz",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}}`,
		},
		{
			name:               "badDataSourceCodeInJSON",
			dataSourceCode:     record1002.DataSource,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0023|Conflicting DATA_SOURCE values 'CUSTOMERS' and 'BOB'"}}}`,
			recordDefinition:   `{"DATA_SOURCE": "BOB", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`,
			recordID:           record1002.ID,
		},
		{
			name:               "badRecordDefinition",
			recordDefinition:   badRecordDefinition,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
		},
		{
			name:               "badRecordID",
			recordID:           badRecordID,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ0024|Conflicting RECORD_ID values 'BadRecordID' and '1001'"}}}`,
		},
		{
			name: "default",
		},
		{
			name:           "nilDataSourceCode",
			dataSourceCode: nilDataSourceCode,
		},
		{
			name:               "nilRecordDefinition",
			recordDefinition:   nilRecordDefinition,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044001","reason":"SENZ3121|JSON Parsing Failure [code=1,offset=0]"}}}`,
		},
		{
			name:     "nilRecordID",
			recordID: nilRecordID,
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
		{
			name:               "withInfo_badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			flags:              senzing.SzWithInfo,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ0023|Conflicting DATA_SOURCE values 'BADDATASOURCECODE' and 'CUSTOMERS'"}}}`,
		},
		{
			name:               "withInfo_badRecordDefinition",
			flags:              senzing.SzWithInfo,
			recordDefinition:   badRecordDefinition,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
		},
		{
			name:               "withInfo_badRecordID",
			flags:              senzing.SzWithInfo,
			recordID:           badRecordID,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ0024|Conflicting RECORD_ID values 'BadRecordID' and '1001'"}}}`,
		},
		{
			name:           "withInfo_nilDataSourceCode",
			dataSourceCode: nilDataSourceCode,
			flags:          senzing.SzWithInfo,
		},
		{
			name:               "withInfo_nilRecordDefinition",
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).AddRecord","error":{"function":"szengine.(*Szengine).AddRecord","error":{"id":"SZSDK60044002","reason":"SENZ3121|JSON Parsing Failure [code=1,offset=0]"}}}`,
			flags:              senzing.SzWithInfo,
			recordDefinition:   nilRecordDefinition,
		},
		{
			name:     "withInfo_nilRecordID",
			flags:    senzing.SzWithInfo,
			recordID: nilRecordID,
		},
	}

	return result
}

func getTestCasesForDeleteRecord() []TestMetadataForDeleteRecord {
	var addendum []TestMetadataForDeleteRecord

	result := []TestMetadataForDeleteRecord{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:               "badDataSourceCode_asSzBadInputError",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:               "badDataSourceCode_asSzErr",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:     "badRecordID",
			recordID: badRecordID,
		},
		{
			name: "default",
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
		{
			name:               "withInfo_badDataSourceCode",
			flags:              senzing.SzWithInfo,
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044005","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:     "withInfo_badRecordID",
			flags:    senzing.SzWithInfo,
			recordID: badRecordID,
		},
	}

	switch {
	case senzingVersion < 40201:
		addendum = []TestMetadataForDeleteRecord{
			{
				name:               "nilDataSourceCode",
				dataSourceCode:     nilDataSourceCode,
				expectedErr:        szerror.ErrSzConfiguration,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}}`,
			},
			{
				name:               "nilRecordID",
				expectedErr:        szerror.ErrSzNotInitialized,
				expectedErrMessage: `{"function": "szengineserver.(*SzEngineServer).DeleteRecord", "error": {"function": "szengine.(*Szengine).DeleteRecord", "error":{"id":"SZSDK60044004","reason":"SENZ0053|RECORD_ID must be provided"}}}`,
				recordID:           nilRecordID,
			},
			{
				name:               "withInfo_nilDataSourceCode",
				flags:              senzing.SzWithInfo,
				dataSourceCode:     nilDataSourceCode,
				expectedErr:        szerror.ErrSzConfiguration,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044005","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}}`,
			},
			{
				name:               "withInfo_nilRecordID",
				expectedErr:        szerror.ErrSzNotInitialized,
				expectedErrMessage: `{"function": "szengineserver.(*SzEngineServer).DeleteRecord", "error": {"function": "szengine.(*Szengine).DeleteRecord", "error":{"id":"SZSDK60044005","reason":"SENZ0053|RECORD_ID must be provided"}}}`,
				flags:              senzing.SzWithInfo,
				recordID:           nilRecordID,
			},
		}

	default:
		addendum = []TestMetadataForDeleteRecord{
			{
				name:               "nilDataSourceCode",
				dataSourceCode:     nilDataSourceCode,
				expectedErr:        szerror.ErrSzUnknownDataSource,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044004","reason":"SENZ2207|Data source code [] does not exist."}}}`,
			},
			{
				name:               "nilRecordID",
				expectedErr:        szerror.ErrSz,
				expectedErrMessage: `{"function": "szengineserver.(*SzEngineServer).DeleteRecord", "error": {"function": "szengine.(*Szengine).DeleteRecord", "error":{"id":"SZSDK60044004","reason":"SENZ2288|No record ID found."}}}`,
				recordID:           nilRecordID,
			},
			{
				name:               "withInfo_nilDataSourceCode",
				flags:              senzing.SzWithInfo,
				dataSourceCode:     nilDataSourceCode,
				expectedErr:        szerror.ErrSzUnknownDataSource,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).DeleteRecord","error":{"function":"szengine.(*Szengine).DeleteRecord","error":{"id":"SZSDK60044005","reason":"SENZ2207|Data source code [] does not exist."}}}`,
			},
			{
				name:               "withInfo_nilRecordID",
				expectedErr:        szerror.ErrSz,
				expectedErrMessage: `{"function": "szengineserver.(*SzEngineServer).DeleteRecord", "error": {"function": "szengine.(*Szengine).DeleteRecord", "error":{"id":"SZSDK60044005","reason":"SENZ2288|No record ID found."}}}`,
				flags:              senzing.SzWithInfo,
				recordID:           nilRecordID,
			},
		}
	}

	result = append(result, addendum...)

	return result
}

func getTestCasesForFindInterestingEntitiesByEntityID() []TestMetadataForFindInterestingEntitiesByEntityID {
	result := []TestMetadataForFindInterestingEntitiesByEntityID{
		{
			name:               "badEntityID",
			entityID:           badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByEntityId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByEntityID","error":{"id":"SZSDK60044010","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilEntityID",
			entityID:           nilEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByEntityId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByEntityID","error":{"id":"SZSDK60044010","reason":"SENZ0037|Unknown resolved entity value '0'"}}}`,
		},
	}

	return result
}

func getTestCasesForFindInterestingEntitiesByRecordID() []TestMetadataForFindInterestingEntitiesByRecordID {
	result := []TestMetadataForFindInterestingEntitiesByRecordID{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByRecordId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:               "badRecordID",
			recordID:           badRecordID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByRecordId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByRecordId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ2207|Data source code [] does not exist."}}}`,
		},
		{
			name:               "nilRecordID",
			recordID:           nilRecordID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindInterestingEntitiesByRecordId","error":{"function":"szengine.(*Szengine).FindInterestingEntitiesByRecordID","error":{"id":"SZSDK60044011","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}}`,
		},
	}

	return result
}

func getTestCasesForFindNetworkByEntityID() []TestMetadataForFindNetworkByEntityID {
	result := []TestMetadataForFindNetworkByEntityID{
		{
			name:               "badBuildOutDegrees",
			buildOutDegrees:    badBuildOutDegrees,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindNetworkByEntityId","error":{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0032|Invalid value of build out degree '-1'"}}}`,
		},
		{
			name:                "badBuildOutMaxEntities",
			buildOutMaxEntities: badBuildOutMaxEntities,
			expectedErr:         szerror.ErrSz,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindNetworkByEntityId","error":{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0029|Invalid value of max entities '-1'"}}}`,
		},
		{
			name:      "badEntityIDs",
			entityIDs: badEntityIDsFunc,
			// IMPROVE: Shouldn't this error?
		},
		{
			name:               "badMaxDegrees",
			maxDegrees:         badMaxDegrees,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindNetworkByEntityId","error":{"function":"szengine.(*Szengine).FindNetworkByEntityID","error":{"id":"SZSDK60044013","reason":"SENZ0031|Invalid value of max degree '-1'"}}}`,
		},
		{
			name: "default",
		},
		{
			name:            "nilBuildOutDegrees",
			buildOutDegrees: nilBuildOutDegrees,
		},
		{
			name:                "nilBuildOutMaxEntities",
			buildOutMaxEntities: nilBuildOutMaxEntities,
		},
		{
			name:       "nilMaxDegrees",
			maxDegrees: nilMaxDegrees,
		},
	}

	return result
}

func getTestCasesForFindNetworkByRecordID() []TestMetadataForFindNetworkByRecordID {
	result := []TestMetadataForFindNetworkByRecordID{
		{
			name:               "badBuildOutDegrees",
			buildOutDegrees:    badBuildOutDegrees,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindNetworkByRecordId","error":{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ0032|Invalid value of build out degree '-1'"}}}`,
		},
		{
			name:               "badBuildOutMaxEntities",
			buildOutDegrees:    badBuildOutMaxEntities,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindNetworkByRecordId","error":{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ0032|Invalid value of build out degree '-1'"}}}`,
		},
		{
			name:               "badDataSourceCode",
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindNetworkByRecordId","error":{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
			recordKeys:         recordKeysBadDataSourceCodeFunc,
		},
		{
			name:               "badMaxDegree",
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindNetworkByRecordId","error":{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ0031|Invalid value of max degree '-1'"}}}`,
			maxDegrees:         badMaxDegrees,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindNetworkByRecordId","error":{"function":"szengine.(*Szengine).FindNetworkByRecordID","error":{"id":"SZSDK60044015","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
			recordKeys:         recordKeysFuncBadRecordIDFunc,
		},
		{
			name: "default",
		},
	}

	return result
}

func getTestCasesForFindPathByEntityID() []TestMetadataForFindPathByEntityID {
	result := []TestMetadataForFindPathByEntityID{
		{
			name:           "avoiding",
			avoidEntityIDs: avoidEntityIDsFunc,
		},
		{
			name:                "avoiding_and_including",
			avoidEntityIDs:      avoidEntityIDsFunc,
			maxDegrees:          1,
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:               "avoiding_badStartEntityID",
			avoidEntityIDs:     avoidEntityIDsFunc,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044021","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
			startEntityID:      badEntityID,
		},
		{
			name:               "badAvoidEntityIDs",
			avoidEntityIDs:     badAvoidEntityIDsFunc,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044021","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
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
			name:                "badRequiredDataSource",
			expectedErr:         szerror.ErrSzBadInput,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044025","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
			requiredDataSources: badRequiredDataSourcesFunc,
		},
		{
			name:               "badStartEntityID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByEntityId","error":{"function":"szengine.(*Szengine).FindPathByEntityID","error":{"id":"SZSDK60044017","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
			startEntityID:      badEntityID,
		},
		{
			name: "default",
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
		{
			name:           "nilAvoidEntityIDs",
			avoidEntityIDs: nilAvoidEntityIDsFunc,
		},
		{
			name:       "nilMaxDegrees",
			maxDegrees: nilMaxDegrees,
		},
		{
			name:                "nilRequiredDataSource",
			requiredDataSources: nilRequiredDataSourcesFunc,
		},
	}

	return result
}

func getTestCasesForFindPathByRecordID() []TestMetadataForFindPathByRecordID {
	result := []TestMetadataForFindPathByRecordID{
		{
			name:            "avoiding",
			avoidRecordKeys: avoidRecordIDsFunc,
			maxDegrees:      1,
		},
		{
			name:                "avoiding_and_including",
			avoidRecordKeys:     avoidRecordIDsFunc,
			maxDegrees:          1,
			requiredDataSources: requiredDataSourcesFunc,
		},
		{
			name:                "avoiding_badStartDataSourceCode",
			avoidRecordKeys:     avoidRecordIDsFunc,
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
			startDataSourceCode: badDataSourceCode,
		},
		{
			name:               "badAvoidRecordKeys",
			avoidRecordKeys:    badAvoidRecordIDsFunc,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
		},
		{
			name:               "badDataRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
			startRecordID:      badRecordID,
		},
		{
			name:                "badDataSourceCode",
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
			startDataSourceCode: badDataSourceCode,
		},
		{
			name:               "badMaxDegrees",
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ0031|Invalid value of max degree '-1'"}}}`,
			maxDegrees:         badMaxDegrees,
		},
		{
			name:               "badRequiredDataSources",
			avoidRecordKeys:    badRequiredDataSourcesFunc,
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044023","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
		},
		{
			name: "default",
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
		{
			name:                "nilDataSourceCode",
			expectedErr:         szerror.ErrSzUnknownDataSource,
			expectedErrMessage:  `{"function":"szengineserver.(*SzEngineServer).FindPathByRecordId","error":{"function":"szengine.(*Szengine).FindPathByRecordID","error":{"id":"SZSDK60044019","reason":"SENZ2207|Data source code [] does not exist."}}}`,
			startDataSourceCode: nilDataSourceCode,
		},
	}

	return result
}

func getTestCasesForGetEntityByEntityID() []TestMetadataForGetEntityByEntityID {
	result := []TestMetadataForGetEntityByEntityID{
		{
			name:               "badEntityID",
			entityID:           badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetEntityByEntityId","error":{"function":"szengine.(*Szengine).GetEntityByEntityID","error":{"id":"SZSDK60044030","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilEntityID",
			entityID:           nilEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetEntityByEntityId","error":{"function":"szengine.(*Szengine).GetEntityByEntityID","error":{"id":"SZSDK60044030","reason":"SENZ0037|Unknown resolved entity value '0'"}}}`,
		},
	}

	return result
}

func getTestCasesForGetEntityByRecordID() []TestMetadataForGetEntityByRecordID {
	result := []TestMetadataForGetEntityByRecordID{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetEntityByRecordId","error":{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetEntityByRecordId","error":{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
			recordID:           badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetEntityByRecordId","error":{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ2207|Data source code [] does not exist."}}}`,
		},
		{
			name:               "nilRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetEntityByRecordId","error":{"function":"szengine.(*Szengine).GetEntityByRecordID","error":{"id":"SZSDK60044032","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}}`,
			recordID:           nilRecordID,
		},
	}

	return result
}

func getTestCasesForGetRecord() []TestMetadataForGetRecord {
	result := []TestMetadataForGetRecord{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetRecord","error":{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetRecord","error":{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
			recordID:           badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetRecord","error":{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ2207|Data source code [] does not exist."}}}`,
		},
		{
			name:               "nilRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetRecord","error":{"function":"szengine.(*Szengine).GetRecord","error":{"id":"SZSDK60044035","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}}`,
			recordID:           nilRecordID,
		},
	}

	return result
}

func getTestCasesForGetVirtualEntityByRecordID() []TestMetadataForGetVirtualEntityByRecordID {
	result := []TestMetadataForGetVirtualEntityByRecordID{
		{
			name:               "badDataSourceCode",
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetVirtualEntityByRecordId","error":{"function":"szengine.(*Szengine).GetVirtualEntityByRecordID","error":{"id":"SZSDK60044038","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
			recordKeys:         badDataSourcesFunc,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetVirtualEntityByRecordId","error":{"function":"szengine.(*Szengine).GetVirtualEntityByRecordID","error":{"id":"SZSDK60044038","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
			recordKeys:         badRecordKeysFunc,
		},
		{
			name: "default",
		},
	}

	return result
}

func getTestCasesForHowEntityByEntityID() []TestMetadataForHowEntityByEntityID {
	result := []TestMetadataForHowEntityByEntityID{
		{
			name:               "badEntityID",
			entityID:           badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).HowEntityByEntityId","error":{"function":"szengine.(*Szengine).HowEntityByEntityID","error":{"id":"SZSDK60044040","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilEntityID",
			entityID:           nilEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).HowEntityByEntityId","error":{"function":"szengine.(*Szengine).HowEntityByEntityID","error":{"id":"SZSDK60044040","reason":"SENZ0037|Unknown resolved entity value '0'"}}}`,
		},
	}

	return result
}

func getTestCasesForGetRecordPreview() []TestMetadataForGetRecordPreview {
	var addendum []TestMetadataForGetRecordPreview

	result := []TestMetadataForGetRecordPreview{
		{
			name: "default",
		},
	}

	switch {
	case senzingVersion < 40200:
		addendum = []TestMetadataForGetRecordPreview{
			{
				name:               "badRecordDefinition",
				expectedErr:        szerror.ErrSzBadInput,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetRecordPreview","error":{"function":"szengine.(*Szengine).GetRecordPreview","error":{"id":"SZSDK60044061","reason":"SENZ0002|Invalid Message"}}}`,
				recordDefinition:   badRecordDefinition,
			},
		}

	default:
		addendum = []TestMetadataForGetRecordPreview{
			{
				name:               "badRecordDefinition",
				expectedErr:        szerror.ErrSzBadInput,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).GetRecordPreview","error":{"function":"szengine.(*Szengine).GetRecordPreview","error":{"id":"SZSDK60044061","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}`,
				recordDefinition:   badRecordDefinition,
			},
		}
	}

	result = append(result, addendum...)

	return result
}

func getTestCasesForProcessRedoRecord() []TestMetadataForProcessRedoRecord {
	var addendum []TestMetadataForProcessRedoRecord

	result := []TestMetadataForProcessRedoRecord{
		{
			name: "default",
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
	}

	switch {
	case senzingVersion < 40201:
		addendum = []TestMetadataForProcessRedoRecord{
			{
				name:               "badRedoRecord",
				expectedErr:        szerror.ErrSzConfiguration,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ProcessRedoRecord","error":{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044044","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}}`,
				redoRecord:         badRedoRecord,
			},
			{
				name:               "nilRedoRecord",
				expectedErr:        szerror.ErrSzBadInput,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ProcessRedoRecord","error":{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044044","reason":"SENZ0007|Empty Message"}}}`,
				redoRecord:         nilRedoRecord,
			},
			{
				name:               "withInfo_badRedoRecord",
				expectedErr:        szerror.ErrSzConfiguration,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ProcessRedoRecord","error":{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044045","reason":"SENZ2136|Error in input mapping, missing required field[DATA_SOURCE]"}}}`,
				flags:              senzing.SzWithInfo,
				redoRecord:         badRedoRecord,
			},
			{
				name:               "withInfo_nilRedoRecord",
				flags:              senzing.SzWithInfo,
				expectedErr:        szerror.ErrSzBadInput,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ProcessRedoRecord","error":{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044045","reason":"SENZ0007|Empty Message"}}}`,
				redoRecord:         nilRedoRecord,
			},
		}

	default:
		addendum = []TestMetadataForProcessRedoRecord{
			{
				name:               "badRedoRecord",
				expectedErr:        szerror.ErrSzUnknownDataSource,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ProcessRedoRecord","error":{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044044","reason":"SENZ2207|Data source code [] does not exist."}}}`,
				redoRecord:         badRedoRecord,
			},
			{
				name:               "nilRedoRecord",
				expectedErr:        szerror.ErrSz,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ProcessRedoRecord","error":{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044044","reason":"SENZ0002|Invalid Message"}}}`,
				redoRecord:         nilRedoRecord,
			},
			{
				name:               "withInfo_badRedoRecord",
				expectedErr:        szerror.ErrSzUnknownDataSource,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ProcessRedoRecord","error":{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044045","reason":"SENZ2207|Data source code [] does not exist."}}}`,
				flags:              senzing.SzWithInfo,
				redoRecord:         badRedoRecord,
			},
			{
				name:               "withInfo_nilRedoRecord",
				flags:              senzing.SzWithInfo,
				expectedErr:        szerror.ErrSz,
				expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ProcessRedoRecord","error":{"function":"szengine.(*Szengine).ProcessRedoRecord","error":{"id":"SZSDK60044045","reason":"SENZ0002|Invalid Message"}}}`,
				redoRecord:         nilRedoRecord,
			},
		}
	}

	result = append(result, addendum...)

	return result
}

func getTestCasesForReevaluateEntity() []TestMetadataForReevaluateEntity {
	result := []TestMetadataForReevaluateEntity{
		{
			name:     "badEntityID",
			entityID: badEntityID,
		},
		{
			name: "default",
		},
		{
			name:     "nilEntityID",
			entityID: nilEntityID,
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
		{
			name:     "withInfo_badEntityID",
			entityID: badEntityID,
			flags:    senzing.SzWithInfo,
		},
		{
			name:     "withInfo_nilEntityID",
			entityID: nilEntityID,
			flags:    senzing.SzWithInfo,
		},
	}

	return result
}

func getTestCasesForReevaluateRecord() []TestMetadataForReevaluateRecord {
	result := []TestMetadataForReevaluateRecord{
		{
			name:               "badDataSourceCode",
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ReevaluateRecord","error":{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044048","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
			dataSourceCode:     badDataSourceCode,
		},
		{
			name:     "badRecordID",
			recordID: badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ReevaluateRecord","error":{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044048","reason":"SENZ2207|Data source code [] does not exist."}}}`,
		},
		{
			name:               "nilRecordID",
			expectedErr:        szerror.ErrSzNotInitialized,
			expectedErrMessage: `{"function": "szengineserver.(*SzEngineServer).ReevaluateRecord", "error": {"function": "szengine.(*Szengine).ReevaluateRecord", "error":{"id":"SZSDK60044048","reason":"SENZ0053|RECORD_ID must be provided"}}}`,
			recordID:           nilRecordID,
		},
		{
			name:  "withInfo",
			flags: senzing.SzWithInfo,
		},
		{
			name:               "withInfo_badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ReevaluateRecord","error":{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044049","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
			flags:              senzing.SzWithInfo,
		},
		{
			name:     "withInfo_badRecordID",
			flags:    senzing.SzWithInfo,
			recordID: badRecordID,
		},
		{
			name:               "withInfo_nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).ReevaluateRecord","error":{"function":"szengine.(*Szengine).ReevaluateRecord","error":{"id":"SZSDK60044049","reason":"SENZ2207|Data source code [] does not exist."}}}`,
			flags:              senzing.SzWithInfo,
		},
		{
			name:               "withInfo_nilRecordID",
			expectedErr:        szerror.ErrSzNotInitialized,
			expectedErrMessage: `{"function": "szengineserver.(*SzEngineServer).ReevaluateRecord", "error": {"function": "szengine.(*Szengine).ReevaluateRecord", "error":{"id":"SZSDK60044049","reason":"SENZ0053|RECORD_ID must be provided"}}}`,
			flags:              senzing.SzWithInfo,
			recordID:           nilRecordID,
		},
	}

	return result
}

func getTestCasesForSearchByAttributes() []TestMetadataForSearchByAttributes {
	result := []TestMetadataForSearchByAttributes{
		{
			name:               "badAttributes",
			attributes:         badAttributes,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).SearchByAttributes","error":{"function":"szengine.(*Szengine).SearchByAttributes","error":{"id":"SZSDK60044053","reason":"SENZ0027|Invalid value for search-attributes"}}}`,
		},
		{
			name:               "badSearchProfile",
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).SearchByAttributes","error":{"function":"szengine.(*Szengine).SearchByAttributes","error":{"id":"SZSDK60044053","reason":"SENZ0088|Unknown search profile value '}{'"}}}`,
			searchProfile:      badSearchProfile,
		},
		{
			name: "default",
		},
		{
			name:               "nilAttributes",
			attributes:         nilAttributes,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).SearchByAttributes","error":{"function":"szengine.(*Szengine).SearchByAttributes","error":{"id":"SZSDK60044053","reason":"SENZ0027|Invalid value for search-attributes"}}}`,
		},
		{
			name:          "nilSearchProfile",
			searchProfile: nilSearchProfile,
		},
	}

	return result
}

func getTestCasesForWhyEntities() []TestMetadataForWhyEntities {
	result := []TestMetadataForWhyEntities{
		{
			name:               "badEnitity1",
			entityID1:          badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyEntities","error":{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
		},
		{
			name:               "badEnitity2",
			entityID2:          badEntityID,
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyEntities","error":{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
		},
		{
			name: "default",
		},
		{
			name:               "nilEnitity1",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyEntities","error":{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '0'"}}}`,
			entityID1:          nilEntityID,
		},
		{
			name:               "nilEnitity2",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyEntities","error":{"function":"szengine.(*Szengine).WhyEntities","error":{"id":"SZSDK60044056","reason":"SENZ0037|Unknown resolved entity value '0'"}}}`,
			entityID2:          nilEntityID,
		},
	}

	return result
}

func getTestCasesForWhyRecordInEntity() []TestMetadataForWhyRecordInEntity {
	result := []TestMetadataForWhyRecordInEntity{
		{
			name:               "badDataSourceCode",
			dataSourceCode:     badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecordInEntity","error":{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:               "badRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecordInEntity","error":{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
			recordID:           badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode",
			dataSourceCode:     nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecordInEntity","error":{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ2207|Data source code [] does not exist."}}}`,
		},
		{
			name:               "nilRecordID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecordInEntity","error":{"function":"szengine.(*Szengine).WhyRecordInEntity","error":{"id":"SZSDK60044058","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}}`,
			recordID:           nilRecordID,
		},
	}

	return result
}

func getTestCasesForWhyRecords() []TestMetadataForWhyRecords {
	result := []TestMetadataForWhyRecords{
		{
			name:               "badDataSourceCode1",
			dataSourceCode1:    badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecords","error":{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:               "badDataSourceCode2",
			dataSourceCode2:    badDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecords","error":{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [BADDATASOURCECODE] does not exist."}}}`,
		},
		{
			name:               "badRecordID1",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecords","error":{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
			recordID1:          badRecordID,
		},
		{
			name:               "badRecordID2",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecords","error":{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[BadRecordID]"}}}`,
			recordID2:          badRecordID,
		},
		{
			name: "default",
		},
		{
			name:               "nilDataSourceCode1",
			dataSourceCode1:    nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecords","error":{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [] does not exist."}}}`,
		},
		{
			name:               "nilDataSourceCode2",
			dataSourceCode2:    nilDataSourceCode,
			expectedErr:        szerror.ErrSzUnknownDataSource,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecords","error":{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ2207|Data source code [] does not exist."}}}`,
		},
		{
			name:               "nilRecordID1",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecords","error":{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}}`,
			recordID1:          nilRecordID,
		},
		{
			name:               "nilRecordID2",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhyRecords","error":{"function":"szengine.(*Szengine).WhyRecords","error":{"id":"SZSDK60044060","reason":"SENZ0033|Unknown record: dsrc[CUSTOMERS], record[]"}}}`,
			recordID2:          nilRecordID,
		},
	}

	return result
}

func getTestCasesForWhySearch() []TestMetadataForWhySearch {
	result := []TestMetadataForWhySearch{
		{
			name:               "badAttributes",
			attributes:         badAttributes,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhySearch","error":{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0027|Invalid value for search-attributes"}}}`,
		},
		{
			name:               "badEntityID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhySearch","error":{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0037|Unknown resolved entity value '-1'"}}}`,
			entityID:           badEntityID,
		},
		{
			name:               "badSearchProfile",
			expectedErr:        szerror.ErrSzBadInput,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhySearch","error":{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0088|Unknown search profile value '}{'"}}}`,
			searchProfile:      badSearchProfile,
		},
		{
			name: "default",
		},
		{
			name:               "nilAttributes",
			attributes:         nilAttributes,
			expectedErr:        szerror.ErrSz,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhySearch","error":{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0027|Invalid value for search-attributes"}}}`,
		},
		{
			name:               "nilEntityID",
			expectedErr:        szerror.ErrSzNotFound,
			expectedErrMessage: `{"function":"szengineserver.(*SzEngineServer).WhySearch","error":{"function":"szengine.(*Szengine).WhySearch","error":{"id":"SZSDK60044064","reason":"SENZ0037|Unknown resolved entity value '0'"}}}`,
			entityID:           nilEntityID,
		},
		{
			name:          "searchProfile",
			searchProfile: "SEARCH",
		},
	}

	return result
}

// ----------------------------------------------------------------------------
// Test data helpers
// ----------------------------------------------------------------------------

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

func badDataSourcesFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]

	return `{"RECORDS": [{"DATA_SOURCE": "` +
		badDataSourceCode +
		`", "RECORD_ID": "` +
		record1.ID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}]}`
}

func badEntityIDsFunc() string {
	return `{"ENTITIES": [{"ENTITY_ID": ` +
		strconv.FormatInt(badEntityID, baseTen) +
		`}]}`
}

func badRecordKeysFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]

	return `{"RECORDS": [{"DATA_SOURCE": "` +
		record1.DataSource +
		`", "RECORD_ID": "` +
		badRecordID +
		`"}, {"DATA_SOURCE": "` +
		record2.DataSource +
		`", "RECORD_ID": "` +
		record2.ID +
		`"}]}`
}

func entityIDsJSON(entityIDs ...int64) string {
	var resultBuilder strings.Builder

	resultBuilder.WriteString(`{"ENTITIES": [`)

	for _, entityID := range entityIDs {
		resultBuilder.WriteString(`{"ENTITY_ID":` + strconv.FormatInt(entityID, baseTen) + `},`)
	}

	result := resultBuilder.String()
	result = result[:len(result)-1] // Remove final comma.
	result += `]}`

	return result
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

func recordKeysJSON(r1DS, r1ID, r2DS, r2ID, r3DS, r3ID string) string {
	return `{"RECORDS": [{"DATA_SOURCE": "` +
		r1DS +
		`", "RECORD_ID": "` +
		r1ID +
		`"}, {"DATA_SOURCE": "` +
		r2DS +
		`", "RECORD_ID": "` +
		r2ID +
		`"}, {"DATA_SOURCE": "` +
		r3DS +
		`", "RECORD_ID": "` +
		r3ID +
		`"}]}`
}

func recordKeysFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]

	return recordKeysJSON(
		record1.DataSource,
		record1.ID,
		record2.DataSource,
		record2.ID,
		record3.DataSource,
		record3.ID,
	)
}

func recordKeysBadDataSourceCodeFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]

	return recordKeysJSON(
		badDataSourceCode,
		record1.ID,
		record2.DataSource,
		record2.ID,
		record3.DataSource,
		record3.ID,
	)
}

func recordKeysFuncBadRecordIDFunc() string {
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]

	return recordKeysJSON(
		record1.DataSource,
		badRecordID,
		record2.DataSource,
		record2.ID,
		record3.DataSource,
		record3.ID,
	)
}

// Return first non-zero length candidate.  Last candidate is default.
func xString(candidates ...string) string {
	var result string
	for _, result = range candidates {
		if result == nilSemaphoreString {
			return ""
		}

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
		if result == nilSemaphoreInt64 {
			return 0
		}

		if result != 0 {
			return result
		}
	}

	return result
}

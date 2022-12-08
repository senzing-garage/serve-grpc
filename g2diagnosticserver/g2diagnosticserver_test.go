package g2diagnosticserver

import (
	"context"
	"testing"

	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	pb "github.com/senzing/go-servegrpc/g2diagnostic"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Internal functions - names begin with lowercase letter
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context) (G2DiagnosticServer, error) {
	var err error = nil
	g2diagnosticServer := G2DiagnosticServer{}

	moduleName := "Test module name"
	verboseLogging := 0 // 0 for no Senzing logging; 1 for logging
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		return g2diagnosticServer, jsonErr
	}

	request := pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err = g2diagnosticServer.Init(ctx, &request)
	return g2diagnosticServer, err
}

func testError(test *testing.T, ctx context.Context, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestGetObject(test *testing.T) {
	ctx := context.TODO()
	_, err := getTestObject(ctx)
	testError(test, ctx, err)
}

// ----------------------------------------------------------------------------
// Test interface functions - names begin with "Test"
// ----------------------------------------------------------------------------

func TestCheckDBPerf(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.CheckDBPerfRequest{
		SecondsToRun: int32(1),
	}
	actual, err := g2diagnosticServer.CheckDBPerf(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestClearLastException(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.ClearLastExceptionRequest{}
	actual, err := g2diagnosticServer.ClearLastException(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestDestroy(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.DestroyRequest{}
	actual, err := g2diagnosticServer.Destroy(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

//func TestEntityListBySize(test *testing.T) {
//	ctx := context.TODO()
//	g2diagnosticServer, _ := getTestObject(ctx)
//	request := pb.GetEntityListBySizeRequest{
//		EntitySize: int32(10),
//	}
//	actual, err := g2diagnosticServer.GetEntityListBySize(ctx, &request)
//	testError(test, ctx, err)
//	test.Log("Actual:", actual)
//
//	entityListBySizeHandle := actual.Result
//	request2 := pb.FetchNextEntityBySizeRequest{
//		EntityListBySizeHandle: entityListBySizeHandle,
//	}
//	actual2, err2 := g2diagnosticServer.FetchNextEntityBySize(ctx, &request2)
//	testError(test, ctx, err2)
//	test.Log("Actual:", actual2)
//
//	request3 := pb.CloseEntityListBySizeRequest{
//		EntityListBySizeHandle: entityListBySizeHandle,
//	}
//	actual3, err3 := g2diagnosticServer.CloseEntityListBySize(ctx, &request3)
//	testError(test, ctx, err3)
//	test.Log("Actual:", actual3)
//}

func TestFindEntitiesByFeatureIDs(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.FindEntitiesByFeatureIDsRequest{
		Features: "{\"ENTITY_ID\":1,\"LIB_FEAT_IDS\":[1,3,4]}",
	}
	actual, err := g2diagnosticServer.FindEntitiesByFeatureIDs(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetAvailableMemoryy(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetAvailableMemoryRequest{}
	actual, err := g2diagnosticServer.GetAvailableMemory(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetDataSourceCounts(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetDataSourceCountsRequest{}
	actual, err := g2diagnosticServer.GetDataSourceCounts(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetDBInfo(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetDBInfoRequest{}
	actual, err := g2diagnosticServer.GetDBInfo(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetEntityDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetEntityDetailsRequest{
		EntityID:                int64(1),
		IncludeInternalFeatures: 1,
	}
	actual, err := g2diagnosticServer.GetEntityDetails(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetEntityResume(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetEntityResumeRequest{
		EntityID: int64(1),
	}
	actual, err := g2diagnosticServer.GetEntityResume(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetEntitySizeBreakdown(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetEntitySizeBreakdownRequest{
		MinimumEntitySize:       int32(1),
		IncludeInternalFeatures: int32(1),
	}
	actual, err := g2diagnosticServer.GetEntitySizeBreakdown(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetFeature(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetFeatureRequest{
		LibFeatID: int64(1),
	}
	actual, err := g2diagnosticServer.GetFeature(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetGenericFeatures(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetGenericFeaturesRequest{
		FeatureType:           "PHONE",
		MaximumEstimatedCount: 10,
	}
	actual, err := g2diagnosticServer.GetGenericFeatures(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetLastException(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetLastExceptionRequest{}
	actual, err := g2diagnosticServer.GetLastException(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetLastExceptionCode(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetLastExceptionCodeRequest{}
	actual, err := g2diagnosticServer.GetLastExceptionCode(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetLogicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetLogicalCoresRequest{}
	actual, err := g2diagnosticServer.GetLogicalCores(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetMappingStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: 1,
	}
	actual, err := g2diagnosticServer.GetMappingStatistics(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetPhysicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetPhysicalCoresRequest{}
	actual, err := g2diagnosticServer.GetPhysicalCores(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetRelationshipDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetRelationshipDetailsRequest{
		RelationshipID:          int64(1),
		IncludeInternalFeatures: 1,
	}
	actual, err := g2diagnosticServer.GetRelationshipDetails(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetResolutionStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetResolutionStatisticsRequest{}
	actual, err := g2diagnosticServer.GetResolutionStatistics(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestGetTotalSystemMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.GetTotalSystemMemoryRequest{}
	actual, err := g2diagnosticServer.GetTotalSystemMemory(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestInit(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		assert.FailNow(test, jsonErr.Error())
	}
	request := pb.InitRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		VerboseLogging: int32(0),
	}
	actual, err := g2diagnosticServer.Init(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestInitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		assert.FailNow(test, jsonErr.Error())
	}
	request := pb.InitWithConfigIDRequest{
		ModuleName:     "Test module name",
		IniParams:      iniParams,
		InitConfigID:   int64(1),
		VerboseLogging: int32(0),
	}
	actual, err := g2diagnosticServer.InitWithConfigID(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

func TestReinit(test *testing.T) {
	ctx := context.TODO()
	g2diagnosticServer, _ := getTestObject(ctx)
	request := pb.ReinitRequest{
		InitConfigID: int64(1),
	}
	actual, err := g2diagnosticServer.Reinit(ctx, &request)
	testError(test, ctx, err)
	test.Log("Actual:", actual)
}

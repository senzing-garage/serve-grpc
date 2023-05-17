package g2diagnosticserver

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type G2DiagnosticServer struct {
	pb.UnimplementedG2DiagnosticServer
	isTrace bool
	logger  logging.LoggingInterface
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2diagnostic package found messages having the format "senzing-6999xxxx".
const ComponentId = 6013

// Log message prefix.
const Prefix = "serve-grpc.g2diagnosticserver."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2diagnostic package.
var IdMessages = map[int]string{
	1:    "Enter " + Prefix + "CheckDBPerf(%+v).",
	2:    "Exit  " + Prefix + "CheckDBPerf(%+v) returned (%s, %v).",
	3:    "Enter " + Prefix + "RegisterObserver(%s).",
	4:    "Exit  " + Prefix + "RegisterObserver(%s) returned (%s, %v).",
	5:    "Enter " + Prefix + "CloseEntityListBySize(%+v).",
	6:    "Exit  " + Prefix + "CloseEntityListBySize(%+v) returned (%v).",
	7:    "Enter " + Prefix + "Destroy(%+v).",
	8:    "Exit  " + Prefix + "Destroy(%+v) returned (%v).",
	9:    "Enter " + Prefix + "FetchNextEntityBySize(%+v).",
	10:   "Exit  " + Prefix + "FetchNextEntityBySize(%+v) returned (%s, %v).",
	11:   "Enter " + Prefix + "FindEntitiesByFeatureIDs(%+v).",
	12:   "Exit  " + Prefix + "FindEntitiesByFeatureIDs(%+v) returned (%s, %v).",
	13:   "Enter " + Prefix + "GetAvailableMemory(%+v).",
	14:   "Exit  " + Prefix + "GetAvailableMemory(%+v) returned (%d, %v).",
	15:   "Enter " + Prefix + "GetDataSourceCounts(%+v).",
	16:   "Exit  " + Prefix + "GetDataSourceCounts(%+v) returned (%s, %v).",
	17:   "Enter " + Prefix + "GetDBInfo(%+v).",
	18:   "Exit  " + Prefix + "GetDBInfo(%+v) returned (%s, %v).",
	19:   "Enter " + Prefix + "GetEntityDetails(%+v).",
	20:   "Exit  " + Prefix + "GetEntityDetails(%+v) returned (%s, %v).",
	21:   "Enter " + Prefix + "GetEntityListBySize(%+v).",
	22:   "Exit  " + Prefix + "GetEntityListBySize(%+v) returned (%v, %v).",
	23:   "Enter " + Prefix + "GetEntityResume(%+v).",
	24:   "Exit  " + Prefix + "GetEntityResume(%+v) returned (%s, %v).",
	25:   "Enter " + Prefix + "GetEntitySizeBreakdown(%+v).",
	26:   "Exit  " + Prefix + "GetEntitySizeBreakdown(%+v) returned (%s, %v).",
	27:   "Enter " + Prefix + "GetFeature(%+v).",
	28:   "Exit  " + Prefix + "GetFeature(%+v) returned (%s, %v).",
	29:   "Enter " + Prefix + "GetGenericFeatures(%+v).",
	30:   "Exit  " + Prefix + "GetGenericFeatures(%+v) returned (%s, %v).",
	31:   "Enter " + Prefix + "UnregisterObserver(%s).",
	32:   "Exit  " + Prefix + "UnregisterObserver(%s) returned (%s, %v).",
	35:   "Enter " + Prefix + "GetLogicalCores(%+v).",
	36:   "Exit  " + Prefix + "GetLogicalCores(%+v) returned (%d, %v).",
	37:   "Enter " + Prefix + "GetMappingStatistics(%+v).",
	38:   "Exit  " + Prefix + "GetMappingStatistics(%+v) returned (%s, %v).",
	39:   "Enter " + Prefix + "GetPhysicalCores(%+v).",
	40:   "Exit  " + Prefix + "GetPhysicalCores(%+v) returned (%d, %v).",
	41:   "Enter " + Prefix + "GetRelationshipDetails(%+v).",
	42:   "Exit  " + Prefix + "GetRelationshipDetails(%+v) returned (%s, %v).",
	43:   "Enter " + Prefix + "GetResolutionStatistics(%+v).",
	44:   "Exit  " + Prefix + "GetResolutionStatistics(%+v) returned (%s, %v).",
	45:   "Enter " + Prefix + "GetTotalSystemMemory(%+v).",
	46:   "Exit  " + Prefix + "GetTotalSystemMemory(%+v) returned (%d, %v).",
	47:   "Enter " + Prefix + "Init(%+v).",
	48:   "Exit  " + Prefix + "Init(%+v) returned (%v).",
	49:   "Enter " + Prefix + "InitWithConfigID(%+v).",
	50:   "Exit  " + Prefix + "InitWithConfigID(%+v) returned (%v).",
	51:   "Enter " + Prefix + "Reinit(%+v).",
	52:   "Exit  " + Prefix + "Reinit(%+v) returned (%v).",
	53:   "Enter " + Prefix + "SetLogLevel(%s).",
	54:   "Exit  " + Prefix + "SetLogLevel(%s) returned (%v).",
	55:   "Enter " + Prefix + "GetObserverOrigin().",
	56:   "Exit  " + Prefix + "GetObserverOrigin() returned (%v).",
	57:   "Enter " + Prefix + "SetObserverOrigin(%s).",
	58:   "Exit  " + Prefix + "SetObserverOrigin(%s) returned (%v).",
	4001: Prefix + "Destroy() not supported in gRPC",
	4002: Prefix + "Init() not supported in gRPC",
	4003: Prefix + "InitWithConfigID() not supported in gRPC",
	5901: "During test setup, call to messagelogger.NewSenzingApiLogger() failed.",
	5902: "During test setup, call to g2eg2engineconfigurationjson.BuildSimpleSystemConfigurationJson() failed.",
	5903: "During test setup, call to g2engine.Init() failed.",
	5904: "During test setup, call to g2engine.PurgeRepository() failed.",
	5905: "During test setup, call to g2engine.Destroy() failed.",
	5906: "During test setup, call to g2config.Init() failed.",
	5907: "During test setup, call to g2config.Create() failed.",
	5908: "During test setup, call to g2config.AddDataSource() failed.",
	5909: "During test setup, call to g2config.Save() failed.",
	5910: "During test setup, call to g2config.Close() failed.",
	5911: "During test setup, call to g2config.Destroy() failed.",
	5912: "During test setup, call to g2configmgr.Init() failed.",
	5913: "During test setup, call to g2configmgr.AddConfig() failed.",
	5914: "During test setup, call to g2configmgr.SetDefaultConfigID() failed.",
	5915: "During test setup, call to g2configmgr.Destroy() failed.",
	5916: "During test setup, call to g2engine.Init() failed.",
	5917: "During test setup, call to g2engine.AddRecord() failed.",
	5918: "During test setup, call to g2engine.Destroy() failed.",
	5920: "During test setup, call to setupSenzingConfig() failed.",
	5921: "During test setup, call to setupPurgeRepository() failed.",
	5922: "During test setup, call to setupAddRecords() failed.",
	5931: "During test setup, call to g2engine.Init() failed.",
	5932: "During test setup, call to g2engine.PurgeRepository() failed.",
	5933: "During test setup, call to g2engine.Destroy() failed.",
}

// Status strings for specific g2diagnostic messages.
var IdStatuses = map[int]string{}

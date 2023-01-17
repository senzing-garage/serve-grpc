package g2diagnosticserver

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing/go-logging/messagelogger"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type G2DiagnosticServer struct {
	pb.UnimplementedG2DiagnosticServer
	isTrace bool
	logger  messagelogger.MessageLoggerInterface
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2diagnostic package found messages having the format "senzing-6999xxxx".
const ProductId = 6013

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2diagnostic package.
var IdMessages = map[int]string{
	1:    "Enter CheckDBPerf(%+v).",
	2:    "Exit  CheckDBPerf(%+v) returned (%s, %v).",
	5:    "Enter CloseEntityListBySize(%+v).",
	6:    "Exit  CloseEntityListBySize(%+v) returned (%v).",
	7:    "Enter Destroy(%+v).",
	8:    "Exit  Destroy(%+v) returned (%v).",
	9:    "Enter FetchNextEntityBySize(%+v).",
	10:   "Exit  FetchNextEntityBySize(%+v) returned (%s, %v).",
	11:   "Enter FindEntitiesByFeatureIDs(%+v).",
	12:   "Exit  FindEntitiesByFeatureIDs(%+v) returned (%s, %v).",
	13:   "Enter GetAvailableMemory(%+v).",
	14:   "Exit  GetAvailableMemory(%+v) returned (%d, %v).",
	15:   "Enter GetDataSourceCounts(%+v).",
	16:   "Exit  GetDataSourceCounts(%+v) returned (%s, %v).",
	17:   "Enter GetDBInfo(%+v).",
	18:   "Exit  GetDBInfo(%+v) returned (%s, %v).",
	19:   "Enter GetEntityDetails(%+v).",
	20:   "Exit  GetEntityDetails(%+v) returned (%s, %v).",
	21:   "Enter GetEntityListBySize(%+v).",
	22:   "Exit  GetEntityListBySize(%+v) returned (%v, %v).",
	23:   "Enter GetEntityResume(%+v).",
	24:   "Exit  GetEntityResume(%+v) returned (%s, %v).",
	25:   "Enter GetEntitySizeBreakdown(%+v).",
	26:   "Exit  GetEntitySizeBreakdown(%+v) returned (%s, %v).",
	27:   "Enter GetFeature(%+v).",
	28:   "Exit  GetFeature(%+v) returned (%s, %v).",
	29:   "Enter GetGenericFeatures(%+v).",
	30:   "Exit  GetGenericFeatures(%+v) returned (%s, %v).",
	35:   "Enter GetLogicalCores(%+v).",
	36:   "Exit  GetLogicalCores(%+v) returned (%d, %v).",
	37:   "Enter GetMappingStatistics(%+v).",
	38:   "Exit  GetMappingStatistics(%+v) returned (%s, %v).",
	39:   "Enter GetPhysicalCores(%+v).",
	40:   "Exit  GetPhysicalCores(%+v) returned (%d, %v).",
	41:   "Enter GetRelationshipDetails(%+v).",
	42:   "Exit  GetRelationshipDetails(%+v) returned (%s, %v).",
	43:   "Enter GetResolutionStatistics(%+v).",
	44:   "Exit  GetResolutionStatistics(%+v) returned (%s, %v).",
	45:   "Enter GetTotalSystemMemory(%+v).",
	46:   "Exit  GetTotalSystemMemory(%+v) returned (%d, %v).",
	47:   "Enter Init(%+v).",
	48:   "Exit  Init(%+v) returned (%v).",
	49:   "Enter InitWithConfigID(%+v).",
	50:   "Exit  InitWithConfigID(%+v) returned (%v).",
	51:   "Enter Reinit(%+v).",
	52:   "Exit  Reinit(%+v) returned (%v).",
	53:   "Enter SetLogLevel(%+v).",
	54:   "Exit  SetLogLevel(%+v) returned (%v).",
	4001: "Destroy() not supported in gRPC",
	4002: "Init() not supported in gRPC",
	4003: "InitWithConfigID() not supported in gRPC",
	4004: "PurgeRepository() not supported in gRPC",
	5901: "setup() call to messagelogger.NewSenzingApiLogger() failed.",
	5902: "setup() call to g2eg2engineconfigurationjson.BuildSimpleSystemConfigurationJson() failed.",
	5903: "setup() call to g2engine.Init() failed.",
	5904: "setup() call to g2engine.PurgeRepository() failed.",
	5905: "setup() call to g2engine.Destroy() failed.",
	5906: "setup() call to g2config.Init() failed.",
	5907: "setup() call to g2config.Create() failed.",
	5908: "setup() call to g2config.AddDataSource() failed.",
	5909: "setup() call to g2config.Save() failed.",
	5910: "setup() call to g2config.Close() failed.",
	5911: "setup() call to g2config.Destroy() failed.",
	5912: "setup() call to g2configmgr.Init() failed.",
	5913: "setup() call to g2configmgr.AddConfig() failed.",
	5914: "setup() call to g2configmgr.SetDefaultConfigID() failed.",
	5915: "setup() call to g2configmgr.Destroy() failed.",
	5916: "setup() call to g2engine.Init() failed.",
	5917: "setup() call to g2engine.AddRecord() failed.",
	5918: "setup() call to g2engine.Destroy() failed.",
	5931: "setup() call to g2engine.Init() failed.",
	5932: "setup() call to g2engine.PurgeRepository() failed.",
	5933: "setup() call to g2engine.Destroy() failed.",
}

// Status strings for specific g2diagnostic messages.
var IdStatuses = map[int]string{}

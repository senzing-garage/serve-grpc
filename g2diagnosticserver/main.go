package g2diagnosticserver

import (
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/go-servegrpc/protobuf/g2diagnostic"
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
	1:    "Enter CheckDBPerf(%d).",
	2:    "Exit  CheckDBPerf(%d) returned (%s, %v).",
	3:    "Enter ClearLastException().",
	4:    "Exit  ClearLastException() returned (%v).",
	5:    "Enter CloseEntityListBySize().",
	6:    "Exit  CloseEntityListBySize() returned (%v).",
	7:    "Enter Destroy().",
	8:    "Exit  Destroy() returned (%v).",
	9:    "Enter FetchNextEntityBySize().",
	10:   "Exit  FetchNextEntityBySize() returned (%s, %v).",
	11:   "Enter FindEntitiesByFeatureIDs(%s).",
	12:   "Exit  FindEntitiesByFeatureIDs(%s) returned (%s, %v).",
	13:   "Enter GetAvailableMemory().",
	14:   "Exit  GetAvailableMemory() returned (%d, %v).",
	15:   "Enter GetDataSourceCounts().",
	16:   "Exit  GetDataSourceCounts() returned (%s, %v).",
	17:   "Enter GetDBInfo().",
	18:   "Exit  GetDBInfo()  returned (%s, %v).",
	19:   "Enter GetEntityDetails(%d, %d).",
	20:   "Exit  GetEntityDetails(%d, %d) returned (%s, %v).",
	21:   "Enter GetEntityListBySize(%d).",
	22:   "Exit  GetEntityListBySize(%d) returned (%v, %v).",
	23:   "Enter GetEntityResume(%d).",
	24:   "Exit  GetEntityResume(%d) returned (%s, %v).",
	25:   "Enter GetEntitySizeBreakdown(%d, %d).",
	26:   "Exit  GetEntitySizeBreakdown(%d, %d) returned (%s, %v).",
	27:   "Enter GetFeature(%d).",
	28:   "Exit  GetFeature(%d) returned (%s, %v).",
	29:   "Enter GetGenericFeatures(%s, %d).",
	30:   "Exit  GetGenericFeatures(%s, %d) returned (%s, %v).",
	31:   "Enter GetLastException().",
	32:   "Exit  GetLastException() returned (%s, %v).",
	33:   "Enter GetLastExceptionCode().",
	34:   "Exit  GetLastExceptionCode() returned (%d, %v).",
	35:   "Enter GetLogicalCores().",
	36:   "Exit  GetLogicalCores() returned (%d, %v).",
	37:   "Enter GetMappingStatistics(%d).",
	38:   "Exit  GetMappingStatistics(%d) returned (%s, %v).",
	39:   "Enter GetPhysicalCores().",
	40:   "Exit  GetPhysicalCores() returned (%d, %v).",
	41:   "Enter GetRelationshipDetails(%d, %d).",
	42:   "Exit  GetRelationshipDetails(%d, %d) returned (%s, %v).",
	43:   "Enter GetResolutionStatistics().",
	44:   "Exit  GetResolutionStatistics() returned (%s, %v).",
	45:   "Enter GetTotalSystemMemory().",
	46:   "Exit  GetTotalSystemMemory() returned (%d, %v).",
	47:   "Enter Init(%s, %s, %d).",
	48:   "Exit  Init(%s, %s, %d) returned (%v).",
	49:   "Enter InitWithConfigID(%s, %s, %d, %d).",
	50:   "Exit  InitWithConfigID(%s, %s, %d, %d) returned (%v).",
	51:   "Enter Reinit(%d).",
	52:   "Exit  Reinit(%d) returned (%v).",
	53:   "Enter SetLogLevel(%v).",
	54:   "Exit  SetLogLevel(%v) returned (%v).",
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

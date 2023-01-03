package g2engineserver

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2engine"
	"github.com/senzing/go-logging/messagelogger"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type G2EngineServer struct {
	pb.UnimplementedG2EngineServer
	isTrace bool
	logger  messagelogger.MessageLoggerInterface
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2engineserver package found messages having the format "senzing-6999xxxx".
const ProductId = 6014

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2engineserver package.
var IdMessages = map[int]string{
	1:    "Enter AddRecord(%+v).",
	2:    "Exit  AddRecord(%+v) returned (%v).",
	3:    "Enter AddRecordWithInfo(%+v).",
	4:    "Exit  AddRecordWithInfo(%+v) returned (%s, %v).",
	5:    "Enter AddRecordWithInfoWithReturnedRecordID(%+v).",
	6:    "Exit  AddRecordWithInfoWithReturnedRecordID(%+v) returned (%s, %s, %v).",
	7:    "Enter AddRecordWithReturnedRecordID(%+v).",
	8:    "Exit  AddRecordWithReturnedRecordID(%+v) returned (%s, %v).",
	9:    "Enter CheckRecord(%+v).",
	10:   "Exit  CheckRecord(%+v) returned (%s, %v).",
	13:   "Enter CloseExport(%+v).",
	14:   "Exit  CloseExport(%+v) returned (%v).",
	15:   "Enter CountRedoRecords(%+v).",
	16:   "Exit  CountRedoRecords(%+v) returned (%d, %v).",
	17:   "Enter DeleteRecord(%+v).",
	18:   "Exit  DeleteRecord(%+v) returned (%v).",
	19:   "Enter DeleteRecordWithInfo(%+v).",
	20:   "Exit  DeleteRecordWithInfo(%+v) returned (%s, %v).",
	21:   "Enter Destroy(%+v).",
	22:   "Exit  Destroy(%+v) returned (%v).",
	23:   "Enter ExportConfigAndConfigID(%+v).",
	24:   "Exit  ExportConfigAndConfigID(%+v) returned (%s, %d, %v).",
	25:   "Enter ExportConfig(%+v).",
	26:   "Exit  ExportConfig(%+v) returned (%s, %v).",
	27:   "Enter ExportCSVEntityReport(%+v).",
	28:   "Exit  ExportCSVEntityReport(%+v) returned (%v, %v).",
	29:   "Enter ExportJSONEntityReport(%+v).",
	30:   "Exit  ExportJSONEntityReport(%+v) returned (%v, %v).",
	31:   "Enter FetchNext(%+v).",
	32:   "Exit  FetchNext(%+v) returned (%s, %v).",
	33:   "Enter FindInterestingEntitiesByEntityID(%+v).",
	34:   "Exit  FindInterestingEntitiesByEntityID(%+v) returned (%s, %v).",
	35:   "Enter FindInterestingEntitiesByRecordID(%+v).",
	36:   "Exit  FindInterestingEntitiesByRecordID(%+v) returned (%s, %v).",
	37:   "Enter FindNetworkByEntityID(%+v).",
	38:   "Exit  FindNetworkByEntityID(%+v) returned (%s, %v).",
	39:   "Enter FindNetworkByEntityID_V2(%+v).",
	40:   "Exit  FindNetworkByEntityID_V2(%+v) returned (%s, %v).",
	41:   "Enter FindNetworkByRecordID(%+v).",
	42:   "Exit  FindNetworkByRecordID(%+v) returned (%s, %v).",
	43:   "Enter FindNetworkByRecordID_V2(%+v).",
	44:   "Exit  FindNetworkByRecordID_V2(%+v) returned (%s, %v).",
	45:   "Enter FindPathByEntityID(%+v).",
	46:   "Exit  FindPathByEntityID(%+v) returned (%s, %v).",
	47:   "Enter FindPathByEntityID_V2(%+v).",
	48:   "Exit  FindPathByEntityID_V2(%+v) returned (%s, %v).",
	49:   "Enter FindPathByRecordID(%+v).",
	50:   "Exit  FindPathByRecordID(%+v) returned (%s, %v).",
	51:   "Enter FindPathByRecordID_V2(%+v).",
	52:   "Exit  FindPathByRecordID_V2(%+v) returned (%s, %v).",
	53:   "Enter FindPathExcludingByEntityID(%+v).",
	54:   "Exit  FindPathExcludingByEntityID(%+v) returned (%s, %v).",
	55:   "Enter FindPathExcludingByEntityID_V2(%+v).",
	56:   "Exit  FindPathExcludingByEntityID_V2(%+v) returned (%s, %v).",
	57:   "Enter FindPathExcludingByRecordID(%+v).",
	58:   "Exit  FindPathExcludingByRecordID(%+v) returned (%s, %v).",
	59:   "Enter FindPathExcludingByRecordID_V2(%+v).",
	60:   "Exit  FindPathExcludingByRecordID_V2(%+v) returned (%v).",
	61:   "Enter FindPathIncludingSourceByEntityID(%+v).",
	62:   "Exit  FindPathIncludingSourceByEntityID(%+v) returned (%s, %v).",
	63:   "Enter FindPathIncludingSourceByEntityID_V2(%+v).",
	64:   "Exit  FindPathIncludingSourceByEntityID_V2(%+v) returned (%s, %v).",
	65:   "Enter FindPathIncludingSourceByRecordID(%+v).",
	66:   "Exit  FindPathIncludingSourceByRecordID(%+v) returned (%s, %v).",
	67:   "Enter FindPathIncludingSourceByRecordID_V2(%+v).",
	68:   "Exit  FindPathIncludingSourceByRecordID_V2(%+v) returned (%s, %v).",
	69:   "Enter GetActiveConfigID(%+v).",
	70:   "Exit  GetActiveConfigID(%+v) returned (%d, %v).",
	71:   "Enter GetEntityByEntityID(%+v).",
	72:   "Exit  GetEntityByEntityID(%+v) returned (%s, %v).",
	73:   "Enter GetEntityByEntityID_V2(%+v).",
	74:   "Exit  GetEntityByEntityID_V2(%+v) returned (%s, %v).",
	75:   "Enter GetEntityByRecordID(%+v).",
	76:   "Exit  GetEntityByRecordID(%+v) returned (%s, %v).",
	77:   "Enter GetEntityByRecordID_V2(%+v).",
	78:   "Exit  GetEntityByRecordID_V2(%+v) returned (%s, %v).",
	83:   "Enter GetRecord(%+v).",
	84:   "Exit  GetRecord(%+v) returned (%s, %v).",
	85:   "Enter GetRecord_V2(%+v).",
	86:   "Exit  GetRecord_V2(%+v) returned (%s, %v).",
	87:   "Enter GetRedoRecord(%+v).",
	88:   "Exit  GetRedoRecord(%+v) returned (%s, %v).",
	89:   "Enter GetRepositoryLastModifiedTime(%+v).",
	90:   "Exit  GetRepositoryLastModifiedTime(%+v) returned (%d, %v).",
	91:   "Enter GetVirtualEntityByRecordID(%+v).",
	92:   "Exit  GetVirtualEntityByRecordID(%+v) returned (%s, %v).",
	93:   "Enter GetVirtualEntityByRecordID_V2(%+v).",
	94:   "Exit  GetVirtualEntityByRecordID_V2(%+v) returned (%s, %v).",
	95:   "Enter HowEntityByEntityID(%+v).",
	96:   "Exit  HowEntityByEntityID(%+v) returned (%s, %v).",
	97:   "Enter HowEntityByEntityID_V2(%+v).",
	98:   "Exit  HowEntityByEntityID_V2(%+v) returned (%s, %v).",
	99:   "Enter Init(%+v).",
	100:  "Exit  Init(%+v) returned (%v).",
	101:  "Enter InitWithConfigID(%+v).",
	102:  "Exit  InitWithConfigID(%+v) returned (%v).",
	103:  "Enter PrimeEngine(%+v).",
	104:  "Exit  PrimeEngine(%+v) returned (%v).",
	105:  "Enter Process(%+v).",
	106:  "Exit  Process(%+v) returned (%v).",
	107:  "Enter ProcessRedoRecord(%+v).",
	108:  "Exit  ProcessRedoRecord(%+v) returned (%s, %v).",
	109:  "Enter ProcessRedoRecordWithInfo(%+v).",
	110:  "Exit  ProcessRedoRecordWithInfo(%+v) returned (%s, %s, %v).",
	111:  "Enter ProcessWithInfo(%+v).",
	112:  "Exit  ProcessWithInfo(%+v) returned (%s, %v).",
	113:  "Enter ProcessWithResponse(%+v).",
	114:  "Exit  ProcessWithResponse(%+v) returned (%s, %v).",
	115:  "Enter ProcessWithResponseResize(%+v).",
	116:  "Exit  ProcessWithResponseResize(%+v) returned (%s, %v).",
	117:  "Enter PurgeRepository(%+v).",
	118:  "Exit  PurgeRepository(%+v) returned (%v).",
	119:  "Enter ReevaluateEntity(%+v).",
	120:  "Exit  ReevaluateEntity(%+v) returned (%v).",
	121:  "Enter ReevaluateEntityWithInfo(%+v).",
	122:  "Exit  ReevaluateEntityWithInfo(%+v) returned (%s, %v).",
	123:  "Enter ReevaluateRecord(%+v).",
	124:  "Exit  ReevaluateRecord(%+v) returned (%v).",
	125:  "Enter ReevaluateRecordWithInfo(%+v).",
	126:  "Exit  ReevaluateRecordWithInfo(%+v) returned (%s, %v).",
	127:  "Enter Reinit(%+v).",
	128:  "Exit  Reinit(%+v) returned (%v).",
	129:  "Enter ReplaceRecord(%+v).",
	130:  "Exit  ReplaceRecord(%+v) returned (%v).",
	131:  "Enter ReplaceRecordWithInfo(%+v).",
	132:  "Exit  ReplaceRecordWithInfo(%+v) returned (%s, %v).",
	133:  "Enter SearchByAttributes(%+v).",
	134:  "Exit  SearchByAttributes(%+v) returned (%s, %v).",
	135:  "Enter SearchByAttributes_V2(%+v).",
	136:  "Exit  SearchByAttributes_V2(%+v) returned (%s, %v).",
	137:  "Enter SetLogLevel(%+v).",
	138:  "Exit  SetLogLevel(%+v) returned (%v).",
	139:  "Enter Stats(%+v).",
	140:  "Exit  Stats(%+v) returned (%s, %v).",
	141:  "Enter WhyEntities(%+v).",
	142:  "Exit  WhyEntities(%+v) returned (%s, %v).",
	143:  "Enter WhyEntities_V2(%+v).",
	144:  "Exit  WhyEntities_V2(%+v) returned (%s, %v).",
	145:  "Enter WhyEntityByEntityID(%+v).",
	146:  "Exit  WhyEntityByEntityID(%+v) returned (%s, %v).",
	147:  "Enter WhyEntityByEntityID_V2(%+v).",
	148:  "Exit  WhyEntityByEntityID_V2(%+v) returned (%s, %v).",
	149:  "Enter WhyEntityByRecordID(%+v).",
	150:  "Exit  WhyEntityByRecordID(%+v) returned (%s, %v).",
	151:  "Enter WhyEntityByRecordID_V2(%+v).",
	152:  "Exit  WhyEntityByRecordID_V2(%+v) returned (%s, %v).",
	153:  "Enter WhyRecords(%+v).",
	154:  "Exit  WhyRecords(%+v) returned (%s, %v).",
	155:  "Enter WhyRecords_V2(%+v).",
	156:  "Exit  WhyRecords_V2(%+v) returned (%s, %v).",
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

// Status strings for specific g2engineserver messages.
var IdStatuses = map[int]string{}

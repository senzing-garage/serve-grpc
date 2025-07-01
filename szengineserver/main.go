package szengineserver

import (
	"errors"

	"github.com/senzing-garage/go-logging/logging"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type SzEngineServer struct {
	szpb.UnimplementedSzEngineServer
	isTrace bool
	logger  logging.Logging
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the szengineserver package found messages having the format "senzing-6999xxxx".
const ComponentID = 6014

// Log message prefix.
const Prefix = "serve-grpc.szengineserver."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the szengineserver package.
var IDMessages = map[int]string{
	1:    "Enter " + Prefix + "AddRecord(%+v).",
	2:    "Exit  " + Prefix + "AddRecord(%+v) returned (%v).",
	3:    "Enter " + Prefix + "AddRecordWithInfo(%+v).",
	4:    "Exit  " + Prefix + "AddRecordWithInfo(%+v) returned (%s, %v).",
	11:   "Enter " + Prefix + "RegisterObserver(%s).",
	12:   "Exit  " + Prefix + "RegisterObserver(%s) returned (%v).",
	13:   "Enter " + Prefix + "CloseExportReport(%+v).",
	14:   "Exit  " + Prefix + "CloseExportReport(%+v) returned (%v).",
	15:   "Enter " + Prefix + "CountRedoRecords(%+v).",
	16:   "Exit  " + Prefix + "CountRedoRecords(%+v) returned (%d, %v).",
	17:   "Enter " + Prefix + "DeleteRecord(%+v).",
	18:   "Exit  " + Prefix + "DeleteRecord(%+v) returned (%v).",
	19:   "Enter " + Prefix + "DeleteRecordWithInfo(%+v).",
	20:   "Exit  " + Prefix + "DeleteRecordWithInfo(%+v) returned (%s, %v).",
	21:   "Enter " + Prefix + "Destroy(%+v).",
	22:   "Exit  " + Prefix + "Destroy(%+v) returned (%v).",
	23:   "Enter " + Prefix + "ExportConfigAndConfigID(%+v).",
	24:   "Exit  " + Prefix + "ExportConfigAndConfigID(%+v) returned (%s, %d, %v).",
	25:   "Enter " + Prefix + "ExportConfig(%+v).",
	26:   "Exit  " + Prefix + "ExportConfig(%+v) returned (%s, %v).",
	27:   "Enter " + Prefix + "ExportCSVEntityReport(%+v).",
	28:   "Exit  " + Prefix + "ExportCSVEntityReport(%+v) returned (%v, %v).",
	29:   "Enter " + Prefix + "ExportJSONEntityReport(%+v).",
	30:   "Exit  " + Prefix + "ExportJSONEntityReport(%+v) returned (%v, %v).",
	31:   "Enter " + Prefix + "FetchNext(%+v).",
	32:   "Exit  " + Prefix + "FetchNext(%+v) returned (%s, %v).",
	33:   "Enter " + Prefix + "FindInterestingEntitiesByEntityID(%+v).",
	34:   "Exit  " + Prefix + "FindInterestingEntitiesByEntityID(%+v) returned (%s, %v).",
	35:   "Enter " + Prefix + "FindInterestingEntitiesByRecordID(%+v).",
	36:   "Exit  " + Prefix + "FindInterestingEntitiesByRecordID(%+v) returned (%s, %v).",
	37:   "Enter " + Prefix + "FindNetworkByEntityID(%+v).",
	38:   "Exit  " + Prefix + "FindNetworkByEntityID(%+v) returned (%s, %v).",
	39:   "Enter " + Prefix + "FindNetworkByEntityID_V2(%+v).",
	40:   "Exit  " + Prefix + "FindNetworkByEntityID_V2(%+v) returned (%s, %v).",
	41:   "Enter " + Prefix + "FindNetworkByRecordID(%+v).",
	42:   "Exit  " + Prefix + "FindNetworkByRecordID(%+v) returned (%s, %v).",
	43:   "Enter " + Prefix + "FindNetworkByRecordID_V2(%+v).",
	44:   "Exit  " + Prefix + "FindNetworkByRecordID_V2(%+v) returned (%s, %v).",
	45:   "Enter " + Prefix + "FindPathByEntityID(%+v).",
	46:   "Exit  " + Prefix + "FindPathByEntityID(%+v) returned (%s, %v).",
	47:   "Enter " + Prefix + "FindPathByEntityID_V2(%+v).",
	48:   "Exit  " + Prefix + "FindPathByEntityID_V2(%+v) returned (%s, %v).",
	49:   "Enter " + Prefix + "FindPathByRecordID(%+v).",
	50:   "Exit  " + Prefix + "FindPathByRecordID(%+v) returned (%s, %v).",
	51:   "Enter " + Prefix + "FindPathByRecordID_V2(%+v).",
	52:   "Exit  " + Prefix + "FindPathByRecordID_V2(%+v) returned (%s, %v).",
	53:   "Enter " + Prefix + "FindPathExcludingByEntityID(%+v).",
	54:   "Exit  " + Prefix + "FindPathExcludingByEntityID(%+v) returned (%s, %v).",
	55:   "Enter " + Prefix + "FindPathExcludingByEntityID_V2(%+v).",
	56:   "Exit  " + Prefix + "FindPathExcludingByEntityID_V2(%+v) returned (%s, %v).",
	57:   "Enter " + Prefix + "FindPathExcludingByRecordID(%+v).",
	58:   "Exit  " + Prefix + "FindPathExcludingByRecordID(%+v) returned (%s, %v).",
	59:   "Enter " + Prefix + "FindPathExcludingByRecordID_V2(%+v).",
	60:   "Exit  " + Prefix + "FindPathExcludingByRecordID_V2(%+v) returned (%v).",
	61:   "Enter " + Prefix + "FindPathIncludingSourceByEntityID(%+v).",
	62:   "Exit  " + Prefix + "FindPathIncludingSourceByEntityID(%+v) returned (%s, %v).",
	63:   "Enter " + Prefix + "FindPathIncludingSourceByEntityID_V2(%+v).",
	64:   "Exit  " + Prefix + "FindPathIncludingSourceByEntityID_V2(%+v) returned (%s, %v).",
	65:   "Enter " + Prefix + "FindPathIncludingSourceByRecordID(%+v).",
	66:   "Exit  " + Prefix + "FindPathIncludingSourceByRecordID(%+v) returned (%s, %v).",
	67:   "Enter " + Prefix + "FindPathIncludingSourceByRecordID_V2(%+v).",
	68:   "Exit  " + Prefix + "FindPathIncludingSourceByRecordID_V2(%+v) returned (%s, %v).",
	69:   "Enter " + Prefix + "GetActiveConfigID(%+v).",
	70:   "Exit  " + Prefix + "GetActiveConfigID(%+v) returned (%d, %v).",
	71:   "Enter " + Prefix + "GetEntityByEntityID(%+v).",
	72:   "Exit  " + Prefix + "GetEntityByEntityID(%+v) returned (%s, %v).",
	73:   "Enter " + Prefix + "GetEntityByEntityID_V2(%+v).",
	74:   "Exit  " + Prefix + "GetEntityByEntityID_V2(%+v) returned (%s, %v).",
	75:   "Enter " + Prefix + "GetEntityByRecordID(%+v).",
	76:   "Exit  " + Prefix + "GetEntityByRecordID(%+v) returned (%s, %v).",
	77:   "Enter " + Prefix + "GetEntityByRecordID_V2(%+v).",
	78:   "Exit  " + Prefix + "GetEntityByRecordID_V2(%+v) returned (%s, %v).",
	79:   "Enter " + Prefix + "UnregisterObserver(%s).",
	80:   "Exit  " + Prefix + "UnregisterObserver(%s) returned (%v).",
	83:   "Enter " + Prefix + "GetRecord(%+v).",
	84:   "Exit  " + Prefix + "GetRecord(%+v) returned (%s, %v).",
	85:   "Enter " + Prefix + "GetRecord_V2(%+v).",
	86:   "Exit  " + Prefix + "GetRecord_V2(%+v) returned (%s, %v).",
	87:   "Enter " + Prefix + "GetRedoRecord(%+v).",
	88:   "Exit  " + Prefix + "GetRedoRecord(%+v) returned (%s, %v).",
	89:   "Enter " + Prefix + "GetRepositoryLastModifiedTime(%+v).",
	90:   "Exit  " + Prefix + "GetRepositoryLastModifiedTime(%+v) returned (%d, %v).",
	91:   "Enter " + Prefix + "GetVirtualEntityByRecordID(%+v).",
	92:   "Exit  " + Prefix + "GetVirtualEntityByRecordID(%+v) returned (%s, %v).",
	93:   "Enter " + Prefix + "GetVirtualEntityByRecordID_V2(%+v).",
	94:   "Exit  " + Prefix + "GetVirtualEntityByRecordID_V2(%+v) returned (%s, %v).",
	95:   "Enter " + Prefix + "HowEntityByEntityID(%+v).",
	96:   "Exit  " + Prefix + "HowEntityByEntityID(%+v) returned (%s, %v).",
	97:   "Enter " + Prefix + "HowEntityByEntityID_V2(%+v).",
	98:   "Exit  " + Prefix + "HowEntityByEntityID_V2(%+v) returned (%s, %v).",
	99:   "Enter " + Prefix + "Init(%+v).",
	100:  "Exit  " + Prefix + "Init(%+v) returned (%v).",
	101:  "Enter " + Prefix + "InitWithConfigID(%+v).",
	102:  "Exit  " + Prefix + "InitWithConfigID(%+v) returned (%v).",
	103:  "Enter " + Prefix + "PrimeEngine(%+v).",
	104:  "Exit  " + Prefix + "PrimeEngine(%+v) returned (%v).",
	107:  "Enter " + Prefix + "ProcessRedoRecord(%+v).",
	108:  "Exit  " + Prefix + "ProcessRedoRecord(%+v) returned (%s, %v).",
	109:  "Enter " + Prefix + "ProcessRedoRecordWithInfo(%+v).",
	110:  "Exit  " + Prefix + "ProcessRedoRecordWithInfo(%+v) returned (%s, %s, %v).",
	119:  "Enter " + Prefix + "ReevaluateEntity(%+v).",
	120:  "Exit  " + Prefix + "ReevaluateEntity(%+v) returned (%v).",
	121:  "Enter " + Prefix + "ReevaluateEntityWithInfo(%+v).",
	122:  "Exit  " + Prefix + "ReevaluateEntityWithInfo(%+v) returned (%s, %v).",
	123:  "Enter " + Prefix + "ReevaluateRecord(%+v).",
	124:  "Exit  " + Prefix + "ReevaluateRecord(%+v) returned (%v).",
	125:  "Enter " + Prefix + "ReevaluateRecordWithInfo(%+v).",
	126:  "Exit  " + Prefix + "ReevaluateRecordWithInfo(%+v) returned (%s, %v).",
	127:  "Enter " + Prefix + "Reinit(%+v).",
	128:  "Exit  " + Prefix + "Reinit(%+v) returned (%v).",
	129:  "Enter " + Prefix + "ReplaceRecord(%+v).",
	130:  "Exit  " + Prefix + "ReplaceRecord(%+v) returned (%v).",
	131:  "Enter " + Prefix + "ReplaceRecordWithInfo(%+v).",
	132:  "Exit  " + Prefix + "ReplaceRecordWithInfo(%+v) returned (%s, %v).",
	133:  "Enter " + Prefix + "SearchByAttributes(%+v).",
	134:  "Exit  " + Prefix + "SearchByAttributes(%+v) returned (%s, %v).",
	135:  "Enter " + Prefix + "SearchByAttributes_V2(%+v).",
	136:  "Exit  " + Prefix + "SearchByAttributes_V2(%+v) returned (%s, %v).",
	137:  "Enter " + Prefix + "SetLogLevel(%s).",
	138:  "Exit  " + Prefix + "SetLogLevel(%s ) returned (%v).",
	139:  "Enter " + Prefix + "Stats(%+v).",
	140:  "Exit  " + Prefix + "Stats(%+v) returned (%s, %v).",
	141:  "Enter " + Prefix + "WhyEntities(%+v).",
	142:  "Exit  " + Prefix + "WhyEntities(%+v) returned (%s, %v).",
	143:  "Enter " + Prefix + "WhyEntities_V2(%+v).",
	144:  "Exit  " + Prefix + "WhyEntities_V2(%+v) returned (%s, %v).",
	153:  "Enter " + Prefix + "WhyRecords(%+v).",
	154:  "Exit  " + Prefix + "WhyRecords(%+v) returned (%s, %v).",
	155:  "Enter " + Prefix + "WhyRecords_V2(%+v).",
	156:  "Exit  " + Prefix + "WhyRecords_V2(%+v) returned (%s, %v).",
	157:  "Enter " + Prefix + "StreamExportCSVEntityReport(%+v).",
	158:  "Exit  " + Prefix + "StreamExportCSVEntityReport(%+v) returned (%s, %v).",
	159:  "Enter " + Prefix + "StreamExportJSONEntityReport(%+v).",
	160:  "Exit  " + Prefix + "StreamExportJSONEntityReport(%+v) returned (%s, %v).",
	161:  "Enter " + Prefix + "GetObserverOrigin().",
	162:  "Exit  " + Prefix + "GetObserverOrigin() returned (%v).",
	163:  "Enter " + Prefix + "SetObserverOrigin(%s).",
	164:  "Exit  " + Prefix + "SetObserverOrigin(%s) returned (%v).",
	165:  "Enter " + Prefix + "PreprocessRecord(%+v).",
	166:  "Exit  " + Prefix + "PreprocessRecord(%+v) returned (%v).",
	167:  "Enter " + Prefix + "WhySearch(%+v).",
	168:  "Exit  " + Prefix + "WhySearch(%+v) returned (%s, %v).",
	601:  "Send  " + Prefix + "StreamExportCSVEntityReport(%+v) item(%s).",
	602:  "Send  " + Prefix + "StreamExportJSONEntityReport(%+v) item(%s).",
	4001: Prefix + "Destroy() not supported in gRPC",
	4002: Prefix + "Init() not supported in gRPC",
	4003: Prefix + "InitWithConfigID() not supported in gRPC",
	5901: "During test setup, call to messagelogger.NewSenzingApiLogger() failed.",
	5902: "During test setup, call to szengineconfigurationjson.BuildSimpleSystemConfigurationJson() failed.",
	5903: "During test setup, call to szengine.Init() failed.",
	5904: "During test setup, call to szdiagnostic.PurgeRepository() failed.",
	5905: "During test setup, call to szengine.Destroy() failed.",
	5906: "During test setup, call to szconfig.Init() failed.",
	5907: "During test setup, call to szconfig.Create() failed.",
	5908: "During test setup, call to szconfig.RegisterDataSource() failed.",
	5909: "During test setup, call to szconfig.Save() failed.",
	5910: "During test setup, call to szconfig.Close() failed.",
	5911: "During test setup, call to szconfig.Destroy() failed.",
	5912: "During test setup, call to szconfigmgr.Init() failed.",
	5913: "During test setup, call to szconfigmgr.AddConfig() failed.",
	5914: "During test setup, call to szconfigmgr.SetDefaultConfigID() failed.",
	5915: "During test setup, call to szconfigmgr.Destroy() failed.",
	5916: "During test setup, call to szengine.Init() failed.",
	5917: "During test setup, call to szengine.AddRecord() failed.",
	5918: "During test setup, call to szengine.Destroy() failed.",
	5920: "During test setup, call to setupSenzingConfig() failed.",
	5921: "During test setup, call to setupPurgeRepository() failed.",
	5922: "During test setup, call to setupAddRecords() failed.",
	5931: "During test setup, call to szengine.Init() failed.",
	5932: "During test setup, call to szdiagnostic.PurgeRepository() failed.",
	5933: "During test setup, call to szengine.Destroy() failed.",
}

// Status strings for specific szengineserver messages.
var IDStatuses = map[int]string{}

var errPackage = errors.New("szengineserver")

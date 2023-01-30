package g2configmgrserver

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/go-logging/messagelogger"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type G2ConfigmgrServer struct {
	pb.UnimplementedG2ConfigMgrServer
	isTrace bool
	logger  messagelogger.MessageLoggerInterface
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2configmgr package found messages having the format "senzing-6999xxxx".
const ProductId = 6012

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2configmgr package.
var IdMessages = map[int]string{
	1:    "Enter AddConfig(%+v).",
	2:    "Exit  AddConfig(%+v) returned (%d, %v).",
	5:    "Enter Destroy(%+v).",
	6:    "Exit  Destroy(%+v) returned (%v).",
	7:    "Enter GetConfig(%+v).",
	8:    "Exit  GetConfig(%+v) returned (%s, %v).",
	9:    "Enter GetConfigList(%+v).",
	10:   "Exit  GetConfigList(%+v) returned (%s, %v).",
	11:   "Enter GetDefaultConfigID(%+v).",
	12:   "Exit  GetDefaultConfigID(%+v) returned (%d, %v).",
	17:   "Enter Init(%+v).",
	18:   "Exit  Init(%+v) returned (%v).",
	19:   "Enter ReplaceDefaultConfigID(%+v).",
	20:   "Exit  ReplaceDefaultConfigID(%+v) returned (%v).",
	21:   "Enter SetDefaultConfigID(%+v).",
	22:   "Exit  SetDefaultConfigID(%+v) returned (%v).",
	23:   "Enter SetLogLevel(%+v).",
	24:   "Exit  SetLogLevel(%+v) returned (%v).",
	4001: "Destroy() not supported in gRPC",
	4002: "Init() not supported in gRPC",
	4003: "InitWithConfigID() not supported in gRPC",
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

// Status strings for specific g2configmgr messages.
var IdStatuses = map[int]string{}

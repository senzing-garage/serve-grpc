package g2configmgrserver

import (
	g2pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type G2ConfigmgrServer struct {
	g2pb.UnimplementedG2ConfigMgrServer
	isTrace bool
	logger  logging.LoggingInterface
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2configmgr package found messages having the format "senzing-6999xxxx".
const ProductId = 6012

// Log message prefix.
const Prefix = "serve-grpc.g2configmgrserver."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2configmgr package.
var IdMessages = map[int]string{
	1:    "Enter " + Prefix + "AddConfig(%+v).",
	2:    "Exit  " + Prefix + "AddConfig(%+v) returned (%d, %v).",
	5:    "Enter " + Prefix + "Destroy(%+v).",
	6:    "Exit  " + Prefix + "Destroy(%+v) returned (%v).",
	7:    "Enter " + Prefix + "GetConfig(%+v).",
	8:    "Exit  " + Prefix + "GetConfig(%+v) returned (%s, %v).",
	9:    "Enter " + Prefix + "GetConfigList(%+v).",
	10:   "Exit  " + Prefix + "GetConfigList(%+v) returned (%s, %v).",
	11:   "Enter " + Prefix + "GetDefaultConfigID(%+v).",
	12:   "Exit  " + Prefix + "GetDefaultConfigID(%+v) returned (%d, %v).",
	17:   "Enter " + Prefix + "Init(%+v).",
	18:   "Exit  " + Prefix + "Init(%+v) returned (%v).",
	19:   "Enter " + Prefix + "ReplaceDefaultConfigID(%+v).",
	20:   "Exit  " + Prefix + "ReplaceDefaultConfigID(%+v) returned (%v).",
	21:   "Enter " + Prefix + "SetDefaultConfigID(%+v).",
	22:   "Exit  " + Prefix + "SetDefaultConfigID(%+v) returned (%v).",
	23:   "Enter " + Prefix + "SetLogLevel(%s).",
	24:   "Exit  " + Prefix + "SetLogLevel(%s) returned (%v).",
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

// Status strings for specific g2configmgr messages.
var IdStatuses = map[int]string{}

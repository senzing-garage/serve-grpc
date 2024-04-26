package szdiagnosticserver

import (
	"github.com/senzing-garage/go-logging/logging"
	pb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type SzDiagnosticServer struct {
	pb.UnimplementedSzDiagnosticServer
	isTrace bool
	logger  logging.LoggingInterface
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the szdiagnostic package found messages having the format "senzing-6999xxxx".
const ComponentId = 6013

// Log message prefix.
const Prefix = "serve-grpc.szdiagnosticserver."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the sz2diagnostic package.
var IdMessages = map[int]string{
	1:    "Enter " + Prefix + "CheckDBPerf(%+v).",
	2:    "Exit  " + Prefix + "CheckDBPerf(%+v) returned (%s, %v).",
	3:    "Enter " + Prefix + "RegisterObserver(%s).",
	4:    "Exit  " + Prefix + "RegisterObserver(%s) returned (%s, %v).",
	7:    "Enter " + Prefix + "Destroy(%+v).",
	8:    "Exit  " + Prefix + "Destroy(%+v) returned (%v).",
	31:   "Enter " + Prefix + "UnregisterObserver(%s).",
	32:   "Exit  " + Prefix + "UnregisterObserver(%s) returned (%s, %v).",
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
	5902: "During test setup, call to g2engineconfigurationjson.BuildSimpleSystemConfigurationJson() failed.",
	5903: "During test setup, call to g2engine.Init() failed.",
	5904: "During test setup, call to g2diagnostic.PurgeRepository() failed.",
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
	5932: "During test setup, call to g2diagnostic.PurgeRepository() failed.",
	5933: "During test setup, call to g2engine.Destroy() failed.",
}

// Status strings for specific szdiagnostic messages.
var IdStatuses = map[int]string{}

package g2configserver

import (
	g2pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type G2ConfigServer struct {
	g2pb.UnimplementedG2ConfigServer
	isTrace bool
	logger  logging.LoggingInterface
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2config package found messages having the format "senzing-6999xxxx".
const ProductId = 6011

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2config package.
var IdMessages = map[int]string{
	1:    "Enter AddDataSource(%+v).",
	2:    "Exit  AddDataSource(%+v) returned (%s, %v).",
	5:    "Enter Close(%+v).",
	6:    "Exit  Close(%+v) returned (%v).",
	7:    "Enter Create(%+v).",
	8:    "Exit  Create(%+v) returned (%v, %v).",
	9:    "Enter DeleteDataSource(%+v).",
	10:   "Exit  DeleteDataSource(%+v) returned (%v).",
	11:   "Enter Destroy(%+v).",
	12:   "Exit  Destroy(%+v) returned (%v).",
	17:   "Enter Init(%+v).",
	18:   "Exit  Init(%+v) returned (%v).",
	19:   "Enter ListDataSources(%+v).",
	20:   "Exit  ListDataSources(%+v) returned (%s, %v).",
	21:   "Enter Load(%+v).",
	22:   "Exit  Load(%+v) returned (%v).",
	23:   "Enter Save(%+v).",
	24:   "Exit  Save(%+v) returned (%s, %v).",
	25:   "Enter SetLogLevel(%+v).",
	26:   "Exit  SetLogLevel(%+v) returned (%v).",
	4001: "Destroy() not supported in gRPC",
	4002: "Init() not supported in gRPC",
	4003: "InitWithConfigID() not supported in gRPC",
	5901: "During setup, call to messagelogger.NewSenzingApiLogger() failed.",
	5902: "During setup, call to g2eg2engineconfigurationjson.BuildSimpleSystemConfigurationJson() failed.",
	5903: "During setup, call to g2engine.Init() failed.",
	5904: "During setup, call to g2engine.PurgeRepository() failed.",
	5905: "During setup, call to g2engine.Destroy() failed.",
	5906: "During setup, call to g2config.Init() failed.",
	5907: "During setup, call to g2config.Create() failed.",
	5908: "During setup, call to g2config.AddDataSource() failed.",
	5909: "During setup, call to g2config.Save() failed.",
	5910: "During setup, call to g2config.Close() failed.",
	5911: "During setup, call to g2config.Destroy() failed.",
	5912: "During setup, call to g2configmgr.Init() failed.",
	5913: "During setup, call to g2configmgr.AddConfig() failed.",
	5914: "During setup, call to g2configmgr.SetDefaultConfigID() failed.",
	5915: "During setup, call to g2configmgr.Destroy() failed.",
	5916: "During setup, call to g2engine.Init() failed.",
	5917: "During setup, call to g2engine.AddRecord() failed.",
	5918: "During setup, call to g2engine.Destroy() failed.",
	5920: "During setup, call to setupSenzingConfig() failed.",
	5921: "During setup, call to setupPurgeRepository() failed.",
	5922: "During setup, call to setupAddRecords() failed.",
	5931: "During setup, call to g2engine.Init() failed.",
	5932: "During setup, call to g2engine.PurgeRepository() failed.",
	5933: "During setup, call to g2engine.Destroy() failed.",
}

// Status strings for specific g2config messages.
var IdStatuses = map[int]string{}

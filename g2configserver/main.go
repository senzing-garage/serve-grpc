package g2configserver

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/go-logging/messagelogger"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type G2ConfigServer struct {
	pb.UnimplementedG2ConfigServer
	isTrace bool
	logger  messagelogger.MessageLoggerInterface
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

// Status strings for specific g2config messages.
var IdStatuses = map[int]string{}

package szconfigserver

import (
	"errors"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type SzConfigServer struct {
	isTrace        bool
	logger         logging.Logging
	logLevelName   string
	observerOrigin string
	observers      []observer.Observer
	szpb.UnimplementedSzConfigServer
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the szconfig package found messages having the format "senzing-6999xxxx".
const ComponentID = 6011

// Log message prefix.
const Prefix = "serve-grpc.szconfigserver."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the szconfig package.
var IDMessages = map[int]string{
	1:    "Enter " + Prefix + "RegisterDataSource(%+v).",
	2:    "Exit  " + Prefix + "RegisterDataSource(%+v) returned (%s, %v).",
	3:    "Enter " + Prefix + "RegisterObserver(%s).",
	4:    "Exit  " + Prefix + "RegisterObserver(%s) returned (%v).",
	5:    "Enter " + Prefix + "Close(%+v).",
	6:    "Exit  " + Prefix + "Close(%+v) returned (%v).",
	7:    "Enter " + Prefix + "Create(%+v).",
	8:    "Exit  " + Prefix + "Create(%+v) returned (%v, %v).",
	9:    "Enter " + Prefix + "UnregisterDataSource(%+v).",
	10:   "Exit  " + Prefix + "UnregisterDataSource(%+v) returned (%v).",
	11:   "Enter " + Prefix + "Destroy(%+v).",
	12:   "Exit  " + Prefix + "Destroy(%+v) returned (%v).",
	13:   "Enter " + Prefix + "UnregisterObserver(%s).",
	14:   "Exit  " + Prefix + "UnregisterObserver(%s) returned (%v).",
	17:   "Enter " + Prefix + "Init(%+v).",
	18:   "Exit  " + Prefix + "Init(%+v) returned (%v).",
	19:   "Enter " + Prefix + "ListDataSources(%+v).",
	20:   "Exit  " + Prefix + "ListDataSources(%+v) returned (%s, %v).",
	21:   "Enter " + Prefix + "Load(%+v).",
	22:   "Exit  " + Prefix + "Load(%+v) returned (%v).",
	23:   "Enter " + Prefix + "Save(%+v).",
	24:   "Exit  " + Prefix + "Save(%+v) returned (%s, %v).",
	25:   "Enter " + Prefix + "SetLogLevel(%s).",
	26:   "Exit  " + Prefix + "SetLogLevel(%s) returned (%v).",
	27:   "Enter " + Prefix + "GetObserverOrigin().",
	28:   "Exit  " + Prefix + "GetObserverOrigin() returned (%v).",
	29:   "Enter " + Prefix + "SetObserverOrigin(%s).",
	30:   "Exit  " + Prefix + "SetObserverOrigin(%s) returned (%v).",
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

// Status strings for specific szconfig messages.
var IDStatuses = map[int]string{}

var errPackage = errors.New("szconfigserver")

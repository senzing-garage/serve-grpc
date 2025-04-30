package szproductserver

import (
	"errors"

	"github.com/senzing-garage/go-logging/logging"
	pb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type SzProductServer struct {
	pb.UnimplementedSzProductServer
	isTrace bool
	logger  logging.Logging
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the szproductserver package found messages having the format "senzing-6999xxxx".
const ComponentID = 6016

// Log message prefix.
const Prefix = "serve-grpc.szproductserver."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the szproductserver package.
var IDMessages = map[int]string{
	1:    "Enter " + Prefix + "RegisterObserver(%s).",
	2:    "Exit  " + Prefix + "RegisterObserver(%s) returned (%v).",
	3:    "Enter " + Prefix + "Destroy(%+v).",
	4:    "Exit  " + Prefix + "Destroy(%+v) returned (%v).",
	5:    "Enter " + Prefix + "UnregisterObserver(%s).",
	6:    "Exit  " + Prefix + "UnregisterObserver(%s) returned (%v).",
	9:    "Enter " + Prefix + "Init(%+v).",
	10:   "Exit  " + Prefix + "Init(%+v) returned (%v).",
	11:   "Enter " + Prefix + "License(%+v).",
	12:   "Exit  " + Prefix + "License(%+v) returned (%s, %v).",
	13:   "Enter " + Prefix + "SetLogLevel(%s).",
	14:   "Exit  " + Prefix + "SetLogLevel(%s) returned (%v).",
	15:   "Enter " + Prefix + "ValidateLicenseFile(%+v).",
	16:   "Exit  " + Prefix + "ValidateLicenseFile(%+v) returned (%s, %v).",
	17:   "Enter " + Prefix + "ValidateLicenseStringBase64(%+v).",
	18:   "Exit  " + Prefix + "ValidateLicenseStringBase64(%+v) returned (%s, %v).",
	19:   "Enter " + Prefix + "Version(%+v).",
	20:   "Exit  " + Prefix + "Version(%+v) returned (%s, %v).",
	21:   "Enter " + Prefix + "GetObserverOrigin().",
	22:   "Exit  " + Prefix + "GetObserverOrigin() returned (%v).",
	23:   "Enter " + Prefix + "SetObserverOrigin(%s).",
	24:   "Exit  " + Prefix + "SetObserverOrigin(%s) returned (%v).",
	4001: Prefix + "Destroy() not supported in gRPC",
	4002: Prefix + "Init() not supported in gRPC",
	4003: Prefix + "InitWithConfigID() not supported in gRPC",
	5901: "During test setup, call to messagelogger.NewSenzingApiLogger() failed.",
	5902: "During test setup, call to szengineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars() failed.",
	5903: "During test setup, call to szengine.Init() failed.",
	5904: "During test setup, call to szdiagnostic.PurgeRepository() failed.",
	5905: "During test setup, call to szengine.Destroy() failed.",
	5906: "During test setup, call to szconfig.Init() failed.",
	5907: "During test setup, call to szconfig.Create() failed.",
	5908: "During test setup, call to szconfig.AddDataSource() failed.",
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

// Status strings for specific szproductserver messages.
var IDStatuses = map[int]string{}

var errPackage = errors.New("szproductserver")

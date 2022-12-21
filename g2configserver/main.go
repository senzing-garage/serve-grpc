package g2configserver

import (
	"github.com/senzing/go-logging/messagelogger"
	pb "github.com/senzing/go-servegrpc/protobuf/g2config"
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
	1:  "Enter AddDataSource(%+v).",
	2:  "Exit  AddDataSource(%+v) returned (%s, %v).",
	3:  "Enter Close(%+v).",
	4:  "Exit  Close(%+v) returned (%v).",
	5:  "Enter Create(%+v).",
	6:  "Exit  Create(%+v) returned (%+v, %v).",
	7:  "Enter DeleteDataSource(%+v).",
	8:  "Exit  DeleteDataSource(%+v) returned (%v).",
	9:  "Enter Destroy(%+v).",
	10: "Exit  Destroy(%+v) returned (%v).",
	11: "Enter Init(%+v).",
	12: "Exit  Init(%+v) returned (%v).",
	13: "Enter ListDataSources(%+v).",
	14: "Exit  ListDataSources(%+v) returned (%s, %v).",
	15: "Enter Load(%+v).",
	16: "Exit  Load(%+v) returned (%v).",
	17: "Enter Save(%+v).",
	18: "Exit  Save(%+v) returned (%s, %v).",
	19: "Enter SetLogLevel(%+v).",
	20: "Exit  SetLogLevel(%+v) returned (%v).",
}

// Status strings for specific g2config messages.
var IdStatuses = map[int]string{}

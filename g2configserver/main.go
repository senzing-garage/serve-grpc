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
const ProductId = 6999

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2config package.
var IdMessages = map[int]string{
	1:  "Enter AddDataSource(%v, %s).",
	2:  "Exit  AddDataSource(%v, %s) returned (%s, %v).",
	3:  "Enter ClearLastException().",
	4:  "Exit  ClearLastException() returned (%v).",
	5:  "Enter Close(%v).",
	6:  "Exit  Close(%v) returned (%v).",
	7:  "Enter Create().",
	8:  "Exit  Create() returned (%v, %v).",
	9:  "Enter DeleteDataSource(%v, %s).",
	10: "Exit  DeleteDataSource(%v, %s) returned (%v).",
	11: "Enter Destroy().",
	12: "Exit  Destroy() returned (%v).",
	13: "Enter GetLastException().",
	14: "Exit  GetLastException() returned (%s, %v).",
	15: "Enter GetLastExceptionCode().",
	16: "Exit  GetLastExceptionCode() returned (%d, %v).",
	17: "Enter Init(%s, %s, %d).",
	18: "Exit  Init(%s, %s, %d) returned (%v).",
	19: "Enter ListDataSources(%v).",
	20: "Exit  ListDataSources(%v) returned (%s, %v).",
	21: "Enter Load(%v, %s).",
	22: "Exit  Load(%v, %s) returned (%v).",
	23: "Enter Save(%v).",
	24: "Exit  Save(%v) returned (%s, %v).",
	25: "Enter SetLogLevel(%v).",
	26: "Exit  SetLogLevel(%v) returned (%v).",
}

// Status strings for specific g2config messages.
var IdStatuses = map[int]string{}

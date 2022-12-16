package g2diagnosticserver

import (
	pb "github.com/senzing/go-servegrpc/protobuf/g2diagnostic"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// server is used to implement helloworld.GreeterServer.
type G2DiagnosticServer struct {
	pb.UnimplementedG2DiagnosticServer
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2diagnostic package found messages having the format "senzing-6999xxxx".
const ProductId = 6999

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2diagnostic package.
var IdMessages = map[int]string{
	1:    "Enter AddDataSource(%v, %s).",
	2:    "Exit  AddDataSource(%v, %s) returned (%s, %v).",
	4001: "Call to G2Config_addDataSource(%v, %s) failed. Return code: %d",
}

// Status strings for specific g2diagnostic messages.
var IdStatuses = map[int]string{}

package g2diagnosticservercli

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-servegrpc/g2diagnosticserver"
	pb "github.com/senzing/go-servegrpc/protobuf/g2diagnostic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-9999xxxx".
const ProductId = 9999

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates.
var IdMessages = map[int]string{
	2001: "Server listening at %v",
	4001: "Call to net.Listen(tcp, %s) failed.",
	5001: "Failed to serve.",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}

// Command line arguments.
var (
	port = flag.Int("port", 50051, "The server port")
)

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func main() {

	// Configure the logging standard library.

	log.SetFlags(0)
	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)

	// Parse input options.

	flag.Parse()

	// Set up socket listener.

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Log(4001, *port, err)
	}

	// Create server.

	grpcServer := grpc.NewServer()
	pb.RegisterG2DiagnosticServer(grpcServer, &g2diagnosticserver.G2DiagnosticServer{})
	reflection.Register(grpcServer)
	logger.Log(2001, listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		logger.Log(5001, err)
	}
}

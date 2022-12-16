package grpcserver

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-servegrpc/g2diagnosticserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// GrpcServerImpl is the default implementation of the GrpcServer interface.
type GrpcServerImpl struct {
	isTrace bool
	logger  messagelogger.MessageLoggerInterface
	Port    int
	Host    string
	// FIXME:  All variables go in here. Cobra will be used outside to
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Command line arguments.
var (
	port = flag.Int("port", 50051, "The server port")
)

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func (grpcServer *GrpcServerImpl) Serve() error {
	var err error = nil

	// Configure the logging standard library.

	log.SetFlags(0)
	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)

	// Parse input options.

	flag.Parse()

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		logger.Log(5902, err)
	}
	logger.Log(2001, iniParams)

	// Set up socket listener.

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Log(4001, *port, err)
	}

	// Create server.

	grpcServer := grpc.NewServer()
	pb.RegisterG2DiagnosticServer(grpcServer, &g2diagnosticserver.G2DiagnosticServer{})
	reflection.Register(grpcServer)
	logger.Log(2002, listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		logger.Log(5001, err)
	}

	return err
}

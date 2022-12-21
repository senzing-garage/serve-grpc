package grpcserver

import (
	"fmt"
	"net"

	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-servegrpc/g2configmgrserver"
	"github.com/senzing/go-servegrpc/g2configserver"
	"github.com/senzing/go-servegrpc/g2diagnosticserver"
	"github.com/senzing/go-servegrpc/g2engineserver"
	"github.com/senzing/go-servegrpc/g2productserver"
	"github.com/senzing/go-servegrpc/protobuf/g2config"
	"github.com/senzing/go-servegrpc/protobuf/g2configmgr"
	"github.com/senzing/go-servegrpc/protobuf/g2diagnostic"
	"github.com/senzing/go-servegrpc/protobuf/g2engine"
	"github.com/senzing/go-servegrpc/protobuf/g2product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// GrpcServerImpl is the default implementation of the GrpcServer interface.
type GrpcServerImpl struct {
	Port     int
	LogLevel messagelogger.Level
	// FIXME:  All variables go in here. Cobra will be used outside to set the variables.
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func (grpcServer *GrpcServerImpl) Serve() error {
	var err error = nil
	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, grpcServer.LogLevel)

	// Log SENZING_ENGINE_CONFIGURATION_JSON for user to see. (However, it's not used.)

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		logger.Log(5902, err)
	}
	logger.Log(2001, iniParams)

	// Set up socket listener.

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServer.Port))
	if err != nil {
		logger.Log(4001, grpcServer.Port, err)
	}
	logger.Log(2002, listener.Addr())

	// Create server.

	aGrpcServer := grpc.NewServer()

	// Configure server.

	g2config.RegisterG2ConfigServer(aGrpcServer, &g2configserver.G2ConfigServer{})
	g2configmgr.RegisterG2ConfigMgrServer(aGrpcServer, &g2configmgrserver.G2ConfigmgrServer{})
	g2diagnostic.RegisterG2DiagnosticServer(aGrpcServer, &g2diagnosticserver.G2DiagnosticServer{})
	g2engine.RegisterG2EngineServer(aGrpcServer, &g2engineserver.G2EngineServer{})
	g2product.RegisterG2ProductServer(aGrpcServer, &g2productserver.G2ProductServer{})
	reflection.Register(aGrpcServer)

	// Run server.

	err = aGrpcServer.Serve(listener)
	if err != nil {
		logger.Log(5001, err)
	}

	return err
}

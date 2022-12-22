package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/senzing/go-logging/logger"
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
	EnableG2config     bool
	EnableG2configmgr  bool
	EnableG2diagnostic bool
	EnableG2engine     bool
	EnableG2product    bool
	Port               int
	LogLevel           logger.Level
	// FIXME:  All variables go in here. Cobra will be used outside to set the variables.
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func (grpcServer *GrpcServerImpl) Serve(ctx context.Context) error {
	var err error = nil
	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, grpcServer.LogLevel)

	// Log entry parameters.

	logger.Log(2000, grpcServer)

	// Determine which services to start. If no services are explicitly set, then all services are started.

	if !grpcServer.EnableG2config && !grpcServer.EnableG2configmgr && !grpcServer.EnableG2diagnostic && !grpcServer.EnableG2engine && !grpcServer.EnableG2product {
		logger.Log(2002)
		grpcServer.EnableG2config = true
		grpcServer.EnableG2configmgr = true
		grpcServer.EnableG2diagnostic = true
		grpcServer.EnableG2engine = true
		grpcServer.EnableG2product = true
	}

	// Set up socket listener.

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServer.Port))
	if err != nil {
		logger.Log(4001, grpcServer.Port, err)
	}
	logger.Log(2003, listener.Addr())

	// Create server.

	aGrpcServer := grpc.NewServer()

	// Configure server.

	if grpcServer.EnableG2config {
		server := &g2configserver.G2ConfigServer{}
		server.SetLogLevel(ctx, grpcServer.LogLevel)
		g2config.RegisterG2ConfigServer(aGrpcServer, server)
	}

	if grpcServer.EnableG2configmgr {
		server := &g2configmgrserver.G2ConfigmgrServer{}
		server.SetLogLevel(ctx, grpcServer.LogLevel)
		g2configmgr.RegisterG2ConfigMgrServer(aGrpcServer, server)
	}

	if grpcServer.EnableG2diagnostic {
		server := &g2diagnosticserver.G2DiagnosticServer{}
		server.SetLogLevel(ctx, grpcServer.LogLevel)
		g2diagnostic.RegisterG2DiagnosticServer(aGrpcServer, server)
	}

	if grpcServer.EnableG2engine {
		server := &g2engineserver.G2EngineServer{}
		server.SetLogLevel(ctx, grpcServer.LogLevel)
		g2engine.RegisterG2EngineServer(aGrpcServer, server)
	}

	if grpcServer.EnableG2product {
		server := &g2productserver.G2ProductServer{}
		server.SetLogLevel(ctx, grpcServer.LogLevel)
		g2product.RegisterG2ProductServer(aGrpcServer, server)
	}

	reflection.Register(aGrpcServer)

	// Run server.

	err = aGrpcServer.Serve(listener)
	if err != nil {
		logger.Log(5001, err)
	}

	return err
}

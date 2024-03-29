package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/senzing-garage/g2-sdk-proto/go/g2config"
	"github.com/senzing-garage/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing-garage/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing-garage/g2-sdk-proto/go/g2engine"
	"github.com/senzing-garage/g2-sdk-proto/go/g2product"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/observerpb"
	"github.com/senzing-garage/serve-grpc/g2configmgrserver"
	"github.com/senzing-garage/serve-grpc/g2configserver"
	"github.com/senzing-garage/serve-grpc/g2diagnosticserver"
	"github.com/senzing-garage/serve-grpc/g2engineserver"
	"github.com/senzing-garage/serve-grpc/g2productserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// GrpcServerImpl is the default implementation of the GrpcServer interface.
type GrpcServerImpl struct {
	EnableAll                      bool
	EnableG2config                 bool
	EnableG2configmgr              bool
	EnableG2diagnostic             bool
	EnableG2engine                 bool
	EnableG2product                bool
	logger                         logging.LoggingInterface
	LogLevelName                   string
	ObserverOrigin                 string
	Observers                      []observer.Observer
	ObserverUrl                    string
	Port                           int
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int64
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging -------------------------------------------------------------------------

// Get the Logger singleton.
func (grpcServer *GrpcServerImpl) getLogger() logging.LoggingInterface {
	var err error = nil
	if grpcServer.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 3},
		}
		grpcServer.logger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return grpcServer.logger
}

// Log message.
func (grpcServer *GrpcServerImpl) log(messageNumber int, details ...interface{}) {
	grpcServer.getLogger().Log(messageNumber, details...)
}

// --- Observing --------------------------------------------------------------

func (grpcServer *GrpcServerImpl) createGrpcObserver(ctx context.Context, parsedUrl url.URL) (observer.Observer, error) {
	_ = ctx
	var err error
	var result observer.Observer

	port := DefaultGrpcObserverPort
	if len(parsedUrl.Port()) > 0 {
		port = parsedUrl.Port()
	}
	target := fmt.Sprintf("%s:%s", parsedUrl.Hostname(), port)

	// TODO: Allow specification of options from ObserverUrl/parsedUrl
	grpcOptions := grpc.WithTransportCredentials(insecure.NewCredentials())

	grpcConnection, err := grpc.Dial(target, grpcOptions)
	if err != nil {
		return result, err
	}
	result = &observer.ObserverGrpc{
		GrpcClient: observerpb.NewObserverClient(grpcConnection),
		Id:         "serve-grpc",
	}
	return result, err
}

// --- Enabling services ---------------------------------------------------------------

// Add G2Config service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2config(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2configserver.G2ConfigServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = g2configserver.GetSdkG2config().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if err != nil {
		panic(err)
	}
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			err = server.RegisterObserver(ctx, observer)
			if err != nil {
				panic(err)
			}
		}
	}
	if len(grpcServer.ObserverOrigin) > 0 {
		server.SetObserverOrigin(ctx, grpcServer.ObserverOrigin)
	}
	g2config.RegisterG2ConfigServer(serviceRegistrar, server)
}

// Add G2Configmgr service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2configmgr(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2configmgrserver.G2ConfigmgrServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = g2configmgrserver.GetSdkG2configmgr().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if err != nil {
		panic(err)
	}
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			err = server.RegisterObserver(ctx, observer)
			if err != nil {
				panic(err)
			}
		}
	}
	if len(grpcServer.ObserverOrigin) > 0 {
		server.SetObserverOrigin(ctx, grpcServer.ObserverOrigin)
	}
	g2configmgr.RegisterG2ConfigMgrServer(serviceRegistrar, server)
}

// Add G2Diagnostic service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2diagnostic(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2diagnosticserver.G2DiagnosticServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = g2diagnosticserver.GetSdkG2diagnostic().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if err != nil {
		panic(err)
	}
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			err = server.RegisterObserver(ctx, observer)
			if err != nil {
				panic(err)
			}
		}
	}
	if len(grpcServer.ObserverOrigin) > 0 {
		server.SetObserverOrigin(ctx, grpcServer.ObserverOrigin)
	}
	g2diagnostic.RegisterG2DiagnosticServer(serviceRegistrar, server)
}

// Add G2Engine service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2engine(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2engineserver.G2EngineServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = g2engineserver.GetSdkG2engine().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if err != nil {
		panic(err)
	}
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			err = server.RegisterObserver(ctx, observer)
			if err != nil {
				panic(err)
			}
		}
	}
	if len(grpcServer.ObserverOrigin) > 0 {
		server.SetObserverOrigin(ctx, grpcServer.ObserverOrigin)
	}
	g2engine.RegisterG2EngineServer(serviceRegistrar, server)
}

// Add G2Product service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2product(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2productserver.G2ProductServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = g2productserver.GetSdkG2product().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if err != nil {
		panic(err)
	}
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			err = server.RegisterObserver(ctx, observer)
			if err != nil {
				panic(err)
			}
		}
	}
	if len(grpcServer.ObserverOrigin) > 0 {
		server.SetObserverOrigin(ctx, grpcServer.ObserverOrigin)
	}
	g2product.RegisterG2ProductServer(serviceRegistrar, server)
}

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func (grpcServer *GrpcServerImpl) Serve(ctx context.Context) error {

	// Log entry parameters.

	grpcServer.log(2000, grpcServer)

	// Initialize observing.

	var anObserver observer.Observer
	if len(grpcServer.ObserverUrl) > 0 {
		parsedUrl, err := url.Parse(grpcServer.ObserverUrl)
		if err != nil {
			return err
		}
		switch parsedUrl.Scheme {
		case "grpc":
			anObserver, err = grpcServer.createGrpcObserver(ctx, *parsedUrl)
			if err != nil {
				return err
			}
		}
		if anObserver != nil {
			grpcServer.Observers = append(grpcServer.Observers, anObserver)
		}
	}

	// Set up socket listener.

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServer.Port))
	if err != nil {
		grpcServer.log(4001, grpcServer.Port, err)
	}
	grpcServer.log(2003, listener.Addr())

	// Create server.

	aGrpcServer := grpc.NewServer()

	// Register services with gRPC server.

	if grpcServer.EnableAll || grpcServer.EnableG2config {
		grpcServer.enableG2config(ctx, aGrpcServer)
	}
	if grpcServer.EnableAll || grpcServer.EnableG2configmgr {
		grpcServer.enableG2configmgr(ctx, aGrpcServer)
	}
	if grpcServer.EnableAll || grpcServer.EnableG2diagnostic {
		grpcServer.enableG2diagnostic(ctx, aGrpcServer)
	}
	if grpcServer.EnableAll || grpcServer.EnableG2engine {
		grpcServer.enableG2engine(ctx, aGrpcServer)
	}
	if grpcServer.EnableAll || grpcServer.EnableG2product {
		grpcServer.enableG2product(ctx, aGrpcServer)
	}

	// Enable reflection.

	reflection.Register(aGrpcServer)

	// Run server.

	err = aGrpcServer.Serve(listener)
	if err != nil {
		grpcServer.log(5001, err)
	}

	return err
}

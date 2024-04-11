package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/observerpb"
	"github.com/senzing-garage/serve-grpc/szconfigmanagerserver"
	"github.com/senzing-garage/serve-grpc/szconfigserver"
	"github.com/senzing-garage/serve-grpc/szdiagnosticserver"
	"github.com/senzing-garage/serve-grpc/szengineserver"
	"github.com/senzing-garage/serve-grpc/szproductserver"
	"github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	"github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// GrpcServerImpl is the default implementation of the GrpcServer interface.
type GrpcServerImpl struct {
	EnableAll             bool
	EnableSzConfig        bool
	EnableSzConfigManager bool
	EnableSzDiagnostic    bool
	EnableSzEngine        bool
	EnableSzProduct       bool
	logger                logging.LoggingInterface
	LogLevelName          string
	ObserverOrigin        string
	Observers             []observer.Observer
	ObserverUrl           string
	Port                  int
	SenzingSettings       string
	SenzingInstanceName   string
	SenzingVerboseLogging int64
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
	server := &szconfigserver.SzConfigServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = szconfigserver.GetSdkSzConfig().Initialize(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, grpcServer.SenzingVerboseLogging)
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
	szconfig.RegisterG2ConfigServer(serviceRegistrar, server)
}

// Add G2Configmgr service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2configmgr(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &szconfigmanagerserver.SzConfigManagerServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = szconfigmanagerserver.GetSdkSzConfigManager().Init(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, grpcServer.SenzingVerboseLogging)
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
	szconfigmanager.RegisterG2ConfigMgrServer(serviceRegistrar, server)
}

// Add G2Diagnostic service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2diagnostic(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &szdiagnosticserver.SzDiagnosticServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = szdiagnosticserver.GetSdkSzDiagnostic().Init(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, grpcServer.SenzingVerboseLogging)
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
	szdiagnostic.RegisterG2DiagnosticServer(serviceRegistrar, server)
}

// Add G2Engine service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2engine(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &szengineserver.SzEngineServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = szengineserver.GetSdkG2engine().Initialize(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, grpcServer.SenzingVerboseLogging)
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
	// szengine.RegisterG2EngineServer(serviceRegistrar, server)
}

// Add G2Product service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2product(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &szproductserver.SzProductServer{}
	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}
	err = szproductserver.GetSdkSzProduct().Init(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, grpcServer.SenzingVerboseLogging)
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
	// szproduct.RegisterG2ProductServer(serviceRegistrar, server)
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

	if grpcServer.EnableAll || grpcServer.EnableSzConfig {
		grpcServer.enableG2config(ctx, aGrpcServer)
	}
	if grpcServer.EnableAll || grpcServer.EnableSzConfigManager {
		grpcServer.enableG2configmgr(ctx, aGrpcServer)
	}
	if grpcServer.EnableAll || grpcServer.EnableSzDiagnostic {
		grpcServer.enableG2diagnostic(ctx, aGrpcServer)
	}
	if grpcServer.EnableAll || grpcServer.EnableSzEngine {
		grpcServer.enableG2engine(ctx, aGrpcServer)
	}
	if grpcServer.EnableAll || grpcServer.EnableSzProduct {
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

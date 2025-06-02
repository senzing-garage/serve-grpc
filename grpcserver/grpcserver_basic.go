package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-helpers/settingsparser"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/observerpb"
	"github.com/senzing-garage/init-database/initializer"
	"github.com/senzing-garage/serve-grpc/szconfigmanagerserver"
	"github.com/senzing-garage/serve-grpc/szconfigserver"
	"github.com/senzing-garage/serve-grpc/szdiagnosticserver"
	"github.com/senzing-garage/serve-grpc/szengineserver"
	"github.com/senzing-garage/serve-grpc/szproductserver"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	"github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicGrpcServer is the default implementation of the GrpcServer interface.
type BasicGrpcServer struct {
	AvoidServing          bool
	EnableAll             bool
	EnableSzConfig        bool
	EnableSzConfigManager bool
	EnableSzDiagnostic    bool
	EnableSzEngine        bool
	EnableSzProduct       bool
	grpcserver            *grpc.Server
	GrpcServerOptions     []grpc.ServerOption
	isInitialized         bool
	logger                logging.Logging
	LogLevelName          string
	ObserverOrigin        string
	Observers             []observer.Observer
	ObserverURL           string
	Port                  int
	SenzingInstanceName   string
	SenzingSettings       string
	SenzingVerboseLogging int64
}

const OptionCallerSkip = 3

// ----------------------------------------------------------------------------
// Public methods
// ----------------------------------------------------------------------------

func (grpcServer *BasicGrpcServer) GetGRPCServer() *grpc.Server {
	return grpcServer.grpcserver
}

func (grpcServer *BasicGrpcServer) Initialize(ctx context.Context) error {
	var err error

	// Log entry parameters.

	grpcServer.log(2000, grpcServer)

	// Initialize observing.

	if len(grpcServer.ObserverURL) > 0 {
		err = grpcServer.setupObserver(ctx)
		if err != nil {
			return err
		}
	}

	// Special database processing.

	err = initializeDatabase(ctx, grpcServer.SenzingSettings)
	if err != nil {
		return err
	}

	// Create server.

	grpcServer.grpcserver = grpc.NewServer(grpcServer.GrpcServerOptions...)

	// Register services with gRPC server.

	grpcServer.enableServices(ctx, grpcServer.grpcserver)

	// Enable reflection.

	reflection.Register(grpcServer.grpcserver)

	grpcServer.isInitialized = true

	return wraperror.Errorf(err, wraperror.NoMessage)
}

func (grpcServer *BasicGrpcServer) Serve(ctx context.Context) error {
	var err error

	_ = ctx

	if !grpcServer.isInitialized {
		return wraperror.Errorf(
			errForPackage,
			"grpcserver.Serve is not initialized. BasicGrpcServer.Initialize() must be called first.",
		)
	}

	// Set up socket listener.

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServer.Port))
	if err != nil {
		grpcServer.log(4001, grpcServer.Port, err)
	}
	defer listener.Close()

	// Run server.

	if !grpcServer.AvoidServing {
		grpcServer.log(2003, listener.Addr())
		err = grpcServer.grpcserver.Serve(listener)
	} else {
		grpcServer.log(2004)
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

// --- Logging -------------------------------------------------------------------------

// Get the Logger singleton.
func (grpcServer *BasicGrpcServer) getLogger() logging.Logging {
	var err error

	if grpcServer.logger == nil {
		options := []interface{}{
			logging.OptionCallerSkip{Value: OptionCallerSkip},
			logging.OptionMessageFields{Value: []string{"id", "text", "reason", "errors", "details"}},
		}

		grpcServer.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}

	return grpcServer.logger
}

// Log message.
func (grpcServer *BasicGrpcServer) log(messageNumber int, details ...interface{}) {
	grpcServer.getLogger().Log(messageNumber, details...)
}

// --- Observing --------------------------------------------------------------

func (grpcServer *BasicGrpcServer) createGrpcObserver(
	ctx context.Context,
	parsedURL url.URL,
) (observer.Observer, error) {
	_ = ctx

	var err error

	var result observer.Observer

	port := DefaultGrpcObserverPort
	if len(parsedURL.Port()) > 0 {
		port = parsedURL.Port()
	}

	target := fmt.Sprintf("%s:%s", parsedURL.Hostname(), port)

	// IMPROVE: Allow specification of options from ObserverUrl/parsedUrl
	grpcOptions := grpc.WithTransportCredentials(insecure.NewCredentials())

	grpcConnection, err := grpc.NewClient(target, grpcOptions)
	if err != nil {
		return result, wraperror.Errorf(err, "grpc.NewClient")
	}

	result = &observer.GrpcObserver{
		GrpcClient: observerpb.NewObserverClient(grpcConnection),
		ID:         "serve-grpc",
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

// --- Enabling services ---------------------------------------------------------------

func (grpcServer *BasicGrpcServer) enableServices(ctx context.Context, aGrpcServer *grpc.Server) {
	if grpcServer.EnableAll || grpcServer.EnableSzConfig {
		grpcServer.enableSzConfig(ctx, aGrpcServer)
	}

	if grpcServer.EnableAll || grpcServer.EnableSzConfigManager {
		grpcServer.enableSzConfigManager(ctx, aGrpcServer)
	}

	if grpcServer.EnableAll || grpcServer.EnableSzDiagnostic {
		grpcServer.enableSzDiagnostic(ctx, aGrpcServer)
	}

	if grpcServer.EnableAll || grpcServer.EnableSzEngine {
		grpcServer.enableSzEngine(ctx, aGrpcServer)
	}

	if grpcServer.EnableAll || grpcServer.EnableSzProduct {
		grpcServer.enableSzProduct(ctx, aGrpcServer)
	}
}

// Add SzConfig service to gRPC server.
func (grpcServer *BasicGrpcServer) enableSzConfig(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &szconfigserver.SzConfigServer{}

	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}

	err = szconfigserver.GetSdkSzConfigManager().
		Initialize(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, grpcServer.SenzingVerboseLogging)
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

	szconfig.RegisterSzConfigServer(serviceRegistrar, server)
}

// Add SzConfigManager service to gRPC server.
func (grpcServer *BasicGrpcServer) enableSzConfigManager(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &szconfigmanagerserver.SzConfigManagerServer{}

	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}

	err = szconfigmanagerserver.GetSdkSzConfigManager().
		Initialize(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, grpcServer.SenzingVerboseLogging)
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

	szconfigmanager.RegisterSzConfigManagerServer(serviceRegistrar, server)
}

// Add SzDiagnostic service to gRPC server.
func (grpcServer *BasicGrpcServer) enableSzDiagnostic(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &szdiagnosticserver.SzDiagnosticServer{}

	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}

	err = szdiagnosticserver.GetSdkSzDiagnostic().
		Initialize(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, senzing.SzInitializeWithDefaultConfiguration, grpcServer.SenzingVerboseLogging)
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

	szdiagnostic.RegisterSzDiagnosticServer(serviceRegistrar, server)
}

// Add SzEngine service to gRPC server.
func (grpcServer *BasicGrpcServer) enableSzEngine(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &szengineserver.SzEngineServer{}

	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}

	err = szengineserver.GetSdkSzEngine().
		Initialize(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, senzing.SzInitializeWithDefaultConfiguration, grpcServer.SenzingVerboseLogging)
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

	szengine.RegisterSzEngineServer(serviceRegistrar, server)
}

// Add SzProduct service to gRPC server.
func (grpcServer *BasicGrpcServer) enableSzProduct(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &szproductserver.SzProductServer{}

	err := server.SetLogLevel(ctx, grpcServer.LogLevelName)
	if err != nil {
		panic(err)
	}

	err = szproductserver.GetSdkSzProduct().
		Initialize(ctx, grpcServer.SenzingInstanceName, grpcServer.SenzingSettings, grpcServer.SenzingVerboseLogging)
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

	szproduct.RegisterSzProductServer(serviceRegistrar, server)
}

func (grpcServer *BasicGrpcServer) setupObserver(ctx context.Context) error {
	var (
		anObserver observer.Observer
		err        error
	)

	parsedURL, err := url.Parse(grpcServer.ObserverURL)
	if err != nil {
		return wraperror.Errorf(err, "url.Parse: %s", grpcServer.ObserverURL)
	}

	if parsedURL.Scheme == "grpc" {
		anObserver, err = grpcServer.createGrpcObserver(ctx, *parsedURL)
		if err != nil {
			return wraperror.Errorf(err, "grpcServer.createGrpcObserver: %v", parsedURL)
		}
	}

	if anObserver != nil {
		grpcServer.Observers = append(grpcServer.Observers, anObserver)
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Special database initialization.
func initializeDatabase(ctx context.Context, senzingSettings string) error {
	var err error

	parsedSenzingSettings, err := settingsparser.New(senzingSettings)
	if err != nil {
		return wraperror.Errorf(err, "New")
	}

	databaseURIs, err := parsedSenzingSettings.GetDatabaseURIs(ctx)
	if err != nil {
		return wraperror.Errorf(err, "GetDatabaseURIs")
	}

	if len(databaseURIs) >= 1 {
		databaseURI := databaseURIs[0]
		if strings.HasPrefix(databaseURI, "sqlite3://") {
			err = initializeSqlite(ctx, senzingSettings, databaseURIs)
		}
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

func initializeSqlite(ctx context.Context, senzingSettings string, databaseURIs []string) error {
	var err error

	databaseURI := databaseURIs[0]

	parsedDatabaseURL, err := url.Parse(databaseURI)
	if err != nil {
		return wraperror.Errorf(err, "url.Parse: %s", databaseURI)
	}

	queryParameters := parsedDatabaseURL.Query()
	if (queryParameters.Get("mode") == "memory") && (queryParameters.Get("cache") == "shared") {
		initializer := &initializer.BasicInitializer{
			DatabaseURLs:          databaseURIs,
			ObserverOrigin:        viper.GetString(option.ObserverOrigin.Arg),
			ObserverURL:           viper.GetString(option.ObserverURL.Arg),
			SenzingInstanceName:   viper.GetString(option.EngineInstanceName.Arg),
			SenzingLogLevel:       viper.GetString(option.LogLevel.Arg),
			SenzingSettings:       senzingSettings,
			SenzingVerboseLogging: viper.GetInt64(option.EngineLogLevel.Arg),
		}

		err = initializer.Initialize(ctx)
		if err != nil {
			return wraperror.Errorf(err, "Initialize")
		}
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

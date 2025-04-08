package grpcserver_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-grpc/grpcserver"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/stretchr/testify/require"
)

var (
	localLogger logging.Logging
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setupSenzingConfig(ctx context.Context, instanceName string, settings string, verboseLogging int64) {
	now := time.Now()

	szAbstractFactory := &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	panicOnError(err)

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	panicOnError(err)

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range datasourceNames {
		_, err := szConfig.AddDataSource(ctx, dataSourceCode)
		panicOnError(err)
	}

	configComment := fmt.Sprintf("Created by szconfigmanagerserver_test at %s", now.UTC())
	configDefinition, err := szConfig.Export(ctx)
	panicOnError(err)

	configID, err := szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	panicOnError(err)

	err = szAbstractFactory.Reinitialize(ctx, configID)
	panicOnError(err)
}

func setupPurgeRepository(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	err := szDiagnostic.Initialize(
		ctx,
		instanceName,
		settings,
		senzing.SzInitializeWithDefaultConfiguration,
		verboseLogging,
	)
	if err != nil {
		return localLogger.NewError(5903, err)
	}

	err = szDiagnostic.PurgeRepository(ctx)
	if err != nil {
		return localLogger.NewError(5904, err)
	}

	err = szDiagnostic.Destroy(ctx)
	if err != nil {
		return localLogger.NewError(5905, err)
	}
	return err
}

func setup() error {
	var err error
	ctx := context.TODO()
	moduleName := "Test module name"
	verboseLogging := int64(0)

	iniParams, err := settings.BuildSimpleSettingsUsingEnvVars()
	panicOnError(err)

	// Add Data Sources to Senzing configuration.

	setupSenzingConfig(ctx, moduleName, iniParams, verboseLogging)

	// Purge repository.

	err = setupPurgeRepository(ctx, moduleName, iniParams, verboseLogging)
	panicOnError(err)

	return err
}

func teardown() error {
	var err error
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestGrpcServerImpl_Serve(test *testing.T) {
	_ = test
	ctx := context.TODO()

	observer1 := &observer.NullObserver{
		ID: "Observer 1",
	}

	logLevelName := "INFO"
	osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")
	if len(osenvLogLevel) > 0 {
		logLevelName = osenvLogLevel
	}

	senzingsettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)

	grpcServer := &grpcserver.BasicGrpcServer{
		AvoidServing:        true,
		EnableAll:           true,
		LogLevelName:        logLevelName,
		Observers:           []observer.Observer{observer1},
		ObserverOrigin:      "Test Observer origin",
		ObserverURL:         "grpc://localhost:1234",
		Port:                8258,
		SenzingInstanceName: "Test gRPC Server",
		SenzingSettings:     senzingsettings,
	}
	err = grpcServer.Serve(ctx)
	require.NoError(test, err)

}

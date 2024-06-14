package grpcserver

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

var (
	localLogger logging.LoggingInterface
)

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

func setupSenzingConfig(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	now := time.Now()

	szConfig := &szconfig.Szconfig{}
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return localLogger.NewError(5906, err)
	}

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return localLogger.NewError(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range datasourceNames {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return localLogger.NewError(5908, err)
		}
	}

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return localLogger.NewError(5909, err)
	}

	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		return localLogger.NewError(5910, err)
	}

	err = szConfig.Destroy(ctx)
	if err != nil {
		return localLogger.NewError(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return localLogger.NewError(5912, err)
	}

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return localLogger.NewError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configID)
	if err != nil {
		return localLogger.NewError(5914, err)
	}

	err = szConfigManager.Destroy(ctx)
	if err != nil {
		return localLogger.NewError(5915, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	err := szDiagnostic.Initialize(ctx, instanceName, settings, senzing.SzInitializeWithDefaultConfiguration, verboseLogging)
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
	var err error = nil
	ctx := context.TODO()
	moduleName := "Test module name"
	verboseLogging := int64(0)

	localLogger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages)
	if err != nil {
		panic(err)
	}

	iniParams, err := settings.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		return localLogger.NewError(5902, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.NewError(5920, err)
	}

	// Purge repository.

	err = setupPurgeRepository(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.NewError(5921, err)
	}

	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestGrpcServerImpl_Serve(test *testing.T) {
	// ctx := context.TODO()

	// observer1 := &observer.ObserverNull{
	// 	Id: "Observer 1",
	// }

	// senzingsettings, err := ssettings.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// grpcServer := &GrpcServerImpl{
	// 	LogLevel:                       logger.LevelInfo,
	// 	Observers:                      []observer.Observer{observer1},
	// 	Port:                           8258,
	// 	Senzingsettings: senzingsettings,
	// 	SenzingModuleName:              "Test gRPC Server",
	// }
	// grpcServer.Serve(ctx)
}

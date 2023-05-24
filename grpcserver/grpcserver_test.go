package grpcserver

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/senzing/g2-sdk-go-base/g2config"
	"github.com/senzing/g2-sdk-go-base/g2configmgr"
	"github.com/senzing/g2-sdk-go-base/g2engine"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-common/truthset"
	"github.com/senzing/go-logging/logging"
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

func setupSenzingConfig(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	now := time.Now()

	aG2config := &g2config.G2config{}
	err := aG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.NewError(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
	if err != nil {
		return localLogger.NewError(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := aG2config.AddDataSource(ctx, configHandle, datasource.Json)
		if err != nil {
			return localLogger.NewError(5908, err)
		}
	}

	configStr, err := aG2config.Save(ctx, configHandle)
	if err != nil {
		return localLogger.NewError(5909, err)
	}

	err = aG2config.Close(ctx, configHandle)
	if err != nil {
		return localLogger.NewError(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return localLogger.NewError(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &g2configmgr.G2configmgr{}
	err = aG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.NewError(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return localLogger.NewError(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return localLogger.NewError(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return localLogger.NewError(5915, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	aG2engine := &g2engine.G2engine{}
	err := aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.NewError(5903, err)
	}

	err = aG2engine.PurgeRepository(ctx)
	if err != nil {
		return localLogger.NewError(5904, err)
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return localLogger.NewError(5905, err)
	}
	return err
}

func setup() error {
	var err error = nil
	ctx := context.TODO()
	moduleName := "Test module name"
	verboseLogging := 0

	localLogger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages)
	if err != nil {
		panic(err)
	}

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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

	// senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// grpcServer := &GrpcServerImpl{
	// 	LogLevel:                       logger.LevelInfo,
	// 	Observers:                      []observer.Observer{observer1},
	// 	Port:                           8258,
	// 	SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	// 	SenzingModuleName:              "Test gRPC Server",
	// }
	// grpcServer.Serve(ctx)
}

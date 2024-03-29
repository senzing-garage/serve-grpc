/*
 */
package cmd

import (
	"context"
	"os"

	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/engineconfiguration"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/serve-grpc/grpcserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Short string = "Start a gRPC server for the Senzing SDK API"
	Use   string = "serve-grpc"
	Long  string = `
Start a gRPC server for the Senzing SDK API.
For more information, visit https://github.com/senzing-garage/serve-grpc
    `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.Configuration,
	option.DatabaseUrl,
	option.EnableAll,
	option.EnableG2config,
	option.EnableG2configmgr,
	option.EnableG2diagnostic,
	option.EnableG2engine,
	option.EnableG2product,
	option.EngineConfigurationJson,
	option.EngineLogLevel,
	option.EngineModuleName,
	option.GrpcPort,
	option.GrpcUrl,
	option.HttpPort,
	option.LogLevel,
	option.ObserverOrigin,
	option.ObserverUrl,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.Background()

	senzingEngineConfigurationJson, err := engineconfiguration.BuildAndVerifySenzingEngineConfigurationJson(ctx, viper.GetViper())
	if err != nil {
		return err
	}

	grpcserver := &grpcserver.GrpcServerImpl{
		EnableAll:                      viper.GetBool(option.EnableAll.Arg),
		EnableG2config:                 viper.GetBool(option.EnableG2config.Arg),
		EnableG2configmgr:              viper.GetBool(option.EnableG2configmgr.Arg),
		EnableG2diagnostic:             viper.GetBool(option.EnableG2diagnostic.Arg),
		EnableG2engine:                 viper.GetBool(option.EnableG2engine.Arg),
		EnableG2product:                viper.GetBool(option.EnableG2product.Arg),
		LogLevelName:                   viper.GetString(option.LogLevel.Arg),
		ObserverOrigin:                 viper.GetString(option.ObserverOrigin.Arg),
		ObserverUrl:                    viper.GetString(option.ObserverUrl.Arg),
		Port:                           viper.GetInt(option.GrpcPort.Arg),
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingModuleName:              viper.GetString(option.EngineModuleName.Arg),
		SenzingVerboseLogging:          viper.GetInt64(option.EngineLogLevel.Arg),
	}
	return grpcserver.Serve(ctx)
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}

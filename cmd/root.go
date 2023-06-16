/*
 */
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/senzing-tools/cmdhelper"
	"github.com/senzing/senzing-tools/envar"
	"github.com/senzing/senzing-tools/help"
	"github.com/senzing/senzing-tools/option"
	"github.com/senzing/serve-grpc/grpcserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Short string = "Start a gRPC server for the Senzing SDK API"
	Use   string = "serve-grpc"
	Long  string = `
Start a gRPC server for the Senzing SDK API.
For more information, visit https://github.com/Senzing/serve-grpc
    `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextBools = []cmdhelper.ContextBool{
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableAll, false),
		Envar:   envar.EnableAll,
		Help:    help.EnableAll,
		Option:  option.EnableAll,
	},
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableG2config, false),
		Envar:   envar.EnableG2config,
		Help:    help.EnableG2config,
		Option:  option.EnableG2config,
	},
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableG2configmgr, false),
		Envar:   envar.EnableG2configmgr,
		Help:    help.EnableG2configmgr,
		Option:  option.EnableG2configmgr,
	},
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableG2diagnostic, false),
		Envar:   envar.EnableG2diagnostic,
		Help:    help.EnableG2diagnostic,
		Option:  option.EnableG2diagnostic,
	},
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableG2engine, false),
		Envar:   envar.EnableG2engine,
		Help:    help.EnableG2engine,
		Option:  option.EnableG2engine,
	},
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableG2product, false),
		Envar:   envar.EnableG2product,
		Help:    help.EnableG2product,
		Option:  option.EnableG2product,
	},
}

var ContextInts = []cmdhelper.ContextInt{
	{
		Default: cmdhelper.OsLookupEnvInt(envar.EngineLogLevel, 0),
		Envar:   envar.EngineLogLevel,
		Help:    help.EngineLogLevel,
		Option:  option.EngineLogLevel,
	},
	{
		Default: cmdhelper.OsLookupEnvInt(envar.GrpcPort, 8258),
		Envar:   envar.GrpcPort,
		Help:    help.GrpcPort,
		Option:  option.GrpcPort,
	},
}

var ContextStrings = []cmdhelper.ContextString{
	{
		Default: cmdhelper.OsLookupEnvString(envar.Configuration, ""),
		Envar:   envar.Configuration,
		Help:    help.Configuration,
		Option:  option.Configuration,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.DatabaseUrl, ""),
		Envar:   envar.DatabaseUrl,
		Help:    help.DatabaseUrl,
		Option:  option.DatabaseUrl,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.EngineConfigurationJson, ""),
		Envar:   envar.EngineConfigurationJson,
		Help:    help.EngineConfigurationJson,
		Option:  option.EngineConfigurationJson,
	},
	{
		Default: fmt.Sprintf("serve-grpc-%d", time.Now().Unix()),
		Envar:   envar.EngineModuleName,
		Help:    help.EngineModuleName,
		Option:  option.EngineModuleName,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.LogLevel, "INFO"),
		Envar:   envar.LogLevel,
		Help:    help.LogLevel,
		Option:  option.LogLevel,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.ObserverOrigin, ""),
		Envar:   envar.ObserverOrigin,
		Help:    help.ObserverOrigin,
		Option:  option.ObserverOrigin,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.ObserverUrl, ""),
		Envar:   envar.ObserverUrl,
		Help:    help.ObserverUrl,
		Option:  option.ObserverUrl,
	},
}

var ContextVariables = &cmdhelper.ContextVariables{
	Bools:   ContextBools,
	Ints:    ContextInts,
	Strings: ContextStrings,
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, *ContextVariables)
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
	cmdhelper.PreRun(cobraCommand, args, Use, *ContextVariables)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.Background()

	logLevelName := viper.GetString(option.LogLevel)

	senzingEngineConfigurationJson := viper.GetString(option.EngineConfigurationJson)
	if len(senzingEngineConfigurationJson) == 0 {
		senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString(option.DatabaseUrl))
		if err != nil {
			return err
		}
	}

	grpcserver := &grpcserver.GrpcServerImpl{
		EnableAll:                      viper.GetBool(option.EnableAll),
		EnableG2config:                 viper.GetBool(option.EnableG2config),
		EnableG2configmgr:              viper.GetBool(option.EnableG2configmgr),
		EnableG2diagnostic:             viper.GetBool(option.EnableG2diagnostic),
		EnableG2engine:                 viper.GetBool(option.EnableG2engine),
		EnableG2product:                viper.GetBool(option.EnableG2product),
		ObserverOrigin:                 viper.GetString(option.ObserverOrigin),
		ObserverUrl:                    viper.GetString(option.ObserverUrl),
		Port:                           viper.GetInt(option.GrpcPort),
		LogLevelName:                   logLevelName,
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingModuleName:              viper.GetString(option.EngineModuleName),
		SenzingVerboseLogging:          viper.GetInt(option.EngineLogLevel),
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

/*
 */
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/senzing-tools/constant"
	"github.com/senzing/senzing-tools/envar"
	"github.com/senzing/senzing-tools/help"
	"github.com/senzing/senzing-tools/helper"
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

var ContextBools = []struct {
	Dfault bool
	Envar  string
	Help   string
	Option string
}{
	{
		Dfault: false,
		Envar:  envar.EnableG2config,
		Help:   help.EnableG2config,
		Option: option.EnableG2config,
	},
	{
		Dfault: false,
		Envar:  envar.EnableG2configmgr,
		Help:   help.EnableG2configmgr,
		Option: option.EnableG2configmgr,
	},
	{
		Dfault: false,
		Envar:  envar.EnableG2diagnostic,
		Help:   help.EnableG2diagnostic,
		Option: option.EnableG2diagnostic,
	},
	{
		Dfault: false,
		Envar:  envar.EnableG2engine,
		Help:   help.EnableG2engine,
		Option: option.EnableG2engine,
	},
	{
		Dfault: false,
		Envar:  envar.EnableG2product,
		Help:   help.EnableG2product,
		Option: option.EnableG2product,
	},
}

var ContextInts = []struct {
	Dfault int
	Envar  string
	Help   string
	Option string
}{
	{
		Dfault: 0,
		Envar:  envar.EngineLogLevel,
		Help:   help.EngineLogLevel,
		Option: option.EngineLogLevel,
	},
	{
		Dfault: 8258,
		Envar:  envar.GrpcPort,
		Help:   help.GrpcPort,
		Option: option.GrpcPort,
	},
}

var ContextStrings = []struct {
	Dfault string
	Envar  string
	Help   string
	Option string
}{
	{
		Dfault: "",
		Envar:  envar.Configuration,
		Help:   help.Configuration,
		Option: option.Configuration,
	},
	{
		Dfault: "",
		Envar:  envar.DatabaseUrl,
		Help:   help.DatabaseUrl,
		Option: option.DatabaseUrl,
	},
	{
		Dfault: "",
		Envar:  envar.EngineConfigurationJson,
		Help:   help.EngineConfigurationJson,
		Option: option.EngineConfigurationJson,
	},
	{
		Dfault: fmt.Sprintf("serve-grpc-%d", time.Now().Unix()),
		Envar:  envar.EngineModuleName,
		Help:   help.EngineModuleName,
		Option: option.EngineModuleName,
	},
	{
		Dfault: "INFO",
		Envar:  envar.LogLevel,
		Help:   help.LogLevel,
		Option: option.LogLevel,
	},
	{
		Dfault: "",
		Envar:  envar.ObserverOrigin,
		Help:   help.ObserverOrigin,
		Option: option.ObserverOrigin,
	},
	{
		Dfault: "",
		Envar:  envar.ObserverUrl,
		Help:   help.ObserverUrl,
		Option: option.ObserverUrl,
	},
}

var ContextStringSlices = []struct {
	Dfault []string
	Envar  string
	Help   string
	Option string
}{}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	for _, contextBool := range ContextBools {
		RootCmd.Flags().Bool(contextBool.Option, contextBool.Dfault, fmt.Sprintf(contextBool.Help, contextBool.Envar))
	}
	for _, contextInt := range ContextInts {
		RootCmd.Flags().Int(contextInt.Option, contextInt.Dfault, fmt.Sprintf(contextInt.Help, contextInt.Envar))
	}
	for _, contextString := range ContextStrings {
		RootCmd.Flags().String(contextString.Option, contextString.Dfault, fmt.Sprintf(contextString.Help, contextString.Envar))
	}
	for _, contextStringSlice := range ContextStringSlices {
		RootCmd.Flags().StringSlice(contextStringSlice.Option, contextStringSlice.Dfault, fmt.Sprintf(contextStringSlice.Help, contextStringSlice.Envar))
	}
}

// If a configuration file is present, load it.
func loadConfigurationFile(cobraCommand *cobra.Command) {
	configuration := ""
	configFlag := cobraCommand.Flags().Lookup(option.Configuration)
	if configFlag != nil {
		configuration = configFlag.Value.String()
	}
	if configuration != "" { // Use configuration file specified as a command line option.
		viper.SetConfigFile(configuration)
	} else { // Search for a configuration file.

		// Determine home directory.

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Specify configuration file name.

		viper.SetConfigName("serve-grpc")
		viper.SetConfigType("yaml")

		// Define search path order.

		viper.AddConfigPath(home + "/.senzing-tools")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/senzing-tools")
	}

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Applying configuration file:", viper.ConfigFileUsed())
	}
}

// Configure Viper with user-specified options.
func loadOptions(cobraCommand *cobra.Command) {
	var err error = nil
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(constant.SetEnvPrefix)

	// Bools

	for _, contextVar := range ContextBools {
		viper.SetDefault(contextVar.Option, contextVar.Dfault)
		err = viper.BindPFlag(contextVar.Option, cobraCommand.Flags().Lookup(contextVar.Option))
		if err != nil {
			panic(err)
		}
	}

	// Ints

	for _, contextVar := range ContextInts {
		viper.SetDefault(contextVar.Option, contextVar.Dfault)
		err = viper.BindPFlag(contextVar.Option, cobraCommand.Flags().Lookup(contextVar.Option))
		if err != nil {
			panic(err)
		}
	}

	// Strings

	for _, contextVar := range ContextStrings {
		viper.SetDefault(contextVar.Option, contextVar.Dfault)
		err = viper.BindPFlag(contextVar.Option, cobraCommand.Flags().Lookup(contextVar.Option))
		if err != nil {
			panic(err)
		}
	}

	// StringSlice

	for _, contextVar := range ContextStringSlices {
		viper.SetDefault(contextVar.Option, contextVar.Dfault)
		err = viper.BindPFlag(contextVar.Option, cobraCommand.Flags().Lookup(contextVar.Option))
		if err != nil {
			panic(err)
		}
	}
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
	loadConfigurationFile(cobraCommand)
	loadOptions(cobraCommand)
	cobraCommand.SetVersionTemplate(constant.VersionTemplate)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.TODO()

	logLevelName := viper.GetString(option.LogLevel)

	senzingEngineConfigurationJson := viper.GetString(option.EngineConfigurationJson)
	if len(senzingEngineConfigurationJson) == 0 {
		senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString(option.DatabaseUrl))
		if err != nil {
			return err
		}
	}

	grpcserver := &grpcserver.GrpcServerImpl{
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
	err = grpcserver.Serve(ctx)
	return err
}

// Used in construction of cobra.Command
func Version() string {
	return helper.MakeVersion(githubVersion, githubIteration)
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

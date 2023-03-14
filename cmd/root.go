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
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/senzing-tools/constant"
	"github.com/senzing/senzing-tools/envar"
	"github.com/senzing/senzing-tools/helper"
	"github.com/senzing/senzing-tools/option"
	"github.com/senzing/servegrpc/grpcserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfiguration           string = ""
	defaultDatabaseUrl             string = ""
	defaultEnableG2config          bool   = false
	defaultEnableG2configmgr       bool   = false
	defaultEnableG2diagnostic      bool   = false
	defaultEnableG2engine          bool   = false
	defaultEnableG2product         bool   = false
	defaultEngineConfigurationJson string = ""
	defaultEngineLogLevel          int    = 0
	defaultGrpcPort                int    = 8258
	defaultLogLevel                string = "INFO"
)

var (
	buildIteration          string = "0"
	buildVersion            string = "0.3.7"
	defaultEngineModuleName string = fmt.Sprintf("initdatabase-%d", time.Now().Unix())
)

// If a configuration file is present, load it.
func loadConfigurationFile(cobraCommand *cobra.Command) {
	configuration := cobraCommand.Flags().Lookup(option.Configuration).Value.String()
	if configuration != "" { // Use configuration file specified as a command line option.
		viper.SetConfigFile(configuration)
	} else { // Search for a configuration file.

		// Determine home directory.

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Specify configuration file name.

		viper.SetConfigName("servegrpc")
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
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(constant.SetEnvPrefix)

	// Bools

	boolOptions := map[string]bool{
		option.EnableG2config:     defaultEnableG2config,
		option.EnableG2configmgr:  defaultEnableG2configmgr,
		option.EnableG2diagnostic: defaultEnableG2diagnostic,
		option.EnableG2engine:     defaultEnableG2engine,
		option.EnableG2product:    defaultEnableG2product,
	}
	for optionKey, optionValue := range boolOptions {
		viper.SetDefault(optionKey, optionValue)
		viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
	}

	// Ints

	intOptions := map[string]int{
		option.EngineLogLevel: defaultEngineLogLevel,
		option.GrpcPort:       defaultGrpcPort,
	}
	for optionKey, optionValue := range intOptions {
		viper.SetDefault(optionKey, optionValue)
		viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
	}

	// Strings

	stringOptions := map[string]string{
		option.Configuration:           defaultConfiguration,
		option.DatabaseUrl:             defaultDatabaseUrl,
		option.EngineConfigurationJson: defaultEngineConfigurationJson,
		option.EngineModuleName:        defaultEngineModuleName,
		option.LogLevel:                defaultLogLevel,
	}
	for optionKey, optionValue := range stringOptions {
		viper.SetDefault(optionKey, optionValue)
		viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
	}
}

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:   "servegrpc",
	Short: "Start a gRPC server for the Senzing SDK API",
	Long: `
Start a gRPC server for the Senzing SDK API.
For more information, visit https://github.com/Senzing/servegrpc
	`,
	PreRun: func(cobraCommand *cobra.Command, args []string) {
		loadConfigurationFile(cobraCommand)
		loadOptions(cobraCommand)
		cobraCommand.SetVersionTemplate(constant.VersionTemplate)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		ctx := context.TODO()

		logLevel, ok := logger.TextToLevelMap[viper.GetString(option.LogLevel)]
		if !ok {
			logLevel = logger.LevelInfo
		}

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
			Port:                           viper.GetInt(option.GrpcPort),
			LogLevel:                       logLevel,
			SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
			SenzingModuleName:              viper.GetString(option.EngineModuleName),
			SenzingVerboseLogging:          viper.GetInt(option.EngineLogLevel),
		}
		err = grpcserver.Serve(ctx)
		return err
	},
	Version: helper.MakeVersion(buildVersion, buildIteration),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Since init() is always invoked, define command line parameters.
func init() {
	RootCmd.Flags().Bool(option.EnableG2config, defaultEnableG2config, fmt.Sprintf("Enable G2Config service [%s]", envar.EnableG2config))
	RootCmd.Flags().Bool(option.EnableG2configmgr, defaultEnableG2configmgr, fmt.Sprintf("Enable G2ConfigMgr service [%s]", envar.EnableG2configmgr))
	RootCmd.Flags().Bool(option.EnableG2diagnostic, defaultEnableG2diagnostic, fmt.Sprintf("Enable G2Diagnostic service [%s]", envar.EnableG2diagnostic))
	RootCmd.Flags().Bool(option.EnableG2engine, defaultEnableG2engine, fmt.Sprintf("Enable G2Config service [%s]", envar.EnableG2engine))
	RootCmd.Flags().Bool(option.EnableG2product, defaultEnableG2product, fmt.Sprintf("Enable G2Config service [%s]", envar.EnableG2product))
	RootCmd.Flags().Int(option.EngineLogLevel, defaultEngineLogLevel, fmt.Sprintf("Log level for Senzing Engine [%s]", envar.EngineLogLevel))
	RootCmd.Flags().Int(option.GrpcPort, defaultGrpcPort, fmt.Sprintf("Port used to serve gRPC [%s]", envar.GrpcPort))
	RootCmd.Flags().String(option.Configuration, defaultConfiguration, fmt.Sprintf("Path to configuration file [%s]", envar.Configuration))
	RootCmd.Flags().String(option.DatabaseUrl, defaultDatabaseUrl, fmt.Sprintf("URL of database to initialize [%s]", envar.DatabaseUrl))
	RootCmd.Flags().String(option.EngineConfigurationJson, defaultEngineConfigurationJson, fmt.Sprintf("JSON string sent to Senzing's init() function [%s]", envar.EngineConfigurationJson))
	RootCmd.Flags().String(option.EngineModuleName, defaultEngineModuleName, fmt.Sprintf("Identifier given to the Senzing engine [%s]", envar.EngineModuleName))
	RootCmd.Flags().String(option.LogLevel, defaultLogLevel, fmt.Sprintf("Log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [%s]", envar.LogLevel))
}

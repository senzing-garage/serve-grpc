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

func makeVersion(version string, iteration string) string {
	result := ""
	if buildIteration == "0" {
		result = version
	} else {
		result = fmt.Sprintf("%s-%s", buildVersion, buildIteration)
	}
	return result
}

func loadConfigurationFile(cobraCommand *cobra.Command) {
	configuration := cobraCommand.Flags().Lookup(constant.Configuration).Value.String()

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
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func loadOptions(cobraCommand *cobra.Command) {

	// Integrate with Viper.

	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(constant.SetEnvPrefix)

	// Define flags in Viper.

	viper.SetDefault(constant.DatabaseUrl, defaultDatabaseUrl)
	viper.BindPFlag(constant.DatabaseUrl, cobraCommand.Flags().Lookup(constant.DatabaseUrl))

	viper.SetDefault(constant.EnableG2config, defaultEnableG2config)
	viper.BindPFlag(constant.EnableG2config, cobraCommand.Flags().Lookup(constant.EnableG2config))

	viper.SetDefault(constant.EnableG2configmgr, defaultEnableG2configmgr)
	viper.BindPFlag(constant.EnableG2configmgr, cobraCommand.Flags().Lookup(constant.EnableG2configmgr))

	viper.SetDefault(constant.EnableG2diagnostic, defaultEnableG2diagnostic)
	viper.BindPFlag(constant.EnableG2diagnostic, cobraCommand.Flags().Lookup(constant.EnableG2diagnostic))

	viper.SetDefault(constant.EnableG2engine, defaultEnableG2engine)
	viper.BindPFlag(constant.EnableG2engine, cobraCommand.Flags().Lookup(constant.EnableG2engine))

	viper.SetDefault(constant.EnableG2product, defaultEnableG2product)
	viper.BindPFlag(constant.EnableG2product, cobraCommand.Flags().Lookup(constant.EnableG2product))

	viper.SetDefault(constant.EngineConfigurationJson, defaultEngineConfigurationJson)
	viper.BindPFlag(constant.EngineConfigurationJson, cobraCommand.Flags().Lookup(constant.EngineConfigurationJson))

	viper.SetDefault(constant.EngineLogLevel, defaultEngineLogLevel)
	viper.BindPFlag(constant.EngineLogLevel, cobraCommand.Flags().Lookup(constant.EngineLogLevel))

	viper.SetDefault(constant.EngineModuleName, defaultEngineModuleName)
	viper.BindPFlag(constant.EngineModuleName, cobraCommand.Flags().Lookup(constant.EngineModuleName))

	viper.SetDefault(constant.GrpcPort, defaultGrpcPort)
	viper.BindPFlag(constant.GrpcPort, cobraCommand.Flags().Lookup(constant.GrpcPort))

	viper.SetDefault(constant.LogLevel, defaultLogLevel)
	viper.BindPFlag(constant.LogLevel, cobraCommand.Flags().Lookup(constant.LogLevel))
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "servegrpc",
	Short: "Start a gRPC server for the Senzing SDK API",
	Long: `
Start a gRPC server for the Senzing SDK API.
For more information, visit https://github.com/Senzing/servegrpc
	`,
	PreRun: func(cobraCommand *cobra.Command, args []string) {
		fmt.Println(">>>>> servegrpc.PreRun")
		loadConfigurationFile(cobraCommand)
		loadOptions(cobraCommand)
		// Set version template.

		// versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
		cobraCommand.SetVersionTemplate(constant.VersionTemplate)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		ctx := context.TODO()

		logLevel, ok := logger.TextToLevelMap[viper.GetString(constant.LogLevel)]
		if !ok {
			logLevel = logger.LevelInfo
		}

		senzingEngineConfigurationJson := viper.GetString(constant.EngineConfigurationJson)
		if len(senzingEngineConfigurationJson) == 0 {
			senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString(constant.DatabaseUrl))
			if err != nil {
				return err
			}
		}

		grpcserver := &grpcserver.GrpcServerImpl{
			EnableG2config:                 viper.GetBool(constant.EnableG2config),
			EnableG2configmgr:              viper.GetBool(constant.EnableG2configmgr),
			EnableG2diagnostic:             viper.GetBool(constant.EnableG2diagnostic),
			EnableG2engine:                 viper.GetBool(constant.EnableG2engine),
			EnableG2product:                viper.GetBool(constant.EnableG2product),
			Port:                           viper.GetInt(constant.GrpcPort),
			LogLevel:                       logLevel,
			SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
			SenzingModuleName:              viper.GetString(constant.EngineModuleName),
			SenzingVerboseLogging:          viper.GetInt(constant.EngineLogLevel),
		}
		err = grpcserver.Serve(ctx)
		return err
	},
	Version: makeVersion(buildVersion, buildIteration),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	fmt.Println(">>>>> servegrpc.init()")
	RootCmd.Flags().Bool(constant.EnableG2config, defaultEnableG2config, fmt.Sprintf("Enable G2Config service [%s]", envar.EnableG2config))
	RootCmd.Flags().Bool(constant.EnableG2configmgr, defaultEnableG2configmgr, fmt.Sprintf("Enable G2ConfigMgr service [%s]", envar.EnableG2configmgr))
	RootCmd.Flags().Bool(constant.EnableG2diagnostic, defaultEnableG2diagnostic, fmt.Sprintf("Enable G2Diagnostic service [%s]", envar.EnableG2diagnostic))
	RootCmd.Flags().Bool(constant.EnableG2engine, defaultEnableG2engine, fmt.Sprintf("Enable G2Config service [%s]", envar.EnableG2engine))
	RootCmd.Flags().Bool(constant.EnableG2product, defaultEnableG2product, fmt.Sprintf("Enable G2Config service [%s]", envar.EnableG2product))
	RootCmd.Flags().Int(constant.EngineLogLevel, defaultEngineLogLevel, fmt.Sprintf("Log level for Senzing Engine [%s]", envar.EngineLogLevel))
	RootCmd.Flags().Int(constant.GrpcPort, defaultGrpcPort, fmt.Sprintf("Port used to serve gRPC [%s]", envar.GrpcPort))
	RootCmd.Flags().String(constant.Configuration, defaultConfiguration, fmt.Sprintf("Path to configuration file [%s]", envar.Configuration))
	RootCmd.Flags().String(constant.DatabaseUrl, defaultDatabaseUrl, fmt.Sprintf("URL of database to initialize [%s]", envar.DatabaseUrl))
	RootCmd.Flags().String(constant.EngineConfigurationJson, defaultEngineConfigurationJson, fmt.Sprintf("JSON string sent to Senzing's init() function [%s]", envar.EngineConfigurationJson))
	RootCmd.Flags().String(constant.EngineModuleName, defaultEngineModuleName, fmt.Sprintf("Identifier given to the Senzing engine [%s]", envar.EngineModuleName))
	RootCmd.Flags().String(constant.LogLevel, defaultLogLevel, fmt.Sprintf("Log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [%s]", envar.LogLevel))
}

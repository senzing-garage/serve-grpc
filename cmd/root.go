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
	"github.com/senzing/servegrpc/grpcserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildIteration                 string = "0"
	buildVersion                   string = "0.3.7"
	defaultConfigurationFile       string = ""
	defaultDatabaseUrl             string = ""
	defaultEngineConfigurationJson string = ""
	defaultEngineLogLevel          int    = 0
	defaultEngineModuleName        string = fmt.Sprintf("initdatabase-%d", time.Now().Unix())
	defaultGrpcPort                int    = 8258
	defaultLogLevel                string = "INFO"
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
		viper.AutomaticEnv()
		viper.SetDefault("configuration", defaultDatabaseUrl)
		viper.BindPFlag("configuration", cobraCommand.Flags().Lookup("configuration"))

		if viper.GetString("configuration") != "" {
			// Use config file from the flag.
			viper.SetConfigFile(viper.GetString("configuration"))
		} else {

			// Determine home directory.

			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			// Configuration file: servegrpc.yaml

			viper.SetConfigName("servegrpc")
			viper.SetConfigType("yaml")

			// Search path order.

			viper.AddConfigPath(home + "/.senzing-tools")
			viper.AddConfigPath(home)
			viper.AddConfigPath("/etc/senzing-tools")
		}

		// If a config file is found, read it in.

		if err := viper.ReadInConfig(); err == nil {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}

		// Integrate with Viper.

		replacer := strings.NewReplacer("-", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.SetEnvPrefix("SENZING_TOOLS")

		// Define flags in Viper.

		viper.SetDefault("enable-g2config", false)
		viper.BindPFlag("enable-g2config", cobraCommand.Flags().Lookup("enable-g2config"))

		viper.SetDefault("enable-g2configmgr", false)
		viper.BindPFlag("enable-g2configmgr", cobraCommand.Flags().Lookup("enable-g2configmgr"))

		viper.SetDefault("enable-g2diagnostic", false)
		viper.BindPFlag("enable-g2diagnostic", cobraCommand.Flags().Lookup("enable-g2diagnostic"))

		viper.SetDefault("enable-g2engine", false)
		viper.BindPFlag("enable-g2engine", cobraCommand.Flags().Lookup("enable-g2engine"))

		viper.SetDefault("enable-g2product", false)
		viper.BindPFlag("enable-g2product", cobraCommand.Flags().Lookup("enable-g2product"))

		viper.SetDefault("engine-log-level", defaultEngineLogLevel)
		viper.BindPFlag("engine-log-level", cobraCommand.Flags().Lookup("engine-log-level"))

		viper.SetDefault("grpc-port", defaultGrpcPort)
		viper.BindPFlag("grpc-port", cobraCommand.Flags().Lookup("grpc-port"))

		viper.SetDefault("database-url", defaultDatabaseUrl)
		viper.BindPFlag("database-url", cobraCommand.Flags().Lookup("database-url"))

		viper.SetDefault("engine-configuration-json", defaultEngineConfigurationJson)
		viper.BindPFlag("engine-configuration-json", cobraCommand.Flags().Lookup("engine-configuration-json"))

		viper.SetDefault("engine-module-name", defaultEngineModuleName)
		viper.BindPFlag("engine-module-name", cobraCommand.Flags().Lookup("engine-module-name"))

		viper.SetDefault("log-level", defaultLogLevel)
		viper.BindPFlag("log-level", cobraCommand.Flags().Lookup("log-level"))

		// Set version template.

		versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
		cobraCommand.SetVersionTemplate(versionTemplate)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		ctx := context.TODO()

		logLevel, ok := logger.TextToLevelMap[viper.GetString("log-level")]
		if !ok {
			logLevel = logger.LevelInfo
		}

		senzingEngineConfigurationJson := viper.GetString("engine-configuration-json")
		if len(senzingEngineConfigurationJson) == 0 {
			senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString("database-url"))
			if err != nil {
				return err
			}
		}

		grpcserver := &grpcserver.GrpcServerImpl{
			EnableG2config:                 viper.GetBool("enable-g2config"),
			EnableG2configmgr:              viper.GetBool("enable-g2configmgr"),
			EnableG2diagnostic:             viper.GetBool("enable-g2diagnostic"),
			EnableG2engine:                 viper.GetBool("enable-g2engine"),
			EnableG2product:                viper.GetBool("enable-g2product"),
			Port:                           viper.GetInt("grpc-port"),
			LogLevel:                       logLevel,
			SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
			SenzingModuleName:              viper.GetString("engine-module-name"),
			SenzingVerboseLogging:          viper.GetInt("engine-log-level"),
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
	RootCmd.Flags().Bool("enable-g2config", false, "Enable G2Config service [SENZING_TOOLS_ENABLE_G2CONFIG]")
	RootCmd.Flags().Bool("enable-g2configmgr", false, "Enable G2ConfigMgr service [SENZING_TOOLS_ENABLE_G2CONFIGMGR]")
	RootCmd.Flags().Bool("enable-g2diagnostic", false, "Enable G2Diagnostic service [SENZING_TOOLS_ENABLE_G2DIAGNOSTIC]")
	RootCmd.Flags().Bool("enable-g2engine", false, "Enable G2Config service [SENZING_TOOLS_ENABLE_G2ENGINE]")
	RootCmd.Flags().Bool("enable-g2product", false, "Enable G2Config service [SENZING_TOOLS_ENABLE_G2PRODUCT]")
	RootCmd.Flags().Int("engine-log-level", defaultEngineLogLevel, "Log level for Senzing Engine [SENZING_TOOLS_ENGINE_LOG_LEVEL]")
	RootCmd.Flags().Int("grpc-port", defaultGrpcPort, "Port used to serve gRPC [SENZING_TOOLS_GRPC_PORT]")
	RootCmd.Flags().String("configuration", defaultConfigurationFile, "Path to configuration file [SENZING_TOOLS_CONFIGURATION]")
	RootCmd.Flags().String("database-url", defaultDatabaseUrl, "URL of database to initialize [SENZING_TOOLS_DATABASE_URL]")
	RootCmd.Flags().String("engine-configuration-json", defaultEngineConfigurationJson, "JSON string sent to Senzing's init() function [SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]")
	RootCmd.Flags().String("engine-module-name", defaultEngineModuleName, "Identifier given to the Senzing engine [SENZING_TOOLS_ENGINE_MODULE_NAME]")
	RootCmd.Flags().String("log-level", defaultLogLevel, "Log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [SENZING_TOOLS_LOG_LEVEL]")
}

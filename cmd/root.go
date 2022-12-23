/*
 */
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/servegrpc/grpcserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	buildVersion   string = "0.0.0"
	buildIteration string = "0"
)

func makeVersion(version string, iteration string) string {
	var result string = ""
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
	Long:  `For more information, visit https://github.com/Senzing/servegrpc`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		ctx := context.TODO()

		logLevel, ok := logger.TextToLevelMap[viper.GetString("log-level")]
		if !ok {
			logLevel = logger.LevelInfo
		}

		grpcserver := &grpcserver.GrpcServerImpl{
			EnableG2config:     viper.GetBool("enable-g2config"),
			EnableG2configmgr:  viper.GetBool("enable-g2configmgr"),
			EnableG2diagnostic: viper.GetBool("enable-g2diagnostic"),
			EnableG2engine:     viper.GetBool("enable-g2engine"),
			EnableG2product:    viper.GetBool("enable-g2product"),
			Port:               viper.GetInt("grpc-port"),
			LogLevel:           logLevel,
		}
		grpcserver.Serve(ctx)
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
	cobra.OnInitialize(initConfig)

	// Define flags for command.

	RootCmd.Flags().BoolP("enable-g2config", "", false, "enable G2Config service [SENZING_TOOLS_ENABLE_G2CONFIG]")
	RootCmd.Flags().BoolP("enable-g2configmgr", "", false, "enable G2ConfigMgr service [SENZING_TOOLS_ENABLE_G2CONFIGMGR]")
	RootCmd.Flags().BoolP("enable-g2diagnostic", "", false, "enable G2Diagnostic service [SENZING_TOOLS_ENABLE_G2DIAGNOSTIC]")
	RootCmd.Flags().BoolP("enable-g2engine", "", false, "enable G2Config service [SENZING_TOOLS_ENABLE_G2ENGINE]")
	RootCmd.Flags().BoolP("enable-g2product", "", false, "enable G2Config service [SENZING_TOOLS_ENABLE_G2PRODUCT]")
	RootCmd.Flags().Int("grpc-port", 8258, "port used to serve gRPC [SENZING_TOOLS_GRPC_PORT]")
	RootCmd.Flags().String("log-level", "INFO", "log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [SENZING_TOOLS_LOG_LEVEL]")

	// Integrate with Viper.

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("SENZING_TOOLS")

	// Define flags in Viper.

	viper.SetDefault("enable-g2config", false)
	viper.BindPFlag("enable-g2config", RootCmd.Flags().Lookup("enable-g2config"))

	viper.SetDefault("enable-g2configmgr", false)
	viper.BindPFlag("enable-g2configmgr", RootCmd.Flags().Lookup("enable-g2configmgr"))

	viper.SetDefault("enable-g2diagnostic", false)
	viper.BindPFlag("enable-g2diagnostic", RootCmd.Flags().Lookup("enable-g2diagnostic"))

	viper.SetDefault("enable-g2engine", false)
	viper.BindPFlag("enable-g2engine", RootCmd.Flags().Lookup("enable-g2engine"))

	viper.SetDefault("enable-g2product", false)
	viper.BindPFlag("enable-g2product", RootCmd.Flags().Lookup("enable-g2product"))

	viper.SetDefault("grpc-port", 8258)
	viper.BindPFlag("grpc-port", RootCmd.Flags().Lookup("grpc-port"))

	viper.SetDefault("log-level", "INFO")
	viper.BindPFlag("log-level", RootCmd.Flags().Lookup("log-level"))

	// Set version template.

	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	RootCmd.SetVersionTemplate(versionTemplate)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".servegrpc" (without extension).
		viper.AddConfigPath(home + "/.senzing-tools")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/senzing-tools")
		viper.SetConfigType("yaml")
		viper.SetConfigName("servegrpc")
	}

	// Read in environment variables that match "SENZING_TOOLS_*" pattern.

	viper.AutomaticEnv()

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

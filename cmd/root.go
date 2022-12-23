/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/servegrpc/grpcserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile            string
	enableG2config     bool         = false
	enableG2configmgr  bool         = false
	enableG2diagnostic bool         = false
	enableG2engine     bool         = false
	enableG2product    bool         = false
	logLevel           logger.Level = logger.LevelInfo
	ok                 bool
	port               int = 8258
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "servegrpc",
	Short: "Start a gRPC server for the Senzing SDK API",
	Long:  `For more information, visit https://github.com/Senzing/servegrpc`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()

		// for _, key := range viper.AllKeys() {
		// 	fmt.Printf(">>> Key: %s = %v\n", key, viper.Get(key))
		// }

		grpcserver := &grpcserver.GrpcServerImpl{
			EnableG2config:     enableG2config,
			EnableG2configmgr:  enableG2configmgr,
			EnableG2diagnostic: enableG2diagnostic,
			EnableG2engine:     enableG2engine,
			EnableG2product:    enableG2product,
			Port:               port,
			LogLevel:           logLevel,
		}
		grpcserver.Serve(ctx)
	},
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

	RootCmd.Flags().BoolP("enable-g2config", "", false, "enable G2Config service [SENZING_TOOLS_ENABLE_G2CONFIG]")
	viper.SetDefault("enable_g2config", false)
	viper.BindPFlag("enable_g2config", RootCmd.Flags().Lookup("enable-g2config"))

	RootCmd.Flags().BoolP("enable-g2configmgr", "", false, "enable G2ConfigMgr service [SENZING_TOOLS_ENABLE_G2CONFIGMGR]")
	viper.SetDefault("enable_g2configmgr", false)
	viper.BindPFlag("enable_g2configmgr", RootCmd.Flags().Lookup("enable-g2configmgr"))

	RootCmd.Flags().BoolP("enable-g2diagnostic", "", false, "enable G2Diagnostic service [SENZING_TOOLS_ENABLE_G2DIAGNOSTIC]")
	viper.SetDefault("enable_g2diagnostic", false)
	viper.BindPFlag("enable_g2diagnostic", RootCmd.Flags().Lookup("enable-g2diagnostic"))

	RootCmd.Flags().BoolP("enable-g2engine", "", false, "enable G2Config service [SENZING_TOOLS_ENABLE_G2ENGINE]")
	viper.SetDefault("enable_g2engine", false)
	viper.BindPFlag("enable_g2engine", RootCmd.Flags().Lookup("enable-g2engine"))

	RootCmd.Flags().BoolP("enable-g2product", "", false, "enable G2Config service [SENZING_TOOLS_ENABLE_G2PRODUCT]")
	viper.SetDefault("enable_g2product", false)
	viper.BindPFlag("enable_g2product", RootCmd.Flags().Lookup("enable-g2product"))

	RootCmd.Flags().Int("grpc-port", 8258, "port used to serve gRPC [SENZING_TOOLS_GRPC_PORT]")
	viper.SetDefault("grpc_port", 8258)
	viper.BindPFlag("grpc_port", RootCmd.Flags().Lookup("grpc-port"))

	RootCmd.Flags().String("log-level", "INFO", "log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [SENZING_TOOLS_LOG_LEVEL]")
	viper.SetDefault("log_level", "INFO")
	viper.BindPFlag("log_level", RootCmd.Flags().Lookup("log-level"))
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
		viper.SetConfigName("config")
	}

	// if err := viper.ReadInConfig(); err != nil {
	// 	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	// 		// Config file not found; ignore error
	// 	} else {
	// 		// Config file was found but another error was produced
	// 		msglog.Log(2001, err)
	// 	}
	// }

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("senzing_tools")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Set local variables for use in command.

	enableG2config = viper.GetBool("enable_g2config")
	enableG2configmgr = viper.GetBool("enable_g2configmgr")
	enableG2diagnostic = viper.GetBool("enable_g2diagnostic")
	enableG2engine = viper.GetBool("enable_g2engine")
	enableG2product = viper.GetBool("enable_g2product")

	logLevel, ok = logger.TextToLevelMap[viper.GetString("log_level")]
	if !ok {
		logLevel = logger.LevelInfo
	}

	port = viper.GetInt("grpc_port")

}

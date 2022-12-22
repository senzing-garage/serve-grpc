package cmd

import (
	"context"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-servegrpc/grpcserver"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(servegrpcCmd)
}

var servegrpcCmd = &cobra.Command{
	Use:   "servegrpc",
	Short: "Start a gRPC server for the Senzing SDK API",
	Long:  `...long description...`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()
		grpcserver := &grpcserver.GrpcServerImpl{
			LogLevel: logger.LevelTrace,
			Port:     8258,
		}
		grpcserver.Serve(ctx)
	},
}

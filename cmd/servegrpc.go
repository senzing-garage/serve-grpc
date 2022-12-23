package cmd

// "context"

// "github.com/senzing/go-logging/logger"
// "github.com/senzing/servegrpc/grpcserver"
// "github.com/spf13/cobra"

// func init() {
// 	rootCmd.AddCommand(servegrpcCmd)
// }

// var servegrpcCmd = &cobra.Command{
// 	Use:   "servegrpc",
// 	Short: "Start a gRPC server for the Senzing SDK API",
// 	Long:  `For more information, visit https://github.com/Senzing/servegrpc`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		ctx := context.TODO()
// 		grpcserver := &grpcserver.GrpcServerImpl{
// 			LogLevel: logger.LevelTrace,
// 			Port:     8258,
// 		}
// 		grpcserver.Serve(ctx)
// 	},
// }

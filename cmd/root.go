/*
 */
package cmd

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-cmdhelping/option/optiontype"
	"github.com/senzing-garage/go-cmdhelping/settings"
	"github.com/senzing-garage/serve-grpc/grpcserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

var serverCertificatePath = option.ContextVariable{
	Arg:     "server-certificate-path",
	Default: option.OsLookupEnvString("SENZING_TOOLS_SERVER_CERTIFICATE_PATH", ""),
	Envar:   "SENZING_TOOLS_SERVER_CERTIFICATE_PATH",
	Help:    "Path to server certificate file. [%s]",
	Type:    optiontype.String,
}

var serverKeyPath = option.ContextVariable{
	Arg:     "server-key-path",
	Default: option.OsLookupEnvString("SENZING_TOOLS_SERVER_KEY_PATH", ""),
	Envar:   "SENZING_TOOLS_SERVER_KEY_PATH",
	Help:    "Path to server key file. [%s]",
	Type:    optiontype.String,
}

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	serverCertificatePath,
	serverKeyPath,
	option.AvoidServe,
	option.Configuration,
	option.DatabaseURL,
	option.EnableAll,
	option.EnableSzConfig,
	option.EnableSzConfigManager,
	option.EnableSzDiagnostic,
	option.EnableSzEngine,
	option.EnableSzProduct,
	option.EngineInstanceName,
	option.EngineLogLevel,
	option.EngineSettings,
	option.GrpcPort,
	option.GrpcURL,
	option.HTTPPort,
	option.LogLevel,
	option.ObserverOrigin,
	option.ObserverURL,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

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
	var err error
	ctx := context.Background()

	senzingSettings, err := settings.BuildAndVerifySettings(ctx, viper.GetViper())
	if err != nil {
		return err
	}

	// Add gRPC server options.

	var grpcServerOptions []grpc.ServerOption

	serverCertificatePathValue := viper.GetString(serverCertificatePath.Arg)
	serverKeyPathValue := viper.GetString(serverKeyPath.Arg)
	if serverCertificatePathValue != "" && serverKeyPathValue != "" {
		certificate, err := tls.LoadX509KeyPair(serverCertificatePathValue, serverKeyPathValue)
		if err != nil {
			return err
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{certificate},
			ClientAuth:   tls.NoClientCert,
		}
		// transportCredentials := credentials.NewServerTLSFromCert(&certificate)
		transportCredentials := credentials.NewTLS(tlsConfig)
		grpcServerOption := grpc.Creds(transportCredentials)
		grpcServerOptions = append(grpcServerOptions, grpcServerOption)
	}

	grpcserver := &grpcserver.BasicGrpcServer{
		AvoidServing:          viper.GetBool(option.AvoidServe.Arg),
		EnableAll:             viper.GetBool(option.EnableAll.Arg),
		EnableSzConfig:        viper.GetBool(option.EnableSzConfig.Arg),
		EnableSzConfigManager: viper.GetBool(option.EnableSzConfigManager.Arg),
		EnableSzDiagnostic:    viper.GetBool(option.EnableSzDiagnostic.Arg),
		EnableSzEngine:        viper.GetBool(option.EnableSzEngine.Arg),
		EnableSzProduct:       viper.GetBool(option.EnableSzProduct.Arg),
		GrpcServerOptions:     grpcServerOptions,
		LogLevelName:          viper.GetString(option.LogLevel.Arg),
		ObserverOrigin:        viper.GetString(option.ObserverOrigin.Arg),
		ObserverURL:           viper.GetString(option.ObserverURL.Arg),
		Port:                  viper.GetInt(option.GrpcPort.Arg),
		SenzingSettings:       senzingSettings,
		SenzingInstanceName:   viper.GetString(option.EngineInstanceName.Arg),
		SenzingVerboseLogging: viper.GetInt64(option.EngineLogLevel.Arg),
	}
	return grpcserver.Serve(ctx)
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}

// func loadTLSCredentials(certFile string, keyFile string) (credentials.TransportCredentials, error) {
// 	certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
// 	if err != nil {
// 		return nil, err
// 	}
// config := &tls.Config{
// 	Certificates: []tls.Certificate{certificate},
// 	ClientAuth:   tls.NoClientCert,
// }
// 	return credentials.NewServerTLSFromCert(&certificate), nil
// }

// func loadTLSCredentialsX() (credentials.TransportCredentials, error) {
// 	// Load certificate of the CA who signed client's certificate
// 	pemClientCA, err := ioutil.ReadFile(clientCACertFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	certPool := x509.NewCertPool()
// 	if !certPool.AppendCertsFromPEM(pemClientCA) {
// 		return nil, fmt.Errorf("failed to add client CA's certificate")
// 	}

// 	// Load server's certificate and private key
// 	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create the credentials and return it
// 	config := &tls.Config{
// 		Certificates: []tls.Certificate{serverCert},
// 		ClientAuth:   tls.RequireAndVerifyClientCert,
// 		ClientCAs:    certPool,
// 	}

// 	return credentials.NewTLS(config), nil
// }

/*
 */
package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
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

var clientCaCertificatePath = option.ContextVariable{
	Arg:     "client-ca-certificate-path",
	Default: option.OsLookupEnvString("SENZING_TOOLS_CLIENT_CA_CERTIFICATE_PATH", ""),
	Envar:   "SENZING_TOOLS_CLIENT_CA_CERTIFICATE_PATH",
	Help:    "Path to client Certificate Authority certificate file. [%s]",
	Type:    optiontype.String,
}

var clientCaCertificatePaths = option.ContextVariable{
	Arg:     "client-ca-certificate-paths",
	Default: []string{},
	Envar:   "SENZING_TOOLS_CLIENT_CA_CERTIFICATE_PATHS",
	Help:    "Paths to client Certificate Authority certificate files. [%s]",
	Type:    optiontype.StringSlice,
}

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
	clientCaCertificatePath,
	clientCaCertificatePaths,
	serverCertificatePath,
	serverKeyPath,
	option.AvoidServe,
	option.DatabaseURL,
	option.EnableAll,
	option.EnableSzConfig,
	option.EnableSzConfigManager,
	option.EnableSzDiagnostic,
	option.EnableSzEngine,
	option.EnableSzProduct,
	option.EngineInstanceName,
	option.EngineLogLevel,
	option.GrpcPort,
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

	// Build Senzing SDK configuration.

	senzingSettings, err := settings.BuildAndVerifySettings(ctx, viper.GetViper())
	if err != nil {
		return err
	}

	// Aggregate gRPC server options.

	var grpcServerOptions []grpc.ServerOption
	serverSideTLSServerOption, err := getServerSideTLSServerOption()
	if err != nil {
		return err
	}
	if serverSideTLSServerOption != nil {
		grpcServerOptions = append(grpcServerOptions, serverSideTLSServerOption)
	}

	// Create and Serve gRPC Server.

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

// If TLS environment variables specified, construct the appropriate gRPC server option.
func getServerSideTLSServerOption() (grpc.ServerOption, error) {
	var err error
	var result grpc.ServerOption
	var clientCAs *x509.CertPool
	var clientAuth tls.ClientAuthType
	var certificates []tls.Certificate

	serverCertificatePathValue := viper.GetString(serverCertificatePath.Arg)
	serverKeyPathValue := viper.GetString(serverKeyPath.Arg)
	if serverCertificatePathValue != "" && serverKeyPathValue != "" {

		// Server-side TLS.

		clientAuth = tls.NoClientCert
		serverCertificate, err := tls.LoadX509KeyPair(serverCertificatePathValue, serverKeyPathValue)
		if err != nil {
			return result, err
		}
		certificates = []tls.Certificate{serverCertificate}

		// Mutual TLS.

		caCertificatePathsValue := viper.GetStringSlice(clientCaCertificatePaths.Arg)
		caCertificatePathValue := viper.GetString(clientCaCertificatePath.Arg)
		if len(caCertificatePathValue) > 0 {
			caCertificatePathsValue = append(caCertificatePathsValue, caCertificatePathValue)
		}

		if len(caCertificatePathsValue) > 0 {
			clientAuth = tls.RequireAndVerifyClientCert
			clientCAs = x509.NewCertPool()
			for _, caCertificatePath := range caCertificatePathsValue {
				pemClientCA, err := os.ReadFile(caCertificatePath)
				if err != nil {
					return result, err
				}
				if !clientCAs.AppendCertsFromPEM(pemClientCA) {
					return result, fmt.Errorf("failed to add client CA's certificate for %s", caCertificatePath)
				}
			}
		}

		// Create TLS configuration.

		tlsConfig := &tls.Config{
			Certificates: certificates,
			ClientAuth:   clientAuth,
			ClientCAs:    clientCAs,
			MinVersion:   tls.VersionTLS12, // See https://pkg.go.dev/crypto/tls#pkg-constants
			MaxVersion:   tls.VersionTLS13,
		}
		transportCredentials := credentials.NewTLS(tlsConfig)
		result = grpc.Creds(transportCredentials)
		return result, err
	}

	// Input error detection.

	if serverCertificatePathValue != "" {
		err = fmt.Errorf("%s is set, but %s is not set. Both need to be set", serverCertificatePath.Envar, serverKeyPath.Envar)
		return result, err
	}
	if serverKeyPathValue != "" {
		err = fmt.Errorf("%s is set, but %s is not set. Both need to be set", serverKeyPath.Envar, serverCertificatePath.Envar)
		return result, err
	}
	return result, err
}

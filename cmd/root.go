/*
 */
package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"math"
	"os"
	"sync"
	"time"

	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-cmdhelping/option/optiontype"
	"github.com/senzing-garage/go-cmdhelping/settings"
	tlshelper "github.com/senzing-garage/go-helpers/tls"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/serve-grpc/grpcserver"
	"github.com/senzing-garage/serve-grpc/httpserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	Short string = "Start a gRPC server for the Senzing SDK API"
	Use   string = "serve-grpc"
	Long  string = `
 Start a gRPC server for the Senzing SDK API.
 For more information, visit https://github.com/senzing-garage/serve-grpc
	 `
)

// For the following, see
// - https://github.com/grpc/grpc-go/blob/master/internal/transport/defaults.go
// - https://pkg.go.dev/google.golang.org/grpc#ServerOption
const (
	defaultKeepaliveEnforcementPolicyMinTimeInSeconds             = 300 // 5 Minutes
	defaultKeepaliveEnforcementPolicyPermitWithoutStream          = false
	defaultKeepaliveServerParameterMaxConnectionAgeGraceInSeconds = 0    // Infinity
	defaultKeepaliveServerParameterMaxConnectionAgeInSeconds      = 0    // Infinity
	defaultKeepaliveServerParameterMaxConnectionIdleInSeconds     = 0    // Infinity
	defaultKeepaliveServerParameterTimeInSeconds                  = 7200 // 2 Hours
	defaultKeepaliveServerParameterTimeoutInSeconds               = 20   // 20 Seconds
	defaultMaxConcurrentStreams                                   = uint32(math.MaxUint32)
	defaultMaxHeaderListSizeInBytes                               = uint32(16 << 20)
	defaultMaxReceiveMessageSizeInBytes                           = math.MaxInt32
	defaultMaxSendMessageSizeInBytes                              = math.MaxInt32
	defaultReadBufferSizeInBytes                                  = 32 * 1024
	defaultWriteBufferSizeInBytes                                 = 32 * 1024
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var clientCaCertificateFile = option.ContextVariable{
	Arg:     "client-ca-certificate-file",
	Default: option.OsLookupEnvString("SENZING_TOOLS_CLIENT_CA_CERTIFICATE_FILE", ""),
	Envar:   "SENZING_TOOLS_CLIENT_CA_CERTIFICATE_FILE",
	Help:    "Path to client Certificate Authority certificate file. [%s]",
	Type:    optiontype.String,
}

var clientCaCertificateFiless = option.ContextVariable{
	Arg:     "client-ca-certificate-files",
	Default: []string{},
	Envar:   "SENZING_TOOLS_CLIENT_CA_CERTIFICATE_FILES",
	Help:    "Paths to client Certificate Authority certificate files. [%s]",
	Type:    optiontype.StringSlice,
}

var enableHTTP = option.ContextVariable{
	Arg:     "enable-http",
	Default: option.OsLookupEnvBool("SENZING_TOOLS_ENABLE_HTTP", false),
	Envar:   "SENZING_TOOLS_ENABLE_HTTP",
	Help:    "Enable GRPC over HTTP service [%s]",
	Type:    optiontype.Bool,
}

var keepaliveEnforcementPolicyMinTimeInSeconds = option.ContextVariable{
	Arg: "keepalive-enforcement-policy-min-time-in-seconds",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_MIN_TIME_IN_SECONDS",
		defaultKeepaliveEnforcementPolicyMinTimeInSeconds,
	),
	Envar: "SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_MIN_TIME_IN_SECONDS",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc/keepalive#EnforcementPolicy. [%s]",
	Type:  optiontype.Int,
}

var keepaliveEnforcementPolicyPermitWithoutStream = option.ContextVariable{
	Arg: "keepalive-enforcement-policy-permit-without-stream",
	Default: option.OsLookupEnvBool(
		"SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_PERMIT_WITHOUT_STREAM",
		defaultKeepaliveEnforcementPolicyPermitWithoutStream,
	),
	Envar: "SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_PERMIT_WITHOUT_STREAM",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc/keepalive#EnforcementPolicy. [%s]",
	Type:  optiontype.Bool,
}

var keepaliveServerParameterMaxConnectionAgeGraceInSeconds = option.ContextVariable{
	Arg: "keepalive-server-parameter-max-connection-age-grace-in-seconds",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_GRACE_IN_SECONDS",
		defaultKeepaliveServerParameterMaxConnectionAgeGraceInSeconds,
	),
	Envar: "SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_GRACE_IN_SECONDS",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc/keepalive#ServerParameters. [%s]",
	Type:  optiontype.Int,
}

var keepaliveServerParameterMaxConnectionAgeInSeconds = option.ContextVariable{
	Arg: "keepalive-server-parameter-max-connection-age-in-seconds",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_IN_SECONDS",
		defaultKeepaliveServerParameterMaxConnectionAgeInSeconds,
	),
	Envar: "SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_IN_SECONDS",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc/keepalive#ServerParameters. [%s]",
	Type:  optiontype.Int,
}

var keepaliveServerParameterMaxConnectionIdleInSeconds = option.ContextVariable{
	Arg: "keepalive-server-parameter-max-connection-idle-in-seconds",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_IDLE_IN_SECONDS",
		defaultKeepaliveServerParameterMaxConnectionIdleInSeconds,
	),
	Envar: "SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_IDLE_IN_SECONDS",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc/keepalive#ServerParameters. [%s]",
	Type:  optiontype.Int,
}

var keepaliveServerParameterTimeInSeconds = option.ContextVariable{
	Arg: "keepalive-server-parameter-time-in-seconds",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIME_IN_SECONDS",
		defaultKeepaliveServerParameterTimeInSeconds,
	),
	Envar: "SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIME_IN_SECONDS",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc/keepalive#ServerParameters. [%s]",
	Type:  optiontype.Int,
}

var keepaliveServerParameterTimeoutInSeconds = option.ContextVariable{
	Arg: "keepalive-server-parameter-timeout-in-seconds",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIMEOUT_IN_SECONDS",
		defaultKeepaliveServerParameterTimeoutInSeconds,
	),
	Envar: "SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIMEOUT_IN_SECONDS",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc/keepalive#ServerParameters. [%s]",
	Type:  optiontype.Int,
}

var maxConcurrentStreams = option.ContextVariable{
	Arg: "max-concurrent-streams",
	Default: option.OsLookupEnvUint32(
		"SENZING_TOOLS_SERVER_MAX_CONCURRENT_STREAMS",
		defaultMaxConcurrentStreams,
	),
	Envar: "SENZING_TOOLS_SERVER_MAX_CONCURRENT_STREAMS",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc#MaxConcurrentStreams. [%s]",
	Type:  optiontype.Uint32,
}

var maxHeaderListSizeInBytes = option.ContextVariable{
	Arg: "max-header-list-size-in-bytes",
	Default: option.OsLookupEnvUint32(
		"SENZING_TOOLS_SERVER_MAX_HEADER_LIST_SIZE_IN_BYTES",
		defaultMaxHeaderListSizeInBytes,
	),
	Envar: "SENZING_TOOLS_SERVER_MAX_HEADER_LIST_SIZE_IN_BYTES",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc#MaxHeaderListSize. [%s]",
	Type:  optiontype.Uint32,
}

var maxReceiveMessageSizeInBytes = option.ContextVariable{
	Arg: "max-receive-message-size-in-bytes",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_MAX_RECEIVE_MESSAGE_SIZE_IN_BYTES",
		defaultMaxReceiveMessageSizeInBytes,
	),
	Envar: "SENZING_TOOLS_SERVER_MAX_RECEIVE_MESSAGE_SIZE_IN_BYTES",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc#MaxRecvMsgSize. [%s]",
	Type:  optiontype.Int,
}

var maxSendMessageSizeInBytes = option.ContextVariable{
	Arg: "max-send-message-size-in-bytes",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_MAX_SEND_MESSAGE_SIZE_IN_BYTES",
		defaultMaxSendMessageSizeInBytes,
	),
	Envar: "SENZING_TOOLS_SERVER_MAX_SEND_MESSAGE_SIZE_IN_BYTES",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc#MaxSendMsgSize. [%s]",
	Type:  optiontype.Int,
}

var readBufferSizeInBytes = option.ContextVariable{
	Arg: "read-buffer-size-in-bytes",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_READ_BUFFER_SIZE_IN_BYTES",
		defaultReadBufferSizeInBytes,
	),
	Envar: "SENZING_TOOLS_SERVER_READ_BUFFER_SIZE_IN_BYTES",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc#ReadBufferSize. [%s]",
	Type:  optiontype.Int,
}

var serverCertificateFile = option.ContextVariable{
	Arg:     "server-certificate-file",
	Default: option.OsLookupEnvString("SENZING_TOOLS_SERVER_CERTIFICATE_FILE", ""),
	Envar:   "SENZING_TOOLS_SERVER_CERTIFICATE_FILE",
	Help:    "Path to server certificate file. [%s]",
	Type:    optiontype.String,
}

var serverKeyFile = option.ContextVariable{
	Arg:     "server-key-file",
	Default: option.OsLookupEnvString("SENZING_TOOLS_SERVER_KEY_FILE", ""),
	Envar:   "SENZING_TOOLS_SERVER_KEY_FILE",
	Help:    "Path to server key file. [%s]",
	Type:    optiontype.String,
}

var serverKeyPassPhrase = option.ContextVariable{
	Arg:     "server-key-passphrase",
	Default: option.OsLookupEnvString("SENZING_TOOLS_SERVER_KEY_PASSPHRASE", ""),
	Envar:   "SENZING_TOOLS_SERVER_KEY_PASSPHRASE",
	Help:    "Pass phrase used to decrypt SENZING_TOOLS_SERVER_KEY_FILE. [%s]",
	Type:    optiontype.String,
}

var writeBufferSizeInBytes = option.ContextVariable{
	Arg: "write-buffer-size-in-bytes",
	Default: option.OsLookupEnvInt(
		"SENZING_TOOLS_SERVER_WRITE_BUFFER_SIZE_IN_BYTES",
		defaultWriteBufferSizeInBytes,
	),
	Envar: "SENZING_TOOLS_SERVER_WRITE_BUFFER_SIZE_IN_BYTES",
	Help:  "See https://pkg.go.dev/google.golang.org/grpc#WriteBufferSize. [%s]",
	Type:  optiontype.Int,
}

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	clientCaCertificateFile,
	clientCaCertificateFiless,
	enableHTTP,
	keepaliveEnforcementPolicyMinTimeInSeconds,
	keepaliveEnforcementPolicyPermitWithoutStream,
	keepaliveServerParameterMaxConnectionAgeGraceInSeconds,
	keepaliveServerParameterMaxConnectionAgeInSeconds,
	keepaliveServerParameterMaxConnectionIdleInSeconds,
	keepaliveServerParameterTimeInSeconds,
	keepaliveServerParameterTimeoutInSeconds,
	maxConcurrentStreams,
	maxHeaderListSizeInBytes,
	maxReceiveMessageSizeInBytes,
	maxSendMessageSizeInBytes,
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
	option.HTTPPort,
	option.LogLevel,
	option.ObserverOrigin,
	option.ObserverURL,
	option.ServerAddress,
	readBufferSizeInBytes,
	serverCertificateFile,
	serverKeyFile,
	serverKeyPassPhrase,
	writeBufferSizeInBytes,
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

// Used in construction of cobra.Command.
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command.
func RunE(_ *cobra.Command, _ []string) error {
	var err error

	ctx := context.Background()

	// Create and Initialize gRPC Server.

	grpcserver, err := buildBasicGrpcServer(ctx)
	if err != nil {
		return wraperror.Errorf(err, "buildBasicGrpcServer")
	}

	err = grpcserver.Initialize(ctx)
	if err != nil {
		return wraperror.Errorf(err, "Initialize")
	}

	// Start services.

	err = startServers(ctx, grpcserver)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// Used in construction of cobra.Command.
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func buildBasicGrpcServer(ctx context.Context) (*grpcserver.BasicGrpcServer, error) {
	var (
		err    error
		result *grpcserver.BasicGrpcServer
	)

	senzingSettings, err := settings.BuildAndVerifySettings(ctx, viper.GetViper())
	if err != nil {
		return result, wraperror.Errorf(err, "BuildAndVerifySettings")
	}

	// Aggregate gRPC server options.

	grpcServerOptions, err := getGrpcServerOptions()
	if err != nil {
		return result, wraperror.Errorf(err, "getGrpcServerOptions")
	}

	// Create Server.

	result = &grpcserver.BasicGrpcServer{
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
		SenzingInstanceName:   viper.GetString(option.EngineInstanceName.Arg),
		SenzingSettings:       senzingSettings,
		SenzingVerboseLogging: viper.GetInt64(option.EngineLogLevel.Arg),
	}

	return result, err
}

func buildBasicHTTPServer(grpcServer *grpc.Server) *httpserver.BasicHTTPServer {
	return &httpserver.BasicHTTPServer{
		AvoidServing:    viper.GetBool(option.AvoidServe.Arg),
		EnableAll:       viper.GetBool(option.EnableAll.Arg),
		EnableGRPC:      viper.GetBool(enableHTTP.Arg),
		GRPCRoutePrefix: "grpc",
		GRPCServer:      grpcServer,
		LogLevelName:    viper.GetString(option.LogLevel.Arg),
		ServerAddress:   viper.GetString(option.ServerAddress.Arg),
		ServerPort:      viper.GetInt(option.HTTPPort.Arg),
	}
}

func buildGrpcServerOption(
	certificates []tls.Certificate,
	clientAuth tls.ClientAuthType,
	clientCAs *x509.CertPool,
) grpc.ServerOption {
	tlsConfig := &tls.Config{
		Certificates: certificates,
		ClientAuth:   clientAuth,
		ClientCAs:    clientCAs,
		MinVersion:   tls.VersionTLS12, // See https://pkg.go.dev/crypto/tls#pkg-constants
		MaxVersion:   tls.VersionTLS13,
	}
	transportCredentials := credentials.NewTLS(tlsConfig)
	result := grpc.Creds(transportCredentials)

	return result
}

func createTLSOption(serverCertificatePathValue string, serverKeyPathValue string) (grpc.ServerOption, error) {
	var (
		err          error
		result       grpc.ServerOption
		clientCAs    *x509.CertPool
		clientAuth   tls.ClientAuthType
		certificates []tls.Certificate
	)

	// Server-side TLS.

	serverKeyPassPhraseValue := viper.GetString(serverKeyPassPhrase.Arg)
	clientAuth = tls.NoClientCert

	serverCertificate, err := tlshelper.LoadX509KeyPair(
		serverCertificatePathValue,
		serverKeyPathValue,
		serverKeyPassPhraseValue,
	)
	if err != nil {
		return result, wraperror.Errorf(err, "LoadX509KeyPair")
	}

	certificates = []tls.Certificate{serverCertificate}

	// Mutual TLS.

	caCertificatePathsValue := viper.GetStringSlice(clientCaCertificateFiless.Arg)
	caCertificatePathValue := viper.GetString(clientCaCertificateFile.Arg)

	if len(caCertificatePathValue) > 0 {
		caCertificatePathsValue = append(caCertificatePathsValue, caCertificatePathValue)
	}

	if len(caCertificatePathsValue) > 0 {
		clientAuth = tls.RequireAndVerifyClientCert
		clientCAs = x509.NewCertPool()

		for _, caCertificatePath := range caCertificatePathsValue {
			pemClientCA, err := os.ReadFile(caCertificatePath)
			if err != nil {
				return result, wraperror.Errorf(err, "os.ReadFile %s", caCertificatePath)
			}

			if !clientCAs.AppendCertsFromPEM(pemClientCA) {
				return result, wraperror.Errorf(
					errPackage,
					"failed to add client CA's certificate for %s",
					caCertificatePath,
				)
			}
		}
	}

	result = buildGrpcServerOption(certificates, clientAuth, clientCAs)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func getGrpcServerOptions() ([]grpc.ServerOption, error) {
	var (
		err    error
		result []grpc.ServerOption
	)

	optionFunctions := []func() (grpc.ServerOption, error){
		getTLSOption,
		getKeepaliveEnforcementPolicyOption,
		getKeepaliveParamsOption,
		getMaxConcurrentStreamsOption,
		getMaxHeaderListSizeOption,
		getMaxRecvMsgSizeOption,
		getMaxSendMsgSizeOption,
		getReadBufferSizeOption,
		getWriteBufferSizeOption,
	}

	for _, optionFunction := range optionFunctions {
		option, err := optionFunction()
		if err != nil {
			return result, err
		}

		if option != nil {
			result = append(result, option)
		}
	}

	return result, err
}

func getKeepaliveEnforcementPolicyOption() (grpc.ServerOption, error) {
	var (
		err    error
		result grpc.ServerOption
	)

	// Option parameters.

	keepaliveEnforcementPolicyMinTimeInSecondsValue := viper.GetInt(keepaliveEnforcementPolicyMinTimeInSeconds.Arg)
	keepaliveEnforcementPolicyPermitWithoutStreamValue := viper.GetBool(
		keepaliveEnforcementPolicyPermitWithoutStream.Arg,
	)

	// Determine if a grpc.ServerOption is needed.

	if (keepaliveEnforcementPolicyMinTimeInSecondsValue != defaultKeepaliveEnforcementPolicyMinTimeInSeconds) ||
		(keepaliveEnforcementPolicyPermitWithoutStreamValue != defaultKeepaliveEnforcementPolicyPermitWithoutStream) { //nolint
		// Create grpc.ServerOption.
		keepaliveEnforcementPolicy := keepalive.EnforcementPolicy{
			MinTime:             time.Duration(keepaliveEnforcementPolicyMinTimeInSecondsValue) * time.Second,
			PermitWithoutStream: keepaliveEnforcementPolicyPermitWithoutStreamValue,
		}
		result = grpc.KeepaliveEnforcementPolicy(keepaliveEnforcementPolicy)
	}

	return result, err
}

func getKeepaliveParamsOption() (grpc.ServerOption, error) {
	var (
		err    error
		result grpc.ServerOption
	)

	// Option parameters.

	keepaliveServerParameterMaxConnectionAgeGraceInSecondsValue := viper.GetInt(
		keepaliveEnforcementPolicyMinTimeInSeconds.Arg,
	)
	keepaliveServerParameterMaxConnectionAgeInSecondsValue := viper.GetInt(
		keepaliveEnforcementPolicyMinTimeInSeconds.Arg,
	)
	keepaliveServerParameterMaxConnectionIdleInSecondsValue := viper.GetInt(
		keepaliveEnforcementPolicyMinTimeInSeconds.Arg,
	)
	keepaliveServerParameterTimeInSecondsValue := viper.GetInt(keepaliveEnforcementPolicyMinTimeInSeconds.Arg)
	keepaliveServerParameterTimeoutInSecondsValue := viper.GetInt(keepaliveEnforcementPolicyMinTimeInSeconds.Arg)

	// Determine if a grpc.ServerOption is needed.

	if (keepaliveServerParameterMaxConnectionAgeGraceInSecondsValue != defaultKeepaliveServerParameterMaxConnectionAgeGraceInSeconds) ||
		(keepaliveServerParameterMaxConnectionAgeInSecondsValue != defaultKeepaliveServerParameterMaxConnectionAgeInSeconds) ||
		(keepaliveServerParameterMaxConnectionIdleInSecondsValue != defaultKeepaliveServerParameterMaxConnectionIdleInSeconds) ||
		(keepaliveServerParameterTimeInSecondsValue != defaultKeepaliveServerParameterTimeInSeconds) ||
		(keepaliveServerParameterTimeoutInSecondsValue != defaultKeepaliveServerParameterTimeoutInSeconds) {
		// Create grpc.ServerOption.
		keepaliveServerParameters := keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(keepaliveServerParameterMaxConnectionIdleInSecondsValue) * time.Second,
			MaxConnectionAge:  time.Duration(keepaliveServerParameterMaxConnectionAgeInSecondsValue) * time.Second,
			MaxConnectionAgeGrace: time.Duration(
				keepaliveServerParameterMaxConnectionAgeGraceInSecondsValue,
			) * time.Second,
			Time:    time.Duration(keepaliveServerParameterTimeInSecondsValue) * time.Second,
			Timeout: time.Duration(keepaliveServerParameterTimeoutInSecondsValue) * time.Second,
		}
		result = grpc.KeepaliveParams(keepaliveServerParameters)
	}

	return result, err
}

func getMaxConcurrentStreamsOption() (grpc.ServerOption, error) {
	var (
		err    error
		result grpc.ServerOption
	)

	maxConcurrentStreamsValue := viper.GetUint32(maxConcurrentStreams.Arg)

	if maxConcurrentStreamsValue != defaultMaxConcurrentStreams {
		result = grpc.MaxConcurrentStreams(maxConcurrentStreamsValue)
	}

	return result, err
}

func getMaxHeaderListSizeOption() (grpc.ServerOption, error) {
	var (
		err    error
		result grpc.ServerOption
	)

	maxHeaderListSizeInBytesValue := viper.GetUint32(maxHeaderListSizeInBytes.Arg)

	if maxHeaderListSizeInBytesValue != defaultMaxHeaderListSizeInBytes {
		result = grpc.MaxHeaderListSize(maxHeaderListSizeInBytesValue)
	}

	return result, err
}

func getMaxRecvMsgSizeOption() (grpc.ServerOption, error) {
	var (
		err    error
		result grpc.ServerOption
	)

	maxReceiveMessageSizeInBytesValue := viper.GetInt(maxReceiveMessageSizeInBytes.Arg)

	if maxReceiveMessageSizeInBytesValue != defaultMaxReceiveMessageSizeInBytes {
		result = grpc.MaxRecvMsgSize(maxReceiveMessageSizeInBytesValue)
	}

	return result, err
}

func getMaxSendMsgSizeOption() (grpc.ServerOption, error) {
	var (
		err    error
		result grpc.ServerOption
	)

	maxSendMessageSizeInBytesValue := viper.GetInt(maxSendMessageSizeInBytes.Arg)

	if maxSendMessageSizeInBytesValue != defaultMaxSendMessageSizeInBytes {
		result = grpc.MaxSendMsgSize(maxSendMessageSizeInBytesValue)
	}

	return result, err
}

func getReadBufferSizeOption() (grpc.ServerOption, error) {
	var (
		err    error
		result grpc.ServerOption
	)

	readBufferSizeInBytesValue := viper.GetInt(readBufferSizeInBytes.Arg)

	if readBufferSizeInBytesValue != defaultReadBufferSizeInBytes {
		result = grpc.ReadBufferSize(readBufferSizeInBytesValue)
	}

	return result, err
}

func getWriteBufferSizeOption() (grpc.ServerOption, error) {
	var (
		err    error
		result grpc.ServerOption
	)

	writeBufferSizeInBytesValue := viper.GetInt(writeBufferSizeInBytes.Arg)

	if writeBufferSizeInBytesValue != defaultWriteBufferSizeInBytes {
		result = grpc.WriteBufferSize(writeBufferSizeInBytesValue)
	}

	return result, err
}

// Get the appropriate GRPC server options.
func getTLSOption() (grpc.ServerOption, error) {
	var (
		err    error
		result grpc.ServerOption
	)

	serverCertificatePathValue := viper.GetString(serverCertificateFile.Arg)
	serverKeyPathValue := viper.GetString(serverKeyFile.Arg)

	// Determine if TLS is requested.

	switch {
	case serverCertificatePathValue != "" && serverKeyPathValue != "":
		result, err = createTLSOption(serverCertificatePathValue, serverKeyPathValue)
	case serverCertificatePathValue != "" && serverKeyPathValue == "":
		err = wraperror.Errorf(
			errPackage,
			"%s is set, but %s is not set. Both need to be set",
			serverCertificateFile.Envar,
			serverKeyFile.Envar,
		)
	case serverCertificatePathValue == "" && serverKeyPathValue != "":
		err = wraperror.Errorf(
			errPackage,
			"%s is set, but %s is not set. Both need to be set",
			serverKeyFile.Envar,
			serverCertificateFile.Envar,
		)
	default:
		err = nil
	}

	return result, err
}

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}

func startServers(ctx context.Context, grpcserver *grpcserver.BasicGrpcServer) error {
	var (
		err       error
		waitGroup sync.WaitGroup
	)

	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		err = grpcserver.Serve(ctx)
		if err != nil {
			panic(wraperror.Errorf(err, "grpcserver.Serve"))
		}
	}()

	if viper.GetBool(enableHTTP.Arg) {
		httpserver := buildBasicHTTPServer(grpcserver.GetGRPCServer())

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			err = httpserver.Serve(ctx)
			if err != nil {
				panic(wraperror.Errorf(err, "httpserver.Serve"))
			}
		}()
	}

	waitGroup.Wait()

	return wraperror.Errorf(err, wraperror.NoMessage)
}

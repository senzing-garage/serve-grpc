package cmd_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/serve-grpc/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

func Test_RootCmd_Execute(test *testing.T) {
	args := []string{"--avoid-serving"}
	setArgs(cmd.RootCmd, args)
	err := cmd.RootCmd.Execute()
	require.NoError(test, err)
}

func Test_PreRun(_ *testing.T) {
	args := []string{"command-name", "--help"}
	cmd.PreRun(cmd.RootCmd, args)
}

func Test_RunE(test *testing.T) {
	test.Setenv("SENZING_TOOLS_AVOID_SERVING", "true")

	os.Args = []string{}
	err := cmd.RunE(cmd.RootCmd, []string{})
	require.NoError(test, err)
}

func Test_RootCmd_Execute_tls_bad_server_certificate_file(test *testing.T) {
	args := []string{
		"--avoid-serving",
		"--server-certificate-file",
		"",
		"--server-key-file",
		"../testdata/certificates/server/private_key.pem",
	}

	setArgs(cmd.RootCmd, args)
	err := cmd.RootCmd.Execute()
	require.Error(test, err)
}

func Test_RootCmd_Execute_tls_bad_server_key_file(test *testing.T) {
	_ = test
	args := []string{
		"--avoid-serving",
		"--server-certificate-file",
		"../testdata/certificates/server/certificate.pem",
		"--server-key-file",
		"",
	}
	setArgs(cmd.RootCmd, args)
	err := cmd.RootCmd.Execute()
	require.Error(test, err)
}

func Test_RootCmd_Execute_tls(test *testing.T) {
	_ = test
	args := []string{
		"--avoid-serving",
		"--server-certificate-file",
		"../testdata/certificates/server/certificate.pem",
		"--server-key-file",
		"../testdata/certificates/server/private_key.pem",
	}
	setArgs(cmd.RootCmd, args)
	err := cmd.RootCmd.Execute()
	require.NoError(test, err)
}

func Test_CompletionCmd(test *testing.T) {
	_ = test
	err := cmd.CompletionCmd.Execute()
	require.NoError(test, err)
	err = cmd.CompletionCmd.RunE(cmd.CompletionCmd, []string{})
	require.NoError(test, err)
}

func Test_DocsCmd(test *testing.T) {
	_ = test
	err := cmd.DocsCmd.Execute()
	require.NoError(test, err)
	err = cmd.DocsCmd.RunE(cmd.DocsCmd, []string{})
	require.NoError(test, err)
}

func Test_DocsAction_badDir(test *testing.T) {
	var buffer bytes.Buffer

	badDir := "/tmp/no/directory/exists"
	err := cmd.DocsAction(&buffer, badDir)
	require.Error(test, err)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Hack from https://github.com/spf13/cobra/issues/2079
func setArgs(cmd *cobra.Command, args []string) {
	if cmd.Flags().Parsed() {
		cmd.Flags().Visit(func(pf *pflag.Flag) {
			if err := pf.Value.Set(pf.DefValue); err != nil {
				panic(wraperror.Errorf(err, "reset argument[%s] value error", pf.Name))
			}
		})
	}

	cmd.SetArgs(args)
}

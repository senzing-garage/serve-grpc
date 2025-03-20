package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

func Test_Execute(test *testing.T) {
	_ = test
	setArgs(RootCmd, []string{})
	os.Args = []string{"command-name", "--avoid-serving"}
	Execute()
}

func Test_Execute_completion(test *testing.T) {
	_ = test
	setArgs(RootCmd, []string{})
	os.Args = []string{"command-name", "completion"}
	Execute()
}

func Test_Execute_docs(test *testing.T) {
	_ = test
	setArgs(RootCmd, []string{})
	os.Args = []string{"command-name", "docs"}
	Execute()
}

func Test_Execute_help(test *testing.T) {
	_ = test
	setArgs(RootCmd, []string{})
	args := []string{"--help"}
	RootCmd.SetArgs(args)
	err := RootCmd.Execute()
	require.NoError(test, err)
}

func Test_RootCmd(test *testing.T) {
	_ = test
	setArgs(RootCmd, []string{})
	os.Args = []string{"command-name", "--avoid-serving"}
	err := RootCmd.Execute()
	require.NoError(test, err)
	err = RootCmd.RunE(RootCmd, []string{})
	require.NoError(test, err)
}

func Test_RootCmd_Execute(test *testing.T) {
	_ = test
	args := []string{"--avoid-serving"}
	setArgs(RootCmd, args)
	err := RootCmd.Execute()
	require.NoError(test, err)
}

func Test_PreRun(test *testing.T) {
	_ = test
	args := []string{"command-name", "--help"}
	PreRun(RootCmd, args)
}

func Test_RunE(test *testing.T) {
	test.Setenv("SENZING_TOOLS_AVOID_SERVING", "true")
	os.Args = []string{}
	err := RunE(RootCmd, []string{})
	require.NoError(test, err)
}

func Test_completionCmd(test *testing.T) {
	_ = test
	err := completionCmd.Execute()
	require.NoError(test, err)
	err = completionCmd.RunE(completionCmd, []string{})
	require.NoError(test, err)
}

func Test_docsCmd(test *testing.T) {
	_ = test
	err := docsCmd.Execute()
	require.NoError(test, err)
	err = docsCmd.RunE(docsCmd, []string{})
	require.NoError(test, err)
}

func Test_RootCmd_Execute_tls_bad_server_certificate_file(test *testing.T) {
	_ = test
	args := []string{
		"--avoid-serving",
		"--server-certificate-file",
		"",
		"--server-key-file",
		"../testdata/certificates/server/private_key.pem",
	}
	setArgs(RootCmd, args)
	err := RootCmd.Execute()
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
	setArgs(RootCmd, args)
	err := RootCmd.Execute()
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
	setArgs(RootCmd, args)
	err := RootCmd.Execute()
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

func Test_completionAction(test *testing.T) {
	var buffer bytes.Buffer
	err := completionAction(&buffer)
	require.NoError(test, err)
}

func Test_docsAction_badDir(test *testing.T) {
	var buffer bytes.Buffer
	badDir := "/tmp/no/directory/exists"
	err := docsAction(&buffer, badDir)
	require.Error(test, err)
}

// Hack from https://github.com/spf13/cobra/issues/2079
func setArgs(cmd *cobra.Command, args []string) {
	if cmd.Flags().Parsed() {
		cmd.Flags().Visit(func(pf *pflag.Flag) {
			if err := pf.Value.Set(pf.DefValue); err != nil {
				panic(fmt.Errorf("reset argument[%s] value error %w", pf.Name, err))
			}
		})
	}
	cmd.SetArgs(args)
}

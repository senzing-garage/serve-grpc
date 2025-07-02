//go:build linux

package cmd_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/senzing-garage/serve-grpc/cmd"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

func Test_RootCmd(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "--avoid-serving"}
	err := cmd.RootCmd.Execute()
	require.NoError(test, err)
	err = cmd.RootCmd.RunE(cmd.RootCmd, []string{})
	require.NoError(test, err)
}

func Test_Execute(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "--avoid-serving"}

	cmd.Execute()
}

func Test_Execute_completion(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "completion"}

	cmd.Execute()
}

func Test_Execute_docs(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "docs"}

	cmd.Execute()
}

func Test_Execute_help(test *testing.T) {
	_ = test
	args := []string{"--help"}
	cmd.RootCmd.SetArgs(args)
	err := cmd.RootCmd.Execute()
	require.NoError(test, err)
}

func Test_RootCmd_Execute_tls_encrypted_key(test *testing.T) {
	_ = test
	args := []string{
		"--avoid-serving",
		"--server-certificate-file",
		"../testdata/certificates/server/certificate.pem",
		"--server-key-file",
		"../testdata/certificates/server/private_key_encrypted.pem",
		"--server-key-passphrase",
		"Passw0rd",
	}
	cmd.RootCmd.SetArgs(args)
	err := cmd.RootCmd.Execute()
	require.NoError(test, err)
}

func Test_DocsAction(test *testing.T) {
	var buffer bytes.Buffer

	err := cmd.DocsAction(&buffer, "/tmp")
	require.NoError(test, err)
}

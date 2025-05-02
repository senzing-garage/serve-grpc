//go:build windows

package cmd_test

import (
	"bytes"
	"testing"

	"github.com/senzing-garage/serve-grpc/cmd"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

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
	require.Error(test, err)
}

func Test_DocsAction(test *testing.T) {
	var buffer bytes.Buffer
	err := cmd.DocsAction(&buffer, "C:\\Temp")
	require.NoError(test, err)
}

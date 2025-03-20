//go:build linux

package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

func Test_Execute(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "--avoid-serving"}
	Execute()
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
	RootCmd.SetArgs(args)
	err := RootCmd.Execute()
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

func Test_docsAction(test *testing.T) {
	var buffer bytes.Buffer
	err := docsAction(&buffer, "/tmp")
	require.NoError(test, err)
}

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

func Test_Execute_completion(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "completion"}
	Execute()
}

func Test_Execute_docs(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "docs"}
	Execute()
}

func Test_Execute_help(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "--help"}
	Execute()
}

func Test_Execute_tls(test *testing.T) {
	_ = test
	// os.Args = []string{
	// 	"command-name",
	// 	"--avoid-serving",
	// 	"--server-certificate-path",
	// 	"../testdata/certificates/server/certificate.pem",
	// 	"--server-key-path",
	// 	"../testdata/certificates/server/private_key.pem"}
	os.Args = []string{"command-name", "--avoid-serving"}
	Execute()
}

// func Test_Execute_tls_bad_no_server_key_path(test *testing.T) {
// 	_ = test
// 	os.Args = []string{
// 		"command-name",
// 		"--avoid-serving",
// 		"--server-certificate-path",
// 		"../testdata/certificates/server/certificate.pem"}
// 	Execute()
// }

// func Test_Execute_tls_bad_no_server_certificate_path(test *testing.T) {
// 	_ = test
// 	os.Args = []string{
// 		"command-name",
// 		"--avoid-serving",
// 		"--server-key-path",
// 		"../testdata/certificates/server/Xprivate_key.pem"}
// 	Execute()
// }

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

func Test_RunE_tls(test *testing.T) {
	// test.Setenv("SENZING_TOOLS_AVOID_SERVING", "true")
	cmdLineArgs := []string{
		"command-name",
		"--avoid-serving",
		"--server-certificate-path",
		"../testdata/certificates/server/certificate.pem",
		"--server-key-path",
		"../testdata/certificates/server/private_key.pem"}
	err := RunE(RootCmd, cmdLineArgs)
	require.NoError(test, err)
}

func Test_RootCmd(test *testing.T) {
	_ = test
	err := RootCmd.Execute()
	require.NoError(test, err)
	err = RootCmd.RunE(RootCmd, []string{})
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

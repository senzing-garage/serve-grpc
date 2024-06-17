//go:build windows

package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

func Test_docsAction(test *testing.T) {
	var buffer bytes.Buffer
	err := docsAction(&buffer, "C:\\Temp")
	require.NoError(test, err)
}

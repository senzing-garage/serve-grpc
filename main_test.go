package main

import (
	"os"
	"testing"
)

func TestMain(test *testing.T) {
	test.Parallel()

	os.Args = []string{"command-name", "--avoid-serving"}

	main()
}

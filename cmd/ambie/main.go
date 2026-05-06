// AMBIE command-line interface.
//
// Install:
//
//	curl -fsSL https://raw.githubusercontent.com/ambie-ai/ambie-cli/main/install.sh | sh
//
// or download a binary from
//
//	https://github.com/ambie-ai/ambie-cli/releases/latest
//
// Get an API key at https://ambie.ai/signup.
package main

import (
	"fmt"
	"os"

	"github.com/ambie-ai/ambie-cli/internal/cli"
)

func main() {
	if err := cli.Root().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

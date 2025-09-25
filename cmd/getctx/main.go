package main

import (
	"fmt"
	"getctx/cmd/getctx/cli"
	"os"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

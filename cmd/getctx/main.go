package main

import (
	"fmt"
	"os"

	"github.com/kacperzielinskidev/getctx/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

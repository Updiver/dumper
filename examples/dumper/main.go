package main

import (
	"fmt"
	"os"

	cmd "github.com/updiver/dumper/examples/dumper/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "dumper error: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"os"

	"github.com/yourepena/stamperbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

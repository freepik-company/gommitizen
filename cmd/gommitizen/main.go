package main

import (
	"os"

	"github.com/freepik-company/gommitizen/internal/pkg/cmd"
)

func main() {
	root := cmd.Root()

	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

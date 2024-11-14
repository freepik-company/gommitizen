package main

import (
	"os"

	"github.com/freepik-company/gommitizen/internal/app/gen-docs/docs"
	"github.com/freepik-company/gommitizen/internal/pkg/cmd"
)

func main() {
	root := cmd.Root()

	err := docs.GenMarkdown(root, os.Stdout)
	if err != nil {
		os.Exit(1)
	}
}

package main

import (
	"log/slog"
	"os"

	"github.com/freepik-company/gommitizen/internal/pkg/cmd"
	"github.com/freepik-company/gommitizen/pkg/docs"
)

const (
	templateFile = "README.md.template"
)

func main() {
	root := cmd.Root()

	err := docs.GenMarkdown(root, os.Stdout, templateFile)
	if err != nil {
		slog.Error("Error generating markdown", "error", err)
		os.Exit(1)
	}
}

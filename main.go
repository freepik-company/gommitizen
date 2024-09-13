package main

import (
	"gommitizen/cmd"
	"gommitizen/internal/prettylog"
	"log/slog"
)

func main() {
	logger := slog.New(prettylog.NewHandler(&slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	cmd.Execute()
}

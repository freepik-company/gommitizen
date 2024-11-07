package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/prettylogconsole"
	"github.com/freepik-company/gommitizen/internal/version"
)

const (
	cmdRootDirPath = "directory"
	cmdRootDebug   = "debug"
)

func Root() *cobra.Command {
	var dirPath string
	var debug bool

	root := &cobra.Command{
		Use:     "gommitizen",
		Version: version.GetVersion(),
		Short:   "A commitizen implementation for Go with multi-project support",
		Long: `A commitizen implementation for Go with multi-project support.
It only supports the conventional commits specification: https://www.conventionalcommits.org/en/v1.0.0/
Currently it only supports the bump command, but it will support the commit command soon.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			dirPath, err = normalizePath(dirPath)
			if err != nil {
				slog.Error(fmt.Sprintf("normalising folders: %v", err))
				os.Exit(1)
			}

			level := slog.LevelInfo
			if debug {
				level = slog.LevelDebug
			}

			logger := slog.New(prettylogconsole.NewHandler(&slog.HandlerOptions{
				AddSource: false,
				Level:     level,
			}))
			slog.SetDefault(logger)
		},
	}

	root.PersistentFlags().StringVarP(&dirPath, "directory", "d", "", "Select a directory to run the command")
	root.PersistentFlags().BoolVar(&debug, cmdRootDebug, false, "Enable debug")

	root.AddCommand(initCmd())
	root.AddCommand(bumpCmd())
	root.AddCommand(getCmd())

	return root
}

func normalizePath(dirPath string) (string, error) {
	if len(dirPath) > 0 {
		if isRelativeDirPath(dirPath) {
			return toAbsoluteDirPath(dirPath)
		} else {
			return dirPath, nil
		}
	}
	return getCurrentDirPath()
}

func getCurrentDirPath() (string, error) {
	dirPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current path: %v", err)
	}
	return dirPath, nil
}

func toAbsoluteDirPath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", fmt.Errorf("error converting to absolute path: %v", err)
	}
	return absPath, nil
}

func isRelativeDirPath(path string) bool {
	return !filepath.IsAbs(path)
}

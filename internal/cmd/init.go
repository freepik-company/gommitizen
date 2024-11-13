package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/config"
	"github.com/freepik-company/gommitizen/internal/git"
)

const (
	cmdInitBumpChangelog = "bump-changelog"
)

func initCmd() *cobra.Command {
	var prefix string
	var updateChangelogOnBump bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Start a repository to use gommitizen",
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(cmdRootDirPath).Value.String()
			initRun(dirPath, prefix, updateChangelogOnBump)
		},
	}

	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Select a prefix for the version file")
	cmd.Flags().BoolVar(&updateChangelogOnBump, cmdInitBumpChangelog, false, "Update changelog on bump")

	return cmd
}

func initRun(dirPath, prefix string, updateChangelogOnBump bool) {
	commit, err := git.GetFirstCommit()
	if err != nil {
		slog.Error(fmt.Sprintf("first commit: %v", err))
		os.Exit(1)
	}

	config := config.NewConfigVersion(dirPath, "0.0.0", commit, prefix)
	config.UpdateChangelogOnBump = updateChangelogOnBump

	err = config.Save()
	if err != nil {
		slog.Error(fmt.Sprintf("config: %v", err))
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("Initializing gommitizen in %s", config.GetFilePath()))
}

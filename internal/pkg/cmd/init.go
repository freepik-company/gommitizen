package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/app/gommitizen/config"
	"github.com/freepik-company/gommitizen/internal/app/gommitizen/git"
)

const (
	initBumpFlagName = "bump-changelog"
)

func initCmd() *cobra.Command {
	var alias string
	var updateChangelogOnBump bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Start a repository to use gommitizen",
		Long: `Initialize the repository to use gommitizen. It will create a file with the version of the project and 
the first commit of the project.`,
		Example: "# To initialize the versioning of a project, run: \n" +
			"gommitizen init -d <directory> -p <prefix>`\n" +
			"# This will create a .version.json file in the given directory with the version 0.0.0.",
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(rootDirPathFlagName).Value.String()
			initRun(dirPath, alias, updateChangelogOnBump)
		},
	}

	cmd.Flags().StringVarP(&alias, "alias", "a", "", "Set a alias for the version file")
	cmd.Flags().BoolVar(&updateChangelogOnBump, initBumpFlagName, false, "Update changelog on bump")

	return cmd
}

func initRun(dirPath, alias string, updateChangelogOnBump bool) {
	commit, err := git.GetFirstCommit()
	if err != nil {
		slog.Error(fmt.Sprintf("first commit: %v", err))
		os.Exit(1)
	}

	config := config.NewConfigVersion(dirPath, "0.0.0", commit, alias)
	config.UpdateChangelogOnBump = updateChangelogOnBump

	err = config.Save()
	if err != nil {
		slog.Error(fmt.Sprintf("config: %v", err))
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("Initializing gommitizen in %s", config.GetFilePath()))
}

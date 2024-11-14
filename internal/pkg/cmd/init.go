package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/app/gommitizen/config"
	"github.com/freepik-company/gommitizen/internal/app/gommitizen/git"
)

func initCmd() *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Start a repository to use gommitizen",
		Long: `Initialize the repository to use gommitizen. It will create a file with the version of the project and 
the first commit of the project.`,
		Example: "To initialize the versioning of a project, run: \n" +
			"```bash\n" +
			"gommitizen init -d <directory> -p <prefix>`\n" +
			"```\n" +
			"This will create a .version.json file in the given directory with the version 0.0.0.",
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(cmdRootDirPath).Value.String()
			initRun(dirPath, prefix)
		},
	}

	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Select a prefix for the version file")

	return cmd
}

func initRun(dirPath, prefix string) {
	commit, err := git.GetFirstCommit()
	if err != nil {
		slog.Error(fmt.Sprintf("first commit: %v", err))
		os.Exit(1)
	}

	config := config.NewConfigVersion(dirPath, "0.0.0", commit, prefix)
	err = config.Save()
	if err != nil {
		slog.Error(fmt.Sprintf("config: %v", err))
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("Initializing gommitizen in %s", config.GetFilePath()))
}
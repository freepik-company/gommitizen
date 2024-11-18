package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	configObj "github.com/freepik-company/gommitizen/internal/app/gommitizen/config"
)

const (
	getPrefixFlagName = "prefix"
	getOutputFlagName = "output"
)

func getCmd() *cobra.Command {
	var prefix, output string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Give a list of projects, their versions and other information",
		Long: `Show information about the projects in the repository. It can show the version, the prefix, the commit 
information and all the information saved in the config file.`,
		Example: "To show all information in yaml format, run:\n" +
			"```bash\n" +
			"gommitizen get all -o yaml\n" +
			"```\n" +
			"To show the version of the projects in plain format, run:\n" +
			"```bash\n" +
			"gommitizen get version -o plain\n" +
			"```\n" +
			"or just:\n" +
			"```bash\n" +
			"gommitizen get version\n" +
			"```\n",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if output != "json" && output != "yaml" && output != "plain" {
				return fmt.Errorf("invalid output format: %s, supported values: json, yaml, plain", output)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				return
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&output, getOutputFlagName, "o", "plain", "Select the output format {json, yaml, plain}")
	cmd.PersistentFlags().StringVarP(&prefix, getPrefixFlagName, "p", "", "A prefix to look for a project to show information")

	cmd.AddCommand(getAllCmd())
	cmd.AddCommand(getVersionCmd())
	cmd.AddCommand(getPrefixCmd())
	cmd.AddCommand(getCommitCmd())

	return cmd
}

func getAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Get all projects information",
		Long: `Get all the information of the projects in the repository. It will show the version, the prefix, the commit
information and all the information saved in the config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(rootDirPathFlagName).Value.String()
			prefix := cmd.Parent().Flag(getPrefixFlagName).Value.String()
			output := cmd.Parent().Flag(getOutputFlagName).Value.String()
			projectsRun(dirPath, prefix, output, nil)
		},
	}
}

func getVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get the version of the projects",
		Long:  `Get the version of the projects in the repository. It will show the version of the projects and the prefix.`,
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(rootDirPathFlagName).Value.String()
			prefix := cmd.Parent().Flag(getPrefixFlagName).Value.String()
			output := cmd.Parent().Flag(getOutputFlagName).Value.String()
			projectsRun(dirPath, prefix, output, []string{"Version", "TagPrefix"})
		},
	}
}

func getPrefixCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prefix",
		Short: "Get the prefix of the projects",
		Long:  `Get the prefix of the projects in the repository.`,
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(rootDirPathFlagName).Value.String()
			prefix := cmd.Parent().Flag(getPrefixFlagName).Value.String()
			output := cmd.Parent().Flag(getOutputFlagName).Value.String()
			projectsRun(dirPath, prefix, output, []string{"TagPrefix"})
		},
	}
}

func getCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commit",
		Short: "Get the commit information of the projects",
		Long:  `Get the commit information of the projects in the repository.`,
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(rootDirPathFlagName).Value.String()
			prefix := cmd.Parent().Flag(getPrefixFlagName).Value.String()
			output := cmd.Parent().Flag(getOutputFlagName).Value.String()
			projectsRun(dirPath, prefix, output, []string{"Commit", "TagPrefix"})
		},
	}
}

func projectsRun(dirPath string, prefix string, output string, filter []string) {
	var configVersionPaths []string
	var err error

	if prefix == "" {
		configVersionPaths, err = configObj.FindConfigVersionFilePath(dirPath)
		if err != nil {
			slog.Error(fmt.Sprintf("finding config version file path: %v", err))
			os.Exit(1)
		}
	} else {
		configVersionPaths, err = configObj.FindConfigVersionFilePathByPrefix(dirPath, prefix)
		if err != nil {
			slog.Error(fmt.Sprintf("finding config version file path by prefix: %v", err))
			os.Exit(1)
		}
	}

	if len(configVersionPaths) == 0 {
		slog.Info("No projects found")
		os.Exit(0)
	}

	var configVersions []*configObj.ConfigVersion
	for _, configVersionPath := range configVersionPaths {
		configVersionFile, err := configObj.ReadConfigVersion(configVersionPath)
		if err != nil {
			slog.Error(fmt.Sprintf("reading configVersionFile version: %v", err))
			continue
		}
		configVersions = append(configVersions, configVersionFile)
	}

	str, err := configObj.PrintConfigVersions(configVersions, filter, output)
	if err != nil {
		slog.Error(fmt.Sprintf("printing config versions: %v", err))
		os.Exit(1)
	}
	slog.Info(str)
}

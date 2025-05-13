package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/app/gommitizen/config"
)

const (
	getAliasFlagName  = "alias"
	getOutputFlagName = "output"
)

func getCmd() *cobra.Command {
	var alias, output string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Give a list of projects, their versions and other information",
		Long: `Show information about the projects in the repository. It can show the version, the alias, the commit 
information and all the information saved in the config file.`,
		Example: "# To show all information in yaml format, run:\n" +
			"gommitizen get all -o yaml\n" +
			"# To show the version of the projects in plain format, run:\n" +
			"gommitizen get version -o plain\n" +
			"# or just:\n" +
			"gommitizen get version\n",
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

	cmd.PersistentFlags().StringVarP(&output, getOutputFlagName, "o", "plain", "select the output format {json, yaml, plain}")
	cmd.PersistentFlags().StringVarP(&alias, getAliasFlagName, "a", "", "a alias to look for a project to show information")

	cmd.AddCommand(getAllCmd())
	cmd.AddCommand(getVersionCmd())
	cmd.AddCommand(getAliasCmd())
	cmd.AddCommand(getCommitCmd())

	return cmd
}

func getAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Get all projects information",
		Long: `Get all the information of the projects in the repository. It will show the version, the alias, the commit
information and all the information saved in the config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(rootDirPathFlagName).Value.String()
			alias := cmd.Parent().Flag(getAliasFlagName).Value.String()
			output := cmd.Parent().Flag(getOutputFlagName).Value.String()
			projectsRun(dirPath, alias, output, nil)
		},
	}
}

func getVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get the version of the projects",
		Long:  `Get the version of the projects in the repository. It will show the version of the projects and the alias.`,
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(rootDirPathFlagName).Value.String()
			alias := cmd.Parent().Flag(getAliasFlagName).Value.String()
			output := cmd.Parent().Flag(getOutputFlagName).Value.String()
			projectsRun(dirPath, alias, output, []string{"Version", "Alias"})
		},
	}
}

func getAliasCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "alias",
		Short: "Get the alias of the projects",
		Long:  `Get the alias of the projects in the repository.`,
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(rootDirPathFlagName).Value.String()
			alias := cmd.Parent().Flag(getAliasFlagName).Value.String()
			output := cmd.Parent().Flag(getOutputFlagName).Value.String()
			projectsRun(dirPath, alias, output, []string{"Alias"})
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
			alias := cmd.Parent().Flag(getAliasFlagName).Value.String()
			output := cmd.Parent().Flag(getOutputFlagName).Value.String()
			projectsRun(dirPath, alias, output, []string{"Commit", "Alias"})
		},
	}
}

func projectsRun(dirPath string, alias string, output string, filter []string) {
	var configVersionPaths []string
	var err error

	if alias == "" {
		configVersionPaths, err = config.FindConfigVersionFilePath(dirPath)
		if err != nil {
			slog.Error(fmt.Sprintf("finding config version file path: %v", err))
			os.Exit(1)
		}
	} else {
		configVersionPaths, err = config.FindConfigVersionFilePathByAlias(dirPath, alias)
		if err != nil {
			slog.Error(fmt.Sprintf("finding config version file path by alias: %v", err))
			os.Exit(1)
		}
	}

	if len(configVersionPaths) == 0 {
		slog.Info("No projects found")
		os.Exit(0)
	}

	var configVersions []*config.ConfigVersion
	for _, configVersionPath := range configVersionPaths {
		configVersionFile, err := config.ReadConfigVersion(configVersionPath)
		if err != nil {
			slog.Error(fmt.Sprintf("reading configVersionFile version: %v", err))
			continue
		}
		configVersions = append(configVersions, configVersionFile)
	}

	str, err := config.PrintConfigVersions(configVersions, filter, output)
	if err != nil {
		slog.Error(fmt.Sprintf("printing config versions: %v", err))
		os.Exit(1)
	}

	// Print directly to stdout for structured formats to allow piping to tools like yq
	if output == "json" || output == "yaml" {
		fmt.Println(str)
	} else {
		slog.Info(str)
	}
}

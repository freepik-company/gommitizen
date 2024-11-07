package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/config"
)

const (
	cmdGetPrefix = "prefix"
	cmdGetOutput = "output"
)

func getCmd() *cobra.Command {
	var prefix, output string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Give a list of projects, their versions and other information",
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

	cmd.PersistentFlags().StringVarP(&output, cmdGetOutput, "o", "plain", "select the output format {json, yaml, plain}")
	cmd.PersistentFlags().StringVarP(&prefix, cmdGetPrefix, "p", "", "select a prefix to look for projects. Don't use with --directory")

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
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(cmdRootDirPath).Value.String()
			prefix := cmd.Parent().Flag(cmdGetPrefix).Value.String()
			output := cmd.Parent().Flag(cmdGetOutput).Value.String()
			projectsRun(dirPath, prefix, output, nil)
		},
	}
}

func getVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get the version of the projects",
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(cmdRootDirPath).Value.String()
			prefix := cmd.Parent().Flag(cmdGetPrefix).Value.String()
			output := cmd.Parent().Flag(cmdGetOutput).Value.String()
			projectsRun(dirPath, prefix, output, []string{"Version", "TagPrefix"})
		},
	}
}

func getPrefixCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prefix",
		Short: "Get the prefix of the projects",
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(cmdRootDirPath).Value.String()
			prefix := cmd.Parent().Flag(cmdGetPrefix).Value.String()
			output := cmd.Parent().Flag(cmdGetOutput).Value.String()
			projectsRun(dirPath, prefix, output, []string{"TagPrefix"})
		},
	}
}

func getCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commit",
		Short: "Get the commit information of the projects",
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(cmdRootDirPath).Value.String()
			prefix := cmd.Parent().Flag(cmdGetPrefix).Value.String()
			output := cmd.Parent().Flag(cmdGetOutput).Value.String()
			projectsRun(dirPath, prefix, output, []string{"Commit", "TagPrefix"})
		},
	}
}

func projectsRun(dirPath string, prefix string, output string, filter []string) {
	var configVersionPaths []string
	var err error

	if prefix == "" {
		configVersionPaths, err = config.FindConfigVersionFilePath(dirPath)
		if err != nil {
			slog.Error(fmt.Sprintf("finding config version file path: %v", err))
			os.Exit(1)
		}
	} else {
		configVersionPaths, err = config.FindConfigVersionFilePathByPrefix(dirPath, prefix)
		if err != nil {
			slog.Error(fmt.Sprintf("finding config version file path by prefix: %v", err))
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
	slog.Info(str)
}

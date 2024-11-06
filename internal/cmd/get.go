package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/config"
)

type projectsOpts struct {
	directory     string
	projectPrefix string
	outputFormat  string
}

var opts = projectsOpts{}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Give a list of projects, their versions and other information",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				return
			}
		},
	}

	cmd.AddCommand(getAllCmd())
	cmd.AddCommand(getVersionCmd())
	cmd.AddCommand(getPrefixCmd())
	cmd.AddCommand(getCommitCmd())

	cmd.PersistentFlags().StringVarP(&opts.outputFormat, "output", "o", "plain", "select the output format {json, yaml, plain}")
	cmd.PersistentFlags().StringVarP(&opts.directory, "directory", "d", "", "select a project directory to retrieve the project information")
	cmd.PersistentFlags().StringVarP(&opts.projectPrefix, "prefix", "p", "", "select a prefix to look for projects. Don't use with --directory")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if opts.outputFormat != "json" && opts.outputFormat != "yaml" && opts.outputFormat != "plain" {
			return fmt.Errorf("invalid output format: %s, supported values: json, yaml, plain", opts.outputFormat)
		}

		return nil
	}

	return cmd
}

func getAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Get all projects information",
		Run: func(cmd *cobra.Command, args []string) {
			projectsRun(opts, nil)
		},
	}
}

func getVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get the version of the projects",
		Run: func(cmd *cobra.Command, args []string) {
			projectsRun(opts, []string{"Version", "TagPrefix"})
		},
	}
}

func getPrefixCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prefix",
		Short: "Get the prefix of the projects",
		Run: func(cmd *cobra.Command, args []string) {
			projectsRun(opts, []string{"TagPrefix"})
		},
	}
}

func getCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commit",
		Short: "Get the commit information of the projects",
		Run: func(cmd *cobra.Command, args []string) {
			projectsRun(opts, []string{"Commit", "TagPrefix"})
		},
	}
}

func projectsRun(opts projectsOpts, filter []string) {
	var configVersionPaths []string

	if opts.projectPrefix == "" {
		nDirPath, err := config.NormalizePath(opts.directory)
		if err != nil {
			slog.Error(fmt.Sprintf("normalising folders: %v", err))
			os.Exit(1)
		}

		configVersionPaths, err = config.FindConfigVersionFilePath(nDirPath)
		if err != nil {
			slog.Error(fmt.Sprintf("finding config version file path: %v", err))
			os.Exit(1)
		}
	} else {
		nDirPath, err := config.NormalizePath(opts.directory)
		if err != nil {
			slog.Error(fmt.Sprintf("normalising folders: %v", err))
			os.Exit(1)
		}

		configVersionPaths, err = config.FindConfigVersionFilePathByPrefix(nDirPath, opts.projectPrefix)
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

	str, err := config.PrintConfigVersions(configVersions, filter, opts.outputFormat)
	if err != nil {
		slog.Error(fmt.Sprintf("printing config versions: %v", err))
		os.Exit(1)
	}
	slog.Info(str)
}

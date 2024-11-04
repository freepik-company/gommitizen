package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/config"
)

type projectsOpts struct {
	directory       string
	projectPrefix   string
	showVersionOnly bool
	showProjectPath bool
	outputFormat    string
}

func projectsCmd() *cobra.Command {
	var opts = projectsOpts{}

	cmd := &cobra.Command{
		Use:   "projects",
		Short: "Give a list of projects, their versions and other information",
		Run: func(cmd *cobra.Command, args []string) {
			projectsRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.directory, "directory", "d", "", "select a project directory to retrieve the project information")
	cmd.Flags().StringVarP(&opts.projectPrefix, "prefix", "p", "", "select a prefix to look for projects. Don't use with --directory")
	cmd.Flags().BoolVarP(&opts.showVersionOnly, "version", "v", false, "show only the version of the projects")
	cmd.Flags().BoolVarP(&opts.showProjectPath, "path", "P", false, "show only the path of the projects")
	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "plain", "select the output format {json, yaml, plain}")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if opts.directory != "" && opts.projectPrefix != "" {
			return fmt.Errorf("flags --directory and --prefix are mutually exclusive")
		}

		if opts.outputFormat != "json" && opts.outputFormat != "yaml" && opts.outputFormat != "plain" {
			return fmt.Errorf("invalid output format: %s, supported values: json, yaml, plain", opts.outputFormat)
		}

		if (opts.showVersionOnly || opts.showProjectPath) && opts.outputFormat != "plain" {
			return fmt.Errorf("flags --version and --path can only be used with plain format")
		}

		return nil
	}

	return cmd
}

func projectsRun(opts projectsOpts) {
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

	var printOption config.PrintPlainOption
	var err error
	if opts.outputFormat == "plain" {
		if opts.showVersionOnly {
			printOption = config.PrintVersionOnly
		} else if opts.showProjectPath {
			printOption = config.PrintPathOnly
		} else {
			printOption = config.PrintAll
		}
		err = config.PrintConfigVersionsPlain(configVersions, printOption)
	} else {
		err = config.PrintConfigVersions(configVersions, opts.outputFormat)
	}

	if err != nil {
		slog.Error(fmt.Sprintf("printing config versions: %v", err))
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"
	"gommitizen/internal/bumpmanager"
	"gommitizen/internal/cmdgit"
	"gommitizen/internal/config"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	validIncrements = []string{"MAJOR", "MINOR", "PATCH"}
)

type bumpOpts struct {
	directory     string
	changelog     bool
	incrementType string
}

func Bump() *cobra.Command {
	opts := bumpOpts{}

	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Make a version bump",
		Run: func(cmd *cobra.Command, args []string) {
			bumpRun(opts.directory, opts.changelog, strings.ToLower(opts.incrementType))
		},
	}

	cmd.Flags().StringVarP(&opts.directory, "directory", "d", "", "Select a project directory to bump")
	cmd.Flags().BoolVarP(&opts.changelog, "changelog", "c", false, "Create CHANGELOG.md")
	cmd.Flags().StringVarP(&opts.incrementType, "increment", "i", "", "Increment version (MAJOR, MINOR, PATCH)")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		increment, _ := cmd.Flags().GetString("increment")
		if increment == "" {
			return nil
		}
		for _, valid := range validIncrements {
			if increment == valid {
				return nil
			}
		}
		return fmt.Errorf(
			"invalid increment value: %s, supported values: %s",
			increment,
			strings.Join(validIncrements, ", "),
		)
	}

	return cmd
}

func bumpRun(dirPath string, changelog bool, incrementType string) {
	if incrementType != "" {
		slog.Info(fmt.Sprintf("Bumping version with increment: %s", incrementType))
	}

	nDirPath, err := config.NormalizePath(dirPath)
	if err != nil {
		slog.Error(fmt.Sprintf("normalising folders: %v", err))
		os.Exit(1)
	}

	configVersionPaths, err := config.FindConfigVersionFilePath(nDirPath)
	if err != nil {
		slog.Error(fmt.Sprintf("find config version paths: %v", err))
		os.Exit(1)
	}

	allModifiedFiles := make([]string, 0)
	allTagVersions := make([]string, 0)
	for _, configVersionPath := range configVersionPaths {
		modifiedFiles, tagVersion, err := bumpByConfig(configVersionPath, changelog, incrementType)
		if err != nil {
			slog.Error(fmt.Sprintf("bump by config: %v", err))
			os.Exit(1)
		}
		if len(modifiedFiles) > 0 {
			allModifiedFiles = append(allModifiedFiles, modifiedFiles...)
			allTagVersions = append(allTagVersions, tagVersion)
		}
	}

	output, err := bumpmanager.BumpCommitAll(allModifiedFiles, allTagVersions)
	if err != nil {
		slog.Error(fmt.Sprintf("bump commit all: %v", err))
		os.Exit(1)
	}

	slog.Info(strings.Join(output, "\n"))
}

func bumpByConfig(configVersionPath string, changelog bool, incrementType string) ([]string, string, error) {
	config, err := config.ReadConfigVersion(configVersionPath)
	if err != nil {
		slog.Info(fmt.Sprintf("Skipping file: %s, %v", configVersionPath, err))
		return []string{}, "", nil
	}

	modifiedFiles := make([]string, 0)
	tagVersion := config.GetTagVersion()

	slog.Info(fmt.Sprintf("Running bump in project %s", config.GetDirPath()))

	commitMessages, err := cmdgit.GetCommitMessages(config.Commit, config.GetDirPath())
	if err != nil {
		return []string{}, "", fmt.Errorf("commit messages: %s", err)
	}
	if incrementType == "" {
		incrementType = bumpmanager.DetermineVersionBump(commitMessages)
	}

	// If the file has been modified, update the version
	if incrementType != "none" {
		newVersion, newVersionStr, err := bumpmanager.IncrementVersion(config.Version, incrementType)
		if err != nil {
			return []string{}, "", fmt.Errorf("increment version: %s", err)
		}

		lastCommit, err := cmdgit.GetLastCommit()
		if err != nil {
			return []string{}, "", fmt.Errorf("last commit: %s", err)
		}

		slog.Info(fmt.Sprintf("%s change, %s -> %s", newVersionStr, config.Version, newVersion))

		modifiedFiles, err = config.UpdateVersion(newVersion, lastCommit)
		if err != nil {
			return []string{}, "", fmt.Errorf("update version: %s", err)
		}

		if changelog {
			// Update the CHANGELOG.md file
			// modifiedFiles = append(modifiedFiles, "ruta changelog")
			slog.Debug("TODO")
		}

		slog.Info("Commit messages:")
		for _, message := range commitMessages {
			slog.Info(fmt.Sprintf(" - %s", message))
		}

		slog.Info("Updated files:")
		for _, file := range modifiedFiles {
			slog.Info(fmt.Sprintf(" - %s", file))
		}

		tagVersion = config.GetTagVersion()
		slog.Info("New tags: " + tagVersion)

		slog.Info(fmt.Sprintf("Updated version in %s", config.GetDirPath()))
	} else {
		slog.Info(fmt.Sprintf("Bump skipped in %s", config.GetDirPath()))
	}
	slog.Info("---")
	return modifiedFiles, tagVersion, nil
}

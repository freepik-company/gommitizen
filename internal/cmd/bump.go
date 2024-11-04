package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/bumpmanager"
	"github.com/freepik-company/gommitizen/internal/changelog"
	"github.com/freepik-company/gommitizen/internal/config"
	"github.com/freepik-company/gommitizen/internal/conventionalcommits"
	"github.com/freepik-company/gommitizen/internal/git"
)

var (
	validIncrements = []string{"MAJOR", "MINOR", "PATCH"}
)

type bumpOpts struct {
	directory       string
	createChangelog bool
	incrementType   string
}

func bumpCmd() *cobra.Command {
	opts := bumpOpts{}

	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Make a version bump",
		Run: func(cmd *cobra.Command, args []string) {
			bumpRun(opts.directory, opts.createChangelog, strings.ToLower(opts.incrementType))
		},
	}

	cmd.Flags().StringVarP(&opts.directory, "directory", "d", "", "select a project directory to bump")
	cmd.Flags().BoolVarP(&opts.createChangelog, "changelog", "c", false, "generate the changelog for the newest version")
	cmd.Flags().StringVarP(&opts.incrementType, "increment", "i", "", "manually specify the desired increment {MAYOR, MINOR, PATCH}")

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

func bumpRun(dirPath string, createChangelog bool, incrementType string) {
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
		modifiedFiles, tagVersion, err := bumpByConfig(configVersionPath, createChangelog, incrementType)
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

func bumpByConfig(configVersionPath string, createChangelog bool, incrementType string) ([]string, string, error) {
	config, err := config.ReadConfigVersion(configVersionPath)
	if err != nil {
		slog.Info(fmt.Sprintf("Skipping file: %s, %v", configVersionPath, err))
		return []string{}, "", nil
	}

	modifiedFiles := make([]string, 0)
	tagVersion := config.GetTagVersion()

	slog.Info(fmt.Sprintf("Running bump in project %s", config.GetDirPath()))

	gitCommits, err := git.GetCommits(config.Commit, config.GetDirPath())
	cvCommits := conventionalcommits.ReadConventionalCommits(gitCommits)
	if err != nil {
		return []string{}, "", fmt.Errorf("commit messages: %s", err)
	}
	if incrementType == "" {
		incrementType = conventionalcommits.DetermineIncrementType(cvCommits)
	}

	// If the file has been modified, update the version
	if incrementType != "none" {
		newVersion, newVersionStr, err := bumpmanager.IncrementVersion(config.Version, incrementType)
		if err != nil {
			return []string{}, "", fmt.Errorf("increment version: %s", err)
		}

		lastCommit, err := git.GetLastCommit()
		if err != nil {
			return []string{}, "", fmt.Errorf("last commit: %s", err)
		}

		slog.Info(fmt.Sprintf("%s change, %s -> %s", newVersionStr, config.Version, newVersion))

		modifiedFiles, err = config.UpdateVersion(newVersion, lastCommit)
		if err != nil {
			return []string{}, "", fmt.Errorf("update version: %s", err)
		}

		if createChangelog {
			changelogFilePath, err := changelog.Apply(config.GetDirPath(), config.Version, cvCommits)
			if err != nil {
				return []string{}, "", fmt.Errorf("update changelog: %s", err)
			}
			modifiedFiles = append(modifiedFiles, changelogFilePath)
		}

		slog.Info("Commit messages:")
		for _, commit := range cvCommits {
			slog.Info(fmt.Sprintf(" - %s", commit))
		}

		slog.Info("Updated files:")
		for _, file := range modifiedFiles {
			slog.Info(fmt.Sprintf(" - %s", file))
		}

		tagVersion = config.GetTagVersion()
		slog.Info("New tags: " + tagVersion)

		slog.Info(fmt.Sprintf("Updated version in %s", config.GetDirPath()))
	} else {
		slog.Info(fmt.Sprintf("bump skipped in %s", config.GetDirPath()))
	}
	slog.Info("---")
	return modifiedFiles, tagVersion, nil
}

package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/app/gommitizen/bumpmanager"
	"github.com/freepik-company/gommitizen/internal/app/gommitizen/changelog"
	"github.com/freepik-company/gommitizen/internal/app/gommitizen/config"
	"github.com/freepik-company/gommitizen/internal/app/gommitizen/conventionalcommits"
	"github.com/freepik-company/gommitizen/internal/app/gommitizen/git"
)

func bumpCmd() *cobra.Command {
	var validIncrements = []string{"MAJOR", "MINOR", "PATCH"}
	var incrementType string
	var createChangelog bool

	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Make a version bump",
		Long:  `Increment the version of the project according to the conventional commits specification.`,
		Example: "To bump the version of a project, run:\n" +
			"```bash\n" +
			"gommitizen bump\n" +
			"```\n" +
			"This will bump the version of all projects in the current directory.\n" +
			"If you want to bump the version of a specific project, run:\n" +
			"```bash\n" +
			"gommitizen bump -d <directory>\n" +
			"```\n" +
			"This will bump the version of the project in the given directory.\n" +
			"If you want to bump the version of projects and generate a changelog, run:\n" +
			"```bash\n" +
			"gommitizen bump -c\n" +
			"```\n" +
			"This will bump the version of the projects and generate a changelog with the changes made since the last " +
			"version.\n" +
			"If you want to bump the version of project to a major version, run:\n" +
			"```bash\n" +
			"gommitizen bump -i major\n" +
			"```\n",
		PreRunE: func(cmd *cobra.Command, args []string) error {
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
		},
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cmd.Root().Flag(rootDirPathFlagName).Value.String()
			bumpRun(dirPath, createChangelog, strings.ToLower(incrementType))
		},
	}

	cmd.Flags().BoolVarP(&createChangelog, "changelog", "c", false, "generate the changelog for the newest version")
	cmd.Flags().StringVarP(&incrementType, "increment", "i", "", "manually specify the desired increment {MAYOR, MINOR, PATCH}")

	return cmd
}

func bumpRun(dirPath string, createChangelog bool, incrementType string) {
	if incrementType != "" {
		slog.Info(fmt.Sprintf("Bumping version with increment: %s", incrementType))
	}

	configVersionPaths, err := config.FindConfigVersionFilePath(dirPath)
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
	gitTag := config.GetGitTag()

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
		// Running pre-bump scripts
		err = config.RunHook("pre-bump")
		if err != nil {
			return []string{}, "", fmt.Errorf("pre bump scripts: %s", err)
		}

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

		// Running post-bump scripts
		err = config.RunHook("post-bump")
		if err != nil {
			return []string{}, "", fmt.Errorf("post bump scripts: %s", err)
		}

		if createChangelog {
			// Running pre-changelog scripts
			err = config.RunHook("pre-changelog")
			if err != nil {
				return []string{}, "", fmt.Errorf("pre changelog scripts: %s", err)
			}

			slog.Info("Generating changelog...")
			changelogFilePath, err := changelog.Apply(config.GetDirPath(), config.Version, cvCommits)
			if err != nil {
				return []string{}, "", fmt.Errorf("update changelog: %s", err)
			}
			modifiedFiles = append(modifiedFiles, changelogFilePath)

			// Running post-changelog scripts
			err = config.RunHook("post-changelog")
			if err != nil {
				return []string{}, "", fmt.Errorf("post changelog scripts: %s", err)
			}
		}

		slog.Info("Commit messages:")
		for _, commit := range cvCommits {
			slog.Info(fmt.Sprintf(" - %s", commit))
		}

		slog.Info("Updated files:")
		for _, file := range modifiedFiles {
			slog.Info(fmt.Sprintf(" - %s", file))
		}

		gitTag = config.GetGitTag()
		slog.Info("New tags: " + gitTag)

		slog.Info(fmt.Sprintf("Updated version in %s", config.GetDirPath()))
	} else {
		slog.Info(fmt.Sprintf("bump skipped in %s", config.GetDirPath()))
	}
	slog.Info("---")
	return modifiedFiles, gitTag, nil
}

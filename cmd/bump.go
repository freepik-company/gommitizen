package cmd

import (
	"fmt"
	"gommitizen/version"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var projectDir string
var changelog bool
var incrementType string
var validIncrements = []string{"MAJOR", "MINOR", "PATCH"}

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Make a version bump",
	Run: func(cmd *cobra.Command, args []string) {
		incrementType, _ = cmd.Flags().GetString("increment")
		if incrementType != "" {
			fmt.Printf("Bumping version with increment: %s\n", incrementType)
		}

		if projectDir == "" {
			fmt.Printf("\n# Run bump in all projects\n\n")
			bumpVersion()
			return
		} else {
			bumpProjectVersion(projectDir)
			return
		}
	},
}

func init() {
	bumpCmd.Flags().StringVarP(&projectDir, "directory", "d", "", "Select a project directory to bump")
	bumpCmd.Flags().BoolVarP(&changelog, "changelog", "c", false, "Create CHANGELOG.md")
	bumpCmd.Flags().StringP("increment", "i", "", "Increment version (MAJOR, MINOR, PATCH)")
	bumpCmd.PreRunE = validateIncrementArgs

	rootCmd.AddCommand(bumpCmd)
}

func validateIncrementArgs(cmd *cobra.Command, args []string) error {
	increment, _ := cmd.Flags().GetString("increment")

	if increment == "" {
		return nil
	}

	for _, valid := range validIncrements {
		if increment == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid increment value: %s", increment)
}

func bumpProjectVersion(project string) {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error obtaining current directory:", err)
		os.Exit(1)
	}

	filePath := filepath.Join(rootDir, project, ".version.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("The file %s does not exist\n", filePath)
		os.Exit(1)
	}

	if bumpRun(rootDir, filePath) != nil {
		fmt.Println("Error running bump:", err)
		os.Exit(1)
	}
}

// Run the bump command for all .version.json files in the current directory and its subdirectories
func bumpVersion() {
	// Get the current directory
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error obtaining current directory:", err)
		os.Exit(1)
	}

	// Find all .version.json files in the current directory and its subdirectories
	fileList, err := version.FindFCVersionFiles(rootDir)
	if err != nil {
		fmt.Println("Error finding .version.json files:", err)
		os.Exit(1)
	}

	if len(fileList) == 0 {
		fmt.Println("Files .version.json not found")
		os.Exit(1)
	}

	// Loop over the found files
	for _, filePath := range fileList {
		err := bumpRun(rootDir, filePath)
		if err != nil {
			fmt.Println("Error running bump:", err)
			continue
		}
	}
}

// Run the bump command for a .version.json file
func bumpRun(rootDir string, filePath string) error {
	var err error
	var relativePath string

	// Get the relative path to the current directory
	relativePath, err = filepath.Rel(rootDir, filePath)
	if err != nil {
		return fmt.Errorf("Error obtaining relative path: %s", err)
	}

	// Print the start message
	fmt.Printf("## Running bump in project %s\n\n", filepath.Dir(relativePath))

	// Read the version data
	config := version.LoadVersionData(filePath)
	err = config.RetrieveRepositoryData()
	if err != nil {
		return fmt.Errorf("Error retrieving repository data: %s", err)
	}

	config.SetUpdateChangelog(changelog) // Set the update changelog flag (default: false)

	// Check if files have been modified in Git
	modified, err := config.IsSomeFileModified()
	if err != nil {
		return fmt.Errorf("Error checking if some file has been modified in Git: %s", err)
	}

	// If the file has been modified, update the version
	if modified {
		currentVersion := config.GetVersion()
		newVersion := currentVersion
		if incrementType != "" {
			newVersion, err = config.IncrementVersion(incrementType)
			if err != nil {
				return fmt.Errorf("Error incrementing version: %s", err)
			}
		} else {
			newVersion, err = config.UpdateVersion()
			if err != nil {
				return fmt.Errorf("Error updating version: %s", err)
			}
		}

		if newVersion == currentVersion {
			fmt.Printf("There is no update of config in %s\n", relativePath)
		} else {
			fmt.Printf("Updated version in %s\n", relativePath)
		}
	} else {
		fmt.Printf("Bump skipped in %s\n", relativePath)
	}
	fmt.Printf("\n")

	return nil
}

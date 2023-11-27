package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gommitizen/internal"
	"gommitizen/version"
	"strings"
)

var directory, prefix string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Start a repository to use gommitizen",
	//	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(directory) == "" {
			directory = internal.GetCurrentDirectory()
		}

		initialize(directory, prefix)
	},
}

func init() {
	initCmd.Flags().StringVarP(&directory, "directory", "d", "", "Select a project directory to initialize")
	initCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Select a prefix for the version file")

	rootCmd.AddCommand(initCmd)
}

func initialize(path, prefix string) {
	// check if path ends with .version.json, if not, append it
	if !strings.HasSuffix(version.ConfigFileName, path) {
		path = path + "/" + version.ConfigFileName
	}

	config := version.NewVersionData(
		"0.0.0",
		version.DefaultCommit,
		path,
		prefix,
	)

	config.Save()

	fmt.Println("Initializing gommitizen in", config.GetFilePath())
}

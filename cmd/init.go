package cmd

import (
	"fmt"
	"gommitizen/version"
	"os"

	"github.com/spf13/cobra"
)

var directory string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Start a repository to use gommitizen",
	//	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if directory != "" {
			initialize(directory)
		} else {
			fmt.Println("A directory must be specified")
		}
	},
}

func init() {
	initCmd.Flags().StringVarP(&directory, "directory", "d", "", "Select a project directory to initialize")

	rootCmd.AddCommand(initCmd)
}

func initialize(path string) {
	// check directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(("The directory does not exist"))
		os.Exit(1)
	}

	config := version.NewVersionData()
	err := config.Initialize(path)
	if err != nil {
		fmt.Println("Error initializing repository:", err)
		os.Exit(1)
	}
}

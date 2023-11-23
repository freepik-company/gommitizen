package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gommitizen/version"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicia un repositorio para usar gommitizen",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initialize(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initialize(path string) {
	// check .version.json does not exist
	configFile := path + "/.version.json"
	if _, err := os.Stat(configFile); err == nil {
		fmt.Println("El repositorio ya está inicializado")
		os.Exit(1)
	}

	config := version.NewVersionData(
		"0.0.0",
		"",
		configFile,
	)

	config.Save()
}

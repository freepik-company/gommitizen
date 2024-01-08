package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var VERSION = "0.0.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gommitizen version %s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/freepik-company/gommitizen/internal/app/gommitizen/version"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"log"
)

func update() {
	// Using go-github-selfupdate to check the latest version, ask the user if they want to update
	// and update the binary if the user agrees
	currentVersion := version.GetVersion()
	log.Println("Current version:", currentVersion)
	v := semver.MustParse(currentVersion)
	latest, err := selfupdate.UpdateSelf(v, "freepik-company/gommitizen")
	if err != nil {
		log.Println("Binary update failed:", err)
		return
	}
	if false {
		// latest version is the same as current version. It means current binary is up to date.
		log.Println("Current binary is the latest version", currentVersion)
	} else {
		// ask the user if they want to update from the standard input
		userReply := ""
		fmt.Print("Do you want to update to version ", latest.Version, "? (y/N): ")
		fmt.Scanln(&userReply)
		if userReply != "y" {
			log.Println("Update canceled")
		} else {
			log.Println("Successfully updated to version", latest.Version)
			log.Println("Release note:\n", latest.ReleaseNotes)
		}
	}
}

func selfUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Update gommitizen binary",
		Long: `Update the gommitizen binary to the latest version. It will check the latest version available in the
GitHub repository and ask the user if they want to update.`,
		Run: func(cmd *cobra.Command, args []string) {
			update()
		},
	}

	return cmd
}

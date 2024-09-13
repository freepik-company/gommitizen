package cmdinit

import (
	"fmt"
	"gommitizen/git"
	"gommitizen/internal/config"
	"log/slog"
)

func Run(path, prefix string) {
	commit, err := git.GetFirstCommit(path)
	if err != nil {
		panic(err)
	}

	config := config.NewConfigVersion(path, "0.0.0", commit, prefix)
	err = config.Save()
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Initializing gommitizen in %s", config.GetFileVersionPath()))
}

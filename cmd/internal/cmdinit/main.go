package cmdinit

import (
	"fmt"
	"gommitizen/git"
	"gommitizen/internal/config"
	"log/slog"
)

func Run(path, prefix string) {
	nPath, err := config.NormalizePath(path)
	if err != nil {
		panic(err)
	}

	commit, err := git.GetFirstCommit(nPath)
	if err != nil {
		panic(err)
	}

	config := config.NewConfigVersion(nPath, "0.0.0", commit, prefix)
	err = config.Save()
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Initializing gommitizen in %s", config.GetFileVersionPath()))
}

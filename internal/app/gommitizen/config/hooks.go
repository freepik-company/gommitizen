package config

import (
	"fmt"
	"log/slog"
	"os/exec"
)

func runHook(v ConfigVersion, hookName string, cmdStr string) error {
	if len(cmdStr) == 0 {
		slog.Debug(fmt.Sprintf("hook %s is empty", hookName))
		return nil
	}

	slog.Debug(fmt.Sprintf("running hook %s: %s", hookName, cmdStr))

	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Env = append(cmd.Env,
		"GOMMITIZEN_DIRPATH="+v.GetDirPath(),
		"GOMMITIZEN_FILEPATH="+v.GetFilePath(),
		"GOMMITIZEN_GITTAG="+v.GetGitTag(),
		"GOMMITIZEN_ALIAS="+v.Alias,
		"GOMMITIZEN_VERSION="+v.Version,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("run hook %s: %v", hookName, err)
	}

	if len(output) > 0 {
		slog.Info(fmt.Sprintf("\n\033[32mHook %s output:\n%s\033[0m", hookName, string(output)))
	} else {
		slog.Info(fmt.Sprintf("\n\033[32mLaunch hook %s\n\033[0m", hookName))
	}

	return nil
}

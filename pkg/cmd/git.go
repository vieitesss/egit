package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func EnsureGitInstalled() {
	if _, err := exec.LookPath("git"); err != nil {
		panic(fmt.Sprintf("'git' is not installed: %s", err))
	}
}

func GitCmd(a string) error {
	EnsureGitInstalled()

	args := strings.Split(a, " ")
	return ExecCmd("git", args...)
}

func GitCmdOut(args ...string) (string, error) {
	EnsureGitInstalled()

	return ExecCmdOut("git", args...)
}

func IsGitRepo() bool {
	if err := GitCmd("rev-parse --git-dir"); err != nil {
		return false
	}

	return true
}


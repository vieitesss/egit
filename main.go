package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func execCmdOut(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func execCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func ensureGitInstalled() {
	if _, err := exec.LookPath("git"); err != nil {
		panic(fmt.Sprintf("'git' is not installed: %s", err))
	}
}

func gitCmd(a string) error {
	ensureGitInstalled()

	args := strings.Split(a, " ")
	return execCmd("git", args...)
}

func gitCmdOut(a string) (string, error) {
	ensureGitInstalled()

	args := strings.Split(a, " ")
	return execCmdOut("git", args...)
}

func isGitRepo() bool {
	if err := gitCmd("rev-parse --git-dir"); err != nil {
		return false
	}

	return true
}

func main() {
	out, err := gitCmdOut("status --porcelain")
	if err != nil {
		fmt.Errorf("error: %v", err)
	}
	fmt.Printf("%v", out)
}

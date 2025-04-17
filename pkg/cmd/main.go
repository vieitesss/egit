package cmd

import (
	"os/exec"
)

func ExecCmdOut(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func ExecCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

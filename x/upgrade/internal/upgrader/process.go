package upgrader

import (
	"os"
	"os/exec"
)

func restart() error {
	return StartProcess(os.Args[0], os.Args[1])
}

func StartProcess(path string, arg string) error {

	cmd := exec.Command(path, arg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Process.Release()
	return err
}

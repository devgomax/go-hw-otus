package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	for name, value := range env {
		if value.NeedRemove {
			if err := os.Unsetenv(name); err != nil {
				log.Fatal(err)
			}
			continue
		}

		if err := os.Setenv(name, value.Value); err != nil {
			log.Fatal(err)
		}
	}

	if err := command.Run(); err != nil {
		return command.ProcessState.ExitCode()
	}

	return 0
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("update environment variable", func(t *testing.T) {
		varName := "1"
		env := Environment{
			varName: EnvValue{Value: "test", NeedRemove: false},
		}

		err := os.Setenv(varName, varName)
		defer os.Unsetenv(varName)

		require.NoError(t, err)

		code := RunCmd([]string{"echo", "Hello World"}, env)
		require.Equal(t, 0, code)

		require.Equal(t, "test", os.Getenv(varName))
	})

	t.Run("remove environment variables", func(t *testing.T) {
		tests := []struct {
			name       string
			varName    string
			old        string
			new        string
			needRemove bool
		}{
			{
				name:       "with value and remove=true",
				varName:    "1",
				old:        "1",
				new:        "test",
				needRemove: true,
			},
			{
				name:       "without value and remove=true",
				varName:    "2",
				old:        "2",
				new:        "",
				needRemove: true,
			},
			{
				name:       "without value and remove=false",
				varName:    "3",
				old:        "3",
				new:        "",
				needRemove: false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := os.Setenv(tc.varName, tc.old)
				defer os.Unsetenv(tc.varName)

				require.NoError(t, err)

				code := RunCmd(
					[]string{"echo", "Hello World"},
					Environment{
						tc.varName: EnvValue{Value: tc.new, NeedRemove: tc.needRemove},
					},
				)

				require.Equal(t, 0, code)
				require.Zero(t, os.Getenv(tc.varName))
			})
		}
	})

	t.Run("returns the same exit code as child util", func(t *testing.T) {
		tests := []struct {
			exitCode int
		}{
			{exitCode: 2},
			{exitCode: 11},
			{exitCode: 43},
		}

		for _, tc := range tests {
			t.Run(fmt.Sprintf("code %v", tc.exitCode), func(t *testing.T) {
				args := []string{"sh", "-c", fmt.Sprintf("exit %v", tc.exitCode)}
				cmd := exec.Command(args[0], args[1:]...)

				err := cmd.Run()
				require.Error(t, err)
				require.Equal(t, tc.exitCode, cmd.ProcessState.ExitCode())

				code := RunCmd(args, nil)
				require.Equal(t, cmd.ProcessState.ExitCode(), code)
			})
		}
	})
}

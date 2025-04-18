package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const envData = "testdata/env/"

func TestReadDir(t *testing.T) {
	t.Run("child directories are not processed", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(envData, "")
		defer os.RemoveAll(tempDir)

		require.NoError(t, err)

		env, err := ReadDir(envData)
		require.NoError(t, err)

		require.NotContains(t, env, tempDir)
	})

	t.Run("all files are processed and set into the env map", func(t *testing.T) {
		env, err := ReadDir(envData)
		require.NoError(t, err)

		dirEntries, err := os.ReadDir(envData)
		require.NoError(t, err)

		for _, entry := range dirEntries {
			if entry.IsDir() {
				require.NotContains(t, env, entry.Name())
			} else {
				require.Contains(t, env, entry.Name())
			}
		}
	})

	t.Run(`files with "=" in the name are not processed`, func(t *testing.T) {
		tempFile, err := os.CreateTemp(envData, "=")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		require.NoError(t, err)

		env, err := ReadDir(envData)
		require.NoError(t, err)

		require.NotContains(t, env, strings.TrimPrefix(tempFile.Name(), envData))
	})
}

package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, dirEntry := range dirEntries {
		envValue, inErr := GetEnvValue(dir, dirEntry)
		if inErr != nil {
			return nil, inErr
		}

		if envValue == nil {
			continue
		}

		env[dirEntry.Name()] = *envValue
	}

	return env, nil
}

// GetEnvValue processes specified os.DirEntry and returns *EnvValue.
func GetEnvValue(dir string, dirEntry os.DirEntry) (*EnvValue, error) {
	if dirEntry.IsDir() {
		return nil, nil
	}

	if strings.Contains(dirEntry.Name(), "=") {
		return nil, nil
	}

	file, err := os.Open(path.Join(dir, dirEntry.Name()))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	line, err := reader.ReadString('\n')
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	line = strings.ReplaceAll(line, "\x00", "\n")
	line = strings.TrimRightFunc(line, unicode.IsSpace)

	if line == "" {
		return &EnvValue{Value: "", NeedRemove: true}, nil
	}

	return &EnvValue{Value: line}, nil
}

package main

import (
	"log"
	"os"
)

func main() {
	dir, cmd := os.Args[1], os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(RunCmd(cmd, env))
}

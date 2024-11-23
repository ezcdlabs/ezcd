/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/ezcdlabs/ezcd/cmd/cli/cmd"
)

var version = "dev" // default version

func main() {
	serviceLoader := NewEnvServiceLoader()
	command := cmd.NewRootCmd(version, serviceLoader)
	err := command.Execute()

	if err != nil {
		os.Exit(1)
	}
}

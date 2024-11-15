/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/ezcdlabs/ezcd/cmd/cli/cmd"

var version = "dev" // default version

func main() {
	cmd.Version = version
	cmd.Execute()
}

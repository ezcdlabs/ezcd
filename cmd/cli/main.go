/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ezcdlabs/ezcd/cmd/cli/cmd"
	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/ezcdlabs/ezcd/pkg/ezcd_postgres"
)

var version = "dev" // default version

func main() {
	serviceLoader := NewEnvServiceLoader()
	command := cmd.NewRootCmd(version, serviceLoader)
	err := command.Execute()

	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}

func NewEnvServiceLoader() envServiceLoader {
	return envServiceLoader{}
}

type envServiceLoader struct{}

func (envServiceLoader) Load() (ezcd.Ezcd, error) {
	ezcdDatabaseUrl := os.Getenv("EZCD_DATABASE_URL")

	if ezcdDatabaseUrl == "" {
		return nil, fmt.Errorf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	ezcdService := ezcd.NewEzcdService(ezcd_postgres.NewPostgresDatabase(ezcdDatabaseUrl))

	return ezcdService, nil
}

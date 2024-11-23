package main

import (
	"fmt"
	"os"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/ezcdlabs/ezcd/pkg/ezcd_postgres"
)

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

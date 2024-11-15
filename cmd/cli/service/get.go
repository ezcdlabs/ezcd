package service

import (
	"log"
	"os"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/ezcdlabs/ezcd/pkg/ezcd_postgres"
)

var ezcdService ezcd.Ezcd

func Get() ezcd.Ezcd {
	if ezcdService != nil {
		return ezcdService
	}

	ezcdDatabaseUrl := os.Getenv("EZCD_DATABASE_URL")

	if ezcdDatabaseUrl == "" {
		log.Fatalf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	ezcdService = ezcd.NewEzcdService(ezcd_postgres.NewPostgresDatabase(ezcdDatabaseUrl))

	return ezcdService
}

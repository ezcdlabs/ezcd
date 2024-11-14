package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/ezcdlabs/ezcd/pkg/ezcd_postgres"
)

var version = "dev" // default version

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ezcd-cli <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "hello":
		fmt.Println("Hello, CLI!")
	case "--version":
		fmt.Printf("Version: %s\n", version)

	case "create-project":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ezcd-cli create-project <name>")
			os.Exit(1)
		}

		ezcdDatabaseUrl := os.Getenv("EZCD_DATABASE_URL")
		if ezcdDatabaseUrl == "" {
			log.Fatalf("database connection string is required, please set EZCD_DATABASE_URL")
		}
		ezcdService := ezcd.NewEzcdService(ezcd_postgres.NewPostgresDatabase(ezcdDatabaseUrl))

		name := os.Args[2]
		project, err := ezcdService.CreateProject(name)
		if err != nil {
			fmt.Printf("Failed to create project: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", project.ID)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

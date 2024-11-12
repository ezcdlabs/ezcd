package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ezcdlabs/ezcd/cmd/server/public"
	"github.com/ezcdlabs/ezcd/pkg/ezcd"
)

var version = "dev" // default version

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("Version: %s\n", version)
		return
	}
	ezcdDatabaseUrl := os.Getenv("EZCD_DATABASE_URL")
	if ezcdDatabaseUrl == "" {
		log.Fatalf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	// TODO: remove me
	log.Printf("Using database url %s", ezcdDatabaseUrl)

	// health check that pings the database
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		type HealthResponse struct {
			Status  string `json:"status"`
			Message string `json:"message,omitempty"`
		}

		if _, err := ezcd.GetProjects(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(HealthResponse{Status: "unhealthy", Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthResponse{Status: "healthy"})
	})

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, API!!!")
	})

	http.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		projects, err := ezcd.GetProjects()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching projects: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(projects); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		projectID := strings.TrimPrefix(r.URL.Path, "/api/projects/")
		if projectID == "" {
			http.Error(w, "Project ID is required", http.StatusBadRequest)
			return
		}

		project, err := ezcd.GetProject(projectID)
		if err != nil {
			if errors.Is(err, ezcd.ErrProjectNotFound) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "Project not found"})
			} else {
				http.Error(w, fmt.Sprintf("Error fetching project: %v", err), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(project); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		}
	})

	if public.IsEmbedded {
		fileserver := http.FileServer(http.FS(public.Box))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				fullPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
				if _, err := fs.Stat(public.Box, fullPath); err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						r.URL.Path = "/"
					} else {
						http.Error(w, "Internal server error", http.StatusInternalServerError)
						return
					}
				}
			}
			fileserver.ServeHTTP(w, r)
		})
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "This is the ezcd server running in dev mode, no web frontend is embedded. Please run the web frontend separately.")
		})
	}

	port := ":3923"
	if envPort := os.Getenv("EZCD_PORT"); envPort != "" {
		port = ":" + envPort
	}

	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}

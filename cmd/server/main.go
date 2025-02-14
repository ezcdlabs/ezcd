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
	"time"

	"github.com/ezcdlabs/ezcd/cmd/server/public"
	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/ezcdlabs/ezcd/pkg/ezcd_postgres"
)

var version = "dev" // default version

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func mapToAnnotatedProject(p ezcd.Project) Project {
	return Project{
		ID:   p.ID,
		Name: p.Name,
	}
}

func mapToAnnotatedProjects(ps []ezcd.Project) []Project {
	projects := make([]Project, len(ps))
	for i, p := range ps {
		projects[i] = mapToAnnotatedProject(p)
	}
	return projects
}

type Commit struct {
	Hash                       string     `json:"hash"`
	AuthorName                 string     `json:"authorName"`
	AuthorEmail                string     `json:"authorEmail"`
	Message                    string     `json:"message"`
	Date                       time.Time  `json:"date"`
	CommitStageStartedAt       *time.Time `json:"commitStageStartedAt"`
	CommitStageCompletedAt     *time.Time `json:"commitStageCompletedAt"`
	CommitStageStatus          string     `json:"commitStageStatus"`
	AcceptanceStageStartedAt   *time.Time `json:"acceptanceStageStartedAt"`
	AcceptanceStageCompletedAt *time.Time `json:"acceptanceStageCompletedAt"`
	AcceptanceStageStatus      string     `json:"acceptanceStageStatus"`
	DeployStartedAt            *time.Time `json:"deployStartedAt"`
	DeployCompletedAt          *time.Time `json:"deployCompletedAt"`
	DeployStatus               string     `json:"deployStatus"`
	LeadTimeCompletedAt        *time.Time `json:"leadTimeCompletedAt"`
}

func mapToAnnotatedCommit(c ezcd.Commit) Commit {
	return Commit{
		Hash:                       c.Hash,
		Message:                    c.Message,
		Date:                       c.Date,
		AuthorName:                 c.AuthorName,
		AuthorEmail:                c.AuthorEmail,
		CommitStageStartedAt:       c.CommitStageStartedAt,
		CommitStageCompletedAt:     c.CommitStageCompletedAt,
		CommitStageStatus:          c.CommitStageStatus.String(),
		AcceptanceStageStartedAt:   c.AcceptanceStageStartedAt,
		AcceptanceStageCompletedAt: c.AcceptanceStageCompletedAt,
		AcceptanceStageStatus:      c.AcceptanceStageStatus.String(),
		DeployStartedAt:            c.DeployStartedAt,
		DeployCompletedAt:          c.DeployCompletedAt,
		DeployStatus:               c.DeployStatus.String(),
		LeadTimeCompletedAt:        c.LeadTimeCompletedAt,
	}
}

func mapToAnnotatedCommits(cs []ezcd.Commit) []Commit {
	commits := make([]Commit, len(cs))
	for i, c := range cs {
		commits[i] = mapToAnnotatedCommit(c)
	}
	return commits
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("Version: %s\n", version)
		return
	}
	ezcdDatabaseUrl := os.Getenv("EZCD_DATABASE_URL")
	if ezcdDatabaseUrl == "" {
		log.Fatalf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	ezcdService := ezcd.NewEzcdService(ezcd_postgres.NewPostgresDatabase(ezcdDatabaseUrl))

	// TODO: remove me
	log.Printf("Using database url %s", ezcdDatabaseUrl)

	// health check that pings the database
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		type HealthResponse struct {
			Status  string `json:"status"`
			Message string `json:"message,omitempty"`
		}

		if _, err := ezcdService.GetProjects(); err != nil {
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
		projects, err := ezcdService.GetProjects()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching projects: %v", err), http.StatusInternalServerError)
			return
		}
		annotatedProjects := mapToAnnotatedProjects(projects)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(annotatedProjects); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/projects/{projectID}/commits", func(w http.ResponseWriter, r *http.Request) {
		projectID := r.PathValue("projectID")
		if projectID == "" {
			http.Error(w, "Project ID is required", http.StatusBadRequest)
			return
		}

		commits, err := ezcdService.GetCommits(projectID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching commits: %v", err), http.StatusInternalServerError)
			return
		}
		annnotatedCommits := mapToAnnotatedCommits(commits)
		w.Header().Set("Content-Type", "application/json")
		// write an etag header so that browsers know whether the data has changed:
		w.Header().Set("ETag", fmt.Sprintf("%d", len(annnotatedCommits)))

		if err := json.NewEncoder(w).Encode(annnotatedCommits); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/projects/{projectID}", func(w http.ResponseWriter, r *http.Request) {
		projectID := r.PathValue("projectID")
		if projectID == "" {
			http.Error(w, "Project ID is required", http.StatusBadRequest)
			return
		}

		project, err := ezcdService.GetProject(projectID)
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

		annotatedProject := mapToAnnotatedProject(*project)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(annotatedProject); err != nil {
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

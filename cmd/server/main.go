package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ezcdlabs/ezcd/cmd/server/public"
)

var version = "dev" // default version

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("Version: %s\n", version)
		return
	}

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, API!!!")
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

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}

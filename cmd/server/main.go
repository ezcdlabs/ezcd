package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    if len(os.Args) > 1 && os.Args[1] == "--version" {
        fmt.Println("V0")
        return
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, Server!")
    })

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("could not start server: %s\n", err)
    }
}
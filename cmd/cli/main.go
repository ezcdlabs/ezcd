package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: ezcd-cli <command>")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "hello":
        fmt.Println("Hello, CLI!")
    case "--version":
        fmt.Println("V0!")
    default:
        fmt.Printf("Unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: vps <command>")
        os.Exit(1)
    }

    cmd := os.Args[1]
    switch cmd {
    case "init":
        fmt.Println("Initializing Hyperlane project...")
    case "up":
        fmt.Println("Creating VM...")
    case "list":
        fmt.Println("Listing VMs...")
    case "ssh":
        if len(os.Args) < 3 {
            fmt.Println("Usage: vps ssh <name>")
            os.Exit(1)
        }
        fmt.Printf("Connecting to VM %s...\n", os.Args[2])
    default:
        fmt.Println("Unknown command:", cmd)
    }
}

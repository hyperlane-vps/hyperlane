package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("Hyperlane agent starting...")
    for {
        fmt.Println("Reporting state to control plane...")
        time.Sleep(5 * time.Second)
    }
}

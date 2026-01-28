package main

import (
	"log"
	"time"

	"github.com/<your-username>/hyperlane/internal/controlplane"
	"github.com/<your-username>/hyperlane/internal/state"
)

func main() {
	store, err := state.NewStateStore("hyperlane.db")
	if err != nil {
		log.Fatalf("Failed to open state store: %v", err)
	}
	defer store.Close()

	// Connect to local agent
	agent, err := controlplane.NewAgentClient("localhost:50051", "certs/ca.crt")
	if err != nil {
		log.Fatalf("Failed to connect to agent: %v", err)
	}
	defer agent.Close()

	reconciler := controlplane.NewReconciler(store, 5*time.Second, agent)
	go reconciler.Start()

	// Keep control plane running
	select {}
}

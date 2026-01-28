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

	reconciler := controlplane.NewReconciler(store, 5*time.Second)
	go reconciler.Start()

	// Keep control plane running
	select {}
}

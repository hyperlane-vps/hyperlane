package controlplane

import (
	"fmt"
	"time"

	"github.com/<your-username>/hyperlane/internal/state"
)

// Reconciler periodically ensures that observed state matches desired state
type Reconciler struct {
	Store      *state.StateStore
	PollPeriod time.Duration
	StopChan   chan struct{}
}

// NewReconciler returns a new reconciler
func NewReconciler(store *state.StateStore, pollPeriod time.Duration) *Reconciler {
	return &Reconciler{
		Store:      store,
		PollPeriod: pollPeriod,
		StopChan:   make(chan struct{}),
	}
}

// Start begins the reconciliation loop
func (r *Reconciler) Start() {
	fmt.Println("Reconciliation loop started")
	ticker := time.NewTicker(r.PollPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.reconcile()
		case <-r.StopChan:
			fmt.Println("Reconciliation loop stopped")
			return
		}
	}
}

// Stop signals the reconciler to stop
func (r *Reconciler) Stop() {
	close(r.StopChan)
}

// reconcile fetches desired and observed states and applies minimal actions
func (r *Reconciler) reconcile() {
	desiredVMs, err := r.Store.ListDesired()
	if err != nil {
		fmt.Println("Error fetching desired state:", err)
		return
	}

	observedVMs, err := r.Store.ListObserved()
	if err != nil {
		fmt.Println("Error fetching observed state:", err)
		return
	}

	obsMap := make(map[string]state.VM)
	for _, vm := range observedVMs {
		obsMap[vm.Name] = vm
	}

	for _, d := range desiredVMs {
		o, exists := obsMap[d.Name]
		if !exists {
			fmt.Printf("VM %s not observed yet. Would provision now.\n", d.Name)
			// TODO: trigger agent to create VM
			continue
		}

		if d.State != o.State {
			fmt.Printf("VM %s state mismatch (desired: %s, observed: %s). Reconciling...\n",
				d.Name, d.State, o.State)
			// TODO: trigger agent to align state
		} else {
			fmt.Printf("VM %s is in desired state: %s\n", d.Name, d.State)
		}
	}
}

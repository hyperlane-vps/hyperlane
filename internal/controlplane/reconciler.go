package controlplane

import (
	"fmt"
	"time"

	"github.com/<your-username>/hyperlane/internal/state"
)

type Reconciler struct {
	Store      *state.StateStore
	PollPeriod time.Duration
	StopChan   chan struct{}
	Agent      *AgentClient
}

func NewReconciler(store *state.StateStore, pollPeriod time.Duration, agent *AgentClient) *Reconciler {
	return &Reconciler{
		Store:      store,
		PollPeriod: pollPeriod,
		StopChan:   make(chan struct{}),
		Agent:      agent,
	}
}

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

func (r *Reconciler) Stop() {
	close(r.StopChan)
}

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

		// VM not yet observed → create it
		if !exists || o.State == state.Requested {
			fmt.Printf("Reconciling VM %s: creating...\n", d.Name)
			r.Agent.CreateVM(d.Name, d.CPU, d.RAM, d.Image)

			// Update observed state as provisioning
			r.Store.SetObservedState(state.VM{
				Name:      d.Name,
				CPU:       d.CPU,
				RAM:       d.RAM,
				Image:     d.Image,
				State:     state.Provisioning,
				UpdatedAt: time.Now(),
			})
			continue
		}

		// VM exists but state mismatch → take corrective action
		if d.State != o.State {
			fmt.Printf("VM %s state mismatch (desired: %s, observed: %s)\n", d.Name, d.State, o.State)
			switch d.State {
			case state.Running:
				r.Agent.CreateVM(d.Name, d.CPU, d.RAM, d.Image)
			case state.Stopped:
				r.Agent.StopVM(d.Name)
			case state.Destroyed:
				r.Agent.DestroyVM(d.Name)
			default:
				fmt.Println("Unhandled desired state:", d.State)
			}
		} else {
			fmt.Printf("VM %s is in desired state: %s\n", d.Name, d.State)
		}
	}
}

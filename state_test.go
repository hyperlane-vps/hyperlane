package state

import (
	"os"
	"testing"
)

func TestStateStore(t *testing.T) {
	path := "test.db"
	defer os.Remove(path)

	store, err := NewStateStore(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	vm := VM{
		Name:  "test-vm",
		CPU:   2,
		RAM:   2048,
		Image: "ubuntu-22.04",
		State: Requested,
	}

	// Desired state
	if err := store.SetDesiredState(vm); err != nil {
		t.Fatalf("SetDesiredState failed: %v", err)
	}

	got, err := store.GetDesiredState(vm.Name)
	if err != nil {
		t.Fatalf("GetDesiredState failed: %v", err)
	}
	if got.Name != vm.Name || got.CPU != vm.CPU {
		t.Fatalf("unexpected VM retrieved: %+v", got)
	}

	// Observed state
	vm.State = Running
	if err := store.SetObservedState(vm); err != nil {
		t.Fatalf("SetObservedState failed: %v", err)
	}

	gotObs, err := store.GetObservedState(vm.Name)
	if err != nil {
		t.Fatalf("GetObservedState failed: %v", err)
	}
	if gotObs.State != Running {
		t.Fatalf("expected state Running, got %v", gotObs.State)
	}
}

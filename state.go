package state

import (
	"encoding/json"
	"errors"
	"time"

	bolt "go.etcd.io/bbolt"
)

// VMState represents the lifecycle state of a VM
type VMState string

const (
	Requested    VMState = "requested"
	Provisioning VMState = "provisioning"
	Running      VMState = "running"
	Stopped      VMState = "stopped"
	Failed       VMState = "failed"
	Destroyed    VMState = "destroyed"
)

// VM represents a single virtual machine
type VM struct {
	Name        string    `json:"name"`
	CPU         int       `json:"cpu"`
	RAM         int       `json:"ram"`
	Image       string    `json:"image"`
	State       VMState   `json:"state"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastError   string    `json:"last_error,omitempty"`
}

// StateStore persists desired and observed VM state
type StateStore struct {
	db *bolt.DB
}

// NewStateStore opens or creates the BoltDB database
func NewStateStore(path string) (*StateStore, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// Ensure buckets exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("desired"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("observed"))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &StateStore{db: db}, nil
}

// Close closes the underlying database
func (s *StateStore) Close() error {
	return s.db.Close()
}

// SetDesiredState stores the desired state of a VM
func (s *StateStore) SetDesiredState(vm VM) error {
	vm.UpdatedAt = time.Now()
	data, err := json.Marshal(vm)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("desired"))
		return b.Put([]byte(vm.Name), data)
	})
}

// GetDesiredState retrieves the desired state of a VM
func (s *StateStore) GetDesiredState(name string) (*VM, error) {
	var vm VM
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("desired"))
		data := b.Get([]byte(name))
		if data == nil {
			return errors.New("desired state not found")
		}
		return json.Unmarshal(data, &vm)
	})
	if err != nil {
		return nil, err
	}
	return &vm, nil
}

// SetObservedState stores the observed state of a VM
func (s *StateStore) SetObservedState(vm VM) error {
	vm.UpdatedAt = time.Now()
	data, err := json.Marshal(vm)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("observed"))
		return b.Put([]byte(vm.Name), data)
	})
}

// GetObservedState retrieves the observed state of a VM
func (s *StateStore) GetObservedState(name string) (*VM, error) {
	var vm VM
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("observed"))
		data := b.Get([]byte(name))
		if data == nil {
			return errors.New("observed state not found")
		}
		return json.Unmarshal(data, &vm)
	})
	if err != nil {
		return nil, err
	}
	return &vm, nil
}

// ListDesired returns all desired VMs
func (s *StateStore) ListDesired() ([]VM, error) {
	var vms []VM
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("desired"))
		return b.ForEach(func(_, v []byte) error {
			var vm VM
			if err := json.Unmarshal(v, &vm); err != nil {
				return err
			}
			vms = append(vms, vm)
			return nil
		})
	})
	return vms, err
}

// ListObserved returns all observed VMs
func (s *StateStore) ListObserved() ([]VM, error) {
	var vms []VM
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("observed"))
		return b.ForEach(func(_, v []byte) error {
			var vm VM
			if err := json.Unmarshal(v, &vm); err != nil {
				return err
			}
			vms = append(vms, vm)
			return nil
		})
	})
	return vms, err
}

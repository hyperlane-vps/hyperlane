package control

import (
	"context"
	"log"
	"time"

	pb "github.com/<your-username>/hyperlane/internal/agent/proto"
)

type Reconciler struct {
	State *ClusterState
	Agent pb.VMServiceClient
}

func (r *Reconciler) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.reconcileOnce(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (r *Reconciler) reconcileOnce(ctx context.Context) {
	report, err := r.Agent.ReportState(ctx, &pb.Empty{})
	if err != nil {
		log.Println("Failed to fetch state:", err)
		return
	}

	// Update observed state
	for _, vm := range report.Vms {
		r.State.Observed[vm.Name] = ObservedVM{
			Name:  vm.Name,
			State: vm.State,
			Node:  vm.NodeId,
		}
	}

	// Compare desired vs observed
	for name, desired := range r.State.Desired {
		obs, exists := r.State.Observed[name]

		if !exists {
			log.Printf("[reconcile] VM %s missing → creating\n", name)
			r.Agent.CreateVM(ctx, &pb.CreateVMRequest{
				Name:  desired.Name,
				Cpu:   uint32(desired.Cpu),
				Ram:   uint32(desired.Ram),
				Image: desired.Image,
			})
			continue
		}

		if obs.State != "running" {
			log.Printf("[reconcile] VM %s not running (%s) → restarting\n", name, obs.State)
			r.Agent.CreateVM(ctx, &pb.CreateVMRequest{
				Name:  desired.Name,
				Cpu:   uint32(desired.Cpu),
				Ram:   uint32(desired.Ram),
				Image: desired.Image,
			})
		}
	}
}

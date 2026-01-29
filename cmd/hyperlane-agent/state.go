package main

import (
	"context"
	"log"

	pb "github.com/<your-username>/hyperlane/internal/agent/proto"
	"github.com/libvirt/libvirt-go"
)

func (a *VMAgent) ReportState(ctx context.Context) (*pb.StateReport, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	domains, err := conn.ListAllDomains(0)
	if err != nil {
		return nil, err
	}

	report := &pb.StateReport{}

	for _, dom := range domains {
		name, _ := dom.GetName()
		info, _ := dom.GetInfo()

		state := "unknown"
		switch info.State {
		case libvirt.DOMAIN_RUNNING:
			state = "running"
		case libvirt.DOMAIN_SHUTOFF:
			state = "stopped"
		case libvirt.DOMAIN_CRASHED:
			state = "crashed"
		}

		report.Vms = append(report.Vms, &pb.VMInfo{
			Name:  name,
			State: state,
			Cpu:   uint32(info.NrVirtCpu),
			Ram:   uint32(info.Memory / 1024),
		})
	}

	return report, nil
}

package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/<your-username>/hyperlane/internal/agent/proto"
)

// VMAgent executes VM lifecycle commands
type VMAgent struct{}

// CreateVM provisions a VM from a ZFS snapshot and boots it
func (a *VMAgent) CreateVM(name string, cpu, ram int, image string) error {
	fmt.Printf("Creating VM %s from image %s (CPU=%d, RAM=%dMB)\n", name, image, cpu, ram)

	// 1. Clone ZFS snapshot
	cloneCmd := exec.Command("zfs", "clone", fmt.Sprintf("images/%s@base", image), fmt.Sprintf("vms/%s", name))
	out, err := cloneCmd.CombinedOutput()
	if err != nil {
		log.Printf("ZFS clone failed: %s\n%s", err, out)
		return err
	}

	// 2. Define libvirt domain XML
	domainXML := fmt.Sprintf(`
<domain type='kvm'>
  <name>%s</name>
  <memory unit='MiB'>%d</memory>
  <vcpu>%d</vcpu>
  <os>
    <type arch='x86_64'>hvm</type>
    <boot dev='hd'/>
  </os>
  <devices>
    <disk type='file' device='disk'>
      <driver name='qemu' type='qcow2'/>
      <source file='/var/lib/libvirt/images/%s.qcow2'/>
      <target dev='vda' bus='virtio'/>
    </disk>
    <interface type='network'>
      <source network='default'/>
    </interface>
    <graphics type='vnc' port='-1'/>
  </devices>
</domain>
`, name, ram, cpu, name)

	// 3. Write domain XML to temp file
	tmpFile := fmt.Sprintf("/tmp/%s.xml", name)
	if err := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' > %s", domainXML, tmpFile)).Run(); err != nil {
		log.Println("Failed to write domain XML:", err)
		return err
	}

	// 4. Define and start the domain
	if err := exec.Command("virsh", "define", tmpFile).Run(); err != nil {
		log.Println("Failed to define domain:", err)
		return err
	}

	if err := exec.Command("virsh", "start", name).Run(); err != nil {
		log.Println("Failed to start domain:", err)
		return err
	}

	fmt.Printf("VM %s created and started successfully!\n", name)
	return nil
}

// StopVM stops a running VM
func (a *VMAgent) StopVM(name string) error {
	fmt.Printf("Stopping VM %s\n", name)
	cmd := exec.Command("virsh", "shutdown", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Failed to shutdown VM: %s\n%s", err, out)
		return err
	}
	return nil
}

// DestroyVM destroys a VM and cleans up its ZFS clone
func (a *VMAgent) DestroyVM(name string) error {
	fmt.Printf("Destroying VM %s\n", name)

	// Stop VM first
	if err := a.StopVM(name); err != nil {
		log.Println("Warning: failed to stop VM before destroy:", err)
	}

	// Destroy libvirt domain
	if err := exec.Command("virsh", "undefine", name).Run(); err != nil {
		log.Println("Failed to undefine domain:", err)
	}

	// Destroy ZFS clone
	if err := exec.Command("zfs", "destroy", fmt.Sprintf("vms/%s", name)).Run(); err != nil {
		log.Println("Failed to destroy ZFS clone:", err)
	}

	fmt.Printf("VM %s destroyed successfully\n", name)
	return nil
}

// ReportState stub (to integrate with gRPC ReportState later)
func (a *VMAgent) ReportState() []*proto.VMInfo {
	// TODO: list VMs and states
	return []*proto.VMInfo{}
}

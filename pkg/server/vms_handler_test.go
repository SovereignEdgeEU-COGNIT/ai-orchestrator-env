package server

import (
	"testing"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestAddVM(t *testing.T) {
	client, server, done := prepareTests(t)

	vm, err := client.GetVM("vm1")
	assert.Nil(t, err)
	assert.Nil(t, vm)

	vm = &core.VM{VMID: "vm1", Hostname: "test_vm_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = client.AddVM(vm)
	assert.Nil(t, err)

	vm, err = client.GetVM("vm1")
	assert.Nil(t, err)
	assert.NotNil(t, vm)
	assert.Equal(t, "vm1", vm.VMID)

	server.Shutdown()
	<-done
}

func TestGetVM(t *testing.T) {
	client, server, done := prepareTests(t)

	vm, err := client.GetVM("vm1")
	assert.Nil(t, err)
	assert.Nil(t, vm)

	vm = &core.VM{VMID: "vm1", Hostname: "test_vm_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = client.AddVM(vm)
	assert.Nil(t, err)

	vm, err = client.GetVM("vm1")
	assert.Nil(t, err)
	assert.NotNil(t, vm)
	assert.Equal(t, "vm1", vm.VMID)

	server.Shutdown()
	<-done
}

func TestGetVMs(t *testing.T) {
	client, server, done := prepareTests(t)

	vms, err := client.GetVMs()
	assert.Nil(t, err)
	assert.Nil(t, vms)

	vm := &core.VM{VMID: "host1", Hostname: "test_vm_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = client.AddVM(vm)
	assert.Nil(t, err)

	vm = &core.VM{VMID: "host2", Hostname: "test_vm2_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = client.AddVM(vm)
	assert.Nil(t, err)

	vms, err = client.GetVMs()
	assert.Nil(t, err)
	assert.NotNil(t, vms)
	assert.Len(t, vms, 2)

	server.Shutdown()
	<-done
}

func TestRemoveVM(t *testing.T) {
	client, server, done := prepareTests(t)

	vm := &core.VM{VMID: "vm1", Hostname: "test_vm_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err := client.AddVM(vm)
	assert.Nil(t, err)

	vm, err = client.GetVM("vm1")
	assert.Nil(t, err)
	assert.NotNil(t, vm)
	assert.Equal(t, "vm1", vm.VMID)

	err = client.RemoveVM("vm1")
	assert.Nil(t, err)

	vm, err = client.GetVM("vm1")
	assert.Nil(t, err)
	assert.Nil(t, vm)

	server.Shutdown()
	<-done
}

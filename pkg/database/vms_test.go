package database

import (
	"testing"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestAddVM(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)
	defer db.Close()

	vm, err := db.GetVM("test_vm_id")
	assert.Nil(t, err)
	assert.Nil(t, vm)

	vm = &core.VM{VMID: "test_vm_id", Hostname: "test_vm_name", CurrentCPU: 0, CurrentMemory: 0}
	err = db.AddVM(vm)
	assert.Nil(t, err)

	vms, err := db.GetVMs()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(vms))

	vm, err = db.GetVM("test_vm_id")
	assert.Nil(t, err)
	assert.NotNil(t, vm)
	assert.Equal(t, "test_vm_id", vm.VMID)
}

func TestVMStateMetric(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)
	defer db.Close()

	vm := &core.VM{VMID: "test_vm1_id", Hostname: "test_vm1_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddVM(vm)
	assert.Nil(t, err)

	vm = &core.VM{VMID: "test_vm2_id", Hostname: "test_vm2_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddVM(vm)
	assert.Nil(t, err)

	vms, err := db.GetVMs()
	assert.Nil(t, err)

	counter := 0
	for _, vm := range vms {
		counter += vm.StateID
	}
	assert.Equal(t, 3, counter)

	// If we remove a vm and add another one, the state ID should be reused
	err = db.RemoveVM("test_vm1_id")
	assert.Nil(t, err)

	vm = &core.VM{VMID: "test_vm3_id", Hostname: "test_vm3_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddVM(vm)
	assert.Nil(t, err)

	vms, err = db.GetVMs()
	assert.Nil(t, err)

	counter = 0
	for _, vm := range vms {
		counter += vm.StateID
	}
	assert.Equal(t, 3, counter)
}

func TestSetVMResources(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)
	defer db.Close()

	vm := &core.VM{VMID: "test_vm_id", Hostname: "test_vm_name", CurrentCPU: 0, CurrentMemory: 0}
	err = db.AddVM(vm)
	assert.Nil(t, err)

	err = db.SetVMResources(vm.VMID, int64(1), int64(2))
	assert.Nil(t, err)

	hosts, err := db.GetVMs()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(hosts))
	assert.Equal(t, int64(1), hosts[0].CurrentCPU)
	assert.Equal(t, int64(2), hosts[0].CurrentMemory)
}

func TestRemoveVM(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)
	defer db.Close()

	vm := &core.VM{VMID: "test_host_id", Hostname: "test_vm_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddVM(vm)
	assert.Nil(t, err)

	vms, err := db.GetVMs()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(vms))

	err = db.RemoveVM(vm.VMID)
	assert.Nil(t, err)

	vms, err = db.GetVMs()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(vms))
}

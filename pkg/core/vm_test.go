package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVMToJSON(t *testing.T) {
	vm := &VM{VMID: "1", Hostname: "hostname"}

	jsonStr, err := vm.ToJSON()
	assert.Nil(t, err)

	_, err = ConvertJSONToVM(jsonStr + "error")
	assert.NotNil(t, err)

	vm2, err := ConvertJSONToVM(jsonStr)
	assert.Nil(t, err)
	assert.True(t, vm2.Equals(vm))
}

func TestVMArrayToJSON(t *testing.T) {
	var vms []*VM

	vm1 := &VM{VMID: "1", Hostname: "hostname1"}
	vm2 := &VM{VMID: "2", Hostname: "hostname2"}

	vms = append(vms, vm1)
	vms = append(vms, vm2)

	jsonStr, err := ConvertVMArrayToJSON(vms)
	assert.Nil(t, err)

	_, err = ConvertJSONToVMArray(jsonStr + "error")
	assert.NotNil(t, err)

	vms2, err := ConvertJSONToVMArray(jsonStr)
	assert.Nil(t, err)
	assert.True(t, IsVMArraysEqual(vms, vms2))
}

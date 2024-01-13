package core

import "encoding/json"

type VM struct {
	VMID          string `json:"vmid"`
	StateID       int    `json:"stateid"`
	Hostname      string `json:"hostname"`
	CurrentCPU    int64  `json:"current_cpu"`
	CurrentMemory int64  `json:"current_memory"`
	Deployed      bool   `json:"deployed"`
	HostID        string `json:"hostid"`
	HostStateID   int    `json:"hoststateid"`
}

func ConvertJSONToVM(jsonString string) (*VM, error) {
	var vm *VM
	err := json.Unmarshal([]byte(jsonString), &vm)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func ConvertJSONToVMArray(jsonString string) ([]*VM, error) {
	var vms []*VM

	err := json.Unmarshal([]byte(jsonString), &vms)
	if err != nil {
		return vms, err
	}

	return vms, nil
}

func ConvertVMArrayToJSON(vms []*VM) (string, error) {
	jsonBytes, err := json.Marshal(vms)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (vm *VM) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(vm)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (vm *VM) Equals(vm2 *VM) bool {
	if vm2 == nil {
		return false
	}

	if vm.VMID == vm2.VMID &&
		vm.StateID == vm2.StateID &&
		vm.Hostname == vm2.Hostname &&
		vm.CurrentCPU == vm2.CurrentCPU &&
		vm.CurrentMemory == vm2.CurrentMemory &&
		vm.Deployed == vm2.Deployed &&
		vm.HostID == vm2.HostID &&
		vm.HostStateID == vm2.HostStateID {
		return true
	}

	return false
}

func IsVMArraysEqual(vms []*VM, vms2 []*VM) bool {
	if len(vms) != len(vms2) {
		return false
	}

	for i := 0; i < len(vms); i++ {
		if !vms[i].Equals(vms2[i]) {
			return false
		}
	}

	return true
}

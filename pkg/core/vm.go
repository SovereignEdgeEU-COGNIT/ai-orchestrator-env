package core

import "encoding/json"

type VM struct {
	VMID     string `json:"vmid"`
	Hostname string `json:"hostname"`
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
		vm.Hostname == vm2.Hostname {
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

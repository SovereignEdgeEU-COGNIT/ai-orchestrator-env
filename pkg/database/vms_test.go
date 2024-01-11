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

	vm := &core.VM{VMID: "test_vm_id", Hostname: "test_host_name"}
	err = db.AddVM(vm)
	assert.Nil(t, err)

	vms, err := db.GetVMs()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(vms))
}

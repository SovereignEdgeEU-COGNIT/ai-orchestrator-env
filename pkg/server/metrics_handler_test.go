package server

import (
	"testing"
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestAddHostMetrics(t *testing.T) {
	client, server, done := prepareTests(t)

	host := &core.Host{HostID: "host1", Hostname: "test_host_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err := client.AddHost(host)
	assert.Nil(t, err)

	host, err = client.GetHost("host1")
	assert.Nil(t, err)
	assert.Equal(t, host.CurrentCPU, int64(0))
	assert.Equal(t, host.CurrentMemory, int64(0))

	now := time.Now()
	err = client.AddMetric("host1", core.HostType, &core.Metric{Timestamp: now, CPU: 1, Memory: 10})
	assert.Nil(t, err)

	oneSecondAgo := now.Add(-time.Second)
	metrics, err := client.GetMetrics("host1", core.HostType, oneSecondAgo, 1)
	assert.Nil(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, metrics[0].CPU, int64(1))
	assert.Equal(t, metrics[0].Memory, int64(10))

	host, err = client.GetHost("host1")
	assert.Nil(t, err)
	assert.Equal(t, host.CurrentCPU, int64(1))
	assert.Equal(t, host.CurrentMemory, int64(10))

	server.Shutdown()
	<-done
}

func TestAddVMMetrics(t *testing.T) {
	client, server, done := prepareTests(t)

	vm := &core.VM{VMID: "vm1", Hostname: "test_vm_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err := client.AddVM(vm)
	assert.Nil(t, err)

	vm, err = client.GetVM("vm1")
	assert.Nil(t, err)
	assert.Equal(t, vm.CurrentCPU, int64(0))
	assert.Equal(t, vm.CurrentMemory, int64(0))

	now := time.Now()
	err = client.AddMetric("vm1", core.VMType, &core.Metric{Timestamp: now, CPU: 1, Memory: 10})
	assert.Nil(t, err)

	oneSecondAgo := now.Add(-time.Second)
	metrics, err := client.GetMetrics("vm1", core.VMType, oneSecondAgo, 1)
	assert.Nil(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, metrics[0].CPU, int64(1))
	assert.Equal(t, metrics[0].Memory, int64(10))

	vm, err = client.GetVM("vm1")
	assert.Nil(t, err)
	assert.Equal(t, vm.CurrentCPU, int64(1))
	assert.Equal(t, vm.CurrentMemory, int64(10))

	server.Shutdown()
	<-done
}

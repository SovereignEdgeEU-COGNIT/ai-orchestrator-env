package database

import (
	"testing"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestAddHost(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)
	defer db.Close()

	host, err := db.GetHost("test_host_id")
	assert.Nil(t, err)
	assert.Nil(t, host)

	host = &core.Host{HostID: "test_host_id", Hostname: "test_host_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddHost(host)
	assert.Nil(t, err)

	hosts, err := db.GetHosts()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(hosts))

	host, err = db.GetHost("test_host_id")
	assert.Nil(t, err)
	assert.NotNil(t, host)
	assert.Equal(t, "test_host_id", host.HostID)
}

func TestHostStateMetric(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)
	defer db.Close()

	host := &core.Host{HostID: "test_host1_id", Hostname: "test_host1_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddHost(host)
	assert.Nil(t, err)

	host = &core.Host{HostID: "test_host2_id", Hostname: "test_host2_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddHost(host)
	assert.Nil(t, err)

	hosts, err := db.GetHosts()
	assert.Nil(t, err)

	counter := 0
	for _, host := range hosts {
		counter += host.StateID
	}
	assert.Equal(t, 3, counter)

	// If we remove a host and add another one, the state ID should be reused
	err = db.RemoveHost("test_host1_id")
	assert.Nil(t, err)

	host = &core.Host{HostID: "test_host3_id", Hostname: "test_host3_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddHost(host)
	assert.Nil(t, err)

	hosts, err = db.GetHosts()
	assert.Nil(t, err)

	counter = 0
	for _, host := range hosts {
		counter += host.StateID
	}
	assert.Equal(t, 3, counter)
}

func TestSetHostResources(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)
	defer db.Close()

	host := &core.Host{HostID: "test_host_id", Hostname: "test_host_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddHost(host)
	assert.Nil(t, err)

	err = db.SetHostResources(host.HostID, int64(1), int64(2))
	assert.Nil(t, err)

	hosts, err := db.GetHosts()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(hosts))
	assert.Equal(t, int64(1), hosts[0].CurrentCPU)
	assert.Equal(t, int64(2), hosts[0].CurrentMemory)
}

func TestRemoveHost(t *testing.T) {
	db, err := PrepareTests()
	assert.Nil(t, err)
	defer db.Close()

	host := &core.Host{HostID: "test_host_id", Hostname: "test_host_name", CurrentCPU: 0.0, CurrentMemory: 0.0}
	err = db.AddHost(host)
	assert.Nil(t, err)

	hosts, err := db.GetHosts()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(hosts))

	err = db.RemoveHost(host.HostID)
	assert.Nil(t, err)

	hosts, err = db.GetHosts()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(hosts))
}

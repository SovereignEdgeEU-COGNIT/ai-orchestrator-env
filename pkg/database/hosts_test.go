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

	host := &core.Host{HostID: "test_host_id", Hostname: "test_host_name"}
	err = db.AddHost(host)
	assert.Nil(t, err)
}

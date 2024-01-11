package server

import (
	"testing"
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestAddMetrics(t *testing.T) {
	client, server, done := prepareTests(t)

	now := time.Now()
	err := client.AddMetric("host1", core.HostType, &core.Metric{Timestamp: now, CPU: 1, Memory: 10})
	assert.Nil(t, err)

	oneSecondAgo := now.Add(-time.Second)
	metrics, err := client.GetMetrics("host1", core.HostType, oneSecondAgo, 1)
	assert.Nil(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, metrics[0].CPU, int64(1))
	assert.Equal(t, metrics[0].Memory, int64(10))

	server.Shutdown()
	<-done
}

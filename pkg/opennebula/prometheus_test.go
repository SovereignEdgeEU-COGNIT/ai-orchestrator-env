package opennebula

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const prometheusURL = "http://localhost:9090"

func TestGetHostIDs(t *testing.T) {
	hostIDs, err := GetHostIDs(prometheusURL)
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, hostIDs, "hostIDs should not be empty")
}

func TestGetVMIDs(t *testing.T) {
	vmIDs, err := GetVMIDs(prometheusURL)
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, vmIDs, "vmIDs should not be empty")
}

func TestGetHostCPU(t *testing.T) {
	cpuLoad, err := GetHostCPU(prometheusURL, "4")
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, cpuLoad > 0, "cpuLoad should be greater than 0")
}

func TestGetHostBusyCPU(t *testing.T) {
	cpuLoad, err := GetHostCPUBusy(prometheusURL, "4")
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, cpuLoad > 0, "cpuLoad should be greater than 0")
}

func TestGetHostUsedMem(t *testing.T) {
	usedMem, err := GetHostUsedMem(prometheusURL, "4")
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, usedMem > 0, "usedMem should be greater than 0")
}

func TestGetHostAvailMem(t *testing.T) {
	availMem, err := GetHostAvailMem(prometheusURL, "4")
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, availMem > 0, "availMem should be greater than 0")
}

func TestGetHostNetTransmit(t *testing.T) {
	netTrans, err := GetHostNetTrans(prometheusURL, "4")
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, netTrans > 0, "netTrans should be greater than 0")
}

func TestGetHostNetRecv(t *testing.T) {
	netRecv, err := GetHostNetRecv(prometheusURL, "4")
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, netRecv > 0, "netTrans should be greater than 0")
}

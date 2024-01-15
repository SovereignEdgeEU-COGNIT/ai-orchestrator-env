package opennebula

import (
	"fmt"
	"testing"
)

const prometheusURL = "http://localhost:9090"

func TestGetHostIDs(t *testing.T) {
	hostIDs, err := GetHostIDs(prometheusURL)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("hostIDs:", hostIDs)
}

func TestGetVMIDs(t *testing.T) {
	vmIDs, err := GetVMIDs(prometheusURL)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("vmIDs:", vmIDs)
}

func TestGetHost(t *testing.T) {
	vmIDs, err := GetHost(prometheusURL, "2")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("vmIDs:", vmIDs)
}

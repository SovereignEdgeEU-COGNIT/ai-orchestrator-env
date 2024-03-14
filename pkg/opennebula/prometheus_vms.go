package opennebula

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
)

func GetVMIDs(prometheusURL string) ([]string, error) {
	var vmsIDs []string

	query := `opennebula_vm_state`
	r, err := QueryPrometheus(prometheusURL, query)
	if err != nil {
		return vmsIDs, err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return vmsIDs, err
	}

	for _, result := range resp.Data.Result {
		vmsIDs = append(vmsIDs, result.Metric.OneVMID)
	}

	return vmsIDs, nil
}

func MapVMHostIDs(prometheusURL string, vmMap map[string]*core.VM) error {

	query := `opennebula_vm_host_id`
	r, err := QueryPrometheus(prometheusURL, query)
	if err != nil {
		return err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return err
	}

	for _, result := range resp.Data.Result {

		hostIdStr, ok := result.Value[1].(string)

		if !ok {
			println("Could not convert hostId to string")
			continue
		}

		if hostId, err := strconv.ParseInt(hostIdStr, 10, 64); err != nil || hostId < 1 {
			println("Could not convert hostId to int or hostId < 1, vmId: ", result.Metric.OneVMID, " hostId: ", hostId, " err: ", err)
			continue
		}

		// if result.Metric.OneVMID in map then set hostID
		if vm, ok := vmMap[result.Metric.OneVMID]; !ok {
			println("VMID not found in map, ", result.Metric.OneVMID)
			continue
		} else {
			vm.HostID = hostIdStr
		}

	}

	return nil
}

func GetVMsCPUUsage(prometheusURL string, vmMap map[string]*core.VM) error {
	return fmt.Errorf("GetVMsCPUUsage not implemented")
}

func GetVMsCPUTotal(prometheusURL string, vmMap map[string]*core.VM) error {

	//return error, still missing API call
	return fmt.Errorf("GetVMsCPUTotal not implemented")

}

func GetVMsDiskRead(prometheusURL string, vmMap map[string]*core.VM) error {

	query := `sum by(one_vm_id)(increase(opennebula_libvirt_block_rd_bytes[40s]))`
	r, err := QueryPrometheus(prometheusURL, query)
	if err != nil {
		return err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return err
	}

	for _, result := range resp.Data.Result {

		diskReadStr, ok := result.Value[1].(string)
		if !ok {
			println("Failed to convert GetVMsDiskRead to string")
		}

		diskReadBytes, err := strconv.ParseFloat(diskReadStr, 64)

		if err != nil {
			println("Failed to convert GetVMsDiskRead to float64")
		}

		// if result.Metric.OneVMID in map then set hostID
		if vm, ok := vmMap[result.Metric.OneVMID]; !ok {
			println("VMID not found in map, ", result.Metric.OneVMID)
			continue
		} else {
			vm.DiskRead = diskReadBytes / 1024 / 1024
		}
	}

	return nil
}

func GetVMsDiskWrite(prometheusURL string, vmMap map[string]*core.VM) error {

	query := `sum by(one_vm_id)(increase(opennebula_libvirt_block_wr_bytes[40s]))`
	r, err := QueryPrometheus(prometheusURL, query)
	if err != nil {
		return err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return err
	}

	for _, result := range resp.Data.Result {

		diskWriteStr, ok := result.Value[1].(string)
		if !ok {
			println("Failed to convert GetVMsDiskWrite to string")
		}

		diskWriteBytes, err := strconv.ParseFloat(diskWriteStr, 64)

		if err != nil {
			println("Failed to convert GetVMsDiskWrite to float64")
		}

		// if result.Metric.OneVMID in map then set hostID
		if vm, ok := vmMap[result.Metric.OneVMID]; !ok {
			println("VMID not found in map, ", result.Metric.OneVMID)
			continue
		} else {
			vm.DiskWrite = diskWriteBytes / 1024 / 1024
		}
	}

	return nil
}

func GetVMsMemUsage(prometheusURL string, vmMap map[string]*core.VM) error {

	query := `opennebula_libvirt_memory_total_bytes`
	r, err := QueryPrometheus(prometheusURL, query)
	if err != nil {
		return err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return err
	}

	for _, result := range resp.Data.Result {

		memUsageStr, ok := result.Value[1].(string)
		if !ok {
			println("Failed to convert GetVMsMemUsage to string")
		}

		memUsageBytes, err := strconv.ParseFloat(memUsageStr, 64)

		if err != nil {
			println("Failed to convert GetVMsMemUsage to float64")
		}

		// if result.Metric.OneVMID in map then set hostID
		if vm, ok := vmMap[result.Metric.OneVMID]; !ok {
			println("VMID not found in map, ", result.Metric.OneVMID)
			continue
		} else {
			vm.UsageMemory = memUsageBytes / 1024 / 1024
		}
	}

	return nil
}

func GetVMsMemTotal(prometheusURL string, vmMap map[string]*core.VM) error {

	query := `opennebula_libvirt_memory_maximum_bytes`
	r, err := QueryPrometheus(prometheusURL, query)
	if err != nil {
		return err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return err
	}

	for _, result := range resp.Data.Result {

		memTotalStr, ok := result.Value[1].(string)
		if !ok {
			println("Failed to convert GetVMsMemTotal to string")
		}

		memTotalBytes, err := strconv.ParseFloat(memTotalStr, 64)

		if err != nil {
			println("Failed to convert GetVMsMemTotal to float64")
		}

		// if result.Metric.OneVMID in map then set hostID
		if vm, ok := vmMap[result.Metric.OneVMID]; !ok {
			println("VMID not found in map, ", result.Metric.OneVMID)
			continue
		} else {
			vm.TotalMemory = memTotalBytes / 1024 / 1024
		}
	}

	return nil
}

func GetVMsNetRx(prometheusURL string, vmMap map[string]*core.VM) error {

	query := `sum by(one_vm_id)(increase(opennebula_libvirt_net_rx_total_bytes[40s]))`
	r, err := QueryPrometheus(prometheusURL, query)
	if err != nil {
		return err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return err
	}

	for _, result := range resp.Data.Result {

		netRxStr, ok := result.Value[1].(string)
		if !ok {
			println("Failed to convert GetVMsNetRx to string")
		}

		netRxBytes, err := strconv.ParseFloat(netRxStr, 64)

		if err != nil {
			println("Failed to convert GetVMsNetRx to float64")
		}

		// if result.Metric.OneVMID in map then set hostID
		if vm, ok := vmMap[result.Metric.OneVMID]; !ok {
			println("VMID not found in map, ", result.Metric.OneVMID)
			continue
		} else {
			vm.NetRX = netRxBytes / 1024 / 1024
		}
	}

	return nil
}

func GetVMsNetTx(prometheusURL string, vmMap map[string]*core.VM) error {

	query := `sum by(one_vm_id)(increase(opennebula_libvirt_net_tx_total_bytes[40s]))`
	r, err := QueryPrometheus(prometheusURL, query)
	if err != nil {
		return err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return err
	}

	for _, result := range resp.Data.Result {

		netTxStr, ok := result.Value[1].(string)
		if !ok {
			println("Failed to convert GetVMsNetTx to string")
		}

		netTxBytes, err := strconv.ParseFloat(netTxStr, 64)

		if err != nil {
			println("Failed to convert GetVMsNetTx to float64")
		}

		// if result.Metric.OneVMID in map then set hostID
		if vm, ok := vmMap[result.Metric.OneVMID]; !ok {
			println("VMID not found in map, ", result.Metric.OneVMID)
			continue
		} else {
			vm.NetTX = netTxBytes / 1024 / 1024
		}
	}

	return nil
}

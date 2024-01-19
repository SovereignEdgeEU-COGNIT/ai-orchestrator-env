package core

import "encoding/json"

type Host struct {
	HostID      string  `json:"hostid"`
	StateID     int     `json:"stateid"`
	TotalCPU    int64   `json:"total_cpu"`
	TotalMemory int64   `json:"total_memory"`
	UsageCPU    float64 `json:"usage_cpu"`
	UsageMemory int64   `json:"usage_memory"`
	VMs         int     `json:"vms"`
}

func ConvertJSONToHost(jsonString string) (*Host, error) {
	var host *Host
	err := json.Unmarshal([]byte(jsonString), &host)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func ConvertJSONToHostArray(jsonString string) ([]*Host, error) {
	var hosts []*Host

	err := json.Unmarshal([]byte(jsonString), &hosts)
	if err != nil {
		return hosts, err
	}

	return hosts, nil
}

func ConvertHostArrayToJSON(hosts []*Host) (string, error) {
	jsonBytes, err := json.Marshal(hosts)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (host *Host) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(host)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (host *Host) Equals(host2 *Host) bool {
	if host2 == nil {
		return false
	}

	if host.HostID == host2.HostID &&
		host.StateID == host2.StateID &&
		host.TotalCPU == host2.TotalCPU &&
		host.TotalMemory == host2.TotalMemory &&
		host.UsageCPU == host2.UsageCPU &&
		host.UsageMemory == host2.UsageMemory {
		return true
	}

	return false
}

func IsHostArraysEqual(hosts []*Host, hosts2 []*Host) bool {
	if len(hosts) != len(hosts2) {
		return false
	}

	for i := 0; i < len(hosts); i++ {
		if !hosts[i].Equals(hosts2[i]) {
			return false
		}
	}

	return true
}

package core

import "encoding/json"

type Host struct {
	HostID        string `json:"hostid"`
	Hostname      string `json:"hostname"`
	CurrentCPU    int64  `json:"current_cpu"`
	CurrentMemory int64  `json:"current_memory"`
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
		host.Hostname == host2.Hostname &&
		host.CurrentCPU == host2.CurrentCPU &&
		host.CurrentMemory == host2.CurrentMemory {
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

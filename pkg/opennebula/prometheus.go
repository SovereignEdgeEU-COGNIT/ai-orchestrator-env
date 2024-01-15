package opennebula

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type PrometheusResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Metric Metric        `json:"metric"`
	Value  []interface{} `json:"value"`
}

type Metric struct {
	Name      string `json:"__name__"`
	Instance  string `json:"instance"`
	Job       string `json:"job"`
	OneHostID string `json:"one_host_id"`
	OneVMID   string `json:"one_vm_id"`
}

func queryPrometheus(prometheusURL, query string) ([]byte, error) {
	fullURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, url.QueryEscape(query))

	resp, err := http.Get(fullURL)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func GetHostIDs(prometheusURL string) ([]string, error) {
	var hostIDs []string

	query := `opennebula_host_state`
	r, err := queryPrometheus(prometheusURL, query)
	if err != nil {
		return hostIDs, err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return hostIDs, err
	}

	for _, result := range resp.Data.Result {
		hostIDs = append(hostIDs, result.Metric.OneHostID)
	}

	return hostIDs, nil
}

func GetVMIDs(prometheusURL string) ([]string, error) {
	var vmsIDs []string

	query := `opennebula_vm_state`
	r, err := queryPrometheus(prometheusURL, query)
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

func GetHost(prometheusURL, hostID string) ([]string, error) {
	var vmsIDs []string

	//query := `label_replace(opennebula_host_mem_total_bytes, "one_host_id", "$1", "instance", "unique_regex_pattern_here")`
	//	query := `{__name__=~"opennebula_host_mem_total_bytes|opennebula_host_mem_usage_bytes", one_host_id=~".+"}`
	//query := `{__name__=~"opennebula_host_mem_total_bytes|opennebula_host_mem_usage_bytes|opennebula_host_cpu_total_ratio", one_host_id=~".+"}`
	//query := `{__name__=~"opennebula_host_mem_total_bytes|opennebula_host_mem_usage_bytes|opennebula_host_cpu_total_ratio|opennebula_host_cpu_total_ratio", one_host_id=~".+"}`
	//	query := `{__name__=~"opennebula_host_mem_total_bytes|opennebula_host_mem_usage_bytes|opennebula_host_cpu_total_ratio|opennebula_host_cpu_usage_ratio", one_host_id=~".+"}`

	//query := `(((count(count(node_cpu_seconds_total{one_host_id="4"}) by (cpu))) - avg(sum by (mode)(rate(node_cpu_seconds_total{mode='idle',one_host_id="4"}[60s])))) * 100) / count(count(node_cpu_seconds_total{one_host_id="4"}) by (cpu))`
	//query := `node_cpu_seconds_total{one_host_id="4"}`
	//	query := `sum by (cpu, instance) (node_cpu_seconds_total{one_host_id="4"})`

	//query := `avg by (mode) (node_cpu_seconds_total{one_host_id="4"})` // Working Cumulative

	query := `avg by (mode) (rate(node_cpu_seconds_total{one_host_id="4"}[40s]))` // Rate

	r, err := queryPrometheus(prometheusURL, query)
	if err != nil {
		return vmsIDs, err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return vmsIDs, err
	}

	fmt.Println(string(r))
	// for _, result := range resp.Data.Result {
	// 	vmsIDs = append(vmsIDs, result.Metric.OneVMID)
	// }

	return vmsIDs, nil
}

// func GetHostTotalMem(prometheusURL, hostID string) ([]string, error) {
// 	var vmsIDs []string
//
// 	opennebula_host_mem_total_bytes{one_host_id="4"}
// 	query := `opennebula_host_vms{one_host_id="4"}`
// 	r, err := queryPrometheus(prometheusURL, query)
// 	if err != nil {
// 		return vmsIDs, err
// 	}
//
// 	var resp PrometheusResponse
// 	err = json.Unmarshal(r, &resp)
// 	if err != nil {
// 		return vmsIDs, err
// 	}
//
// 	fmt.Println(string(r))
// 	// for _, result := range resp.Data.Result {
// 	// 	vmsIDs = append(vmsIDs, result.Metric.OneVMID)
// 	// }
//
// 	return vmsIDs, nil
// }

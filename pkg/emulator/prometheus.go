package emulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
	OneHostID                       string `json:"one_host_id"`
	OneVMID                         string `json:"one_vm_id"`
	Name                            string `json:"__name__"`
	DockerComposeConfigHash         string `json:"container_label_com_docker_compose_config_hash"`
	DockerComposeContainerNumber    string `json:"container_label_com_docker_compose_container_number"`
	DockerComposeOneoff             string `json:"container_label_com_docker_compose_oneoff"`
	DockerComposeProject            string `json:"container_label_com_docker_compose_project"`
	DockerComposeProjectConfigFiles string `json:"container_label_com_docker_compose_project_config_files"`
	DockerComposeProjectWorkingDir  string `json:"container_label_com_docker_compose_project_working_dir"`
	DockerComposeService            string `json:"container_label_com_docker_compose_service"`
	DockerComposeVersion            string `json:"container_label_com_docker_compose_version"`
	CPU                             string `json:"cpu"`
	ID                              string `json:"id"`
	Image                           string `json:"image"`
	Instance                        string `json:"instance"`
	Job                             string `json:"job"`
	NameLabel                       string `json:"name"`
}

type HostMetric struct {
	Name        string    `json:"name"`
	Timestamp   time.Time `json:"timestamp"`
	CPURate     float64   `json:"cpurate"`
	MemoryUsage int64     `json:"memoryusage"`
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

func GetFlavourMetricForHost(prometheusURL, host string) (*HostMetric, error) {
	timestamp, cpuRate, err := getCPURateForHost(prometheusURL, host)
	if err != nil {
		return nil, err
	}

	timestamp, memoryUsage, err := getMemoryUsageForHost(prometheusURL, host)
	if err != nil {
		return nil, err
	}

	return &HostMetric{
		Name:        host,
		Timestamp:   timestamp,
		CPURate:     cpuRate,
		MemoryUsage: memoryUsage,
	}, nil
}

func getMemoryUsageForHost(prometheusURL, host string) (time.Time, int64, error) {
	query := `container_memory_usage_bytes{name="` + host + `"}`

	var timestamp time.Time
	var memoryUsage int64

	r, err := queryPrometheus(prometheusURL, query)
	if err != nil {
		return timestamp, memoryUsage, err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return timestamp, memoryUsage, err
	}

	if len(resp.Data.Result) == 0 {
		return timestamp, memoryUsage, errors.New("No metrics found for " + host + " host")
	}

	if len(resp.Data.Result) > 1 {
		return timestamp, memoryUsage, errors.New("Only one metric result is allowed for " + host + " host")
	}

	var result = resp.Data.Result[0]
	if result.Metric.NameLabel == host {
		var memoryUsageStr string
		var ok bool
		var floatTimestamp float64
		if len(result.Value) == 2 {
			floatTimestamp, ok = result.Value[0].(float64)
			if !ok {
				return timestamp, memoryUsage, errors.New("Failed to parse timestamp metric")
			}

			memoryUsageStr, ok = result.Value[1].(string)
			if !ok {
				return timestamp, memoryUsage, errors.New("Failed to parse rate metric")
			}
		} else {
			return timestamp, memoryUsage, errors.New("Invalid number of values for " + host + " host")
		}

		memoryUsage, err = strconv.ParseInt(memoryUsageStr, 10, 64)
		if err != nil {
			return timestamp, memoryUsage, err
		}

		seconds := int64(floatTimestamp)
		nanoseconds := int64((floatTimestamp - float64(seconds)) * 1e9)
		timestamp = time.Unix(seconds, nanoseconds)
	}

	return timestamp, memoryUsage, nil
}

func getCPURateForHost(prometheusURL, host string) (time.Time, float64, error) {
	query := `rate(container_cpu_user_seconds_total{name="` + host + `"}[30s])`

	var timestamp time.Time
	var cpuRate float64

	r, err := queryPrometheus(prometheusURL, query)
	if err != nil {
		return timestamp, cpuRate, err
	}

	var resp PrometheusResponse
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return timestamp, cpuRate, err
	}

	if len(resp.Data.Result) == 0 {
		return timestamp, cpuRate, errors.New("No metrics found for " + host + " host")
	}

	if len(resp.Data.Result) > 1 {
		return timestamp, cpuRate, errors.New("Only one metric result is allowed for " + host + " host")
	}

	var result = resp.Data.Result[0]
	if result.Metric.NameLabel == host {
		var cpuRateStr string
		var ok bool
		var floatTimestamp float64
		if len(result.Value) == 2 {
			floatTimestamp, ok = result.Value[0].(float64)
			if !ok {
				return timestamp, cpuRate, errors.New("Failed to parse timestamp metric")
			}

			cpuRateStr, ok = result.Value[1].(string)
			if !ok {
				return timestamp, cpuRate, errors.New("Failed to parse rate metric")
			}
		} else {
			return timestamp, cpuRate, errors.New("Invalid number of values for " + host + " host")
		}

		cpuRate, err = strconv.ParseFloat(cpuRateStr, 64)
		if err != nil {
			return timestamp, cpuRate, err
		}

		seconds := int64(floatTimestamp)
		nanoseconds := int64((floatTimestamp - float64(seconds)) * 1e9)
		timestamp = time.Unix(seconds, nanoseconds)
	}

	return timestamp, cpuRate, nil
}

func GetTotalCPU(prometheusURL, host string) (float64, error) {
	// TODO The query machine_cpu_cores return the info below, but we need to map instance to host
	// machine_cpu_cores{boot_id="13156795-2d5e-4745-8dca-5267d7027ba4", instance="cadvisor:8080", job="cadvisor1", machine_id="28cf4e1904244c5d92a327cf54597d65", system_uuid="28cf4e19-0424-4c5d-92a3-27cf54597d65"}
	// machine_cpu_cores{boot_id="d9461eb7-1525-4328-8bb6-3fe40d8cf28c", instance="194.28.122.123:8080", job="cadvisor2", machine_id="1f1d300b5e074b768d1f1bcb3bf12220", system_uuid="1f1d300b-5e07-4b76-8d1f-1bcb3bf12220"}
	return 5.0, nil
}

func GetTotalMemory(prometheusURL, host string) (int64, error) {
	// TODO query machine_memory_bytes
	return int64(15348236288), nil
}

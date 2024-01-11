package client

import (
	"errors"
	"strconv"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/go-resty/resty/v2"
)

type EnvClient struct {
	restyClient *resty.Client
	host        string
	port        int
	protocol    string
}

func CreateEnvClient(host string, port int, insecure bool) *EnvClient {
	client := &EnvClient{}
	client.restyClient = resty.New()

	client.host = host
	client.port = port

	client.protocol = "https"
	if insecure {
		client.protocol = "http"
	}

	return client
}

func checkStatus(statusCode int, body string) error {
	if statusCode != 200 {
		return errors.New(body)
	}

	return nil
}

func (client *EnvClient) AddMetric(metric *core.Metric) error {
	jsonString, err := metric.ToJSON()
	if err != nil {
		return err
	}

	resp, err := client.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonString).
		Post(client.protocol + "://" + client.host + ":" + strconv.Itoa(client.port) + "/metrics")
	if err != nil {
		return err
	}

	err = checkStatus(resp.StatusCode(), string(resp.Body()))
	if err != nil {
		return err
	}

	return nil
}


func (client *EnvClient) GetEvents(hostID string , metricType int, sinceUnixNano int64, count) ([]*core.Metric, error) {
	resp, err := client.restyClient.R().
		SetHeader("APIKEY", client.apiKey).
		Get(client.protocol + "://" + client.host + ":" + strconv.Itoa(client.port) + "/metrics?hostid=" + hostID + "&metrictype=" + strconv.Itoa(metricType) + "&since=" + strconv.FormatInt(sinceUnixNano, 10) + "&count=" + strconv.Itoa(count))
	if err != nil {
		return nil, err
	}

	err = checkStatus(resp.StatusCode(), string(resp.Body()))
	if err != nil {
		return nil, err
	}

	respBodyString := string(resp.Body())

	seismogramPackages, err := core.ConvertJSONToMetricArray(respBodyString)
	if err != nil {
		return nil, err
	}

	return seismogramPackages, nil
}

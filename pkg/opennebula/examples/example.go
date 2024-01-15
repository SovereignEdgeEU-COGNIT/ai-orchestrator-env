package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Prometheus server URL
	prometheusURL := "http://localhost:9090"

	// PromQL query
	//query := `opennebula_host_vms{one_host_id="3"}`
	query := `opennebula_host_state`

	// Construct the full URL for the query
	fullURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, query)

	// Make the HTTP request
	resp, err := http.Get(fullURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Print the response body
	fmt.Println(string(body))
}

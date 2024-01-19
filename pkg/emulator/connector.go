package emulator

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/client"
	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type HostInfo struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
	Port int    `json:"port"`
}

type EmulatedHost struct {
	HostInfo HostInfo `json:"host"`
	Flavors  []string `json:"flavors"`
}

type EmulatorConnector struct {
	ctrlPlaneHost  string
	ctrolPlanePort int
	envServerHost  string
	envServerPort  int
	prometheusHost string
	prometheusPort int
	client         *resty.Client
	envClient      *client.EnvClient
}

func CreateEmulatorConnector(ctrlPlaneHost string, ctrlPlanePort int, envServerHost string, envServerPort int, prometheusHost string, prometheusPort int) *EmulatorConnector {
	log.WithFields(log.Fields{
		"CtrlPlaneHost":  ctrlPlaneHost,
		"CtrlPlanePort":  ctrlPlanePort,
		"EnvServerHost":  envServerHost,
		"EnvServerPort":  envServerPort,
		"PrometheusHost": prometheusHost,
		"PrometheusPort": prometheusPort}).
		Info("Creating emulator connector")

	return &EmulatorConnector{
		ctrlPlaneHost:  ctrlPlaneHost,
		ctrolPlanePort: ctrlPlanePort,
		envServerHost:  envServerHost,
		envServerPort:  envServerPort,
		prometheusHost: prometheusHost,
		prometheusPort: prometheusPort,
		client:         resty.New(),
		envClient:      client.CreateEnvClient(envServerHost, envServerPort, true)}
}

func checkStatus(statusCode int, body string) error {
	if statusCode != 200 {
		return errors.New(body)
	}

	return nil
}

func (c *EmulatorConnector) sync() error {
	prometheusURL := "http://" + c.prometheusHost + ":" + strconv.Itoa(c.prometheusPort)

	resp, err := c.client.R().
		Get("http://" + c.ctrlPlaneHost + ":" + strconv.Itoa(c.ctrolPlanePort) + "/hosts/flavor")
	if err != nil {
		return err
	}

	err = checkStatus(resp.StatusCode(), string(resp.Body()))
	if err != nil {
		return err
	}

	var emulatedHosts []EmulatedHost
	err = json.Unmarshal(resp.Body(), &emulatedHosts)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error unmarshalling response")
		return err
	}

	hosts, err := c.envClient.GetHosts()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error fetching hosts from env server")
		return err
	}

	hostMaps := make(map[string]bool)
	for _, host := range hosts {
		hostMaps[host.HostID] = true
	}

	vms, err := c.envClient.GetVMs()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error fetching VMs from env server")
		return err
	}

	vmMaps := make(map[string]bool)
	for _, emulatedHost := range emulatedHosts {
		for _, flavor := range emulatedHost.Flavors {
			vmMaps[flavor] = true
		}
	}

	// Add missing hosts
	for _, emulatedHost := range emulatedHosts {
		if _, ok := hostMaps[emulatedHost.HostInfo.Name]; !ok {
			log.WithFields(log.Fields{"host": emulatedHost.HostInfo.Name}).Info("Adding host to env server")
			host, err := c.envClient.GetHost(emulatedHost.HostInfo.Name)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Error fetching host from env server")
				return err
			}
			if host == nil {
				totalCPU, err := GetTotalCPU(prometheusURL, emulatedHost.HostInfo.Name)
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("Error fetching total CPU")
					return err
				}

				totalMemory, err := GetTotalMemory(prometheusURL, emulatedHost.HostInfo.Name)
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("Error fetching total memory")
					return err
				}

				host := &core.Host{HostID: emulatedHost.HostInfo.Name, TotalCPU: totalCPU, TotalMemory: totalMemory}

				err = c.envClient.AddHost(host)
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("Error adding host to env server")
					return err
				}
			}
		}
	}

	// Remove hosts that are not in the emulator
	for _, host := range hosts {
		if _, ok := hostMaps[host.HostID]; !ok {
			log.WithFields(log.Fields{"host": host.HostID}).Info("Removing host from env server")
			err = c.envClient.RemoveHost(host.HostID)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Error removing host from env server")
				return err
			}

		}
	}

	// Add VMs
	for _, emulatedHost := range emulatedHosts {
		for _, flavor := range emulatedHost.Flavors {
			vm, err := c.envClient.GetVM(flavor)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Error fetching VM from env server")
				return err
			}
			if vm == nil {
				vm := &core.VM{VMID: flavor}
				err = c.envClient.AddVM(vm)
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("Error adding VM to env server")
					return err
				}

				err = c.envClient.Bind(vm.VMID, emulatedHost.HostInfo.Name)
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("Error binding VM to host")
					return err
				}

				log.WithFields(log.Fields{"VMID": flavor, "HostID": vm.HostID}).Info("Adding VM to env server")
			}
		}
	}

	// Remove VMs that are not found in the emulator
	for _, vm := range vms {
		if _, ok := vmMaps[vm.VMID]; !ok {
			log.WithFields(log.Fields{"vm": vm.VMID}).Info("Removing unbinding VM")
			err = c.envClient.Unbind(vm.VMID, vm.HostID)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Error unbinding VM")
				return err
			}
		}
	}

	return nil
}

func (c *EmulatorConnector) fetchMetrics() error {
	prometheusURL := "http://" + c.prometheusHost + ":" + strconv.Itoa(c.prometheusPort)

	vms, err := c.envClient.GetVMs()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error fetching VMs from env server")
		return err
	}

	hosts, err := c.envClient.GetHosts()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error fetching hosts from env server")
		return err
	}

	for _, host := range hosts {
		hostMetric, err := GetFlavourMetricForHost(prometheusURL, host.HostID)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Error fetching flavour metric")
			return err
		}
		err = c.envClient.AddMetric(host.HostID, core.HostType, &core.Metric{Timestamp: hostMetric.Timestamp, CPU: hostMetric.CPURate, Memory: hostMetric.MemoryUsage})
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Error adding metric to env server")
			return err
		}
	}

	for _, vm := range vms {
		hostMetric, err := GetFlavourMetricForHost(prometheusURL, vm.HostID) // Assume one VM per host
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Error fetching flavour metric")
			return err
		}
		err = c.envClient.AddMetric(vm.VMID, core.VMType, &core.Metric{Timestamp: hostMetric.Timestamp, CPU: hostMetric.CPURate, Memory: hostMetric.MemoryUsage})
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Error adding metric to env server")
			return err
		}
	}

	return nil
}

func (c *EmulatorConnector) Start() {
	for {
		log.WithFields(log.Fields{"CtrlPlaneHost": c.ctrlPlaneHost, "CtrlPlanePort": c.ctrolPlanePort, "EnvServerHost": c.envServerHost, "EnvServerPort": c.envServerPort, "PrometheusHost": c.prometheusHost, "PrometheusPort": c.prometheusPort}).Info("Syncing")
		err := c.sync()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Error syncing data")
		}

		err = c.fetchMetrics()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Error fetching metrics")
		}

		time.Sleep(1 * time.Second)
	}
}

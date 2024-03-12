package opennebula

import (
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/client"
	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	log "github.com/sirupsen/logrus"
)

type Monitor struct {
	prometheusURL string
	stopFlag      bool
	client        *client.EnvClient
}

func newMonitor(prometheusURL string) *Monitor {
	m := &Monitor{}
	m.prometheusURL = prometheusURL
	m.client = client.CreateEnvClient("localhost", 50080, true)
	return m
}

func (m *Monitor) runForever() {
	go func() {
		for {
			if m.stopFlag {
				return
			}

			hostIDs, err := GetHostIDs(m.prometheusURL)
			if err != nil {
				log.Error("Failed to get host IDs: ", err)
			}

			time.Sleep(1 * time.Second)

			oneHosts := make(map[string]bool)
			var cpuTotal float64
			for _, hostID := range hostIDs {
				// Add hosts to the envserver that are present on opennebula
				hostFromServer, err := m.client.GetHost(hostID)
				if err != nil {
					log.Debug("Failed to get host from server: ", err)
				}

				oneHosts[hostID] = true
				if hostFromServer == nil {
					host := &core.Host{HostID: hostID}
					err = m.client.AddHost(host)
					if err != nil {
						log.Debug("Failed to add host to server: ", err)
					}
				}

				// Collect metric data
				now := time.Now()

				cpuTotal, err = GetHostCPU(m.prometheusURL, hostID)
				if err != nil {
					log.Debug("Failed to get host CPU: ", err)
				}
				usedMem, err := GetHostUsedMem(m.prometheusURL, hostID)
				if err != nil {
					log.Debug("Failed to get host used memory: ", err)
				}
				diskRead, err := GetHostDiskRead(m.prometheusURL, hostID)
				if err != nil {
					log.Debug("Failed to get host available memory: ", err)
				}
				diskWrite, err := GetHostDiskWrite(m.prometheusURL, hostID)
				if err != nil {
					log.Debug("Failed to get host available memory: ", err)
				}
				netRX, err := GetHostNetRX(m.prometheusURL, hostID)
				if err != nil {
					log.Debug("Failed to get host available memory: ", err)
				}
				netTX, err := GetHostNetRX(m.prometheusURL, hostID)
				if err != nil {
					log.Debug("Failed to get host available memory: ", err)
				}

				// Add metric to the envserver
				err = m.client.AddMetric(hostID, core.HostType, &core.Metric{Timestamp: now, CPU: cpuTotal, Memory: usedMem, DiskRead: diskRead, DiskWrite: diskWrite, NetTX: netTX, NetRX: netRX})
				if err != nil {
					log.Error("Failed to add metric to server: ", err)
				}
			}

			serverHosts, err := m.client.GetHosts()
			if err != nil {
				log.Error("Failed to get hosts from server: ", err)
			}

			// TODO: This code is untested
			for _, serverHost := range serverHosts {
				_, ok := oneHosts[serverHost.HostID]
				if !ok {
					err := m.client.RemoveHost(serverHost.HostID)
					if err != nil {
						log.Error("Failed to remove host from server: ", err)
					}
				}
			}
		}
	}()
}

func (m *Monitor) stop() {
	m.stopFlag = true
}

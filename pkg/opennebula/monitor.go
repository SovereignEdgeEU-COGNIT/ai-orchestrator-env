package opennebula

import (
	"fmt"
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/client"
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

			for _, hostID := range hostIDs {
				fmt.Println("Host ID: ", hostID)
				cpuTotal, err := GetHostCPU(m.prometheusURL, hostID)
				if err != nil {
					log.Error("Failed to get host CPU: ", err)
				}
				fmt.Println(" CPU Total: ", cpuTotal)

				// hostFromServer, err := m.client.GetHost(hostID)
				// if err!=nil {
				// }

			}
			fmt.Println()
		}
	}()
}

func (m *Monitor) stop() {
	m.stopFlag = true
}

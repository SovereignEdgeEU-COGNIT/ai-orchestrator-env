package cli

import (
	"os"
	"strconv"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/client"
	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	hostsCmd.AddCommand(addHostCmd)
	hostsCmd.AddCommand(getHostsCmd)
	rootCmd.AddCommand(hostsCmd)

	addHostCmd.Flags().StringVarP(&HostID, "hostid", "", "", "Host Id")
	addHostCmd.MarkFlagRequired("hostid")

	addHostCmd.Flags().StringVarP(&TotalCPU, "totalcpu", "", "", "Total CPU in millicores")
	addHostCmd.MarkFlagRequired("totalcpu")
	addHostCmd.Flags().StringVarP(&TotalMemory, "totalmem", "", "", "Total memory in bytes")
	addHostCmd.MarkFlagRequired("totalmem")
}

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Manage hosts",
	Long:  "Manage hosts",
}

var addHostCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new host",
	Long:  "Add a new host",
	Run: func(cmd *cobra.Command, args []string) {
		parseEnv()

		if HostID == "" {
			log.Fatal("Host Id is required")
		}

		client := client.CreateEnvClient(ServerHost, ServerPort, Insecure)

		totalCPU, err := strconv.ParseInt(TotalCPU, 10, 64)
		CheckError(err)

		totalMem, err := strconv.ParseInt(TotalMemory, 10, 64)
		CheckError(err)

		err = client.AddHost(&core.Host{HostID: HostID, TotalCPU: totalCPU, TotalMemory: totalMem})
		CheckError(err)

		log.WithFields(log.Fields{
			"HostId":      HostID,
			"TotalCPU":    TotalCPU,
			"TotalMemory": TotalMemory}).
			Info("Host added")
	},
}

var getHostsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all hosts",
	Long:  "List all hosts",
	Run: func(cmd *cobra.Command, args []string) {
		parseEnv()

		client := client.CreateEnvClient(ServerHost, ServerPort, Insecure)

		hosts, err := client.GetHosts()
		CheckError(err)

		if len(hosts) == 0 {
			log.Info("No hosts found")
			os.Exit(0)
		}

		printHostsTable(hosts)
	},
}

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
	vmsCmd.AddCommand(addVMCmd)
	vmsCmd.AddCommand(getVMsCmd)
	vmsCmd.AddCommand(bindCmd)
	rootCmd.AddCommand(vmsCmd)

	addVMCmd.Flags().StringVarP(&VMID, "vmid", "", "", "VM Id")
	addVMCmd.MarkFlagRequired("vmid")

	addVMCmd.Flags().StringVarP(&TotalCPU, "totalcpu", "", "", "Total CPU in millicores")
	addVMCmd.MarkFlagRequired("totalcpu")
	addVMCmd.Flags().StringVarP(&TotalMemory, "totalmem", "", "", "Total memory in bytes")
	addVMCmd.MarkFlagRequired("totalmem")

	bindCmd.Flags().StringVarP(&VMID, "vmid", "", "", "VM Id")
	addVMCmd.MarkFlagRequired("vmid")

	bindCmd.Flags().StringVarP(&HostID, "hostid", "", "", "Host Id")
	addVMCmd.MarkFlagRequired("hostid")
}

var vmsCmd = &cobra.Command{
	Use:   "vms",
	Short: "Manage VMs",
	Long:  "Manage VMs",
}

var addVMCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new VM",
	Long:  "Add a new VM",
	Run: func(cmd *cobra.Command, args []string) {
		parseEnv()

		if VMID == "" {
			log.Fatal("VM Id is required")
		}

		client := client.CreateEnvClient(ServerHost, ServerPort, Insecure)

		totalCPU, err := strconv.ParseInt(TotalCPU, 10, 64)
		CheckError(err)

		totalMem, err := strconv.ParseInt(TotalMemory, 10, 64)
		CheckError(err)

		err = client.AddVM(&core.VM{VMID: VMID, TotalCPU: totalCPU, TotalMemory: totalMem})
		CheckError(err)

		log.WithFields(log.Fields{
			"VMId":        VMID,
			"TotalCPU":    TotalCPU,
			"TotalMemory": TotalMemory}).
			Info("VM added")
	},
}

var getVMsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all VMs",
	Long:  "List all VMs",
	Run: func(cmd *cobra.Command, args []string) {
		parseEnv()

		client := client.CreateEnvClient(ServerHost, ServerPort, Insecure)

		vms, err := client.GetVMs()
		CheckError(err)

		if len(vms) == 0 {
			log.Info("No VM found")
			os.Exit(0)
		}

		printVMsTable(vms)
	},
}

var bindCmd = &cobra.Command{
	Use:   "bind",
	Short: "Bind a VM to a host",
	Long:  "Bind a VM to a host",
	Run: func(cmd *cobra.Command, args []string) {
		parseEnv()

		if VMID == "" {
			log.Fatal("VM Id is required")
		}

		if HostID == "" {
			log.Fatal("Host Id is required")
		}

		client := client.CreateEnvClient(ServerHost, ServerPort, Insecure)

		err := client.Bind(VMID, HostID)
		CheckError(err)

		log.WithFields(log.Fields{
			"VMId":   VMID,
			"HostId": HostID}).Info("VM bound")
	},
}

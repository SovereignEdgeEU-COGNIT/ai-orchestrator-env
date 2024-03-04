package cli

import (
	"strconv"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/internal/table"
	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/muesli/termenv"
)

func printVMsTable(vms []*core.VM) {
	t, theme := createTable(0)

	var cols = []table.Column{
		{ID: "VMID", Name: "VMId", SortIndex: 1},
		{ID: "StateID", Name: "StateId", SortIndex: 2},
		{ID: "Deployed", Name: "Deployed", SortIndex: 3},
		{ID: "HostID", Name: "HostID", SortIndex: 4},
		{ID: "HostStateID", Name: "Host StateID", SortIndex: 5},
		{ID: "TotalCPU", Name: "Total CPU", SortIndex: 6},
		{ID: "TotalMem", Name: "Total Mem", SortIndex: 7},
		{ID: "UsageCPU", Name: "Usage CPU", SortIndex: 8},
		{ID: "UsageMem", Name: "Usage Mem", SortIndex: 9},
	}
	t.SetCols(cols)

	for _, vm := range vms {
		row := []interface{}{
			termenv.String(vm.VMID).Foreground(theme.ColorCyan),
			termenv.String(strconv.Itoa(vm.StateID)).Foreground(theme.ColorViolet),
			termenv.String(strconv.FormatBool(vm.Deployed)).Foreground(theme.ColorBlue),
			termenv.String(vm.HostID).Foreground(theme.ColorBlue),
			termenv.String(strconv.Itoa(vm.HostStateID)).Foreground(theme.ColorViolet),
			termenv.String(strconv.FormatFloat(vm.TotalCPU, 'f', -1, 64)).Foreground(theme.ColorMagenta),
			termenv.String(strconv.FormatInt(vm.TotalMemory, 10)).Foreground(theme.ColorMagenta),
			termenv.String(strconv.FormatFloat(vm.UsageCPU, 'f', -1, 64)).Foreground(theme.ColorGreen),
			termenv.String(strconv.FormatInt(vm.UsageMemory, 10)).Foreground(theme.ColorGreen),
		}
		t.AddRow(row)
	}

	t.Render()
}

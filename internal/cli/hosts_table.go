package cli

import (
	"strconv"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/internal/table"
	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/muesli/termenv"
)

func printHostsTable(hosts []*core.Host) {
	t, theme := createTable(0)

	var cols = []table.Column{
		{ID: "HostID", Name: "HostId", SortIndex: 1},
		{ID: "StateID", Name: "StateId", SortIndex: 2},
		{ID: "TotalCPU", Name: "Total CPU", SortIndex: 3},
		{ID: "TotalMem", Name: "Total Mem", SortIndex: 4},
		{ID: "UsageCPU", Name: "Usage CPU", SortIndex: 5},
		{ID: "UsageMem", Name: "Usage Mem", SortIndex: 6},
		{ID: "VMS", Name: "VMs", SortIndex: 7},
	}
	t.SetCols(cols)

	for _, host := range hosts {
		row := []interface{}{
			termenv.String(host.HostID).Foreground(theme.ColorCyan),
			termenv.String(strconv.Itoa(host.StateID)).Foreground(theme.ColorViolet),
			termenv.String(strconv.FormatInt(host.TotalCPU, 10)).Foreground(theme.ColorMagenta),
			termenv.String(strconv.FormatInt(host.TotalMemory, 10)).Foreground(theme.ColorMagenta),
			termenv.String(strconv.FormatInt(host.UsageCPU, 10)).Foreground(theme.ColorGreen),
			termenv.String(strconv.FormatInt(host.UsageMemory, 10)).Foreground(theme.ColorGreen),
			termenv.String(strconv.Itoa(host.VMs)).Foreground(theme.ColorBlue),
		}
		t.AddRow(row)
	}

	t.Render()
}

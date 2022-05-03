package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/utam0k/wsdbg/pkg/crt"
)

func setUpTable(header []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	return table
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "show all workspaces",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := getAddress()
		if err != nil {
			log.Fatalln(err)
		}

		ns, err := getNamespace()
		if err != nil {
			log.Fatalln(err)
		}

		cc, err := crt.NewClient(addr, ns)
		defer cc.Close()

		if err != nil {
			log.Fatalf("cannot connet to cllient: %v", err)
		}

		wss, err := cc.FetchWsContainers()
		if err != nil {
			log.Fatalln(err)
		}

		header := []string{"ID", "Cpu(%)", "CpuLimit", "Memory(%)", "CgroupPath"}
		table := setUpTable(header)

		for _, ws := range wss {
			pss, err := cc.FetchProcessesInWs(ws)
			if err != nil {
				log.Fatalf("failed to get processes: %v", err)
			}

			var (
				cpuUsage float64
				memUsage float32
			)
			for _, ps := range pss {
				cu, err := ps.CPUPercent()
				if err != nil {
					// if a process cannot found, it means process have already existed.
					continue
				}
				cpuUsage += cu
				mu, err := ps.MemoryPercent()
				if err != nil {
					continue
				}
				memUsage += mu
			}

			table.Append([]string{ws.Id, fmt.Sprintf("%.2f", cpuUsage), ws.CpuMax.String(), fmt.Sprintf("%.2f", memUsage), ws.CgroupPath})
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

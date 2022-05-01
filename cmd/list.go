package cmd

import (
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

		header := []string{"ID", "CpuLimit", "CgroupPath"}
		table := setUpTable(header)

		for _, ws := range wss {
			table.Append([]string{ws.Id, ws.CpuMax.String(), ws.CgroupPath})
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

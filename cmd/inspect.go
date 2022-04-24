package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/utam0k/wsdbg/pkg/crt"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <workspaceId>",
	Short: "get a detail infomation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// wsid := args[0]
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

		for _, ws := range wss {
			fmt.Println(ws.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

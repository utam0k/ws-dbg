package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/utam0k/wsdbg/pkg/crt"
)

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

		for _, ws := range wss {
			fmt.Println(ws.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

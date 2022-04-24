package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/utam0k/wsdbg/pkg/crt"
)

var specCmd = &cobra.Command{
	Use:   "spec <workspaceId>",
	Short: "get a spec",
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
		cc.FetchWsContainers()
	},
}

func init() {
	rootCmd.AddCommand(specCmd)
}

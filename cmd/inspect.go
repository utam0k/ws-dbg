package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/utam0k/wsdbg/pkg/crt"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "show detailed information about a workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetId := args[0]
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
			log.Fatalf("cannot connet to client: %v", err)
		}

		wss, err := cc.FetchWsContainers()
		if err != nil {
			log.Fatalln(err)
		}
		var ws crt.Workspace
		for _, w := range wss {
			if w.Id == targetId {
				ws = w
				break
			}
		}
		if ws.Id == "" {
			log.Fatalf("cannot find %s", targetId)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		err = enc.Encode(ws)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

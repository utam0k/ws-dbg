package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
	"github.com/utam0k/wsdbg/pkg/crt"
)

type WorkspaceInfo struct {
	Workspace crt.Workspace
	Processes []ProcessInfo
}

type ProcessInfo struct {
	Pid            int32
	PPid           int32
	Name           string
	CpuUsage       float64
	MemoryUsage    float32
	IoCountersStat *process.IOCountersStat
}

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

		pss, err := cc.FetchProcessesInWs(ws)
		if err != nil {
			log.Fatalf("failed to get processes: %v", err)
		}

		var psInfo []ProcessInfo
		for _, ps := range pss {
			name, err := ps.Cmdline()
			if err != nil {
				continue
			}
			cpuUsage, err := ps.CPUPercent()
			if err != nil {
				continue
			}
			memUsage, err := ps.MemoryPercent()
			if err != nil {
				continue
			}
			ppid, err := ps.Parent()
			if err != nil {
				continue
			}
			ioCounters, err := ps.IOCounters()
			if err != nil {
				continue
			}

			pi := ProcessInfo{
				Pid:            ps.Pid,
				PPid:           ppid.Pid,
				Name:           name,
				CpuUsage:       cpuUsage,
				MemoryUsage:    memUsage,
				IoCountersStat: ioCounters,
			}
			psInfo = append(psInfo, pi)
		}
		wi := WorkspaceInfo{
			Workspace: ws,
			Processes: psInfo,
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		err = enc.Encode(wi)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

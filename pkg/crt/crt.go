package crt

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
	"github.com/utam0k/wsdbg/pkg/cgroup"
)

type Client interface {
	FetchWsContainers() ([]Workspace, error)
	FetchProcessesInWs(ws Workspace) ([]*process.Process, error)
	Close()
}

func NewClient(addr, namespace string) (Client, error) {
	return connectContainerd(addr, namespace)
}

type Workspace struct {
	Id          string
	ContainerId string
	CgroupPath  string
	CpuMax      cgroup.CpuLimit
}

func (w *Workspace) String() string {
	return fmt.Sprintf("%s\t%s\t%s", w.Id, w.CpuMax.String(), w.CgroupPath)
}

package crt

import (
	"time"

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
	MemoryStat  cgroup.MemoryStat
	Uptime      time.Duration
}

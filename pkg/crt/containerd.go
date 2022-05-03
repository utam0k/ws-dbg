package crt

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/shirou/gopsutil/process"
	"github.com/utam0k/wsdbg/pkg/cgroup"
)

const (
	kubernetesNamespace            = "k8s.io"
	dockerNamespce                 = "moby"
	containerLabelCRIKind          = "io.cri-containerd.kind"
	containerLabelK8sContainerName = "io.kubernetes.container.name"
	containerLabelK8sPodName       = "io.kubernetes.pod.name"
)

func connectContainerd(addr string, ns string) (*ContainerdClient, error) {
	cc, err := containerd.New(addr, containerd.WithDefaultNamespace(ns))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to containerd at %s: %w", addr, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = cc.Version(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to containerd: %w", err)
	}

	return &ContainerdClient{
		Client: cc,
	}, nil
}

type ContainerdClient struct {
	Client *containerd.Client
}

func (cc *ContainerdClient) FetchWsContainers() ([]Workspace, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cs, err := cc.Client.ContainerService().List(ctx)
	if err != nil {
		log.Fatal("cannot container list")
	}

	var wss []Workspace
	for _, c := range cs {
		c, err := cc.Client.LoadContainer(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		spec, _ := c.Spec(ctx)

		envs := make(map[string]string)
		for _, e := range spec.Process.Env {
			pair := strings.Split(e, "=")
			if len(pair) < 2 {
				log.Fatal("environment variable parsing fails.")
			}
			envs[pair[0]] = pair[1]
		}

		wsId := envs["GITPOD_WORKSPACE_ID"]
		ws := Workspace{
			Id:          wsId,
			ContainerId: c.ID(),
		}

		if !isWorkspace(ws) {
			continue
		}

		cpuMax, err := cgroup.ReadCpuLimit(spec.Linux.CgroupsPath)
		if err != nil {
			log.Fatalln(err)
		}
		ws.CgroupPath = spec.Linux.CgroupsPath
		ws.CpuMax = cpuMax

		wss = append(wss, ws)
	}

	return wss, nil
}

func isWorkspace(ws Workspace) bool {
	return ws.Id != ""
}

func (cc *ContainerdClient) FetchProcessesInWs(ws Workspace) ([]*process.Process, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := cc.Client.TaskService().ListPids(ctx, &tasks.ListPidsRequest{ContainerID: ws.ContainerId})
	if err != nil {
		return nil, err
	}

	var processes []*process.Process
	for _, p := range res.Processes {
		ps, err := process.NewProcess(int32(p.Pid))
		if err != nil {
			// if a process cannot found, it means process have already existed.
			continue
		}
		processes = append(processes, ps)
	}
	return processes, nil
}

func (cc *ContainerdClient) Close() {
	if err := cc.Client.Close(); err != nil {
		log.Fatalln(err)
	}
}

package crt

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/containerd/containerd"
	runtimespec "github.com/opencontainers/runtime-spec/specs-go"
)

const (
	kubernetesNamespace            = "k8s.io"
	dockerNamespce                 = "moby"
	containerLabelCRIKind          = "io.cri-containerd.kind"
	containerLabelK8sContainerName = "io.kubernetes.container.name"
	containerLabelK8sPodName       = "io.kubernetes.pod.name"
)

func connectContainerd(socketPath string) (*ContainerdClient, error) {
	cc, err := containerd.New(socketPath, containerd.WithDefaultNamespace(dockerNamespce))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to containerd at %s: %w", socketPath, err)
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

func (cc *ContainerdClient) FetchContainers() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cs, err := cc.Client.ContainerService().List(ctx)
	if err != nil {
		log.Fatal("cannot container list")
	}
	fmt.Printf("containers: %d\n", len(cs))
	for _, c := range cs {
		c, _ := cc.Client.LoadContainer(ctx, c.ID)
		spec, _ := c.Spec(ctx)
		if spec.Linux == nil {
			spec.Linux = &runtimespec.Linux{}
		}
		if spec.Linux.Resources == nil {
			spec.Linux.Resources = &runtimespec.LinuxResources{}
		}
		if spec.Linux.Resources.CPU == nil {
			spec.Linux.Resources.CPU = &runtimespec.LinuxCPU{}
		}

		fmt.Printf("cgroupPath: %v\n", spec.Linux.CgroupsPath)
		fmt.Printf("envs: %v\n", spec.Process.Env)
	}
}

func isWorkspace() bool {
	return true
}

package crt

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/containerd"
)

const (
	kubernetesNamespace            = "k8s.io"
	containerLabelCRIKind          = "io.cri-containerd.kind"
	containerLabelK8sContainerName = "io.kubernetes.container.name"
	containerLabelK8sPodName       = "io.kubernetes.pod.name"
)

func connectContainerd(socketPath string) (*ContainerdClient, error) {
	cc, err := containerd.New(socketPath, containerd.WithDefaultNamespace(kubernetesNamespace))
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
	cc.Client.Containers()
}

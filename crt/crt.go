package crt

type Client interface {
	FetchContainers()
}

func NewClient() (Client, error) {
	return connectContainerd("/var/run/docker/containerd/containerd.sock")
}

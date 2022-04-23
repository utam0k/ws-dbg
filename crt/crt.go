package crt

type Client interface {
	FetchContainers()
}

func NewClient() (Client, error) {
	return connectContainerd("/run/containerd/containerd.sock")
}

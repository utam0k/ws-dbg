package crt

import "fmt"

type Client interface {
	FetchWsContainers() ([]Workspace, error)
	Close()
}

func NewClient(addr, namespace string) (Client, error) {
	return connectContainerd(addr, namespace)
}

type Workspace struct {
	Id         string
	CgroupPath string
}

func (w *Workspace) String() string {
	return fmt.Sprintf("%s\t%s", w.Id, w.CgroupPath)
}

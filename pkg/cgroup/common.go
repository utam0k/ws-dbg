package cgroup

import (
	"os"
	"path/filepath"
)

const CgroupRoot = "/sys/fs/cgroup"

type CgroupVersion int

const (
	V1 CgroupVersion = iota
	V2
)

func IsV2() (bool, error) {
	controllers := filepath.Join(CgroupRoot, "cgroup.controllers")
	_, err := os.Stat(controllers)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

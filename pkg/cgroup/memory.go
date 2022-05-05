package cgroup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
)

type MemoryStat struct {
	Usage resource.Quantity
	Limit resource.Quantity
}

func (m *MemoryStat) String() string {
	var (
		limit string
		usage string
	)

	if m.Limit.Equal(resource.Quantity{}) {
		limit = "max"
	} else {
		m.Limit.RoundUp(resource.Mega)
		limit = m.Limit.String()
	}

	usage = m.Usage.String()

	return fmt.Sprintf("%s/%s", usage, limit)
}

func ReadMemoryStat(cgroupPath string) (MemoryStat, error) {
	isv2, err := IsV2()
	if err != nil {
		return MemoryStat{}, fmt.Errorf("failed to detect the cgroup version: %w", err)
	}

	if isv2 {
		return MemoryStat{}, errors.New("Memory of cgroup v2 not supported")
	} else {
		return readMemoryStatV1(filepath.Join(CgroupRoot, "memory", cgroupPath))
	}
}

func readMemoryStatV1(cgroupPath string) (MemoryStat, error) {
	limit, err := readMemoryLimit(cgroupPath)
	if err != nil {
		return MemoryStat{}, err
	}
	return MemoryStat{
		Limit: limit,
	}, nil
}

func readMemoryLimit(cgroupPath string) (resource.Quantity, error) {
	fn := filepath.Join(cgroupPath, "memory.limit_in_bytes")
	fc, err := os.ReadFile(fn)
	if err != nil {
		return resource.Quantity{}, err
	}
	s := strings.TrimSpace(string(fc))
	if s == "max" {
		return resource.Quantity{}, nil
	}
	return resource.MustParse(s), nil
}

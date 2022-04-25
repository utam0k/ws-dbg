package cgroup

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CpuMax struct {
	Quota  time.Duration
	Period time.Duration
}

func (c *CpuMax) String() string {
	return fmt.Sprintf("%s\t%s", c.Quota.String(), c.Period.String())
}

func ReadCpuMax(cgroupPath string) (CpuMax, error) {
	cpuMax, err := os.ReadFile(filepath.Join(string(cgroupPath), "cpu.max"))
	if err != nil {
		return CpuMax{}, fmt.Errorf("unable to read cpu.max: %w", err)
	}

	parts := strings.Fields(string(cpuMax))
	if len(parts) != 2 {
		return CpuMax{}, fmt.Errorf("cpu.max did not have expected number of fields: %s", parts)
	}

	var quota int64
	if parts[0] == "max" {
		quota = math.MaxInt64
	} else {
		quota, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return CpuMax{}, fmt.Errorf("could not parse quota of %s: %w", parts[0], err)
		}
	}

	period, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return CpuMax{}, fmt.Errorf("could not parse period of %s: %w", parts[1], err)
	}

	return CpuMax{
			Quota:  time.Duration(quota) * time.Microsecond,
			Period: time.Duration(period) * time.Microsecond,
		},
		nil
}

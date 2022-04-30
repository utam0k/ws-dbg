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

type CpuLimit struct {
	Quota  time.Duration
	Period time.Duration
}

func (c *CpuLimit) String() string {
	if c.Quota == time.Duration(math.MaxInt) {
		return fmt.Sprintf("max/%s", c.Period.String())
	}

	return fmt.Sprintf("%s/%s", c.Quota.String(), c.Period.String())
}

func ReadCpuLimit(cgroupPath string) (CpuLimit, error) {
	isv2, err := IsV2()
	if err != nil {
		return CpuLimit{}, fmt.Errorf("failed to detect the cgroup version: %w", err)
	}

	if isv2 {
		return readCpuLimitV2(filepath.Join(CgroupRoot, cgroupPath))
	} else {
		return readCpuLimitV1(filepath.Join(CgroupRoot, "cpu", cgroupPath))
	}
}

func readCpuLimitV2(cgroupPath string) (CpuLimit, error) {
	cpuMax, err := os.ReadFile(filepath.Join(string(cgroupPath), "cpu.max"))
	if err != nil {
		return CpuLimit{}, fmt.Errorf("unable to read cpu.max: %w", err)
	}

	parts := strings.Fields(string(cpuMax))
	if len(parts) != 2 {
		return CpuLimit{}, fmt.Errorf("cpu.max did not have expected number of fields: %s", parts)
	}

	var quota int64
	if parts[0] == "max" {
		quota = math.MaxInt64
	} else {
		quota, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return CpuLimit{}, fmt.Errorf("could not parse quota of %s: %w", parts[0], err)
		}
	}

	period, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return CpuLimit{}, fmt.Errorf("could not parse period of %s: %w", parts[1], err)
	}

	return CpuLimit{
			Quota:  time.Duration(quota) * time.Microsecond,
			Period: time.Duration(period) * time.Microsecond,
		},
		nil
}

func readCpuLimitV1(cgroupPath string) (CpuLimit, error) {
	quota, err := readCfsQuota(cgroupPath)
	if err != nil {
		return CpuLimit{}, fmt.Errorf("unable to read cpu quota: %w", err)
	}

	period, err := readCfsPeriod(cgroupPath)
	if err != nil {
		return CpuLimit{}, fmt.Errorf("unable to read cpu period: %w", err)
	}

	return CpuLimit{
		Quota:  quota,
		Period: period,
	}, nil
}

func readCfsPeriod(cgroupPath string) (time.Duration, error) {
	fn := filepath.Join(cgroupPath, "cpu.cfs_period_us")
	s, err := os.ReadFile(fn)
	if err != nil {
		return 0, err
	}

	p, err := strconv.ParseInt(strings.TrimSpace(string(s)), 10, 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(int64(p)) * time.Microsecond, nil
}

func readCfsQuota(cgroupPath string) (time.Duration, error) {
	fn := filepath.Join(cgroupPath, "cpu.cfs_quota_us")
	s, err := os.ReadFile(fn)
	if err != nil {
		return 0, err
	}

	p, err := strconv.ParseInt(strings.TrimSpace(string(s)), 10, 64)
	if err != nil {
		return 0, err
	}

	if p < 0 {
		return time.Duration(math.MaxInt64), nil
	}
	return time.Duration(p) * time.Microsecond, nil
}

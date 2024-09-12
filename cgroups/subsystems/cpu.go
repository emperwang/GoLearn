package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type CpuSubSystem struct {
}

func (c *CpuSubSystem) Name() string {

	return "cpu"
}

func (c *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {

	if subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true); err == nil {
		log.Infof("cpuset get cgroup path: %s", subsystemCgroupPath)
		if res.CpuShare != "" {
			if err := os.WriteFile(path.Join(subsystemCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
				return fmt.Errorf("cpu set limit error, %v", err)
			}
		}
	} else {
		return err
	}

	return nil
}

func (c *CpuSubSystem) Apply(cgroupPath string, pid int) error {

	if subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false); err == nil {
		if err := os.WriteFile(path.Join(subsystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (c *CpuSubSystem) Remove(cgroupPath string) error {
	if subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsystemCgroupPath)
	}
	return nil
}

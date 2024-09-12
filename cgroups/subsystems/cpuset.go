package subsystems

import (
	"os"
	"path"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type CpusetSubsystem struct {
}

func (c *CpusetSubsystem) Name() string {
	return "cpuset"
}

func (c *CpusetSubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true); err == nil {
		log.Infof("cpuset get cgroup path: %s", subsystemCgroupPath)
		if res.CpuSet != "" {
			err := os.WriteFile(path.Join(subsystemCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644)
			if err != nil {
				log.Errorf("cpuset error: %s", err)
			}
			err = os.WriteFile(path.Join(subsystemCgroupPath, "cpuset.mems"), []byte(strconv.Itoa(0)), 0644)
			return err
		}
	} else {
		return err
	}

	return nil
}

func (c *CpusetSubsystem) Apply(cgroupPath string, pid int) error {
	if subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false); err == nil {
		return os.WriteFile(path.Join(subsystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	}
	return nil
}

func (c *CpusetSubsystem) Remove(cgroupPath string) error {
	if subsystemCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsystemCgroupPath)
	}
	return nil
}

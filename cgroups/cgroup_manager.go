package cgroups

import (
	"tutorial/GoLearn/cgroups/subsystems"

	log "github.com/sirupsen/logrus"
)

type CgroupManager struct {
	// 子目录的名字
	Path string

	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {

	return &CgroupManager{
		Path: path,
	}
}

// 添加进程到cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subsys := range subsystems.Subsystems {
		err := subsys.Apply(c.Path, pid)
		if err != nil {
			log.Errorf("manager apply pid error: %v", err)
		}
	}

	return nil
}

// 设置各个subsystem的cgroup限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subsys := range subsystems.Subsystems {
		err := subsys.Set(c.Path, res)
		if err != nil {
			log.Errorf("set %s limit error: %v", subsys.Name(), err)
		}
	}
	return nil
}

// 释放各个subsystem挂载中的cgroup
func (c *CgroupManager) Destory() error {
	for _, subsys := range subsystems.Subsystems {
		if err := subsys.Remove(c.Path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
		}
	}

	return nil
}

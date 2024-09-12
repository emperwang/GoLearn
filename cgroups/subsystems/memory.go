package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type MemorySubSystem struct {
}

func (m *MemorySubSystem) Name() string {
	return "memory"
}

// 设置 cgrouppath 对应的资源限制
func (m *MemorySubSystem) Set(cgrouppath string, res *ResourceConfig) error {

	if subsysCgroupPath, err := GetCgroupPath(m.Name(), cgrouppath, true); err == nil {
		log.Infof("get cgroup path; %s", subsysCgroupPath)
		if res.MemoryLimit != "" {
			// 设置此cgroup的内存限制
			if err := os.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory limit fail %v", err)
			}
			return nil
		}
	} else {
		return err
	}

	return nil
}

// 把进程加入到cgroup中
func (m *MemorySubSystem) Apply(cgrouppath string, pid int) error {
	if subsystemPath, err := GetCgroupPath(m.Name(), cgrouppath, false); err == nil {
		// 把进程PID 写如到 tasks 文件中
		if err := os.WriteFile(path.Join(subsystemPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup process fail %v", err)
		}
	} else {
		return err
	}

	return nil
}

func (m *MemorySubSystem) Remove(cgroupath string) error {
	if subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupath, false); err == nil {
		// 删除 cgroup 便是删除 cgrouppath对应的目录
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

// 通过 /proc/self/mountinfo 找出挂载了某个 subsystem的hierarchy  cgroup根节点所在目录

func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")

		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

// 得到 cgroup的 absolute path
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {

	cgroupRoot := FindCgroupMountPoint(subsystem)
	log.Infof("cgroup path: %s", cgroupRoot)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
				cgroupPath := path.Join(cgroupRoot, cgroupPath)
				return cgroupPath, nil
			} else {
				return "", fmt.Errorf("create cgroup path err: %s", cgroupPath)
			}
		}
	} else {
		return "", fmt.Errorf("Get cgroupPath error: %v", err)
	}
	return path.Join(cgroupRoot, cgroupPath), nil
}

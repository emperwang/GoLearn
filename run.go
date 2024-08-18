package main

import (
	"os"
	"strings"
	"tutorial/GoLearn/cgroups"
	"tutorial/GoLearn/cgroups/subsystems"
	"tutorial/GoLearn/container"

	log "github.com/sirupsen/logrus"
)

func Run(tty bool, command []string, res *subsystems.ResourceConfig, volume, containerName, imageName string, envSlice []string) {
	// pass environment
	parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName, envSlice)

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// record container info
	containerName, err := container.RecordContainerInfo(parent.Process.Pid, command, containerName, volume)

	if err != nil {
		log.Errorf("record container info error, %s", err)
		return
	}
	// 添加cgroup限制
	cgroupManage := cgroups.NewCgroupManager("mydocker-cgroup")

	defer cgroupManage.Destory()
	cgroupManage.Set(res)

	cgroupManage.Apply(parent.Process.Pid)

	sendInitCommand(command, writePipe)
	if tty {
		parent.Wait()
		container.DeleteWorkSpae(volume, containerName)
		container.DeleteRecordInfo(containerName)
	}
}

// 发送init 命令
func sendInitCommand(cmd []string, writePipe *os.File) {
	command := strings.Join(cmd, " ")
	log.Infof("command all is %s", command)

	writePipe.WriteString(command)
	writePipe.Close()

}

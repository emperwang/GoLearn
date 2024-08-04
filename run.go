package main

import (
	"os"
	"strings"
	"tutorial/GoLearn/cgroups"
	"tutorial/GoLearn/cgroups/subsystems"
	"tutorial/GoLearn/container"

	log "github.com/sirupsen/logrus"
)

func Run(tty bool, command []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 添加cgroup限制
	cgroupManage := cgroups.NewCgroupManager("mydocker-cgroup")

	defer cgroupManage.Destory()
	cgroupManage.Set(res)

	cgroupManage.Apply(parent.Process.Pid)

	sendInitCommand(command, writePipe)
	parent.Wait()
	os.Exit(1)
}

// 发送init 命令
func sendInitCommand(cmd []string, writePipe *os.File) {
	command := strings.Join(cmd, " ")
	log.Infof("command all is %s", command)

	writePipe.WriteString(command)
	writePipe.Close()

}

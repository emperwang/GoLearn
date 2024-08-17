package container

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	_ "tutorial/GoLearn/nsenter"

	log "github.com/sirupsen/logrus"
)

const (
	ENV_EXEC_PID = "mydocker_pid"
	ENV_EXEC_CMD = "mydocker_cmd"
)

func ExecContainer(containerName string, commandArray []string) {
	// 根据 containerName 获取 pid
	pid, err := getPidFromContainerName(containerName)
	if err != nil {
		log.Errorf("retrieve container %s pid error %v", containerName, pid)
		return
	}
	log.Infof("retrieve pid %s from %s", pid, containerName)
	// 拼接命令,便于传递
	cmdStr := strings.Join(commandArray, " ")

	log.Infof("container pid %s, cmd %s", pid, cmdStr)

	// 从自身fork出一个新的process,  参数为 exec
	cmd := exec.Command("/proc/self/exe", "exec")

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)

	if err := cmd.Run(); err != nil {
		log.Errorf("exec container %s error %v", containerName, err)
	}
}

func getPidFromContainerName(containerName string) (string, error) {

	dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)
	configPath := path.Join(dirUrl, ConfigName)

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Errorf("read config file failed, %v", err)
		return "", nil
	}
	var containerInfo ContainerInfo

	if err := json.Unmarshal(bytes, &containerInfo); err != nil {
		log.Errorf("json format error %v", err)
		return "", nil
	}

	return containerInfo.Pid, nil
}

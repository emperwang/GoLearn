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
	pid, err := GetPidFromContainerName(containerName)
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
	// 获取PID对应的环境变量, 也就是容器的环境变量
	containerEnvs := GetEnvsByPid(pid)
	cmd.Env = append(os.Environ(), containerEnvs...)
	if err := cmd.Run(); err != nil {
		log.Errorf("exec container %s error %v", containerName, err)
	}
}

func GetPidFromContainerName(containerName string) (string, error) {

	dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)
	configPath := path.Join(dirUrl, ConfigName)

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Errorf("read config file failed, %v", err)
		return "", err
	}
	var containerInfo ContainerInfo

	if err := json.Unmarshal(bytes, &containerInfo); err != nil {
		log.Errorf("json format error %v", err)
		return "", err
	}

	return containerInfo.Pid, nil
}

// 根据PID 获取其环境变量
func GetEnvsByPid(pid string) []string {

	envPath := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := os.ReadFile(envPath)

	if err != nil {
		log.Errorf("read env file %s err %v", envPath, contentBytes)
		return []string{}
	}

	envs := strings.Split(string(contentBytes), "\u0000")

	return envs
}

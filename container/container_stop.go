package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func StopContainer(containerName string) {
	pid, err := GetPidFromContainerName(containerName)

	if err != nil {
		log.Errorf("get container %s pid error %v", containerName, err)
		return
	}
	pidInt, err := strconv.Atoi(pid)

	if err != nil {
		log.Errorf("invalid pid %s %d", pid, pidInt)
		return
	}

	// stop process
	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("stop container %s error %v", containerName, err)
		return
	}

	// update container info
	containerInfo, err := GetContainerInfoFromContainerName(containerName)

	if err != nil {
		log.Errorf("get container info error %v", err)
	}

	containerInfo.Statue = STOP
	containerInfo.Pid = ""
	containerInfo.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	UpdateContainerInfo(*containerInfo)
}

func GetContainerInfoFromContainerName(containerName string) (*ContainerInfo, error) {

	dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)
	configPath := path.Join(dirUrl, ConfigName)

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Errorf("read config file failed, %v", err)
		return nil, err
	}
	var containerInfo ContainerInfo

	if err := json.Unmarshal(bytes, &containerInfo); err != nil {
		log.Errorf("json format error %v", err)
		return nil, err
	}

	return &containerInfo, nil
}

package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
	"tutorial/GoLearn/utils"

	log "github.com/sirupsen/logrus"
)

type ContainerInfo struct {
	Pid        string `json:"pid"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	CreateTime string `json:"createTime"`
	Statue     string `json:"status"`
}

var (
	RUNNING             string = "running"
	STOP                string = "stop"
	EXIT                string = "exit"
	DefaultInfoLocation string = "/var/run/mydocker/%s/"
	ConfigName          string = "config.json"
)

func RecordContainerInfo(containerPid int, commandArray []string, containerName string) (string, error) {
	id := utils.RandStringBytes(10)

	createTime := time.Now().Format("2024-08-17 00:00:00")

	command := strings.Join(commandArray, "")

	// if no container name
	if containerName == "" {
		containerName = id
	}

	containInfo := &ContainerInfo{
		Id:         id,
		Name:       containerName,
		Command:    command,
		CreateTime: createTime,
		Statue:     RUNNING,
	}

	jsonBytes, err := json.Marshal(containInfo)

	if err != nil {
		log.Errorf("marshal containerinfo fail when record container info. %s", err)
		return containerName, err
	}

	jsonString := string(jsonBytes)

	// get absoulte container info file path
	dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)

	if err := os.MkdirAll(dirUrl, 0644); err != nil {
		log.Errorf("create container info path error %s", err)
		return containerName, err
	}

	filePath := path.Join(dirUrl, ConfigName)
	file, err := os.Create(filePath)

	if err != nil {
		log.Errorf("create container info file fail. %s", err)
		return containerName, err
	}

	defer file.Close()

	_, err = file.WriteString(jsonString)

	if err != nil {

		log.Errorf("write info into file error. %s", err)
		return containerName, err
	}

	return containerName, nil
}

func DeleteRecordInfo(containerName string) {
	dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)

	err := os.RemoveAll(dirUrl)
	if err != nil {
		log.Errorf("delete record info fail %s", err)
	}
}

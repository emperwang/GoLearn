package container

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
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
	UpdateTime string `json:"updateTime"`
}

var (
	RUNNING             string = "running"
	STOP                string = "stop"
	EXIT                string = "exit"
	DefaultInfoLocation string = "/var/run/mydocker/%s/"
	ConfigName          string = "config.json"
	ContainerLogFile    string = "container.log"
)

func RecordContainerInfo(containerPid int, commandArray []string, containerName string) (string, error) {
	log.Infof("container parent pid %d", containerPid)
	id := utils.RandStringBytes(10)

	createTime := time.Now().Format("2006-01-02 15:04:05")

	command := strings.Join(commandArray, "")

	// if no container name
	if containerName == "" {
		containerName = id
	}

	containInfo := &ContainerInfo{
		Id:         id,
		Pid:        strconv.Itoa(containerPid),
		Name:       containerName,
		Command:    command,
		CreateTime: createTime,
		Statue:     RUNNING,
		UpdateTime: createTime,
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
	log.Infof("container info : %s", jsonString)
	_, err = file.WriteString(jsonString)

	if err != nil {

		log.Errorf("write info into file error. %s", err)
		return containerName, err
	}

	return containerName, nil
}

func UpdateContainerInfo(containInfo ContainerInfo) error {
	jsonBytes, err := json.Marshal(containInfo)

	if err != nil {
		return fmt.Errorf("marshal containerinfo fail when record container info. %s", err)
	}
	jsonString := string(jsonBytes)

	// get absoulte container info file path
	dirUrl := fmt.Sprintf(DefaultInfoLocation, containInfo.Name)
	filePath := path.Join(dirUrl, ConfigName)

	log.Infof("container info : %s", jsonString)
	os.WriteFile(filePath, jsonBytes, 0644)

	if err != nil {
		return fmt.Errorf("write info into file error. %s", err)
	}

	return nil
}

func DeleteRecordInfo(containerName string) {
	dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)

	err := os.RemoveAll(dirUrl)
	if err != nil {
		log.Errorf("delete record info fail %s", err)
	}
}

func ListAllContainers() {
	dirUrl := fmt.Sprintf(DefaultInfoLocation, "")
	// get all container dir
	dirs, err := os.ReadDir(dirUrl)

	if err != nil {
		log.Errorf("read container dir faile %s", err)
		return
	}

	var containers []*ContainerInfo

	for _, dir := range dirs {
		// get specific container info
		tmpContainer, err := getContainer(dir)
		if err != nil {
			log.Errorf("encounter error when get container info, %s", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)

	// print info title
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATE_TIME\n")

	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", item.Id, item.Name, item.Pid, item.Statue, item.Command, item.CreateTime)
	}

	if err := w.Flush(); err != nil {
		log.Errorf("flush error %v", err)
	}
}

func getContainer(dir fs.DirEntry) (*ContainerInfo, error) {
	containerDir := fmt.Sprintf(DefaultInfoLocation, dir.Name())

	configPath := path.Join(containerDir, ConfigName)

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Errorf("read container info file failed. %s", err)
		return nil, err
	}
	var container ContainerInfo

	err = json.Unmarshal(bytes, &container)
	if err != nil {
		log.Errorf("container info format invalid, %s", err)
		return nil, err
	}

	return &container, nil
}

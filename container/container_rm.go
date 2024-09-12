package container

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func ContainerRemove(containerName string) {
	containerInfo, err := GetContainerInfoFromContainerName(containerName)

	if err != nil {
		log.Errorf("retrieve container info error %v", err)
		return
	}

	if containerInfo.Statue != STOP {
		log.Errorf("can't remove container %s as it is not stop.", containerName)
		return
	}

	dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirUrl); err != nil {
		log.Errorf("remove %s err %v", containerName, err)
	}
}

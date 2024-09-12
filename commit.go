package main

import (
	"fmt"
	"os/exec"
	"path"
	"tutorial/GoLearn/container"

	log "github.com/sirupsen/logrus"
)

func commitContainer(containerName, imagename string) {
	mntUrl := fmt.Sprintf(container.MntDir, containerName)
	imageTar := path.Join(container.RootDir, imagename+".tar")
	log.Infof("tar path: %s", imageTar)

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntUrl, ".").CombinedOutput(); err != nil {
		log.Errorf("tar folder %s error. %v", mntUrl, err)
	}
}

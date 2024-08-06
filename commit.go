package main

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func commitContainer(imagename string) {
	mntUrl := "/root/docker/mnt"
	imageTar := "/root/docker/" + imagename + ".tar"

	log.Infof("tar path: %s", imageTar)

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntUrl, ".").CombinedOutput(); err != nil {
		log.Errorf("tar folder %s error. %v", mntUrl, err)
	}
}

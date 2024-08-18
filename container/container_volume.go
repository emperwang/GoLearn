package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

func NewWorkSpace(volume, imageName, containerName string) {
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateMountPoint(containerName, imageName)

	//  如果传递了volume参数, 则需要挂载 volume
	if volume != "" {
		volumeUrls := volumeUrlExtract(volume)

		if len(volumeUrls) == 2 && volumeUrls[0] != "" && volumeUrls[1] != "" {
			MountVolume(containerName, volumeUrls)
		} else {
			log.Errorf("volume parameter input not correct. %s", volume)
		}
	}
}

func volumeUrlExtract(volume string) []string {
	var volumeUrls []string
	volumeUrls = strings.Split(volume, ":")
	return volumeUrls
}

func MountVolume(containerName string, volumeUrls []string) {
	// 宿主机目录
	parentUrl := volumeUrls[0]
	baseDir := path.Dir(parentUrl)
	workDir := path.Join(baseDir, "work")
	lowerDir := path.Join(baseDir, "lower")
	log.Infof("parentUrl %s, basedir %s, workDir %s", parentUrl, baseDir, workDir)
	if err := os.MkdirAll(parentUrl, 0777); err != nil {
		log.Errorf("Mkdir parent dir %s error. %v", parentUrl, err)
	}

	if err := os.MkdirAll(workDir, 0777); err != nil {
		log.Errorf("Mkdir workdir %s error. %v", workDir, err)
	}
	if err := os.MkdirAll(lowerDir, 0777); err != nil {
		log.Errorf("Mkdir lowerDir %s error. %v", workDir, err)
	}

	// 容器挂载点
	containerUrl := volumeUrls[1]
	mntUrl := fmt.Sprintf(MntDir, containerName)
	containerVolumeUrl := path.Join(mntUrl, containerUrl)

	if err := os.MkdirAll(containerVolumeUrl, 0777); err != nil {
		log.Errorf("create containerVolume dir %s error. %v", containerVolumeUrl, err)
	}

	// 把宿主机文件目录挂载到 容器目录中
	dirs := "lowerdir=" + lowerDir + ",upperdir=" + parentUrl + ",workdir=" + workDir
	log.Infof("mount volume cmd: mount -t overlay -o %s overlay %s", dirs, containerVolumeUrl)
	cmd := exec.Command("mount", "-t", "overlay", "-o", dirs, "overlay", containerVolumeUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("Mount volume fail. %v", err)
	}
}

// 创建只读层
func CreateReadOnlyLayer(imageName string) {
	unTarFolderUrl := path.Join(RootDir, imageName)
	imageTarUrl := path.Join(RootDir, imageName+".tar")

	exist, err := PathExists(unTarFolderUrl)

	if err != nil {
		log.Infof("Fail to judge whether dir %s exists, %v", unTarFolderUrl, err)
	}

	if !exist {
		if err := os.MkdirAll(unTarFolderUrl, 0644); err != nil {
			log.Errorf("Mkdir dir %s error. %v", unTarFolderUrl, err)
		}
		log.Info("begin create read layer")
		if _, err := exec.Command("tar", "-xvf", imageTarUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			log.Errorf("untar dir %s error %v", unTarFolderUrl, err)
		}
	}
}

// 创建读写层
func CreateWriteLayer(containerName string) {
	writeUrl := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.MkdirAll(writeUrl, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeUrl, err)
	}
}

// 为container 创建 mountPoint
func CreateMountPoint(containerName, imageName string) {
	writeUrl := fmt.Sprintf(WriteLayerUrl, containerName)
	readLower := path.Join(RootDir, imageName)
	workUrl := fmt.Sprintf(WorkDir, containerName)
	mntUrl := fmt.Sprintf(MntDir, containerName)
	// 创建mnt文件夹
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntUrl, err)
	}

	if err := os.MkdirAll(workUrl, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntUrl, err)
	}

	// 把writeLayer目录和busybox目录 mount到 mnt目录
	dirs := "lowerdir=" + readLower + "," + "upperdir=" + writeUrl + ",workdir=" + workUrl
	log.Infof("mount cmd: mount -t overlay -o %s overlay %s", dirs, mntUrl)
	cmd := exec.Command("mount", "-t", "overlay", "-o", dirs, "overlay", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteWorkSpae(volume, containerName string) {
	if volume != "" {
		volumeUrls := volumeUrlExtract(volume)

		if len(volumeUrls) == 2 && volumeUrls[0] != "" && volumeUrls[1] != "" {
			DeleteMountPointWithVolume(containerName, volumeUrls)
		}
	}
	DeleteMountPoint(containerName)
	DeleteWriteLayer(containerName)
}

// 卸载容器中的挂载点
func DeleteMountPointWithVolume(containerName string, volumeUrls []string) {
	mntUrl := fmt.Sprintf(MntDir, containerName)
	// 容器中的挂载点
	containerUrl := path.Join(mntUrl, volumeUrls[1])

	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("umount volume fail. %v", err)
	}

	if _, err := exec.Command("umount", mntUrl).CombinedOutput(); err != nil {
		log.Errorf("Umount mountpoint %s failed. error %v", mntUrl, err)
	}

	if err := os.RemoveAll(mntUrl); err != nil {
		log.Errorf("Remove mountpoint dir %s,  error %v", mntUrl, err)
	}
}

func DeleteMountPoint(containerName string) {
	mntUrl := fmt.Sprintf(MntDir, containerName)
	workUrl := fmt.Sprintf(WorkDir, containerName)
	cmd := exec.Command("umount", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("umount point %s fail. %v", mntUrl, err)
	}

	if _, err := exec.Command("umount", workUrl).CombinedOutput(); err != nil {
		log.Errorf("umount point %s fail. %v", workUrl, err)
	}

	if err := os.RemoveAll(mntUrl); err != nil {
		log.Errorf("Remove mnt dir %s error. %v", mntUrl, err)
	}

	if err := os.RemoveAll(workUrl); err != nil {
		log.Errorf("Remove mnt dir %s error. %v", workUrl, err)
	}
}

func DeleteWriteLayer(containerName string) {
	writeUrl := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.RemoveAll(writeUrl); err != nil {
		log.Errorf("Remove write layer %s fail. %v", writeUrl, err)
	}
}

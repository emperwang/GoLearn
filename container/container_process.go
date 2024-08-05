package container

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

/*
*这是父进程, 也就是当前进程执行的内容.
 1. 这里的 /proc/self/exe 调用中, /proc指的是当前运行进程自己的环境. exe其实就是自己调用了自己.
    使用这种方式对创建出来的进程进行初始化

2. 后面的args是参数, 其中 init 是传递给本进程的第一个参数. 在本例中, 其实就是调用initCommand去初始化环境和资源
3. clond 就是fork一个新进程, 并且使用了 namespace 隔离新创建的进程和外部环境
4. 如果执行了 -ti 参数, 就需要把当前进程的输入输出导入到标准输入输出上
*/
func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("new Pipe error: %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	rootDir := "/root/docker"
	mntDir := "/root/docker/mnt"
	NewWorkSpace(rootDir, mntDir, volume)
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Dir = mntDir
	return cmd, writePipe
}

// 命令通过 pipe 传递
func NewPipe() (*os.File, *os.File, error) {

	read, write, err := os.Pipe()

	if err != nil {
		return nil, nil, err
	}

	return read, write, nil
}

func NewWorkSpace(rootUrl string, mntUrl string, volume string) {
	CreateReadOnlyLayer(rootUrl)
	CreateWriteLayer(rootUrl)

	CreateMountPoint(rootUrl, mntUrl)

	//  如果传递了volume参数, 则需要挂载 volume
	if volume != "" {
		volumeUrls := volumeUrlExtract(volume)

		if len(volumeUrls) == 2 && volumeUrls[0] != "" && volumeUrls[1] != "" {
			MountVolume(rootUrl, mntUrl, volumeUrls)
		} else {
			log.Errorf("volume parameter input not correct. %s", &volume)
		}
	}
}

func volumeUrlExtract(volume string) []string {
	var volumeUrls []string
	volumeUrls = strings.Split(volume, ":")
	return volumeUrls
}

func MountVolume(rootUrl string, mntUrl string, volumeUrls []string) {
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

func CreateReadOnlyLayer(rootUrl string) {
	busyboxUrl := path.Join(rootUrl, "busybox")
	busyboxTarUrl := path.Join(rootUrl, "busybox.tar")

	exist, err := PathExists(busyboxUrl)

	if err != nil {
		log.Infof("Fail to judge whether dir %s exists, %v", busyboxUrl, err)
	}

	if exist == false {
		if err := os.MkdirAll(busyboxUrl, 0777); err != nil {
			log.Errorf("Mkdir dir %s error. %v", busyboxUrl, err)
		}
		log.Info("begin create read layer")
		if _, err := exec.Command("tar", "-xvf", busyboxTarUrl, "-C", busyboxUrl).CombinedOutput(); err != nil {
			log.Errorf("untar dir %s error %v", busyboxTarUrl, err)
		}
	}
}

func CreateWriteLayer(rootUrl string) {
	writeUrl := path.Join(rootUrl, "writeLayer")
	if err := os.MkdirAll(writeUrl, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeUrl, err)
	}
}

func CreateMountPoint(rootUrl string, mntUrl string) {
	writeUrl := path.Join(rootUrl, "writeLayer")
	busyboxUrl := path.Join(rootUrl, "busybox")
	workUrl := path.Join(rootUrl, "work")
	// 创建mnt文件夹
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntUrl, err)
	}
	if err := os.MkdirAll(workUrl, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntUrl, err)
	}

	// 把writeLayer目录和busybox目录 mount到 mnt目录
	dirs := "lowerdir=" + busyboxUrl + "," + "upperdir=" + writeUrl + ",workdir=" + workUrl
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

func DeleteWorkSpae(rootUrl string, mntUrl string, volume string) {
	if volume != "" {
		volumeUrls := volumeUrlExtract(volume)

		if len(volumeUrls) == 2 && volumeUrls[0] != "" && volumeUrls[1] != "" {
			DeleteMountPointWithVolume(rootUrl, mntUrl, volumeUrls)
		}
	}
	DeleteMountPoint(rootUrl, mntUrl)
	DeleteWriteLayer(rootUrl)
}

// 卸载容器中的挂载点
func DeleteMountPointWithVolume(rootUrl string, mntUrl string, volumeUrls []string) {
	// 容器中的挂载点
	containerUrl := path.Join(mntUrl, volumeUrls[1])

	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("umount volume fail. %v", err)
	}
}

func DeleteMountPoint(rootUrl string, mntUrl string) {
	cmd := exec.Command("umount", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("umount point fail. %v", err)
	}

	if err := os.RemoveAll(mntUrl); err != nil {
		log.Errorf("Remove mnt dir %s error. %v", mntUrl, err)
	}
}

func DeleteWriteLayer(rootUrl string) {
	writeUrl := path.Join(rootUrl, "writeLayer")
	workUrl := path.Join(rootUrl, "work")
	if err := os.RemoveAll(writeUrl); err != nil {
		log.Errorf("Remove write layer %s fail. %v", writeUrl, err)
	}

	if err := os.RemoveAll(workUrl); err != nil {
		log.Errorf("Remove work layer %s fail. %v", writeUrl, err)
	}
}

package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

/*
这里的init 环境是在容器内部执行的, 也就是说,代码执行到这里后, 容器已经创建出来.
这是容器执行的第一个进程.

使用mount 先重新挂载 proc 文件系统,以便后面通过ps等命令查看进程情况
*/
func RunContainerInitProcess() error {
	cmdArray := readUserCommnad()

	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error. cmdArray is nil")
	}
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loopup path error %v", err)
		return err
	}

	log.Infof("find path %s", path)

	// 初始化挂载点
	SetupMount()

	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}

func readUserCommnad() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(pipe)

	if err != nil {
		log.Errorf("init read pipe error: %v", err)
		return nil
	}

	msgString := string(msg)
	return strings.Split(msgString, " ")
}

func pivoteRoot(newRoot string) error {
	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	// 为了使当前root和新的root 不在同一个文件系统下, 我们把 newRoot 重新mount 以下.
	// bind mount 是把相同的 内容换了一个挂载点的挂载方法
	if err := syscall.Mount(newRoot, newRoot, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		log.Errorf("Mount rootfs to iteself error: %v", err)
		return fmt.Errorf("Mount rootfs to itself error %v", err)
	}

	// 创建put_old
	putOld := filepath.Join(newRoot, ".pivot_put_root")

	if err := os.Mkdir(putOld, 0777); err != nil {
		return fmt.Errorf("create put_old error: %v", err)
	}

	// 切换到新的rootfs
	err := syscall.PivotRoot(newRoot, putOld)
	if err != nil {
		return fmt.Errorf("pivot root dir error: %v", err)
	}

	// 修改当前工作目录到 根目录

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir to / fail, %v", err)
	}

	putOldDir := filepath.Join("/", ".pivot_put_root")

	if err := syscall.Unmount(putOldDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount put_old failed. %v", err)
	}

	return os.Remove(putOldDir)
}

// init 挂载点  // https://h8o75a2u.mirror.aliyuncs.com
func SetupMount() {
	// current path
	pwd, err := os.Getwd()
	if err != nil {

	}

	log.Infof("current location is %s", pwd)

	pivoteRoot(pwd)

	// proc

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	// tmpfs
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=0755")
}

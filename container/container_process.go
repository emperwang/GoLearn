package container

import (
	"os"
	"os/exec"
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
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}

	cmd := exec.Command("/proc/self/exe", args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}

/*
这里的init 环境是在容器内部执行的, 也就是说,代码执行到这里后, 容器已经创建出来.
这是容器执行的第一个进程.

使用mount 先重新挂载 proc 文件系统,以便后面通过ps等命令查看进程情况
*/
func RunContainerInitProcess(command string, args []string) error {
	log.Infof("init command: %s", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	argv := []string{command}

	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}

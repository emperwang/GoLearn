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
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
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
	cmd.ExtraFiles = []*os.File{readPipe}
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

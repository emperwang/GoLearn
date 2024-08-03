package ns

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// mount -t fstype device dir
func NewNsProcess() {
	// 要运行的命令
	cmd := exec.Command("sh")
	// 设置命令运行时的 属性
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,                // 镜像中的uid值
				HostID:      syscall.Getuid(), //  宿主机的ID
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,                // 镜像中的GID
				HostID:      syscall.Getgid(), // 宿主机 gid
				Size:        1,
			},
		},
	}

	// cmd.SysProcAttr.Credential = &syscall.Credential{
	// 	Uid: uint32(0),
	// 	Gid: uint32(0),
	// }

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

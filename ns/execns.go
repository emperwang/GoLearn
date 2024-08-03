package ns

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

// memory subsystem hierarchy 根目录
const cgroupMemoryHierarchMount = "/sys/fs/cgroup/memory"

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

func NewCgroupFunc() {
	fmt.Println(os.Args[0])
	// 容器进程
	log.Printf("current pid %d", syscall.Getpid())
	cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)

	// cmd.SysProcAttr = &syscall.SysProcAttr{}

	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// if err := cmd.Run(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// cmd = exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	//   执行后 进程内存占用大于设置的值, 被kill
	// fork出来的进程在外部namespace映射的PID
	log.Printf("%d \n", cmd.Process.Pid)
	// 在系统默认memory subsystem中创建 cgroup
	os.Mkdir(path.Join(cgroupMemoryHierarchMount, "testmemorylimit"), 0755)

	// 将容器进程加入到 cgroup中
	os.WriteFile(path.Join(cgroupMemoryHierarchMount, "testmemorylimit", "tasks"),
		[]byte(strconv.Itoa(cmd.Process.Pid)), 0644)

	// 限制memory
	os.WriteFile(path.Join(cgroupMemoryHierarchMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)

	cmd.Process.Wait()

}

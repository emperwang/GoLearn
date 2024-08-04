package subsystems

// 资源配置的结构体.
// 包括内存限制, CPU时间片权重, CPU核心数
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

// interface for set cgroups
type Subsystem interface {
	Name() string                               // 资源名称, 如  memory  cpu
	Set(path string, res *ResourceConfig) error // 设置资源 限制
	Apply(path string, pid int) error           // 把 进程PID添加到资源中
	Remove(path string) error                   // 移除资源限制
}

var (
	Subsystems []Subsystem = []Subsystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
		&CpusetSubsystem{},
	}
)

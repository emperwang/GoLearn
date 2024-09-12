package network

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"tutorial/GoLearn/container"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type Network struct {
	Name    string     // 网络名
	IpRange *net.IPNet // 地址段
	Driver  string     // 驱动名
}

func (n *Network) Dump(dumpPath string) error {
	// check file if exists
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0644)
		} else {
			return err
		}
	}
	// 保存的文件名是网络的名字
	nwPth := path.Join(dumpPath, n.Name)
	// 打开保存的文件用于写入
	newFile, err := os.OpenFile(nwPth, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Errorf("open file error %v", err)
		return err
	}

	defer newFile.Close()

	data, err := json.Marshal(n)

	if err != nil {
		log.Errorf("format json error %v", err)
		return err
	}

	_, err = newFile.WriteString(string(data))

	if err != nil {
		log.Errorf("write network info into file error %v", err)
		return err
	}

	return nil
}

func (n *Network) Load(loadpath string) error {
	newConfigFile, err := os.Open(loadpath)

	if err != nil {
		return err
	}
	defer newConfigFile.Close()
	data, err := io.ReadAll(newConfigFile)

	if err != nil {
		log.Errorf("read network config error : %v", err)
		return err
	}

	err = json.Unmarshal(data, n)

	if err != nil {
		log.Errorf("unmatshal json error %v", err)
	}
	return nil
}

func (n *Network) Remove(dumpPath string) error {
	// 检查文件状态
	if _, err := os.Stat(path.Join(dumpPath, n.Name)); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		// 删除文件
		return os.Remove(path.Join(dumpPath, n.Name))
	}
}

var (
	drivers  = map[string]NetworkDriver{}
	networks = map[string]Network{}
)

func RegisterDriver(name string, driver *NetworkDriver) {
	drivers[name] = *driver
}

// 创建网络
func CreateNetwork(driver, subnet, name string) error {
	// 将字符串转换为 net.IPNet 对象
	_, cidr, _ := net.ParseCIDR(subnet)

	// 通过IPAM 获取网关IP.  获取到网段中第一个IP作为网关IP
	gatewayIP, err := ipAllocator.Allocate(cidr)

	if err != nil {
		return err
	}

	cidr.IP = gatewayIP

	// 调用指定的网络驱动创建网络, 这里的drivers 字典是各个网络驱动的实例字典.
	nw, err := (drivers[driver]).Create(cidr.String(), name)

	if err != nil {
		return err
	}
	// 保存网络信息到文件系统中
	return nw.Dump(ipadmDefaultAllocatorPath)
}

// 连接到容器之前创建的 网络 mydocker run -net testnet -p 8080:80  xxx
func Connect(networkName string, cinfo *container.ContainerInfo) error {
	//  从networks 中获取容器要连接的网络信息, networks中保存 已经创建好的网络
	network, ok := networks[networkName]

	if !ok {
		return fmt.Errorf("no sucnh netowkr: %s", networkName)
	}

	ip, err := ipAllocator.Allocate(network.IpRange)
	if err != nil {
		return err
	}

	ep := &EndPoint{
		ID:          fmt.Sprintf("%s-%s", cinfo.Id, networkName),
		IPAddress:   ip,
		Network:     &network,
		PortMapping: nil,
	}

	// 调用 driver.Connect 方法去配置网络端点
	if err := (drivers[network.Driver]).Connect(&network, ep); err != nil {
		return err
	}
	// 进入到容器网络 namespace 配置容器网络设备的IP地址和路由.

	if err := configEndpointIpAddressAndRoute(ep, cinfo); err != nil {
		return err
	}

	// 配置容器到 宿主机 的端口映射.
	return configPortMapping(ep, cinfo)
}

// 配置端口映射
func configPortMapping(ep *EndPoint, cinfo *container.ContainerInfo) error {

	for _, pmap := range ep.PortMapping {
		portMap := strings.Split(pmap, ":")

		if len(portMap) != 2 {
			log.Errorf("port map config error %s", pmap)
			continue
		}

		/*
			通过 iptables 命令来添加 DNAT规则. 将宿主机的端口请求转发到 容器的地址和端口上
		*/

		iptablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s", portMap[0], ep.IPAddress.String(), portMap[1])

		cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)

		output, err := cmd.Output()

		if err != nil {
			log.Errorf("iptable output %s, %v", output, err)
			continue
		}
	}

	return nil
}

// 配置容器网络端点 的地址和路由
func configEndpointIpAddressAndRoute(ep *EndPoint, cinfo *container.ContainerInfo) error {
	// locate PeerName
	peerLink, err := netlink.LinkByName(ep.Device.PeerName)

	if err != nil {
		return err
	}

	/*
		将容器的网络端点加入到容器的网络空间中
		并使这个函数下面的操作 都在对应的 网络空间中进行
		执行完函数后,  恢复为默认的网络空间
	*/
	defer enterContainerNetns(&peerLink, cinfo)()

	/*
		获取到容器的IP 地址及网段, 用于配置 容器内部接口地址
		如: 容器IP 192.168.1.2  网段为 192.168.1.0/24
		那么这里产出的IP 就是 192.168.1.2/24 用于容器中 veth配置
	*/
	interfaceIP := *ep.Network.IpRange
	interfaceIP.IP = ep.IPAddress

	// 设置容器内 veth的IP
	if err := setInterfaceIP(ep.Device.PeerName, interfaceIP.String()); err != nil {
		return fmt.Errorf("%v, %s", ep.Network, err)
	}

	// 启动 peerName
	if err := setInterfaceUp(ep.Device.PeerName); err != nil {
		return err
	}

	// 设置容器内的外部请求都通过容器内的 veth访问
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	// 相当于  route add -net 0.0.0.0/0 gw bridgeIp  dev  容器内veth
	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw:        ep.Network.IpRange.IP,
		Dst:       cidr,
	}

	// 添加route, 相当于  route add
	if err := netlink.RouteAdd(defaultRoute); err != nil {
		return err
	}

	return nil
}

/*
将容器的网络端点加入到容器的网络空间中
并锁定当前程序所执行的线程,  使当前线程进入到容器的网络空间中
返回值是一个函数指针, 执行这个返回函数才会退出容器的 net namespace. 回归到宿主机的 net namespace
*/
func enterContainerNetns(link *netlink.Link, cinfo *container.ContainerInfo) func() {
	/*
		找到容器的 net namespace.  /proc/[pid]/ns/net 打开这个文件的文件描述符就可以在操作 net namespace
		cinfo.PID 即 容器在宿主机上的映射PID
		对应的 /proc/[pid]/ns/net 就是容器内容的 net namespace
	*/
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", cinfo.Pid), os.O_RDONLY, 0)

	if err != nil {
		fmt.Errorf("err get container net namespace : %v", err)
	}

	// get file description handler
	nsFD := f.Fd()

	/*
		锁定当前程序所执行的线程, 如果不锁定操作系统线程的话,G语言的 goroutine 可能会被调度到 别的线程上去
		就不能保证一直在所需要的网络空间中
	*/
	runtime.LockOSThread()

	// 修改 veth的另一端, 将其移动到 容器的net namespace中
	if err := netlink.LinkSetNsFd(*link, int(nsFD)); err != nil {
		log.Errorf("error set link netns: %v", err)
	}

	// 获取当前的 net namespace. 以便后面从容器的 net namespace中退出,回到原来的net namespace中
	origins, err := netns.Get()

	if err != nil {
		log.Errorf("error get current netns: %v", err)
	}

	// 将当前线程进入到容器的net namespace中
	if err := netns.Set(netns.NsHandle(nsFD)); err != nil {
		log.Errorf("errors set netns : %v", err)
	}

	// 返回到之前的 net namespace中
	// 在容器的 net namespace中.  执行完容器配置后, 调用此函数可以将程序恢复到原来的 net namespace中
	return func() {
		// 恢复到原来的 net namespace
		netns.Set(origins)
		origins.Close()

		// 取消当前程序的线程锁定
		runtime.UnlockOSThread()
		f.Close()
	}
}

// 加载所有的 网络配置信息到 networks中
func Init() error {
	// load driver
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver

	if _, err := os.Stat(ipadmDefaultAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(ipadmDefaultAllocatorPath, 0644)
		} else {
			return err
		}
	}

	// 检查目录中所有文件, 并进行遍历
	filepath.Walk(ipadmDefaultAllocatorPath, func(nwPath string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// load file
		_, nwName := path.Split(nwPath)

		nw := &Network{
			Name: nwName,
		}

		if err := nw.Load(nwPath); err != nil {
			log.Errorf("load network config file error: %v", err)
		}

		networks[nwName] = *nw

		return nil
	})

	return nil
}

// list all network info
func ListNetwork() {

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)

	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")

	for _, nw := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n", nw.Name, nw.IpRange, nw.Driver)
	}

	if err := w.Flush(); err != nil {
		log.Errorf("output info error: %v", err)
	}
}

func DeleteNetwork(networkName string) error {
	nw, ok := networks[networkName]

	if !ok {
		return fmt.Errorf("No such network: %s", networkName)
	}

	//1. 回收IP 地址
	if err := ipAllocator.Release(nw.IpRange, &nw.IpRange.IP); err != nil {
		return fmt.Errorf("Error remove network gateway ip: %v", err)
	}

	// 2. 删除网络设备

	if err := (*drivers[nw.Driver]).Delete(nw); err != nil {
		return fmt.Errorf("Error remove network driver: %v", err)
	}

	// 3. 删除配置

	return nw.Remove(ipadmDefaultAllocatorPath)
}

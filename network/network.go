package network

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"
	"tutorial/GoLearn/container"

	log "github.com/sirupsen/logrus"
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

	return nil
}

var (
	drivers = map[string]*NetworkDriver{}
)

func RegisterDriver(name string, driver *NetworkDriver) {
	drivers[name] = driver
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
	nw, err := (*drivers[driver]).Create(cidr.String(), name)

	if err != nil {
		return err
	}
	// 保存网络信息到文件系统中
	return nw.Dump(ipadmDefaultAllocatorPath)
}

var networks = map[string]Network{}

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
	if err := (*drivers[network.Driver]).Connect(&network, ep); err != nil {
		return err
	}
	// 进入到容器网络 namespace 配置容器网络设备的IP地址和路由.

	if err := configEndpointIpAddressAndRoute(ep, cinfo); err != nil {
		return err
	}

	// 配置容器到 宿主机 的端口映射.
	return configPortMapping(ep, cinfo)
}

// 加载所有的 网络配置信息到 networks中
func Init() error {
	// load driver
	var bridgeDriver = &BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = bridgeDriver

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

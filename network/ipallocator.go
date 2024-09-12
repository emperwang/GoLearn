package network

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	ipadmDefaultAllocatorPath = "/var/run/mydocker/network/ipam/subnet.json"
)

type IPAM struct {
	// 文件存放位置
	SubnetAllocatorPath string

	// 网段 和 位图算法的数组map,  kay是网段,  value为分配的位图数组
	Subnets *map[string]string
}

// 初始化一个默认的 Ipallocator
var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipadmDefaultAllocatorPath,
}

// 加载网段地址分配信息
func (ipam *IPAM) Load() error {

	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	// open and read file
	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	if err != nil {
		return err
	}
	defer subnetConfigFile.Close()
	data, err := io.ReadAll(subnetConfigFile)

	// unsharshal json data
	err = json.Unmarshal(data, ipam.Subnets)

	return err
}

func (ipam *IPAM) Dump() error {
	// check file if exists
	ipadmConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipadmConfigFileDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(ipadmConfigFileDir, 0644)
		} else {
			return err
		}
	}

	subnetConfgiFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("open config file error %v", err)
		return err
	}
	defer subnetConfgiFile.Close()
	jsonData, err := json.Marshal(ipam.Subnets)

	if err != nil {
		log.Errorf("format json data error %v", err)
		return err
	}
	_, err = subnetConfgiFile.WriteString(string(jsonData))

	return err
}

// 在网段中 分配一个可用的IP地址
func (ipam *IPAM) Allocate(subnet *net.IPNet) (net.IP, error) {

	// 存放网段中地址分配信息的 映射
	ipam.Subnets = &map[string]string{}

	// load config info
	err := ipam.Load()

	if err != nil {
		return nil, err
	}

	// net.IPNet.mask.size 返回网段的子网掩码的总长度和网段前面的固定为长度
	// 127.0.0.1/8  mask为  255.0.0.0
	// 那么返回的就是 qianm 255的位长度 和总长度 :  8和32
	one, size := subnet.Mask.Size()

	// 如果之前没有分配过, 则要初始化网段的分配配置
	if _, exists := (*ipam.Subnets)[subnet.String()]; !exists {
		/*
			用 0 填满这个网段的配置.  1 << uint8(size-one) 标识这个网段有多少个可用地址
			size-one 是子网掩码后面的网络位数. 2^(size-one) 标识网段中可用IP数
		*/
		(*ipam.Subnets)[subnet.String()] = strings.Repeat("0", 1<<uint8(size-one))
	}
	var ip net.IP = nil
	//
	for idx := range (*ipam.Subnets)[subnet.String()] {

		// 0 未分配
		if (*ipam.Subnets)[subnet.String()][idx] == '0' {
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])

			ipalloc[idx] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)
		}

		// subnet.IP为一个网段的初始化ip, 如 192.168.0.0/16 初始化IP = 192.168.0.0
		ip = subnet.IP

		// 计算IP 地址
		/*
			如网段 172.16.0.0/12, 数组序列为  65525, 那么在 172.16.0.0 上分别加上 [uint8(65535 >> 24), uint8(65535>>16),uint8(65535>>8), uint8(65535>>0)]
		*/

		for t := uint8(4); t > 0; t -= 1 {
			[]byte(ip)[4-t] += uint8(idx >> ((t - 1) * 8))
		}

		ip[3] += 1
		break
	}
	ipam.Dump()
	return ip, nil
}

// ip address release
func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}

	err := ipam.Load()

	if err != nil {
		return err
	}

	// 计算 IP 地址在 网段中的索引位置
	idx := 0

	releaseIp := ipaddr.To4()

	releaseIp[3] -= 1

	for t := uint8(4); t > 0; t -= 1 {
		idx += int(releaseIp[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}

	// 将分配的位图数组中 索引位置 设置为0
	ipalloc := []byte((*ipam.Subnets)[subnet.String()])

	ipalloc[idx] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)

	ipam.Dump()

	return nil
}

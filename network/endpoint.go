package network

import (
	"net"

	"github.com/vishvananda/netlink"
)

type EndPoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"dev"`         // 网络设备
	IPAddress   net.IP           `json:"ip"`          // 设备使用的IP
	MacAddress  net.HardwareAddr `json:"mac"`         // mac 地址
	PortMapping []string         `json:"portmapping"` // 端口映射信息
	Network     *Network         `json:"network"`
}

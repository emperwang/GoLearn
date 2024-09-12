package network

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type BridgeNetworkDriver struct {
}

func (bd BridgeNetworkDriver) Name() string {

	return "bridge"
}

// 创建网络
func (bd BridgeNetworkDriver) Create(subnet, name string) (*Network, error) {

	ip, ipRange, err := net.ParseCIDR(subnet)

	if err != nil {
		log.Errorf("invalide subnet : %s, %v", subnet, err)
		return nil, err
	}

	ipRange.IP = ip

	n := &Network{
		Name:    name,
		IpRange: ipRange,
	}

	err = bd.InitBridge(n)

	return n, err
}

// 删除网络
func (bd BridgeNetworkDriver) Delete(network Network) error {
	bridgeName := network.Name

	// get Link by Name
	bridge, err := netlink.LinkByName(bridgeName)

	if err != nil {
		log.Errorf("locate netlink error: %v", err)
		return err
	}
	return netlink.LinkDel(bridge)
}

// 链接容器网络端点到网络
func (bd *BridgeNetworkDriver) Connect(network *Network, endpoint *EndPoint) error {
	//
	bridgeName := network.Name

	// 通过接口名 获取到 linux bridge 接口对象和接口属性
	bridge, err := netlink.LinkByName(bridgeName)

	if err != nil {
		return err
	}

	// create veth
	linkAttr := netlink.NewLinkAttrs()
	// 取 endpoint ID 的前5位为 接口名
	linkAttr.Name = endpoint.ID[:5]
	// 设置  veth 的一端挂载到 对应的 linux bridge
	linkAttr.MasterIndex = bridge.Attrs().Index

	// create veth, 通过perrName 配置veth 另外一端的接口名
	endpoint.Device = netlink.Veth{
		LinkAttrs: linkAttr,
		PeerName:  "cif" + endpoint.ID[:5],
	}

	// 创建veth
	if err := netlink.LinkAdd(&endpoint.Device); err != nil {
		return fmt.Errorf("Add endpoint device error: %v", err)
	}

	// 调用 LinkSetUp, 启动 veth. 相当于 ip link set xxx up

	err = netlink.LinkSetUp(&endpoint.Device)

	return err
}

// 从网络上移除容器网络端点
func (bd *BridgeNetworkDriver) Disconnect(network Network, endpoint *EndPoint) error {

	return nil
}

// 初始化bridge 设备
func (bd *BridgeNetworkDriver) InitBridge(n *Network) error {

	// 1. 创建 bridge 设备
	bridgeName := bd.Name()

	if err := createBridgeInterface(bridgeName); err != nil {
		log.Errorf("Error add bridge: %s, error: %v", bridgeName, err)
		return err
	}
	// 2. 设置 Bridge 设备的地址和路由
	gateWayIp := n.IpRange
	gateWayIp.IP = n.IpRange.IP

	if err := setInterfaceIP(bridgeName, gateWayIp.String()); err != nil {
		log.Errorf("error assign address: %s on bridge: %s with an error of %v", gateWayIp.String(), bridgeName, err)
		return err
	}
	// 3. 启动 bridge 设备
	if err := setInterfaceUp(bridgeName); err != nil {
		log.Errorf("Error set bridge up: %s, error: %v", bridgeName, err)
		return err
	}

	// 4. 设置iptables 的SNAT的规则

	if err := setupIPtables(bridgeName, n.IpRange); err != nil {
		log.Errorf("set iptables %s for bridge %s, error: %v", n.IpRange, bridgeName, err)
		return err
	}

	return nil
}

// create Linux Bridge device
func createBridgeInterface(bridgeName string) error {
	//  先检查是否已经存在同名的 bridge设备
	_, err := net.InterfaceByName(bridgeName)

	if err == nil || !strings.Contains(err.Error(), "no such network interface") {
		return err
	}

	// 初始化一个 Link 基础对象,  Link的名字即 Bridge 虚拟设备的名字
	la := netlink.NewLinkAttrs()
	la.Name = bridgeName

	//使用刚才创建的 Link 属性 创建netlink的 Bridge 对象
	br := &netlink.Bridge{LinkAttrs: la}

	// 调用 netlink的 LinkAdd 创建  Bridge 虚拟网络设备
	// LinkAdd 用来创建 网络虚拟设备,  相当于  ip link  add xxx
	if err := netlink.LinkAdd(br); err != nil {
		return fmt.Errorf("Bridge creation failed for bridge %s : %v", bridgeName, err)
	}

	return nil
}

// 设置一个网络接口的IP 地址
func setInterfaceIP(name, rawIp string) error {
	// 通过 netlink 的LinkByName 找到网络接口
	iface, err := netlink.LinkByName(name)

	if err != nil {

		return fmt.Errorf("error occur when get interface: %v", err)
	}

	/*
		由于 netlinkParseIPNet 是对 net.ParseCIDR的一个封装, 因此可以将 net.ParseCIDR 返回值的IP 和 net 整合
	*/

	ipNet, err := netlink.ParseIPNet(rawIp)

	if err != nil {
		return err
	}

	// 通过netink.AddrAdd 给网络接口配置地址.  相当于  ip addr add xxx 命令
	// 同时如果配置了 地址所在网段的信息，例如 192.168.0.0/24  还会配置路由表
	// 192.168.0.0/24 转发到这个网络接口上
	addr := &netlink.Addr{
		IPNet: ipNet,
		Label: "",
		Flags: 0,
		Scope: 0,
	}
	return netlink.AddrAdd(iface, addr)
}

func setInterfaceUp(bridgeName string) error {

	iface, err := netlink.LinkByName(bridgeName)

	if err != nil {
		return err
	}

	// 通过 netlink.LinkSetUp 设置状态为Up, 等价于  ip link set xxx up
	if err := netlink.LinkSetUp(iface); err != nil {
		return fmt.Errorf("Error enabling interface for %s: %v", bridgeName, err)
	}

	return nil
}

// 设置 iptables 对应 bridge 的 MASQUERADE 规则
func setupIPtables(bridgename string, subnet *net.IPNet) error {
	// 由于 Go 没有知道iptables 的库, 所以需要通过命令的方式来配置

	// iptables -t nat -A POSTROUTING -s sourceIp ! -o bridgename -j MASQUERADE

	iptablecmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgename)

	cmd := exec.Command("iptables", strings.Split(iptablecmd, " ")...)

	// execute cmd
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("cmd execute fail output %v, error %v", string(output), err)
	}

	return nil
}

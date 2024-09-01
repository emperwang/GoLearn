package network

// mydocker network create--subnet 192.168.0.0/24--driver bridge testbridgenet
// mydocker run-ti-p 80:80--net testbridgenet xxxx
// 网络驱动
type NetworkDriver interface {
	// 驱动名
	Name() string

	// 创建网络
	Create(subnet, name string) (*Network, error)
	// 删除网络
	Delete(netwoek Network) error

	// 链接容器网络端点到网络
	Connect(network *Network, endpoint *EndPoint) error

	// 从网络上移除容器网络端点
	Disconnect(network Network, endpoint *EndPoint) error
}

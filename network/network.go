package network

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"path"

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

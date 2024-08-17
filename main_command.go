package main

import (
	"fmt"
	"tutorial/GoLearn/cgroups/subsystems"
	"tutorial/GoLearn/container"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// 定义了runCommand的flags, 作用类似于运行命令时使用 -- 来指定参数
var runCommand = cli.Command{
	Name: "run",
	Usage: `create a container with namespace and cgroups limit.
				mydocker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "create volume",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
	},
	/*
		run 命令真正执行的函数
		1. 获取参数是否包含 command
		2. 获取用户指定的 command
		3. 调用 run function 去准备启动容器
	*/
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string

		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		detach := context.Bool("d")

		// tty 和  detach 不能共存
		if tty && detach {
			return fmt.Errorf("ti and d parameter can not bot provided")
		}

		// resource limit
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}
		// 获取volume 参数
		volume := context.String("v")
		// retrieve container name
		containerName := context.String("name")
		Run(tty, cmdArray, resConf, volume, containerName)
		return nil
	},
}

// 这里定义了 initCommand 的具体操作
var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process to run User's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("init coming")
		err := container.RunContainerInitProcess()
		return err
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container into a image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			log.Errorf("Missing container name")
		}

		imageName := context.Args().Get(0)
		// commit container
		commitContainer(imageName)

		return nil
	},
}

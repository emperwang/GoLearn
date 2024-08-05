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

		// resource limit
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}
		// 获取volume 参数
		volume := context.String("v")
		Run(tty, cmdArray, resConf, volume)
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

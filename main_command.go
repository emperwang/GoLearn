package main

import (
	"fmt"
	"os"
	"strconv"
	"tutorial/GoLearn/cgroups/subsystems"
	"tutorial/GoLearn/container"
	"tutorial/GoLearn/network"

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
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "set environment",
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
		// get image name
		imageName := cmdArray[0]
		cmdArray = cmdArray[1:]

		// 获取volume 参数
		volume := context.String("v")
		// retrieve container name
		containerName := context.String("name")

		// 获取要设置的环境变量
		envSlice := context.StringSlice("e")

		Run(tty, cmdArray, resConf, volume, containerName, imageName, envSlice)
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
		containerName := context.Args().Get(0)
		imageName := context.Args().Get(1)
		// commit container
		commitContainer(containerName, imageName)

		return nil
	},
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all container info",
	Action: func(context *cli.Context) error {
		container.ListAllContainers()
		return nil
	},
}

var logCommand = cli.Command{
	Name:  "log",
	Usage: "view container log",
	Action: func(context *cli.Context) error {
		if len(context.Args()) != 1 {
			return fmt.Errorf("please input the container name that you want to view")
		}
		container.ViewContainerLog(context.Args().Get(0))
		return nil
	},
}

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "enter process namespace",
	Action: func(context *cli.Context) error {

		// callback function
		if os.Getenv(container.ENV_EXEC_PID) != "" {
			log.Infof("pid callback pid %s", strconv.Itoa(os.Getgid()))
			return nil
		}

		// 命令格式  mydocker exec containername cmd
		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container name or command")
		}

		containerName := context.Args().Get(0)
		var commandArray []string

		for _, arg := range context.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		container.ExecContainer(containerName, commandArray)
		return nil
	},
}

var stopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop container from terminal",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("please input the container name that you want to stop")
		}

		containerName := context.Args().Get(0)
		container.StopContainer(containerName)
		return nil
	},
}

var removeCommand = cli.Command{
	Name:  "rm",
	Usage: "remove container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("please input the container name that you want to remove")
		}
		containerName := context.Args().Get(0)
		container.ContainerRemove(containerName)
		return nil
	},
}

var networkcomand = cli.Command{
	Name:  "network",
	Usage: "network operation",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a container network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "network driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet cidr",
				},
			},
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}
				network.Init()

				err := network.CreateNetwork(context.String("driver"), context.String("subnet"), context.Args()[0])
				return err
			},
		},
		{
			Name:  "list",
			Usage: "list container network",
			Action: func(context *cli.Context) error {
				network.Init()
				network.ListNetwork()
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove container network",
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}

				network.Init()
				err := network.DeleteNetwork(context.Args()[0])
				return err
			},
		},
	},
}

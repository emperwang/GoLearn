package main

import (
	"fmt"
	"os"
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

		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}

// 这里定义了 initCommand 的具体操作
var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process to run User's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("init coming")
		cmd := context.Args().Get(0)
		log.Infof("cmd %s", cmd)

		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	parent.Wait()
	os.Exit(1)
}

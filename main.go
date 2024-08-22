package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Commands = []*cli.Command{
		TuiCommand,
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Infof("app start failed. %v", err)
	}
}

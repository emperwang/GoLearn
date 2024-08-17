package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	const Usage = `mdocker is a simple container runtime implementation`
	app := cli.NewApp()

	app.Name = "mydocker"
	app.Usage = Usage
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		commitCommand,
		listCommand,
		logCommand,
		execCommand,
		stopCommand,
		removeCommand,
	}
	app.Before = func(ctx *cli.Context) error {
		// Log as Json
		log.SetFormatter(&log.TextFormatter{})
		log.SetReportCaller(true)
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

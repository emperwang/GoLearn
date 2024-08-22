package main

import (
	"com.learn/command/tui"
	cli "github.com/urfave/cli/v2"
)

var TuiCommand = &cli.Command{
	Name:  "tui",
	Usage: "tui usuage",

	Action: func(context *cli.Context) error {
		tui.TuiStarter()
		return nil
	},
}

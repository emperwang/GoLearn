package main

import (
	"com.learn/command/trending"
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

var GhQuery = &cli.Command{
	Name:  "gh",
	Usage: "gihug trending",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "lang",
			Usage:    "specific language to query. f.g: java, python,go,javascript",
			Value:    "go",
			Required: false,
		},
		&cli.StringFlag{
			Name:  "format",
			Usage: "output format , f.g: json, table",
			Value: "table",
		},
	},
	Action: func(context *cli.Context) error {
		inputLang := context.String("lang")
		format := context.String("format")
		trending.GhTrendingQuery(inputLang, format)
		return nil
	},
}

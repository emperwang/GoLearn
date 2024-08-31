package main

import (
	"encoding/json"
	"fmt"
	"os"

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
	},
	Action: func(context *cli.Context) error {
		inputLang := context.String("lang")
		GhTrendingQuery(inputLang)
		return nil
	},
}

func GhTrendingQuery(language string) {
	infos := []*trending.TrendingInfo{}
	switch language {
	case "java":
		infos, _ = trending.JavaDefaultGHTrending.Query()
	case "python":
		infos, _ = trending.PythonDefaultGHTrending.Query()
	case "go":
		infos, _ = trending.GoDefaultGHTrending.Query()
	case "javascript":
		infos, _ = trending.NodeJsDefaultGHTrending.Query()
	default:

	}

	data, _ := json.MarshalIndent(infos, "", " ")

	fmt.Fprintf(os.Stdout, "%s", string(data))
}

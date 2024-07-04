package args

import (
	"flag"
	"fmt"
	"os"
)

func ReadFlags() {
	cpath, _ := os.Getwd()
	pathFLag := flag.String("path", cpath, "give path to search")
	lines := flag.Bool("l", false, "count line or not")

	flag.Parse()

	if len(os.Args) <= 2 {
		fmt.Println("please input variables")
		os.Exit(0)
	}

	positionsVar := flag.Args()

	fmt.Println("get path input: ", *pathFLag)
	fmt.Println("count lines or not: ", *lines)
	fmt.Println("other cmdline variables:", positionsVar)
}

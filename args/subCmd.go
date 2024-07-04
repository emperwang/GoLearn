package args

import (
	"flag"
	"fmt"
	"os"
)

func Subcmd() {
	Query := flag.NewFlagSet("query", flag.ExitOnError)
	param := Query.String("param", "", "query parameter")
	qtimeout := Query.Int("timeout", 10, "timeout value")

	Add := flag.NewFlagSet("add", flag.ExitOnError)
	body := Add.String("body", "", "data to add")
	headers := Add.String("header", "", "header to be add into post request. split with , ")

	if len(os.Args) < 3 {
		fmt.Println("please input valid subcmd: query, add")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "query":
		Query.Parse(os.Args[2:])
		fmt.Println("input parameter is : ", *param)
		fmt.Println("timeout : ", *qtimeout)
	case "add":
		Add.Parse(os.Args[2:])
		fmt.Println("body: ", *body)
		fmt.Println("headers: ", *headers)
	default:
		fmt.Println("invalid parameter")
	}
}

package args

import (
	"fmt"
	"os"
)

func ReadCmdArgs() {
	if len(os.Args) < 7 {
		Usage()
		os.Exit(1)
	}
	// read cmd parameters
	var args []string = os.Args
	var name, host, port string
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-name":
			{
				i += 1
				name = args[i]
			}
		case "-host":
			{
				i += 1
				host = args[i]
			}
		case "-port":
			{
				i += 1
				port = args[i]
			}
		default:
			{
				fmt.Printf("Unknown command, i = ", i, ", arg=%v\n", args[i])
			}
		}
	}

	fmt.Printf("your input parameter: name= %s, host= %s, port= %v\n", name, host, port)
}

func Usage() {
	name := os.Args[0]
	fmt.Println("please input like this: ", name, " -name name -host host -port port")
}

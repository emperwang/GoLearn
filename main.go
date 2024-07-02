package main

import (
	"fmt"
	_ "fmt"
	"os"
	"path"
	_ "tutorial/GoLearn/container"
	"tutorial/GoLearn/fileutils"
	_ "tutorial/GoLearn/statement"
	_ "tutorial/GoLearn/variable"
)

// go proxy http://goproxy.cn
func main() {
	// fmt.Println("hello world")
	//variable.TestVariable()
	//statement.IfStatement()
	//statement.ForStatement()
	//statement.SwitchStatement()
	//container.ArrayFunction()
	//container.SliceContainer()
	//container.MapContainer()

	pwd, _ := os.Getwd()
	mainpath := path.Join(pwd, "main.go")
	fileutils.ReadFile(mainpath)

	fmt.Println("hello world")
}

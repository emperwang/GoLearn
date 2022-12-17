package main

import (
	_"tutorial/GoLearn/variable"
	_"tutorial/GoLearn/statement"
	_"tutorial/GoLearn/container"
	"path"
	"os"
	_"fmt"
	"tutorial/GoLearn/fileutils"
)

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
}
package fileutils

import (
	"fmt"
	"os"
	"testing"
	"path"
	"tutorial/GoLearn/fileutils"
)

func TestReadFile(t *testing.T) {
	pwd, _ := os.Getwd()
	mainpath := path.Join(pwd, "..", "main.go")
	mainpath2 := pwd + "\\..\\main.go"
	fmt.Println("pwd: ", pwd, ",mainpath : " ,mainpath, ", mainpath2:",mainpath2)
	fileutils.ReadFile(mainpath2)
}

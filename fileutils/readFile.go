package fileutils

import (
	"fmt"
	"os"
	"bufio"
	"io"
)

func ReadFile(path string) {
	file, err := os.Open(path)
	exists := os.IsExist(err)
	if exists {
		fmt.Println("file ", path," exists.")
	}
	defer file.Close()
	if err != nil {
		fmt.Println("file ", path," open failed. err = ", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(file)
	line,err := reader.ReadString('\n')

	for err == nil && err != io.EOF {
		fmt.Printf("%s",line)
		line,err = reader.ReadString('\n')
	}
	//last line
	fmt.Printf("%s",line)
}







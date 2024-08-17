package container

import (
	"fmt"
	"io"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

func ViewContainerLog(container string) {
	dirUrl := fmt.Sprintf(DefaultInfoLocation, container)

	logPath := path.Join(dirUrl, ContainerLogFile)

	file, err := os.Open(logPath)

	if err != nil {
		log.Errorf("open file error %v when view container log.", err)
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	fmt.Fprint(os.Stdout, string(content))
}

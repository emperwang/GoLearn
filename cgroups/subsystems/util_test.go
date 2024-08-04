package subsystems

import (
	"strings"
	"testing"
)

func TestStringSplit(t *testing.T) {
	var str = "rw,cpuacct,cpu"
	t.Log("test begin")
	for t1, opt := range strings.Split(str, ",") {
		t.Log(t1)
		t.Log(opt)
	}
}

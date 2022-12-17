package async

import (
	"fmt"
)

// 通过 channel来和主线程交互
func outputByAsync(out chan<- int) {
	for i:=0; i<10; i++ {
		fmt.Println("output",i)
	}
	out<- 10
	close(out)
}


func Calloutput(){
	var info chan int = make(chan int, 1)
	go outputByAsync(info)

	// 从channel中读取数据
	// 当channel关闭后,读完数据会退出for
	for i := range info {
		fmt.Println("outputAsync finished. i = ",i)
	}
}


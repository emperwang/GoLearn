package statement

import (
	"fmt"
)

func ForStatement(){

	for i:=0; i<10;i++ {
		fmt.Println("output by for loop.")
	}

	j := 0

	for j < 5 {
		fmt.Println("output by for loop j")
		j++
	}

	var a [3]int = [3]int{3,4,5}
	for index,val := range a {
		fmt.Println("index=", index, "val=", val)
	}
}


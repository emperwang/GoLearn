package container

import(
	"fmt"
)

func ArrayFunction(){

	// 数组定义
	var arr1 [3]int
	// 添加元素到数组中
	arr1[0] = 1
	arr1[1] = 2
	arr1[2] = 3

	arr2 := [3]int{1,2,3}

	arr3 := [...]int{4, 5, 6}

	var arr4 [3]int = [3]int{8,9,0}

	fmt.Println("arr1=",arr1, ", arr2=",arr2, ", arr3=",arr3, ", arr4=",arr4)

}



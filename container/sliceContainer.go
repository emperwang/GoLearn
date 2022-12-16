package container

import (
	"fmt"
)

func SliceContainer(){

	// 从数组中生成切片
	var arr1 [10]int = [10]int{1,2,3,4,5,6,7,8,9,10}
	slice1 := arr1[0:5]

	fmt.Println("slice1 = ", slice1)

	// 直接声明切片
	slice2 := []int{}		// 空切片
	var slice3 []int

	slice2 = append(slice2, 2, 3,4,5)
	slice3 = append(slice3, 4,5,6,7,8)

	fmt.Println("slice2 = ", slice2,", slice3=",slice3)

	// 通过make创建切面
	slice4 := make([]int, 5, 10)
	slice4 = append(slice4,1,2,3,4,5)
	fmt.Println("slice4 = ", slice4)
}
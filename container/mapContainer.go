package container

import (
	"fmt"
)

func MapContainer(){

	// 定义
	mapname := make(map[string]string)

	mapname["key1"] = "value1"
	mapname["key2"] = "value2"
	mapname["key3"] = "value3"

	fmt.Println("mapname = ", mapname)

	// mapname遍历
	for key,val := range mapname {
		fmt.Println("key = ", key,",value=",val)
	}

	// 遍历key
	for key := range mapname {
		fmt.Println("key = ", key)
	}

	// 遍历value
	for _,val := range mapname {
		fmt.Println("value = ", val)
	}
}
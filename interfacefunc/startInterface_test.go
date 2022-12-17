package interfacefunc

import (
	"fmt"
	"testing"
	"tutorial/GoLearn/interfacefunc"
)

func TestWorkFlow(t *testing.T) {
	var wf []interfacefunc.WorkFlow
	phone,_ := interfacefunc.NewPhone("huawei")

	computer,_ := interfacefunc.NewComputer("xiaomi")

	wf = append(wf, &computer)
	wf = append(wf,&phone)

	for _, wf := range wf {
		wf.Init()
		wf.Start()
		wf.Stop()
		wf.Destory()
	}
	for _, v := range wf {
		fmt.Printf("v type = %T\n", v)
		switch v.(type) {
			case *interfacefunc.Phone: {
				ph := v.(*interfacefunc.Phone)
				ph.Init()
				ph.Start()
				ph.Stop()
				ph.Destory()
			}
			case *interfacefunc.Computer:{
				cm := v.(*interfacefunc.Computer)
				cm.Init()
				cm.Start()
				cm.Stop()
				cm.Destory()
			}
			default : {
				fmt.Println("unknow type")
			}
		}
	}
}
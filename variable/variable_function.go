package variable
import (
	"fmt"
)

func TestVariable(){
	// 定义变量并打印
	var a1 int
	a1 = 10
	fmt.Println("a1 = ", a1)
	var b1, b2, b3 int
	b1 = 10
	b2 = 20
	b3 = 30

	fmt.Println("b1=", b1, ", b2=", b2,", b3=", b3)

	var c1, c2,c3 = 40, 50, 60

	fmt.Println("c1=", c1, ", c2=",c2, ", c3=", c3)
	d1, name, d2 := 70, "tom", 80

	fmt.Println("d1=", d1, ", d2=",d2, ", name=", name)
}


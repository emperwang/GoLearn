package OOP

import (
	"fmt"
)

// 结构体声明
// 结构体字段 仍然遵循 首字符大写是导出符号
type User struct {
	Age int
	Address string
	Man bool
}

// create user
func NewUser(age int, addr string, man bool) User {

	entity := User{
		Age: age,
		Address: addr,
		Man: man,
	}

	return entity
}

// u *User 是指针, 即引用传递, 表示操作的和传递进来的是同一个实例
func (u *User) AddAge(step int) bool{
	u.Age += step
	return true
}

// u User 值传递, 表示函数操作的u和真实传递过来的不是同一个实例
func (u User) SetAge(age int ) bool{
	flag := true
	if age > 150 || age<0 {
		fmt.Println("invalid age")
		flag = false
	}else {
		u.Age = age
	}
	return flag
}



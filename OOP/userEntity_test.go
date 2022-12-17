package OOP

import (
	"fmt"
	"testing"
	"tutorial/GoLearn/OOP"
)

func TestUserNew(t *testing.T) {
	user := OOP.NewUser(50, "beijing", true)
	fmt.Println("age=",user.Age, ", address=",user.Address, ", Isman=", user.Man)
}


func TestUserAddAge(t *testing.T) {
	user := OOP.NewUser(50, "beijing", true)
	user.AddAge(50)
	fmt.Println("user = ", user)
}






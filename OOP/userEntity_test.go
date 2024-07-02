package OOP

import (
	"fmt"
	"testing"
)

func TestUserNew(t *testing.T) {
	user := NewUser(50, "beijing", true)
	fmt.Println("age=", user.Age, ", address=", user.Address, ", Isman=", user.Man)
}

func TestUserAddAge(t *testing.T) {
	user := NewUser(50, "beijing", true)
	user.AddAge(50)
	fmt.Println("user = ", user)
}

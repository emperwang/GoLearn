package statement

import (
	"fmt"
)

func SwitchStatement() {

	var input int
	fmt.Print("please input a number within 100: ")

	fmt.Scanf("%d", &input)

	switch {
		case input>=0 && input < 10:{
			fmt.Println("input number between 0 and 10")
		}
		case input>=10 && input < 20:{
			fmt.Println("input number between 10 and 20")
		}
		case input>=20 && input < 30:{
			fmt.Println("input number between 20 and 30")
		}
		case input>=30 && input < 40:{
			fmt.Println("input number between 30 and 40")
		}
		case input>=40 && input < 50:{
			fmt.Println("input number between 40 and 50")
		}
		case input>=50 && input < 60:{
			fmt.Println("input number between 50 and 60")
		}
		case input>=60 && input < 70:{
			fmt.Println("input number between 60 and 70")
		}
		case input>=70 && input < 80:{
			fmt.Println("input number between 70 and 80")
		}
		case input>=80 && input < 90:{
			fmt.Println("input number between 80 and 90")
		}
		case input>=90 && input < 101:{
			fmt.Println("input number between 90 and 100")
		}
		default : {
			fmt.Println("invalid input number.")
		}
	}

	var temp interface{}

	switch temp.(type) {
		case int: {
			fmt.Println("input type is: int")
		}

		case string: {
			fmt.Println("input type is: String`")
		}

		case bool: {
			fmt.Println("input type is: bool")
		}
		case nil: {
			fmt.Println("nil type")
		}
		default: {
			fmt.Println("unknow type")
		}
	}
}

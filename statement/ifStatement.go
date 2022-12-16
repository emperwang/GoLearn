package statement

import (
	"fmt"
	"os"
)

func IfStatement() {
	var n int
	fmt.Println("please input a number within 100: ")
	_, err := fmt.Scanf("%d", &n)

	if(err != nil) {
		fmt.Println("read failed: ", err)
		os.Exit(1)
	}

	if n >= 0 && n < 10 {
		fmt.Println("n between 0 and 10")
	}else if n >= 10 && n < 20 {
		fmt.Println("n between 10 and 20")
	}else if n >= 20 && n < 30 {
		fmt.Println("n between 20 and 30")
	}else if n >= 30 && n < 40 {
		fmt.Println("n between 30 and 40")
	}else if n >= 40 && n < 50 {
		fmt.Println("n between 40 and 50")
	}else if n >= 50 && n < 60 {
		fmt.Println("n between 50 and 60")
	}else if n >= 60 && n < 70 {
		fmt.Println("n between 60 and 70")
	}else if n >= 70 && n < 80 {
		fmt.Println("n between 70 and 80")
	}else if n >= 80 && n < 90 {
		fmt.Println("n between 80 and 90")
	}else if n >= 90 && n < 101 {
		fmt.Println("n between 90 and 100")
	}

}





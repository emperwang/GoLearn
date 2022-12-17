package statement

import (
	"testing"
	"fmt"
	"tutorial/GoLearn/statement"
	_"os"
)


func TestMin(t *testing.T) {
	res := statement.Func1(100,50)
	if res == 50 {
		t.Logf("TestMin expected")
	}else {
		t.Fatalf("TestMin unexpected")
	}
}


func TestMax(t *testing.T) {
	res := statement.Max(100, 300)
	if res == 300 {
		t.Logf("TestMax expected")
	}else {
		t.Fatalf("TestMax unexpected")
	}
}


func TesstGetMinFunction(t *testing.T){
	addFunc := statement.GetAddFunc()
	addFunc(1)
	res := addFunc(1)
	fmt.Printf("res = %v\n", res)
	if res == 2 {
		t.Logf("addfunc expected")
	}else {
		t.Fatalf("addfunc unexpected")
	}
}


func TestProcessNum(t *testing.T) {
	res := statement.ProcessNum(100,200, statement.Max)
	fmt.Println("res = ", res)

	if res == 200 {
		t.Logf("ProcessNum expected")
	}else {
		t.Fatalf("ProcessNum unexpected")
	}
}


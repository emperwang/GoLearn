package OOP

import (
	"fmt"
	"testing"
)

func TestPupilFuncs(t *testing.T) {
	pupil := Pupil{}
	// call person's function
	pupil.ChangeName("tom")
	pupil.SetSex(true)
	// call Pupil's function
	pupil.SetScore(100)
	pupil.SetGrade(5)

	fmt.Println("pupil = ", pupil)
}

func TestCollageStudent(t *testing.T) {
	stu := CollageStudent{}
	// call person's function
	stu.ChangeName("Jason")
	stu.SetSex(true)
	// call Pupil's function
	stu.SetScore(100)
	stu.SetGrade(5)
	// call student's function
	stu.SetStudentId("1234567890")

	fmt.Println("student = ", stu)
}

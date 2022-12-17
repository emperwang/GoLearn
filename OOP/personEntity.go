package OOP

// person声明
type Person struct {
	name string
	sex  bool		// true -> Man
}

func (p *Person) ChangeName(name string) {
	p.name = name
}

func (p *Person) GetName() string {
	return p.name
}

func (p *Person) GetSex() bool {
	return p.sex
}

func (p *Person)SetSex(sex bool) {
	p.sex = sex
}

// 小学生
// 继承 person
type Pupil struct {
	Person
	score int
	grade int
}


func (p *Pupil) GetScore() int {
	return p.score
}

func (p *Pupil) SetScore(score int) {
	p.score = score
}

func (p *Pupil) GetGrade() int {
	return p.grade
}

func (p *Pupil) SetGrade(grade int) {
	p.grade = grade
}


// 大学生
// 继承于 person 和 pupil
// 多继承
// 此处的继承 同时继承了字段和方法 (field and methods)
type CollageStudent struct {
	Person
	Pupil
	studentId string
}

func (p *CollageStudent) GetStudentId() string {

	return p.studentId
}

func (p *CollageStudent) SetStudentId(studentId string)  {
	p.studentId = studentId
}










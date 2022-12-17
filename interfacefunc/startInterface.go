package interfacefunc

import (
	"fmt"
)

type WorkFlow interface {
	Init() (bool, error)
	Start() (bool, error)
	Stop() (bool, error)
	Destory() (bool, error)
}


type Phone struct {
	name string
}

func NewPhone(name string) (Phone, error) {
	return Phone{name: name}, nil
}

func (p *Phone) Init() (bool, error) {
	fmt.Println("phone :", p.name, " init")
	return true, nil
}

func (p *Phone) Start() (bool, error){
	fmt.Println("phone :", p.name, " start")
	return true, nil
}

func (p *Phone) Stop() (bool, error){
	fmt.Println("phone :", p.name, " stop")
	return true, nil
}

func (p *Phone) Destory() (bool, error){
	fmt.Println("phone :", p.name, " destory")
	return true, nil
}


type Computer struct {
	name string
}

func NewComputer(name string) (Computer, error) {
	return Computer{name: name}, nil
}

func (p *Computer) Init() (bool, error) {
	fmt.Println("computer :", p.name, " init")
	return true, nil
}

func (p *Computer) Start() (bool, error){
	fmt.Println("computer :", p.name, " start")
	return true, nil
}

func (p *Computer) Stop() (bool, error){
	fmt.Println("computer :", p.name, " stop")
	return true, nil
}

func (p *Computer) Destory() (bool, error){
	fmt.Println("computer :", p.name, " destory")
	return true, nil
}


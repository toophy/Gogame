package help

import (
	"fmt"
)

var (
	testcases = map[string]interface{}{
		"Npc.SetName": (*Npc).SetName,
		"Npc.SetSex":  (*Npc).SetSex,
		"errstring":   "Can not call this as a function",
		"errnumeric":  123456789,
	}
	funcs = NewFuncs(100)
)

func init() {

	for k, v := range testcases {
		err := funcs.Bind(k, v)
		if k[:3] == "err" {
			if err == nil {
				fmt.Printf("Bind %s: %s", k, "an error should be paniced.")
			}
		} else {
			if err != nil {
				fmt.Printf("Bind %s: %s", k, err)
			}
		}
	}
}

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

// + , - , length, ...

type XObject struct {
	name string
}

func (x *XObject) SetName(name string) {
	x.name = name
	fmt.Println(x.name)
}

type Npc struct {
	XObject
	sex string
}

func (n *Npc) SetSex(sex string) {
	n.sex = sex
	fmt.Println(n.sex)
}

type Functor struct {
	name string
	arg  []interface{}
}

func (f *Functor) Post(name string, arg ...interface{}) {
	f.name = name
	f.arg = arg
}

func (f *Functor) Call() {
	fmt.Println(f.name, f.arg)
	_, err := funcs.Call(f.name, f.arg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

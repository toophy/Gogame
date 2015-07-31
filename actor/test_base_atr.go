package actor

import (
//"fmt"
)

type BaseAtr struct {
	Name string
}

func (a *Actor) BaseAtr_getName() string {
	v := a.Mdl_check(Amdl_BaseAtr)

	if v != nil {
		return v.(*BaseAtr).Name
	}

	return ""
}

func (a *Actor) BaseAtr_setName(name string) {
	v := a.Mdl_check(Amdl_BaseAtr)

	if v != nil {
		v.(*BaseAtr).Name = name
	}
}

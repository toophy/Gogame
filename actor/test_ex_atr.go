package actor

import (
//"fmt"
)

const (
	Exatr_Null = iota
	Exatr_Hp
	Exatr_Mp
	Exatr_Str
	Exatr_Agi
	Exatr_Ene
	Exatr_Last
)

type ExAtr struct {
	Bases   [Exatr_Last - 1]int64
	Changes [Exatr_Last - 1]int64
}

// ExAtr_getBase
func (a *Actor) ExAtr_getBase(id int) int64 {
	if id <= Exatr_Null || id >= Exatr_Last {
		return 0
	}

	v := a.Mdl_check(Amdl_ExAtr)

	if v != nil {
		return v.(*ExAtr).Bases[id]
	}

	return 0
}

func (a *Actor) ExAtr_getChange(id int) int64 {
	if id <= Exatr_Null || id >= Exatr_Last {
		return 0
	}

	v := a.Mdl_check(Amdl_ExAtr)

	if v != nil {
		return v.(*ExAtr).Changes[id]
	}

	return 0
}

func (a *Actor) ExAtr_change(id int, val int64) bool {
	if id <= Exatr_Null || id >= Exatr_Last {
		return false
	}

	v := a.Mdl_check(Amdl_ExAtr)

	if v != nil {
		v.(*ExAtr).Bases[id] = v.(*ExAtr).Bases[id] - v.(*ExAtr).Changes[id] + val
		v.(*ExAtr).Changes[id] = val
		return true
	}

	return true
}

func (a *Actor) ExAtr_clearChange(id int) bool {
	if id <= Exatr_Null || id >= Exatr_Last {
		return false
	}

	v := a.Mdl_check(Amdl_ExAtr)

	if v != nil {
		v.(*ExAtr).Bases[id] = v.(*ExAtr).Bases[id] - v.(*ExAtr).Changes[id]
		v.(*ExAtr).Changes[id] = 0
		return true
	}

	return false
}

package webext

import (
	. "github.com/tmluthfiana/phonebook/controllers"
	routing "github.com/tmluthfiana/phonebook/modules/routing"
)

func RegisterClass() []interface{} {
	base := new(routing.BaseController)

	ret := []interface{}{}
	ret = append(ret, &Phonebook{base})

	return ret
}

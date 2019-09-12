package webext

import (
	. "github.com/tmluthfiana/phonebook/controllers"
	routing "modules/routing"
)

func RegisterClass() []interface{} {
	base := new(routing.BaseController)

	ret := []interface{}{}
	ret = append(ret, &Phonebook{base})

	return ret
}

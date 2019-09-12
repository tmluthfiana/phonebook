package webext

import (
	. "productcategorymanagement/controllers"
	routing "teja/routing"
)

func RegisterClass() []interface{} {
	base := new(routing.BaseController)

	ret := []interface{}{}
	ret = append(ret, &Product{base})
	ret = append(ret, &Category{base})

	return ret
}

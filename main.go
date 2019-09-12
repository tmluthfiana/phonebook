package main

import (
	routing "modules/routing"
	"net/http"
	w "phonebook/webext"

	_ "github.com/eaciit/dbox/dbc/mongo"
)

func main() {
	routing := routing.NewRouting("phonebook/controllers", w.RegisterClass())

	routing.Get("/product/get", "Product.Get")
	routing.Get("/product/view/{id}", "Product.Get")
	routing.Post("/product/save", "Product.Save")
	routing.Put("/product/edit/{id}", "Product.Save")
	routing.Delete("/product/delete/{id}", "Product.Delete")

	http.ListenAndServe(":3030", routing.Routing())
}

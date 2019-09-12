package main

import (
	w "github.com/tmluthfiana/phonebook/webext"
	routing "modules/routing"
	"net/http"

	_ "github.com/eaciit/dbox/dbc/mongo"
)

func main() {
	routing := routing.NewRouting("phonebook/controllers", w.RegisterClass())

	routing.Get("/phonebook/get", "phonebook.Get")
	routing.Get("/phonebook/view/{id}", "phonebook.Get")
	routing.Post("/phonebook/save", "phonebook.Save")
	routing.Put("/phonebook/edit/{id}", "phonebook.Save")
	routing.Delete("/phonebook/delete/{id}", "phonebook.Delete")

	http.ListenAndServe(":3030", routing.Routing())
}

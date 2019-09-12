package main

import (
	routing "github.com/tmluthfiana/phonebook/modules/routing"
	w "github.com/tmluthfiana/phonebook/webext"
	"net/http"

	_ "github.com/eaciit/dbox/dbc/mongo"
)

func main() {
	routing := routing.NewRouting("phonebook/controllers", w.RegisterClass())

	routing.Get("/phonebook/get", "Phonebook.Get")
	routing.Get("/phonebook/view/{id}", "Phonebook.Get")
	routing.Post("/phonebook/save", "Phonebook.Save")
	routing.Put("/phonebook/edit/{id}", "Phonebook.Save")
	routing.Delete("/phonebook/delete/{id}", "Phonebook.Delete")

	http.ListenAndServe(":3030", routing.Routing())
}

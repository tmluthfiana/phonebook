package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	helper "github.com/tmluthfiana/phonebook/helper"
	model "github.com/tmluthfiana/phonebook/model"
	"net/http"
	"testing"
)

func TestPhonebookGet(t *testing.T) {
	payload := struct {
		Take int
		Skip int
	}{
		Take: 10,
		Skip: 0,
	}

	out, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}

	payloadBuff := bytes.NewBufferString(string(out))

	cli := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:3030/phonebook/get", payloadBuff)
	req.Header["Content-type"] = []string{"application/json"}
	if err != nil {
		t.Error("Failed Create Connection")
	}

	resp, err := cli.Do(req)
	if err != nil {
		t.Error("Failed Call Connection")
	}

	response := helper.Result{}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&response)
	t.Log(response)
	if err != nil {
		t.Error("Failed to decode")
	}
}

func TestPhonebookView(t *testing.T) {
	Id := "5d7a50e24db82327ee59c456"
	payload := struct {
		Take int
		Skip int
	}{
		Take: 10,
		Skip: 0,
	}

	out, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}

	payloadBuff := bytes.NewBufferString(string(out))

	cli := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:3030/phonebook/view/"+Id, payloadBuff)
	req.Header["Content-type"] = []string{"application/json"}
	if err != nil {
		t.Error("Failed Create Connection")
	}

	resp, err := cli.Do(req)
	if err != nil {
		t.Error("Failed Call Connection")
	}

	response := helper.Result{}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&response)
	t.Log(response)
	if err != nil {
		t.Error("Failed to decode")
	}
}

func TestPhonebookSave(t *testing.T) {

	Phone := model.PhoneNumberDetail{
		PhoneNo:   "081317595876",
		ProneType: "Mobile",
		PhoneExt:  "",
	}

	PhoneDetail := []model.PhoneNumberDetail{}
	PhoneDetail = append(PhoneDetail, Phone)

	payload := model.Phonebook{
		FirstName:   "Tias",
		LastName:    "Faluthi",
		PhoneNumber: PhoneDetail,
		Email:       "triasluth@gmail.com",
	}

	// fmt.Println(fmt.Sprintf("payload TES %+v", payload))

	out, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}

	fmt.Println("OUT >> ", string(out))

	payloadBuff := bytes.NewBufferString(string(out))

	cli := http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3030/phonebook/save", payloadBuff)
	req.Header["Content-type"] = []string{"application/json"}
	if err != nil {
		t.Error("Failed Create Connection")
	}

	resp, err := cli.Do(req)
	if err != nil {
		t.Error("Failed Call Connection")
	}

	response := model.Phonebook{}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&response)
	t.Log(response)
	if err != nil {
		t.Error("Failed to decode")
	}

	if response.Id == "" {
		t.Error("Failed to Save")
	}
}

func TestPhonebookEdit(t *testing.T) {
	// for edit current phone number
	// Phone := model.PhoneNumberDetail{
	// 	PhoneNo:   "08223009617",
	// 	ProneType: "Mobile",
	// 	PhoneExt:  "",
	// }
	//PhoneDetail := []model.PhoneNumberDetail{}
	// PhoneDetail = append(PhoneDetail, Phone)

	// for add new phone number
	PhoneDetail := []model.PhoneNumberDetail{
		{
			PhoneNo:   "08223009617",
			ProneType: "Mobile",
			PhoneExt:  "",
		},
		{
			PhoneNo:   "08123009615",
			ProneType: "Mobile",
			PhoneExt:  "",
		},
	}

	payload := model.Phonebook{
		FirstName:   "Agil",
		LastName:    "D",
		PhoneNumber: PhoneDetail,
	}

	out, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}
	payloadBuff := bytes.NewBufferString(string(out))

	cli := http.Client{}
	req, err := http.NewRequest(http.MethodPut, "http://localhost:3030/phonebook/edit/5d7a50e24db82327ee59c456", payloadBuff)
	req.Header["Content-type"] = []string{"application/json"}
	if err != nil {
		t.Error("Failed Create Connection")
	}

	resp, err := cli.Do(req)
	if err != nil {
		t.Error("Failed Call Connection")
	}

	response := model.Phonebook{}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&response)
	t.Log(response)
	if err != nil {
		t.Error("Failed to decode")
	}

	if response.Id == "" {
		t.Error("Failed to Save")
	}
}

func TestPhonebookDelete(t *testing.T) {
	Payload := `{}`
	payloadBuff := bytes.NewBufferString(Payload)

	cli := http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3030/phonebook/delete/5d7ad69b7505711d7a832849", payloadBuff)
	req.Header["Content-type"] = []string{"application/json"}
	if err != nil {
		t.Error("Failed Create Connection")
	}

	resp, err := cli.Do(req)
	if err != nil {
		t.Error("Failed Call Connection")
	}

	response := model.Phonebook{}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&response)
	t.Log(response)
	if err != nil {
		t.Error("Failed to decode")
	}

	if response.Id == "" {
		t.Error("Failed to Delete")
	}
}

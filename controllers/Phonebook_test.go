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
	payload := model.Phonebook{}

	payload.FirstName = "Agil"
	payload.LastName = "D"

	fmt.Println(fmt.Sprintf("payload TES %+v", payload))

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
	payload := model.Phonebook{}

	payload.FirstName = "Agil"
	payload.LastName = "D"

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
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3030/phonebook/delete/5d7a50e24db82327ee59c456", payloadBuff)
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

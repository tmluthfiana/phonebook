package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	helper "phonebook/helper"
	model "phonebook/model"
	"testing"
)

func TestCategoryGet(t *testing.T) {
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

func TestCategoryView(t *testing.T) {
	Id := "5d726cf13b3eb739d0a7271f"
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

func TestCategorySave(t *testing.T) {
	payload := model.Category{}

	payload.Name = "Tes"
	payload.Code = "Tes"

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

	response := model.Product{}
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

func TestCategoryEdit(t *testing.T) {
	payload := model.Category{}

	payload.Name = "Rajman xxx"
	payload.Code = "axsx"

	out, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}
	payloadBuff := bytes.NewBufferString(string(out))

	cli := http.Client{}
	req, err := http.NewRequest(http.MethodPut, "http://localhost:3030/phonebook/edit/5d7277dc3b3eb74388d4f289", payloadBuff)
	req.Header["Content-type"] = []string{"application/json"}
	if err != nil {
		t.Error("Failed Create Connection")
	}

	resp, err := cli.Do(req)
	if err != nil {
		t.Error("Failed Call Connection")
	}

	response := model.Product{}
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

func TestCategoryDelete(t *testing.T) {
	Payload := `{}`
	payloadBuff := bytes.NewBufferString(Payload)

	cli := http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3030/phonebook/delete/5d7277dc3b3eb74388d4f289", payloadBuff)
	req.Header["Content-type"] = []string{"application/json"}
	if err != nil {
		t.Error("Failed Create Connection")
	}

	resp, err := cli.Do(req)
	if err != nil {
		t.Error("Failed Call Connection")
	}

	response := model.Product{}
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

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	helper "phonebook/helper"
	model "phonebook/model"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestGet(t *testing.T) {
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

func TestView(t *testing.T) {
	Id := "5d72709a3b3eb70eb0753e05"
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

func TestSave(t *testing.T) {
	payload := model.Product{}

	payload.Name = "Tes"
	payload.Code = "Tes"
	payload.Category = &model.Category{
		Id:   bson.ObjectIdHex("5d72709a3b3eb70eb0753e05"),
		Code: "xx",
		Name: "xx",
	}

	pricing := model.ProductPricing{
		Price:     989000,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	arPricing := []*model.ProductPricing{}
	arPricing = append(arPricing, &pricing)

	payload.Pricing = arPricing

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

func TestEdit(t *testing.T) {
	payload := model.Product{}

	payload.Name = "Rajman xxx"
	payload.Code = "axsx"
	payload.Category = &model.Category{
		Id:   bson.ObjectIdHex("5d72709a3b3eb70eb0753e05"),
		Code: "xx",
		Name: "xx",
	}

	pricing := model.ProductPricing{
		Price:     989000,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	arPricing := []*model.ProductPricing{}
	arPricing = append(arPricing, &pricing)

	payload.Pricing = arPricing

	fmt.Println(fmt.Sprintf("payload TES %+v", payload))

	out, err := json.Marshal(&payload)
	if err != nil {
		panic(err)
	}
	payloadBuff := bytes.NewBufferString(string(out))

	cli := http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3030/phonebook/edit/5d7276b03b3eb74388d4f286", payloadBuff)
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

func TestDelete(t *testing.T) {
	Payload := `{}`
	payloadBuff := bytes.NewBufferString(Payload)

	cli := http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3030/phonebook/delete/5d7276b03b3eb74388d4f286", payloadBuff)
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

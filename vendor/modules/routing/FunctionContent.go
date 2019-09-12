package routing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

type WeContent struct {
	Writer http.ResponseWriter
	Req    *http.Request
	vars   map[string]string
}

func (f *WeContent) Parse(d interface{}) error {
	v := reflect.ValueOf(d)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("Invalid controller object passed (%s). Controller object should be a pointer", v.Kind())
	}

	r := f.Req
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}

	if err := json.Unmarshal(body, d); err != nil {
		if err := json.NewEncoder(f.Writer).Encode(err); err != nil {
			return err
		}
	}

	return nil
}

func (f *WeContent) VarsGet(k string) (string, error) {
	if v, isexist := f.vars[k]; isexist {
		return v, nil
	}

	return "", errors.New("Not Found")
}

func (f *WeContent) JSON(d interface{}) []byte {
	f.Writer.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(d)
	if err != nil {
		return nil
	}

	return js
}

func (f *WeContent) NotFound(er error) interface{} {
	return f.error(404, er.Error())
}

func (f *WeContent) ServerError(er error) interface{} {
	return f.error(500, er.Error())
}

func (f *WeContent) error(Code int, Message string) []byte {
	f.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	f.Writer.WriteHeader(Code) // unprocessable entity

	return []byte(Message)
}

func (f *WeContent) Return(d []byte) {
	f.Writer.Write(d)
}

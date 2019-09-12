package controllers

import (
	"errors"
	"fmt"
	routing "modules/routing"
	helper "phonebook/helper"
	model "phonebook/model"

	"gopkg.in/mgo.v2/bson"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type Phonebook struct {
	*routing.BaseController
}

func (p *Phonebook) Get(r *routing.WeContent) interface{} {
	frm := struct {
		Id   string
		Take int
		Skip int
		Sort []tk.M
	}{}
	if e := r.Parse(&frm); e != nil {
		return r.ServerError(e)
	}

	v, err := r.VarsGet("id")
	if err == nil {
		frm.Id = v
	}

	qry := tk.M{
		"limit": frm.Take,
		"skip":  frm.Skip,
	}

	var dbFilter []*db.Filter

	if frm.Id != "" {
		dbFilter = append(dbFilter, db.Eq("_id", bson.ObjectIdHex(v)))
	}

	if len(dbFilter) > 0 {
		qry.Set("where", db.And(dbFilter...))
	}

	conn, err := helper.ConnectToDB()
	defer conn.Close()
	if err != nil {
		return r.ServerError(err)
	}

	ctx := orm.New(conn)
	crs, err := ctx.Find(new(model.Phonebook), qry)
	defer crs.Close()
	if err != nil {
		return r.ServerError(err)
	}

	data := make([]model.Phonebook, 0)
	err = crs.Fetch(&data, 0, false)
	if err != nil {
		return r.ServerError(err)
	}

	qTotal := ctx.Connection.NewQuery()
	qTotal.Where(dbFilter...)
	crs, err = qTotal.Aggr(db.AggrSum, 1, "Count").From(new(model.Phonebook).TableName()).Group("").Cursor(nil)
	defer crs.Close()
	if err != nil {
		return r.ServerError(err)
	}

	total := 0
	tkm := tk.M{}
	crs.Fetch(&tkm, 1, false)
	if tkm != nil {
		total = tkm.GetInt("Count")
	}

	res := helper.NewResult()

	if frm.Id != "" {
		if len(data) > 0 {
			return r.JSON(res.SetData(data[0]).SetTotal(total))
		} else {
			return r.NotFound(errors.New("ID not found"))
		}
	}

	return r.JSON(res.SetData(data).SetTotal(total))
}

func (p *Phonebook) Save(r *routing.WeContent) interface{} {
	model := model.Phonebook{}
	if e := r.Parse(&model); e != nil {
		return r.ServerError(e)
	}

	fmt.Println(fmt.Sprintf("xxxx %+v", model))

	v, err := r.VarsGet("id")
	if err == nil {
		model.Id = bson.ObjectIdHex(v)
	}

	if model.Code == "" {
		return r.ServerError(errors.New("Code is required"))
	}

	if model.Name == "" {
		return r.ServerError(errors.New("Name is required"))
	}

	if model.Category == nil {
		return r.ServerError(errors.New("Category Required"))
	}

	if len(model.Pricing) == 0 {
		return r.ServerError(errors.New("Pricing Required"))
	}

	for _, prc := range model.Pricing {
		if prc.Id == "" {
			prc.Id = bson.NewObjectId()
		}
	}

	conn, err := helper.ConnectToDB()
	defer conn.Close()
	if err != nil {
		return r.ServerError(err)
	}

	tk.Printfn("model %+v", model.Category)
	tk.Printfn("model %+v", model)

	err = helper.SaveRecord(&model)
	if err != nil {
		return r.ServerError(err)
	}

	return r.JSON(model)
}

func (c *Phonebook) Delete(r *routing.WeContent) interface{} {
	v, err := r.VarsGet("id")
	if err == nil {
		tk.Printfn("r %+v", v)
	}

	model := model.Phonebook{
		Id: bson.ObjectIdHex(v),
	}
	if err := helper.DeleteRecord(&model); err != nil {
		return r.ServerError(err)
	}

	return r.JSON(model)
}

package model

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type Phonebook struct {
	orm.ModelBase `bson:"-" json:"-"`
	Id            bson.ObjectId `bson:"_id" json:"_id"`
	FirstName     string        `bson:"FirstName" json:"FirstName"`
	LastName      string        `bson:"LastName" json:"LastName"`
	PhoneNumber   []PhoneNumberDetail
	Email         string `bson:"Email" json:"Email"`
	LastAction    string
	Status        string
	CreatedDate   time.Time
	CreatedBy     string
	UpdateDate    time.Time
	UpdateBy      string
}

func (e *Phonebook) PreSave() error {
	if e.Id == "" {
		e.Id = bson.NewObjectId()
		e.CreatedDate = time.Now()
		e.LastAction = "insert"
	} else {
		e.UpdateDate = time.Now()
		e.LastAction = "update"
	}

	return nil
}

type PhoneNumberDetail struct {
	PhoneNo   string `bson:"PhoneNo" json:"PhoneNo"`
	ProneType string `bson:"ProneType" json:"ProneType"`
	PhoneExt  string `bson:"PhoneExt" json:"PhoneExt"`
}

func (e *Phonebook) RecordID() interface{} {
	return e.Id
}

func (m *Phonebook) TableName() string {
	return "Phonebook"
}

package mongo

import (
	"errors"
	"fmt"
	"github.com/eaciit/dbox"
	"github.com/eaciit/errorlib"
	. "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2"
	"reflect"
)

const (
	modCursor = "Cursor"

	QueryResultCursor = "MongoCursor"
	QueryResultPipe   = "MongoPipe"
)

type Cursor struct {
	dbox.Cursor

	ResultType string
	mgoCursor  *mgo.Query
	mgoPipe    *mgo.Pipe
	mgoIter    *mgo.Iter

	count int

	session          *mgo.Session
	isPoolingSession bool
}

func (c *Cursor) Close() {
	if c.mgoIter != nil {
		c.mgoIter.Close()
	}

	if c.session != nil && !c.isPoolingSession {
		c.session.Close()
	}
}

func (c *Cursor) validate() error {
	if c.ResultType == QueryResultPipe {
		if c.mgoPipe == nil {
			return errors.New(fmt.Sprintf("Pipe is nil"))
		}
	} else if c.ResultType == QueryResultCursor {
		if c.mgoCursor == nil {
			return errors.New("Query cursor is nil")
		}
	}

	return nil
}

func (c *Cursor) prepIter() error {
	e := c.validate()
	if e != nil {
		return e
	}
	if c.mgoIter == nil {
		if c.ResultType == QueryResultPipe {
			//fmt.Println("Preparing Pipe Iter")
			c.mgoIter = c.mgoPipe.Iter()
		} else if c.ResultType == QueryResultCursor {
			//fmt.Println("Preparing Query Iter")
			c.mgoIter = c.mgoCursor.Iter()
		}
	}
	return nil
}

func (c *Cursor) Count() int {
	return c.count
}

func (c *Cursor) ResetFetch() error {
	c.mgoIter = nil
	e := c.prepIter()
	if e != nil {
		return errorlib.Error(packageName, modCursor, "ResetFetch", e.Error())
	}
	return nil
}

func (c *Cursor) Fetch(m interface{}, n int, closeWhenDone bool) error {
	if closeWhenDone {
		defer c.Close()
	}

	e := c.prepIter()
	if e != nil {
		return errorlib.Error(packageName, modCursor, "Fetch", e.Error())
	}

	if c.mgoIter == nil {
		return errorlib.Error(packageName, modCursor, "Fetch", "Iter object is not yet initialized")
	}

	if !IsPointer(m) {
		return errorlib.Error(packageName, modCursor, "Fetch", "Model object should be pointer")
	}

	//ds := dbox.NewDataSet(m)
	//	var datas []interface{}
	//	rt := reflect.TypeOf(m)
	//datias := reflect.MakeSlice(rt, 0, 0)
	if n == 0 {
		//datas := []interface{}{}
		e = c.mgoIter.All(m)
		if e != nil {
			return errorlib.Error(packageName, modCursor,
				"Fetch", e.Error())
		}
		//ds.Data = datas
	} else if n == 1 && reflect.ValueOf(m).Elem().Kind() != reflect.Slice {
		//dataBuf := M{}
		//Printf("Record count: %d\n", func() int { count, _ := c.mgoCursor.Count(); return count }())

		bOk := c.mgoIter.Next(m)

		if !bOk {
			errtxt := ""
			if c.mgoIter.Err() == nil {
				errtxt = "Not found"
			} else {
				errtxt = c.mgoIter.Err().Error()
			}
			return errorlib.Error(packageName, modCursor, "Fetch", errtxt)
		}
		/*
			e := Serde(dataBuf, m, "json")
			if e != nil {
				return errorlib.Error(packageName, modCursor, "Fetch", "Fetch serde fail. "+e.Error())
			}
		*/
	} else if n > 1 || reflect.ValueOf(m).Elem().Kind() == reflect.Slice {
		fetched := 0
		fetching := true

		v := reflect.TypeOf(m).Elem().Elem()
		ivs := reflect.MakeSlice(reflect.SliceOf(v), 0, 0)
		for fetching {
			iv := reflect.New(v).Interface()

			tiv := M{}
			bOk := c.mgoIter.Next(&tiv)
			if bOk {
				Serde(tiv, iv, "json")
				ivs = reflect.Append(ivs, reflect.ValueOf(iv).Elem())
				fetched++
				if fetched == n {
					fetching = false
				}
			} else {
				fetching = false
			}
		}
		reflect.ValueOf(m).Elem().Set(ivs)
		// for fetching {
		// 	dataHolder, e := GetEmptySliceElement(m)
		// 	if e != nil {
		// 		return errorlib.Error(packageName, modCursor, "Fetch", e.Error())
		// 	}
		// 	var bOk bool
		// 	if IsPointer(dataHolder) {
		// 		bOk = c.mgoIter.Next(dataHolder)
		// 	} else {
		// 		bOk = c.mgoIter.Next(&dataHolder)
		// 	}
		// 	if bOk {
		// 		AppendSlice(m, dataHolder)
		// 		fetched++
		// 		if fetched == n {
		// 			fetching = false
		// 		}
		// 	} else {
		// 		fetching = false
		// 	}
		// }
	}

	return nil
}

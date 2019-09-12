package mongo

import (
	"github.com/eaciit/dbox"
	"gopkg.in/mgo.v2"

	"regexp"
	"time"

	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
)

const (
	packageName   = "eaciit.dbox.dbc.mongo"
	modConnection = "Connection"
)

type Connection struct {
	dbox.Connection

	session *mgo.Session
}

func init() {
	dbox.RegisterConnector("mongo", NewConnection)
}

func NewConnection(ci *dbox.ConnectionInfo) (dbox.IConnection, error) {
	if ci.Settings == nil {
		ci.Settings = toolkit.M{}
	}
	c := new(Connection)
	c.SetInfo(ci)
	c.SetFb(dbox.NewFilterBuilder(new(FilterBuilder)))
	return c, nil
}

func (c *Connection) Connect() error {
	info := new(mgo.DialInfo)
	ci := c.Info()
	if ci == nil {
		return errorlib.Error(packageName, modConnection, "Connect", "ConnectionInfo is not initialized")
	}
	if ci.UserName != "" {
		info.Username = ci.UserName
		info.Password = ci.Password
		if ci.AuthMechanism == "DEFAULT" {
			info.Source = "admin"
		} else {
			info.Mechanism = ci.AuthMechanism
			info.Source = ""
		}
	}
	info.Addrs = ci.Host
	info.Database = ci.Database

	if ci.Settings == nil {
		ci.Settings = toolkit.M{}
	}

	poollimit := ci.Settings.GetInt("poollimit")
	if poollimit > 0 {
		info.PoolLimit = poollimit
	}

	timeout := ci.Settings.GetInt("timeout")
	if timeout > 0 {
		info.Timeout = time.Duration(timeout) * time.Second
	}

	//sess, e := mgo.Dial(info.Addrs[0])
	sess, e := mgo.DialWithInfo(info)
	if e != nil {
		hlist := ""
		for _, tm0 := range ci.Host {
			hlist += tm0 + ","
		}
		return errorlib.Error(packageName, modConnection,
			"Connect", e.Error()+" "+ci.UserName+"@"+hlist+"/"+ci.Database)
	}
	sess.SetMode(mgo.Monotonic, true)
	// sess.SetPoolLimit(100)
	if len(ci.Host) > 1 {
		sess.SetMode(mgo.Primary, true)
		sess.SetSafe(&mgo.Safe{WMode: "majority"})
	}
	c.session = sess
	return nil
}

func (c *Connection) NewQuery() dbox.IQuery {
	q := new(Query)
	q.SetConnection(c)
	q.SetThis(q)
	return q
}

func (c *Connection) ObjectNames(obj dbox.ObjTypeEnum) []string {
	mgoDb := c.session.DB(c.Info().Database)
	if obj == "" {
		obj = dbox.ObjTypeAll
	}

	astr := []string{}

	if obj == dbox.ObjTypeAll || obj == dbox.ObjTypeTable {
		cols, err := mgoDb.CollectionNames()
		if err != nil {
			return []string{}
		}

		for _, col := range cols {
			if cond, _ := regexp.MatchString("^(.*)((\\.(indexes)|\\.(js)))$", col); !cond {
				astr = append(astr, col)
			}
		}

	}

	if obj == dbox.ObjTypeAll || obj == dbox.ObjTypeProcedure {
		cols := mgoDb.C("system.js")
		res := []toolkit.M{}
		err := cols.Find(nil).All(&res)
		if err != nil {
			toolkit.Printf("%v\n", err.Error())
			return []string{}
		}

		// toolkit.Printf("%v\n", res)
		for _, col := range res {
			astr = append(astr, col["_id"].(string))
		}

	}

	return astr
}

func (c *Connection) Close() {
	if c.session != nil {
		c.session.Close()
	}
}

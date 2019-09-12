package helper

import (
	"strings"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

var GlobalConfig map[string]string = map[string]string{
	"host":      "localhost:27017",
	"database":  "apiapi",
	"username":  "",
	"password":  "",
	"mechanism": "DEFAULT",
}

func ConnectToDB() (db.IConnection, error) {
	hostParam := strings.Split(strings.TrimSpace(GlobalConfig["host"]), "~")
	hostList := []string{}
	for _, host := range hostParam {
		hostList = append(hostList, host)
	}

	dbPassword := DecryptAes128(GlobalConfig["password"], AES128KEY)
	ci := &db.ConnectionInfo{hostList, GlobalConfig["database"], GlobalConfig["username"], dbPassword, GlobalConfig["mechanism"], nil}
	c, e := db.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func SaveRecord(m orm.IModel) error {
	conn, err := ConnectToDB()
	defer conn.Close()
	if err != nil {
		return err
	}

	ctx := orm.New(conn)
	err = ctx.Save(m)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRecord(m orm.IModel) error {
	conn, err := ConnectToDB()
	defer conn.Close()
	if err != nil {
		return nil
	}

	ctx := orm.New(conn)
	err = ctx.Delete(m)
	if err != nil {
		return err
	}

	return nil
}

//GetDataFromDB Get Data from query (pipe)
func GetDataFromDB(pipe []tk.M, tablename string) ([]tk.M, error) {
	result := []tk.M{}

	conn, err := ConnectToDB()
	defer conn.Close()
	if err != nil {
		return result, err
	}
	query := conn.NewQuery()

	if len(pipe) != 0 {
		query.Command("pipe", pipe)
	}

	csrU, err := query.From(tablename).Cursor(nil)
	defer csrU.Close()
	if err != nil {
		return nil, err
	}

	err = csrU.Fetch(&result, 0, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

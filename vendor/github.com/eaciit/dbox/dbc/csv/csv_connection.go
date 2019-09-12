package csv

import (
	"encoding/csv"
	// "fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/eaciit/cast"
	"github.com/eaciit/dbox"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	// "time"
	// "reflect"
)

type TypeOpenFile_Enum int

type headerstruct struct {
	// index 	 int
	name     string
	dataType string
	format   string
}

const (
	TypeOpenFile_Append TypeOpenFile_Enum = iota
	TypeOpenFile_Create
)

const (
	packageName   = "eaciit.dbox.dbc.csv"
	modConnection = "Connection"
)

type Connection struct {
	dbox.Connection

	TypeOpenFile TypeOpenFile_Enum
	ExecOpr      bool
	setNewHeader bool
	isUseHeader  bool
	isMapHeader  bool

	file     *os.File
	tempfile *os.File
	reader   *csv.Reader
	writer   *csv.Writer

	headerColumn []headerstruct
}

func init() {
	dbox.RegisterConnector("csv", NewConnection)
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

/*
	file				string	// File Path => host

	useheader			bool	// field column, harus ada
	comma           	rune 	// field delimiter (set to ',' by NewReader)
	comment          	rune 	// comment character for start of line
	fieldsperrecord  	int  	// number of expected fields per record
	lazyquotes       	bool 	// allow lazy quotes
	trailingcomma    	bool 	// ignored; here for backwards compatibility
	trimleadingspace 	bool
*/
func (c *Connection) Connect() error {
	ci := c.Info()
	if ci == nil {
		return errorlib.Error(packageName, modConnection, "Connect", "ConnectionInfo is not initialized")
	}

	c.isUseHeader = ci.Settings.Get("useheader", false).(bool)
	isNewFile := ci.Settings.Get("newfile", false).(bool)
	c.setNewHeader = false

	if filePath := ci.Host[0]; filePath != "" {
		var err error
		c.file, err = os.Open(filePath)
		if err != nil {
			if isNewFile {
				c.file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
				if err != nil {
					return errorlib.Error(packageName, modConnection, "Connect", "Cannot Create New File")
				}

				c.writer = csv.NewWriter(c.file)
				if delimiter := ci.Settings.Get("delimiter", "").(string); delimiter != "" {
					c.writer.Comma = rune(delimiter[0])
				}

				if c.isUseHeader {
					c.setNewHeader = true
				}
			}

			if !isNewFile || err != nil {
				return errorlib.Error(packageName, modConnection, "Connect", "Cannot Open File "+filePath)
			}
		}
		c.reader = csv.NewReader(c.file)
	} else {
		return errorlib.Error(packageName, modConnection, "Connect", "File is not initialized")
	}

	c.SetReaderParam()
	if !c.setNewHeader {
		c.SetHeaderData()
	}
	c.isMapHeader = false
	if ci.Settings.Has("mapheader") {
		c.isMapHeader = true
		tMapHeader := ci.Settings["mapheader"].([]toolkit.M)
		c.SetHeaderSliceToolkitM(tMapHeader)

		if c.setNewHeader {
			var dataTemp []string

			for _, valHeader := range c.headerColumn {
				dataTemp = append(dataTemp, valHeader.name)
			}

			if len(dataTemp) > 0 {
				c.TypeOpenFile = TypeOpenFile_Append
				c.setNewHeader = false

				c.StartSessionWrite()

				c.writer.Write(dataTemp)
				c.writer.Flush()

				c.EndSessionWrite()
			}
		}
	}

	return nil
}

func (c *Connection) SetReaderParam() {
	ci := c.Info()

	if delimiter := ci.Settings.Get("delimiter", "").(string); delimiter != "" {
		c.reader.Comma = rune(delimiter[0])
	}

	if comment := ci.Settings.Get("comment", "").(string); comment != "" {
		c.reader.Comment = rune(comment[0])
	}

	if fieldsPerRecord := ci.Settings.GetInt("fieldsperrecord"); fieldsPerRecord > 0 {
		c.reader.FieldsPerRecord = fieldsPerRecord
	}

	if lazyQuotes := ci.Settings.Get("lazyquotes", false).(bool); lazyQuotes {
		c.reader.LazyQuotes = lazyQuotes
	}

	if trailingComma := ci.Settings.Get("trailingcomma", false).(bool); trailingComma {
		c.reader.TrailingComma = trailingComma
	}

	if trimLeadingSpace := ci.Settings.Get("trimleadingspace", false).(bool); trimLeadingSpace {
		c.reader.TrailingComma = trimLeadingSpace
	}

}

func (c *Connection) SetHeaderToolkitM(tMapHeader toolkit.M) {
	var tempstruct []headerstruct
	i := 0
	for cols := range tMapHeader {
		ts := headerstruct{}
		ts.name = cols
		ts.dataType = "string"
		tempstruct = append(tempstruct, ts)
		i += 1
	}
	c.headerColumn = tempstruct
}

func (c *Connection) SetHeaderSliceToolkitM(tMapHeader []toolkit.M) {
	var tempstruct []headerstruct
	i := 0
	for _, val := range tMapHeader {
		ts := headerstruct{}
		// ts.index = i
		for cols, dataType := range val {
			ts.name = cols
			ts.dataType = dataType.(string)
		}
		tempstruct = append(tempstruct, ts)
		i += 1
	}
	c.headerColumn = tempstruct
}

func (c *Connection) SetHeaderData() {
	ci := c.Info()
	dateformat := ci.Settings.Get("dateformat", "").(string)

	var tempstruct []headerstruct

	tempData, e := c.reader.Read()
	for i, v := range tempData {
		ts := headerstruct{}
		// ts.index = i
		ts.name = string(i)
		ts.dataType = "string"
		if c.isUseHeader {
			ts.name = v
		}
		tempstruct = append(tempstruct, ts)
	}
	if c.isUseHeader && e != io.EOF {
		tempData, e = c.reader.Read()
	}

	isCheckType := true
	ix := 0
	for isCheckType && e != io.EOF {
		ix += 1
		isCheckType = false

		for i, v := range tempData {
			if v != "" {
				tempstruct[i].dataType = "string"
				matchNumber := false
				matchFloat := false

				matchF1 := false
				matchF2 := false

				//dd.mm.yyyy dd/mm/yyyy dd-mm-yyyy
				//yyyy.mm.dd yyyy/mm/dd yyyy-mm-dd
				// formatDate := "((^(0[0-9]|[0-9]|(1|2)[0-9]|3[0-1])(\\.|\\/|-)(0[0-9]|[0-9]|1[0-2])(\\.|\\/|-)[\\d]{4}$)|(^[\\d]{4}(\\.|\\/|-)(0[0-9]|[0-9]|1[0-2])(\\.|\\/|-)(0[0-9]|[0-9]|(1|2)[0-9]|3[0-1])$))"
				F1 := "(^(0[0-9]|[0-9]|(1|2)[0-9]|3[0-1])(\\.|\\/|-)(0[0-9]|[0-9]|1[0-2])(\\.|\\/|-)[\\d]{4}$)"
				F2 := "(^[\\d]{4}(\\.|\\/|-)(0[0-9]|[0-9]|1[0-2])(\\.|\\/|-)(0[0-9]|[0-9]|(1|2)[0-9]|3[0-1])$)"
				matchF1, _ = regexp.MatchString(F1, v)
				matchF2, _ = regexp.MatchString(F2, v)
				if !matchF1 && !matchF2 && dateformat != "" {
					d := cast.String2Date(v, dateformat)
					if !d.IsZero() {
						tempstruct[i].dataType = "date"
						tempstruct[i].format = dateformat
					}
				}

				x := strings.Index(v, ".")

				if x > 0 {
					matchFloat = true
					v = strings.Replace(v, ".", "", 1)
				}

				matchNumber, _ = regexp.MatchString("^\\d+$", v)

				if matchNumber && string(v[0]) != "0" {
					tempstruct[i].dataType = "int"
					if matchFloat {
						tempstruct[i].dataType = "float"
					}
				}

				if matchF1 || matchF2 {
					tempstruct[i].dataType = "date"
					switch {
					case strings.Contains(v, "."):
						tempstruct[i].format = "dd.MM.YYYY"
						if matchF2 {
							tempstruct[i].format = "YYYY.MM.dd"
						}
					case strings.Contains(v, "-"):
						tempstruct[i].format = "dd-MM-YYYY"
						if matchF2 {
							tempstruct[i].format = "YYYY-MM-dd"
						}
					case strings.Contains(v, "/"):
						tempstruct[i].format = "dd/MM/YYYY"
						if matchF2 {
							tempstruct[i].format = "YYYY/MM/dd"
						}
					}
				}
			}
		}
		for _, v := range tempstruct {
			if v.dataType == "" {
				isCheckType = true
			}
		}

		if isCheckType {
			tempData, _ = c.reader.Read()
		}

		// fmt.Println(ix, "-", isCheckType)
		// fmt.Println(tempstruct)
		if ix > 5 {
			break
		}
	}

	c.headerColumn = tempstruct

	c.file.Close()
	c.file, _ = os.Open(ci.Host[0])
	c.reader = csv.NewReader(c.file)
	c.SetReaderParam()

	if c.isUseHeader {
		tempData, _ = c.reader.Read()
	}
}

func (c *Connection) NewQuery() dbox.IQuery {
	q := new(Query)
	q.SetConnection(c)
	q.SetThis(q)
	return q
}

func (c *Connection) Close() {
	if c.file != nil {
		c.file.Close()
	}

	if c.setNewHeader {
		os.Remove(c.Info().Host[0])
	}
}

func (c *Connection) StartSessionWrite() error {
	c.Close()

	ci := c.Info()
	if ci == nil {
		return errorlib.Error(packageName, modConnection, "Connect", "ConnectionInfo is not initialized")
	}

	if filePath := ci.Host[0]; filePath != "" {
		var err error

		c.file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return errorlib.Error(packageName, modConnection, "SessionWrite", "Cannot Open File "+filePath)
		}

		c.reader = csv.NewReader(c.file)
		c.SetReaderParam()
		dataTemp := make([]string, 0)

		if c.isUseHeader {
			dataTemp, _ = c.reader.Read()
		}

		if c.TypeOpenFile == TypeOpenFile_Create {

			c.tempfile, err = os.OpenFile(filePath+".temp", os.O_RDWR|os.O_CREATE, 0666)
			c.writer = csv.NewWriter(c.tempfile)

			if c.isUseHeader {
				// dataTemp, _ := c.reader.Read()
				c.writer.Write(dataTemp)
				c.writer.Flush()
			}

		} else {
			c.writer = csv.NewWriter(c.file)
		}
	}

	if delimiter := ci.Settings.Get("delimiter", "").(string); delimiter != "" {
		c.writer.Comma = rune(delimiter[0])
	}

	return nil
}

func (c *Connection) EndSessionWrite() error {
	c.Close()
	c.writer = nil
	if c.TypeOpenFile == TypeOpenFile_Create {
		c.tempfile.Close()
		if c.ExecOpr {
			os.Remove(c.Info().Host[0])
			os.Rename(c.Info().Host[0]+".temp", c.Info().Host[0])
		}
	}

	e := c.Connect()
	if e != nil {
		return errorlib.Error(packageName, modConnection, "SessionWrite", "Reopen Read File")
	}

	return nil

}

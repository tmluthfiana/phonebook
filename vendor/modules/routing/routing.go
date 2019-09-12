package routing

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	mux "github.com/gorilla/mux"
)

type HttpMethod string

var (
	get    HttpMethod = "GET"
	post   HttpMethod = "POST"
	put    HttpMethod = "PUT"
	delete HttpMethod = "DELETE"
)

type webContext struct {
	Class  string
	Func   func(*WeContent) interface{}
	Method HttpMethod
}

type Router struct {
	ControllerPath string
	ClassList      []interface{}
	GorillaMux     *mux.Router

	UrlPath map[string]webContext
}

func NewRouting(ControllerPath string, ClassList []interface{}) *Router {
	r := new(Router)

	r.ClassList = ClassList
	r.GorillaMux = mux.NewRouter()
	r.ControllerPath = ControllerPath

	r.Dispatch()

	return r
}

func (rt *Router) ScanningClass(c string) (func(*WeContent) interface{}, error) {
	ar := strings.Split(c, ".")
	if len(ar) != 2 {
		return nil, errors.New("Invalid class name")
	}

	for _, class := range rt.ClassList {
		v := reflect.ValueOf(class)
		controllerName := reflect.Indirect(v).Type().Name()

		if controllerName == ar[0] {
			t := reflect.TypeOf(class)
			methodCount := t.NumMethod()
			for mi := 0; mi < methodCount; mi++ {
				isFnContent := false
				method := t.Method(mi)

				tm := method.Type
				if tm.NumIn() == 2 && tm.In(1).String() == "*routing.WeContent" && ar[1] == method.Name {
					if tm.NumOut() == 1 && tm.Out(0).Kind() == reflect.Interface {
						isFnContent = true
					}
				}

				if isFnContent {
					fnc := v.MethodByName(method.Name).Interface().(func(*WeContent) interface{})

					return fnc, nil
				}

				fmt.Println("method", method.Name)
			}
		}
	}

	return nil, errors.New("Not Found")
}

func (rt *Router) url() *Router {
	if rt.UrlPath == nil {
		rt.UrlPath = map[string]webContext{}
	}

	return rt
}

func (rt *Router) registerController(path string, c string, m HttpMethod) {
	if v, err := rt.ScanningClass(c); err == nil {
		fmt.Println("Path >>", path, "class", c)

		rt.GorillaMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(fmt.Sprintf("%+v", v), path)

			wc := new(WeContent)
			wc.Writer = w
			wc.Req = r
			wc.vars = mux.Vars(r)

			data := v(wc)

			wc.Return(data.([]byte))
		})
	}
}

func (rt *Router) Get(path string, c string) {
	rt.registerController(path, c, get)
}

func (rt *Router) Post(path string, c string) {
	rt.registerController(path, c, post)
}

func (rt *Router) Put(path string, c string) {
	rt.registerController(path, c, put)
}

func (rt *Router) Delete(path string, c string) {
	rt.registerController(path, c, delete)
}

func (rt *Router) Dispatch() *Router {
	// for k, v := range rt.url().UrlPath {
	// 	fmt.Println("Registering", k)

	// }

	return rt
}

func (rt *Router) Routing() *mux.Router {
	return rt.GorillaMux
}

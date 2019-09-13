package knot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var (
	apps = map[string]*App{}
)

type App struct {
	Name              string
	Enable            bool
	LayoutTemplate    string
	ViewsPath         string
	DefaultOutputType OutputType
	UseSSL            bool
	CertificatePath   string
	PrivateKeyPath    string

	controllers map[string]interface{}
	statics     map[string]string
}

func (a *App) Register(c interface{}) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("Unable to register %v, type is not pointer \n", c)
	}

	name := strings.ToLower(reflect.Indirect(v).Type().Name())
	a.Controllers()[name] = c
	return nil
}

func (a *App) Statics() map[string]string {
	if a.statics == nil {
		a.statics = map[string]string{}
	}
	return a.statics
}

func (a *App) Static(prefix, path string) {
	if path == "" {
		delete(a.Statics(), prefix)
		return
	}
	a.Statics()[prefix] = path
}

func (a *App) Controllers() map[string]interface{} {
	if a.controllers == nil {
		a.controllers = map[string]interface{}{}
	}
	return a.controllers
}

func NewApp(name string) *App {
	app := new(App)
	app.Name = name
	app.Enable = true
	app.DefaultOutputType = DefaultOutputType
	return app
}

type AppContainerConfig struct {
	Address string
}

func RegisterApp(app *App) {
	apps[app.Name] = app
}

func GetApp(appname string) *App {
	app, _ := apps[appname]
	return app
}

func getIncludeFiles(dirname string) []string {
	fis, e := ioutil.ReadDir(dirname)
	if e != nil {
		return []string{}
	}

	files := []string{}
	for _, fi := range fis {
		if fi.IsDir() {
			files = append(files, getIncludeFiles(filepath.Join(dirname, fi.Name()))...)
		} else if strings.HasPrefix(fi.Name(), "_") { //--- include is file started with _
			files = append(files, fi.Name())
		}
	}
	return files
}

func StartApp(app *App, address string) *Server {
	return StartAppWithFn(app, address, map[string]FnContent{})
}

func StartAppWithFn(app *App, address string, otherRoutes map[string]FnContent) *Server {
	DefaultOutputType = app.DefaultOutputType
	ks := new(Server)
	ks.Address = address

	//appname := app.Name
	//-- end of regex
	includes := []string{}
	if app.ViewsPath != "" {
		includes = getIncludeFiles(app.ViewsPath)
	}
	ks.Log().Info("Scan application " + app.Name + " for controller registration")
	for _, controller := range app.Controllers() {
		ks.RegisterWithConfig(controller, "", &ResponseConfig{
			//AppName:        appname,
			ViewsPath:      app.ViewsPath,
			LayoutTemplate: app.LayoutTemplate,
			IncludeFiles:   includes,
		})
	}

	for surl, spath := range app.Statics() {
		staticUrlPrefix := "/" + surl
		ks.RouteStatic(staticUrlPrefix, spath)
	}

	if app.UseSSL {
		ks.UseSSL = true
		ks.CertificatePath = app.CertificatePath
		ks.PrivateKeyPath = app.PrivateKeyPath
	}

	ks.Route("/status", statusContainer)
	ks.Route("/stop", stopContainer)

	// register both / and /page which handlers are come from `otherRoutes`
	ks.Route("/", indexContainer(otherRoutes["/"], otherRoutes["page"]))

	registerOtherRoutes(ks, otherRoutes)

	ks.Listen()

	return ks
}

func StartContainer(c *AppContainerConfig) *Server {
	return StartContainerWithFn(c, map[string]FnContent{})
}

func StartContainerWithFn(c *AppContainerConfig, otherRoutes map[string]FnContent) *Server {
	ks := new(Server)
	ks.Address = c.Address

	for k, app := range apps {
		appname := strings.ToLower(k)
		//-- need to handle appname translation in Regex way
		if strings.Contains(appname, " ") {
			appname = strings.Replace(appname, " ", "", 0)
		}
		//-- end of regex
		includes := []string{}
		if app.ViewsPath != "" {
			includes = getIncludeFiles(app.ViewsPath)
		}
		ks.Log().Info("Scan application " + appname + " for controller registration")
		for _, controller := range app.Controllers() {
			ks.RegisterWithConfig(controller, appname, &ResponseConfig{
				AppName:        k,
				ViewsPath:      app.ViewsPath,
				LayoutTemplate: app.LayoutTemplate,
				IncludeFiles:   includes,
			})
		}

		for surl, spath := range app.Statics() {
			staticUrlPrefix := appname + "/" + surl
			ks.RouteStatic(staticUrlPrefix, spath)
		}
	}

	ks.Route("/status", statusContainer)
	ks.Route("/stop", stopContainer)

	// register both / and /page which handlers are come from `otherRoutes`
	ks.Route("/", indexContainer(otherRoutes["/"], otherRoutes["page"]))

	registerOtherRoutes(ks, otherRoutes)

	ks.Listen()

	return ks
}

func stopContainer(r *WebContext) interface{} {
	defer r.Server.Stop()
	return "Knot Server (" + r.Server.Address + ") will be stopped. Bye\n"
}

func statusContainer(r *WebContext) interface{} {
	r.Config.OutputType = OutputHtml

	str := "Knot Server v1.0 (c) Eaciit"
	return str
}

func indexContainer(indexCallback FnContent, pageCallback FnContent) FnContent {
	return FnContent(func(r *WebContext) interface{} {
		regex := regexp.MustCompile("/page/[a-zA-Z0-9_]+(/.*)?$")
		rURL := r.Request.URL.String()

		// if start with /page then use /page handler
		// otherwise, it will be / handler
		if regex.MatchString(rURL) {
			args := strings.Split(strings.Replace(rURL, "/page/", "", -1), "/")
			// the rest param after /page/ stored on header with key `PAGE_ID`
			r.Request.Header.Set("PAGE_ID", args[0])
			if pageCallback != nil {
				return pageCallback(r)
			}
		} else {
			if indexCallback != nil {
				return indexCallback(r)
			}
		}

		// If the pageCallback or indexCallback not provided, then it should return 404
		http.Error(r.Writer, "404 Page not found", 404)
		return nil
	})
}

func registerOtherRoutes(ks *Server, otherRoutes map[string]FnContent) {
	for route, handler := range otherRoutes {
		if strings.ToLower(route) == "prerequest" {
			ks.PreRequest(handler)
			continue
		}

		if strings.ToLower(route) == "postrequest" {
			ks.PostRequest(handler)
			continue
		}

		if !strings.HasPrefix(route, "/") {
			route = fmt.Sprintf("/%s", route)
		}

		// ignore handler from /page and /
		if strings.ToLower(route) == "/page" || route == "/" {
			continue
		}

		ks.Route(route, handler)
	}
}

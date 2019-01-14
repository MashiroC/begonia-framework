package begonia

import (
	"fmt"
	"log"
	"net/http"
)

type RouteAction interface {
	getHandle(method string, uri string) (*Handle, error)
	addHandle(*Handle)
	initialization([]string)
}

const (
	TREE = "tree"
	HASH = "hash"
)

type Application struct {
	Route RouteAction
}

func Init(args ...string) *Application {
	app := &Application{&HashRoute{}}
	app.Route.initialization(args)

	return app
}

func (app *Application) Start(port int) {
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	length := len(r.RequestURI)
	println(r.RequestURI)
	if r.RequestURI[length-1] == '/' {
		r.RequestURI = r.RequestURI[0 : length-1]
	}
	h, err := app.Route.getHandle(r.Method, r.RequestURI)
	if err == nil {
		h.
	}
}

func (impl *Application) SetRouteAction(r RouteAction) {
	impl.Route = r
}

func (impl *Application) SetHashRoute() {
	impl.Route = &route.HashRoute{}
}

//func (impl *Application) SetTreeRoute(){
//	impl.Route=&TreeRoute{}
//}

func (impl *Application) setRouteAction(r RouteAction) {
	if (impl.Route.hasChild()) {
		panic("has child")
	}

	impl.Route = r
}

func (app *Application) AddHandle(h *Handle) {
	length := len(h.Uri)
	if h.Uri[length-1] == '/' {
		h.Uri = h.Uri[0 : length-1]
	}
	app.Route.addHandle(h)
}

func (app *Application) AddController(control interface{}) {
	//TODO:添加控制器
}

func (app *Application) AddBeen(been interface{}) {
	//TODO:添加been
}

func (app *Application) Get(uri string, f func(*Context)) {
	h := &Handle{Uri: uri, Method: "GET", Fun: f}
	app.AddHandle(h)
}

func (app *Application) Post(uri string, f func(*Context)) {
	h := &Handle{Uri: uri, Method: "POST", Fun: f}
	app.AddHandle(h)
}

package begonia

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type RouteAction interface {
	getHandle(method string, uri string) (*Handle, error)
	addHandle(*Handle)
	initialization([]string)
	isEmpty() bool
}

const (
	TREE = "tree"
	HASH = "hash"
)

type application struct {
	Route RouteAction
}

func Init(args ...string) *application {
	app := &application{&HashRoute{}}
	app.Route.initialization(args)

	return app
}

func (app *application) Start(port int) {
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (app *application) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//TODO:处理icon暂时先直接return去掉iconv
	if r.RequestURI=="/favicon.ico"{
		return
	}

	length := len(r.RequestURI)
	if r.RequestURI[length-1] == '/' {
		r.RequestURI = r.RequestURI[0 : length-1]
	}

	uri:=r.RequestURI
	index:=strings.IndexRune(r.RequestURI,'?')
	if index!=-1 {
		uri=uri[:index]
	}

	h, err := app.Route.getHandle(r.Method, uri)

	if err == nil {
		ctx := &Context{R: r, W: w, Header: r.Header}
		param := make(map[string]string)

		//解析参数 组装到map上
		if h.Method == "GET" {
			vars := r.URL.Query()
			for key, value := range vars {
				param[key] = value[0]
			}
		} else {
			contentType := strings.Split(r.Header["Content-nodeType"][0], "/")[1]
			if contentType == "x-www-form-urlencoded" {
				r.ParseForm()
				vars := r.PostForm
				for key, value := range vars {
					param[key] = value[0]
				}
			} else if strings.Contains(contentType, "form-data") {
				r.ParseMultipartForm(32 << 20)
				vars := r.MultipartForm.Value
				for key, value := range vars {
					param[key] = value[0]
				}
			}
		}
		ctx.Param = param
		ctx.this=h
		h.Fun(ctx)
	}else{
		fmt.Println("error")
	}
}

func (app *application) SetRouteAction(r RouteAction) {
	if !app.Route.isEmpty() {
		panic("route is not empty")
	}
	app.Route = r
}

func (app *application) SetRouteHash() {
	app.Route = &HashRoute{}
}

//func (app *application) SetTreeRoute(){
//	app.Route=&TreeRoute{}
//}

func (app *application) AddHandle(h *Handle) {
	length := len(h.Uri)
	if h.Uri[length-1] == '/' {
		h.Uri = h.Uri[0 : length-1]
	}
	app.Route.addHandle(h)
}


//func (app *application) AddController(control interface{}) {
//	//TODO:添加控制器
//}

//func (app *application) AddBeen(been interface{}) {
//	//TODO:添加been
//}

func (app *application) Get(uri string, f func(*Context)) {
	h := &Handle{Uri: uri, Method: "GET", Fun: f}
	app.AddHandle(h)
}

func (app *application) Post(uri string, f func(*Context)) {
	h := &Handle{Uri: uri, Method: "POST", Fun: f}
	app.AddHandle(h)
}

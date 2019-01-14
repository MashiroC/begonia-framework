package begonia

import (
	"fmt"
	"net/http"
	"strings"
)

type Handle struct {
	Uri    string
	Method string
	Fun    func(ctx *Context)
	Header []string
	next *Handle
}

func (h *Handle) getMethod() string {
	return h.Method
}

func (h *Handle) getUri() string {
	return h.Uri

}



func (h *Handle) exec(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{}
	ctx.R = r
	ctx.W = w
	ctx.Header = r.Header
	param := make(map[string]string)

	if h.Method == "GET" {
		vars := r.URL.Query()
		for key, value := range vars {
			param[key] = value[0]
		}
	} else {
		contentType := strings.Split(r.Header["Content-Type"][0], "/")[1]
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

	h.Fun(ctx)
}

func (ctx *Context) ResponseString(r string) {
	_, err := ctx.W.Write([]byte(r))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (ctx *Context) ResponseBytes(b []byte) {
	ctx.W.Write(b)
}

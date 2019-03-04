package begonia

import (
	"fmt"
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

func (ctx *Context) String(r string) {
	_, err := ctx.W.Write([]byte(r))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (ctx *Context) Bytes(b []byte) {
	ctx.W.Write(b)
}

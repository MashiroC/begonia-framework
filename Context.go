package begonia

import "net/http"

type Context struct {
	Param  map[string]string
	Header http.Header
	R      *http.Request
	W      http.ResponseWriter
	this   *Handle
}

func (c *Context) Next() {
	if c.this.next != nil {
		c.this = c.this.next
		c.this.next.Fun(c)
	} else {
		panic("handle not has next")
	}
}

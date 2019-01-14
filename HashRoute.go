package begonia

import (
	"net/http"
)

type HandleMap map[string]*Handle

type Action interface {
	getMethod() string
	getUri() string
	exec(http.ResponseWriter, *http.Request)
}

type HashRoute struct {
	handles        map[string]HandleMap
	handleExistMap map[string]bool
}

func (hash *HashRoute) getHandle(method string, uri string) (h *Handle, err error) {
	handles := hash.handles[method]
	h, ok := handles[uri]
	if ok {
		return
	}
	if _, ok := hash.handleExistMap[uri]; ok {
		err = MethodNotAllowError{}
	} else {
		err = NotFoundError{}
	}
	return
}

func (hash *HashRoute) initialization(args []string) {
	m := make(map[string]HandleMap)
	exist := make(map[string]bool)
	hash.handleExistMap = exist
	hash.handles = m
}

func (hash *HashRoute) addHandle(h *Handle) {
	handleMap := hash.handles[h.Method]
	if handleMap == nil {
		handleMap = make(map[string]*Handle)
	}

	handle, ok := handleMap[h.Uri]

	if ok {
		for handle.next != nil {
			handle = handle.next
		}
		handle.next = h
	} else {
		handleMap[h.Uri] = h
	}
}

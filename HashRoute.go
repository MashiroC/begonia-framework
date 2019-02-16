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

//判断路由是否为空，为空可以更换实现方式。
//比如哈希表换成trie树
func (hash *HashRoute) isEmpty() bool {
	return len(hash.handleExistMap) == 0
}

func (hash *HashRoute) getHandle(method, uri string) (h *Handle, err error) {
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
	handleMap, ok := hash.handles[h.Method]
	if !ok {
		handleMap = make(map[string]*Handle)
		hash.handles[h.Method] = handleMap
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

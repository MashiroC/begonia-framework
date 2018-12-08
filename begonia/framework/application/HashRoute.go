package application

import (
	"net/http"
)

type HashRoute struct {
	routeMap *map[string]Handle
}

func (hash *HashRoute) hasChild() bool {
	return true
}

func (hash *HashRoute) execHandle(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	routeMap := *hash.routeMap
	handle, ok := routeMap[uri]
	if ok {
		if handle.Method == r.Method {
			handle.execFun(w, r)
		} else {
			handle405(w, r)
		}
	} else {
		handle404(w, r)
	}
	//TODO:return
}

func (hash *HashRoute) initialization(args []string) {
	m := make(map[string]Handle)
	hash.routeMap = &m
}

func (hash *HashRoute) addHandle(h Handle) {
	routeMap := *hash.routeMap
	_, ok := routeMap[h.Uri]

	if ok {
		panic("this route is used")
	}

	routeMap[h.Uri] = h

}

func handle403(w http.ResponseWriter, r *http.Request) {

}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("404 not found"))
}

func handle405(w http.ResponseWriter, r *http.Request) {

}

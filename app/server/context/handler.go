package context

import "net/http"

type Handler interface {
	ServeHTTP(c Context, w http.ResponseWriter, r *http.Request)
}

type HandlerFunc func(c Context, w http.ResponseWriter, r *http.Request)

func (hf HandlerFunc) ServeHTTP(c Context, w http.ResponseWriter, r *http.Request) {
	hf(c, w, r)
}

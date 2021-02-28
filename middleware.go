package mux

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func (s *Server) Use(middleware Middleware) {
	s.middlewares = append([]Middleware{middleware}, s.middlewares...)
}

func (s *Server) exec(h http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	//build
	for _, mid := range s.middlewares {
		h = mid(h)
	}
	h(w, r)
}

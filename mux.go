package mux

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	MySelf http.Server
}

type Router struct {
	m map[string]func(http.ResponseWriter, *http.Request)
}

func NewRouter() *Router {
	return &Router{m: make(map[string]func(http.ResponseWriter, *http.Request))}
}
func NewServer(addr string) (*Server, *Router) {
	r := NewRouter()
	s := &Server{}
	s.MySelf.Addr = addr
	s.MySelf.Handler = r
	return s, r
}
func (mainServer *Server) Stop() error {
	if mainServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		e := mainServer.MySelf.Shutdown(ctx)
		return e
	}
	return nil
}
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := router.m[r.URL.String()]; ok {
		h(w, r)
	} else if k, ok := hasPreffixInMap(router.m, r.URL.String()); ok {
		router.m[k](w, r)
	} else {
		fmt.Fprint(w, `<!DOCTYPE html><html><head><title>404</title><meta charset="utf-8"><meta name="viewport" content="width=device-width"></head><body>404 not found</body></html>`)
	}
}
func (r *Router) HandleFunc(url string, f func(http.ResponseWriter, *http.Request)) {
	r.m[url] = f
}
func hasPreffixInMap(m map[string]func(http.ResponseWriter, *http.Request), p string) (string, bool) {
	for k, _ := range m {
		if len(p) >= len(k) && k == p[:len(k)] {
			return k, true
		}
	}
	return "", false
}

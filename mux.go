package mux

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	MySelf *http.Server
	r      *router
}

type router struct {
	m map[string]func(http.ResponseWriter, *http.Request)
}

func newRouter() *router {
	return &router{m: make(map[string]func(http.ResponseWriter, *http.Request))}
}
func NewServer(addr string) *Server {
	s := &Server{}
	r := newRouter()
	s.MySelf = &http.Server{Addr: addr, Handler: r}
	s.r = r
	return s
}
func (mainServer *Server) ListenAndServe() error {
	return mainServer.MySelf.ListenAndServe()
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
func (mainServer *Server) HandleFunc(url string, f func(http.ResponseWriter, *http.Request)) {
	mainServer.r.m[url] = f
}
func (s *Server) Handle(pattern string, h http.Handler) {
	s.r.m[pattern] = h.ServeHTTP
}
func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := rt.m[r.URL.String()]; ok {
		h(w, r)
	} else if k, ok := hasPreffixInMap(rt.m, r.URL.String()); ok {
		rt.m[k](w, r)
	} else {
		fmt.Fprint(w, `<!DOCTYPE html><html><head><title>404</title><meta charset="utf-8"><meta name="viewport" content="width=device-width"></head><body>404 not found</body></html>`)
	}
}
func hasPreffixInMap(m map[string]func(http.ResponseWriter, *http.Request), p string) (string, bool) {
	for k, _ := range m {
		if len(p) >= len(k) && k == p[:len(k)] {
			return k, true
		}
	}
	return "", false
}

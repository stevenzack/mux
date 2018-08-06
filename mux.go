package mux

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	MySelf      *http.Server
	prehandlers []func(http.ResponseWriter, *http.Request)
	r, mr       map[string]func(http.ResponseWriter, *http.Request)
}

func NewServer(addr string) *Server {
	s := &Server{}
	s.MySelf = &http.Server{Addr: addr, Handler: s}
	s.r = make(map[string]func(http.ResponseWriter, *http.Request))
	s.mr = make(map[string]func(http.ResponseWriter, *http.Request))
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
	mainServer.r[url] = f
}
func (s *Server) HandleMultiReqs(url string, f func(http.ResponseWriter, *http.Request)) {
	s.mr[url] = f
}
func (s *Server) Handle(pattern string, h http.Handler) {
	s.r[pattern] = h.ServeHTTP
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, v := range s.prehandlers {
		v(w, r)
	}
	if h, ok := s.r[r.URL.String()]; ok {
		h(w, r)
	} else if k, ok := hasPreffixInMap(s.mr, r.URL.String()); ok {
		s.mr[k](w, r)
	} else {
		fmt.Fprint(w, `<!DOCTYPE html><html><head><title>404</title><meta charset="utf-8"><meta name="viewpos" content="width=device-width"></head><body>404 not found</body></html>`)
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
func (s *Server) AddPrehandler(f func(http.ResponseWriter, *http.Request)) {
	s.prehandlers = append(s.prehandlers, f)
}

package mux

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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

func (s *Server) ListenAndServe() error {
	return s.MySelf.ListenAndServe()
}

func (s *Server) Stop() error {
	if s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		e := s.MySelf.Shutdown(ctx)
		return e
	}
	return nil
}

func (s *Server) HandleFunc(url string, f func(http.ResponseWriter, *http.Request)) {
	s.r[url] = f
}

func (s *Server) HandleHtml(url string, text string) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(text))
	}
}
func (s *Server) HandleJs(url string, text string) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(text))
	}
}
func (s *Server) HandleCss(url string, text string) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(text))
	}
}
func (s *Server) HandleSvg(url string, text string) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write([]byte(text))
	}
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
	url := strings.Split(r.URL.String(), "?")[0]
	if h, ok := s.r[url]; ok {
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

// AddRoutes adds all s2's routes to server
func (s *Server) AddRoutes(s2 *Server) {
	for k, v := range s2.r {
		_, ok := s.r[k]
		if !ok {
			s.r[k] = v
		}
	}

	for k, v := range s2.mr {
		_, ok := s.mr[k]
		if !ok {
			s.mr[k] = v
		}
	}
}

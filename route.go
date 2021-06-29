package mux

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			log.Println("⚠️", r.RequestURI, "⚠️\t", e, string(debug.Stack()))
			http.Error(w, "Server internal error", http.StatusInternalServerError)
		}
	}()

	//prehandler
	for _, v := range s.prehandlers {
		interrupt := v(w, r)
		if interrupt {
			return
		}
	}
	url := strings.Split(r.URL.String(), "?")[0]

	//route
	switch r.Method {
	case http.MethodGet:
		if h, ok := s.get[url]; ok {
			s.exec(h, w, r)
			return
		}
	case http.MethodPost:
		if h, ok := s.post[url]; ok {
			s.exec(h, w, r)
			return
		}
	case http.MethodPut:
		if h, ok := s.put[url]; ok {
			s.exec(h, w, r)
			return
		}
	case http.MethodDelete:
		if h, ok := s.delete[url]; ok {
			s.exec(h, w, r)
			return
		}
	case http.MethodPatch:
		if h, ok := s.patch[url]; ok {
			s.exec(h, w, r)
			return
		}
	}

	if h, ok := s.r[url]; ok {
		s.exec(h, w, r)
	} else if k, ok := hasPreffixInMap(s.mr, r.URL.String()); ok {
		s.exec(s.mr[k], w, r)
	} else {
		s.NotFound(w, r)
	}
}

func (s *Server) HandleNotFound(fn http.HandlerFunc) {
	s.notFound = fn
}

func (s *Server) NotFound(w http.ResponseWriter, r *http.Request) {
	if s.notFound != nil {
		s.notFound(w, r)
		return
	}
	s.Error(w, `not found`, http.StatusNotFound)
}

func (s *Server) Error(w http.ResponseWriter, e string, code int) {
	fmt.Fprint(w, `<!DOCTYPE html><html><head><title>`+strconv.Itoa(code)+` `+e+`</title><meta charset="utf-8"><meta name="viewpos" content="width=device-width"></head><body>`+e+`</body></html>`)
}

func hasPreffixInMap(m map[string]http.HandlerFunc, p string) (string, bool) {
	for k, _ := range m {
		if len(p) >= len(k) && k == p[:len(k)] {
			return k, true
		}
	}
	return "", false
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

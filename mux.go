package mux

import (
	"context"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"time"
)

type Server struct {
	HTTPServer                    *http.Server
	middlewares                   []Middleware
	prehandlers                   []func(http.ResponseWriter, *http.Request) bool
	r, mr                         map[string]http.HandlerFunc
	get, post, put, delete, patch map[string]http.HandlerFunc
}

func NewServer(addr string) *Server {
	s := &Server{}
	s.HTTPServer = &http.Server{Addr: addr, Handler: s}
	s.r = make(map[string]http.HandlerFunc)
	s.mr = make(map[string]http.HandlerFunc)
	s.get = make(map[string]http.HandlerFunc)
	s.post = make(map[string]http.HandlerFunc)
	s.put = make(map[string]http.HandlerFunc)
	s.delete = make(map[string]http.HandlerFunc)
	s.patch = make(map[string]http.HandlerFunc)
	return s
}

func (s *Server) ListenAndServe() error {
	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Stop() error {
	if s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		e := s.HTTPServer.Shutdown(ctx)
		return e
	}
	return nil
}

// AddPrehandler adds prehandler which returns interrupt
func (s *Server) AddPrehandler(f func(http.ResponseWriter, *http.Request) bool) {
	s.prehandlers = append(s.prehandlers, f)
}

func (s *Server) GET(url string, f http.HandlerFunc) {
	s.get[url] = f
}

func (s *Server) POST(url string, f http.HandlerFunc) {
	s.post[url] = f
}

func (s *Server) PUT(url string, f http.HandlerFunc) {
	s.put[url] = f
}

func (s *Server) DELETE(url string, f http.HandlerFunc) {
	s.delete[url] = f
}

func (s *Server) PATCH(url string, f http.HandlerFunc) {
	s.patch[url] = f
}

func (s *Server) HandleFunc(url string, f http.HandlerFunc) {
	s.r[url] = f
}

func (s *Server) ServeBytes(url string, bytes []byte) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(url)))
		w.Header().Set("Cache-Control", "public")
		w.Write(bytes)
	}
}

func (s *Server) ServeFile(uri string, path string) {
	s.r[uri] = func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

func (s *Server) HandleWoff(url string, bytes []byte) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		SetWoffHeader(w)
		w.Write(bytes)
	}
}
func (s *Server) HandleRes(url string, bytes []byte) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(url)))
		w.Header().Set("Cache-Control", "public")
		w.Write(bytes)
	}
}
func (s *Server) HandleHtml(url string, text []byte) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "public")
		w.Write(text)
	}
}

func (s *Server) HandleHtmlFunc(url string, fn func(w http.ResponseWriter, r *http.Request)) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "public")
		fn(w, r)
	}
}

func (s *Server) HandleJs(url string, text []byte) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		w.Header().Set("Cache-Control", "public")
		w.Write(text)
	}
}
func (s *Server) HandleCss(url string, text []byte) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Cache-Control", "public")
		w.Write(text)
	}
}
func (s *Server) HandleSvg(url string, text []byte) {
	s.r[url] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public")
		w.Write(text)
	}
}

func (s *Server) HandleMultiReqs(url string, f http.HandlerFunc) {
	s.mr[url] = f
}

func (s *Server) Handle(pattern string, h http.Handler) {
	s.r[pattern] = h.ServeHTTP
}

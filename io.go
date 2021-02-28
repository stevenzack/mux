package mux

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
)

type ResponseWriterBuf struct {
	Rw     http.ResponseWriter
	Writer io.WriteCloser
}

func (r *ResponseWriterBuf) Header() http.Header {
	return r.Rw.Header()
}

func (r *ResponseWriterBuf) Write(p []byte) (int, error) {
	return r.Writer.Write(p)
}

func (r *ResponseWriterBuf) WriteHeader(code int) {
	r.Rw.WriteHeader(code)
}

func (r *ResponseWriterBuf) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := r.Rw.(http.Hijacker); ok {
		return h.Hijack()
	}
	panic(fmt.Sprintf("gzipWriter.w:%s is not a http.Hijacker", reflect.TypeOf(r.Rw).String()))
}

func ReadBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

func ReadBodyString(r *http.Request) (string, error) {
	b, e := ReadBody(r)
	if e != nil {
		return "", e
	}
	return string(b), nil
}

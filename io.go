package mux

import (
	"io/ioutil"
	"net/http"
)

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

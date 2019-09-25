package mux

import (
	"errors"
	"net/http"
	"strings"
)

func GetURIParam(r *http.Request, index int) (string, error) {
	strs := strings.Split(r.RequestURI, "/")
	if len(strs) <= index+1 {
		return "", errors.New("404 not found")
	}
	return strs[index+1], nil
}

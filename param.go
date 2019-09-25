package mux

import (
	"errors"
	"net/http"
	"strings"
)

func GetURIParam(r *http.Request, index int) (string, error) {
	strs := strings.Split(r.RequestURI, "/")
	if len(strs) <= index {
		return "", errors.New("404 not found")
	}
	return strs[index], nil
}

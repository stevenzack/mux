package mux

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func SetJsHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/javascript")
}

func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/json")
}

func SetCssHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/css")
}

func SetSvgHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
}

func SetContentLength(w http.ResponseWriter, l int) {
	w.Header().Set("Content-Length", strconv.Itoa(l))

}

func SetHtmlHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

func SetWoffHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/font-woff")
}

func WriteHtml(w http.ResponseWriter, html string) {
	SetHtmlHeader(w)
	w.Write([]byte(html))
}

func WriteJSON(w http.ResponseWriter, v interface{}) {
	SetJSONHeader(w)
	b, e := json.Marshal(v)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	w.Write(b)
}

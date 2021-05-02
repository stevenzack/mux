package mux

import (
	"compress/gzip"
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func GzipAutoMiddleware(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept-Encoding")
		content := r.Header.Get("Content-Encoding")
		if strings.Contains(content, "gzip") {
			reader, e := gzip.NewReader(r.Body)
			if e != nil {
				e = fmt.Errorf("reading gzip body failed:%w", e)
				log.Println(e)
				http.Error(w, e.Error(), http.StatusBadRequest)
				return
			}
			r.Body = reader
		}

		if strings.Contains(accept, "gzip") {
			ext := filepath.Ext(r.RequestURI)
			mime := mime.TypeByExtension(ext)
			if r.RequestURI == "/" || strings.HasPrefix(mime, "text/") || ext == ".json" || ext == ".js" {
				w.Header().Set("Content-Encoding", "gzip")
				gw := gzip.NewWriter(w)
				defer gw.Flush()
				defer gw.Close()
				w = &ResponseWriterBuf{
					Rw:     w,
					Writer: gw,
				}
			}
		}
		hf(w, r)
	}
}

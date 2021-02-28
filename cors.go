package mux

import (
	"net/http"
	"strings"
)

type CorsOption struct {
	Methods     []string
	Headers     []string
	Origins     []string
	Credentials bool
}

func CorsMiddleware(opt CorsOption) Middleware {
	for i := range opt.Headers {
		opt.Headers[i] = strings.ToLower(opt.Headers[i])
	}
	if len(opt.Methods) == 0 {
		opt.Methods = append(opt.Methods, http.MethodPost, http.MethodGet, http.MethodOptions)
	}
	methods := strings.Join(opt.Methods, ",")
	headers := strings.Join(opt.Headers, ",")
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodOptions {
				hf(w, r)
				return
			}

			var origin = "*"
			if len(opt.Origins) > 0 {
				originHeader := r.Header.Get("Origin")
				for _, v := range opt.Origins {
					if v == originHeader {
						origin = originHeader
						break
					}
				}

				//not found
				if origin == "*" {
					origin = ""
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", headers)
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			if opt.Credentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}
	}
}

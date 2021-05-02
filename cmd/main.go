package main

import (
	"log"
	"net/http"

	"github.com/StevenZack/mux"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	s := mux.NewServer(":8080")
	s.Use(mux.GzipAutoMiddleware)
	s.HandleMultiReqs("/", home)
	e := s.ListenAndServe()
	if e != nil {
		log.Println(e)
		return
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../notfound.html")
}

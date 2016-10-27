package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/npetzall/badger-go/lib/badge"
	"github.com/npetzall/badger-go/lib/nexus"
)

func createBadge(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	err, data := badge.CreateBadge(q.Get("l"), q.Get("r"), q.Get("c"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", "image/svg+xml")
		w.Write(data)
	}
}

func main() {
	r := mux.NewRouter()
	r.Path("/badge").Handler(http.HandlerFunc(createBadge))
	nexus.Configure(r.PathPrefix(nexus.PathPrefix).Subrouter())

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

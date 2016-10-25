package main

import (
	"net/http"

	"github.com/npetzall/badger-go/lib/badge"
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
	http.HandleFunc("/badge", createBadge)
	http.ListenAndServe(":8080", nil)
}

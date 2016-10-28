package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/npetzall/badger-go/lib/badge"
	"github.com/npetzall/badger-go/lib/nexus"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("badger-go")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func GetLogger() *logging.Logger {
	return log
}

func createBadge(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	err, data := badge.CreateBadge(q.Get("l"), q.Get("r"), q.Get("c"))
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Add("Content-Type", "image/svg+xml")
		w.Write(data)
	}
}

func init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
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
	log.Infof("Starting server at: %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

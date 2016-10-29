package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/npetzall/badger-go/lib/badge"
	"github.com/npetzall/badger-go/lib/logger"
	"github.com/npetzall/badger-go/lib/nexus"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("badger-go")

func createBadge(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	err, data := badge.CreateBadge(q.Get("l"), q.Get("r"), q.Get("c"))
	if err != nil {
		log.Error(err)
		badge.DownstreamError(w, err)
	} else {
		w.Header().Add("Content-Type", "image/svg+xml")
		w.Write(data)
	}
}

func init() {
	logger.Init()
}

func main() {
	flag.Parse()
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

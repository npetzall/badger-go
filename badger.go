package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/npetzall/badger-go/lib/badge"
	"github.com/npetzall/badger-go/lib/logger"
	"github.com/npetzall/badger-go/system/artifactory"
	"github.com/npetzall/badger-go/system/nexus"
	"github.com/npetzall/badger-go/system/sonarqube"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("badger-go")

func createBadge(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	err, data := badge.CreateBadge(q.Get("l"), q.Get("r"), q.Get("c"))
	if err != nil {
		log.Error(err)
		badge.DownstreamError(w)
	} else {
		w.Header().Add("Content-Type", "image/svg+xml")
		w.Write(data)
	}
}

func main() {
	flag.Parse()
	logger.Init()
	r := mux.NewRouter()
	r.Path("/badge").Handler(http.HandlerFunc(createBadge))
	nexus.Configure(r.PathPrefix(nexus.PathPrefix).Subrouter())
	sonarqube.Configure(r.PathPrefix(sonarqube.PathPrefix).Subrouter())
	artifactory.Configure(r.PathPrefix(artifactory.PathPrefix).Subrouter())

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Infof("Starting server at: %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

package nexus

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/npetzall/badger-go/lib/badge"
	"github.com/npetzall/badger-go/lib/requests"
	logging "github.com/op/go-logging"
)

const defNexusURL = "http://localhost:8081/nexus"

var nexusURLFlag string

var log = logging.MustGetLogger("nexus")

var nexusURL = defNexusURL

const PathPrefix = "/nexus"

func latestVersion(w http.ResponseWriter, r *http.Request) {
	req, rerr := http.NewRequest("GET", nexusURL+"/service/local/artifact/maven/resolve?v=LATEST&"+r.URL.RawQuery, nil)
	if rerr != nil {
		log.Error(rerr)
		badge.DownstreamError(w, rerr)
		return
	}
	req.Header.Add("Accept", "application/json")
	v, e := requests.Request(req, requests.JsonFirstMatch("$.data.version+"))
	if e != nil {
		log.Error(e)
		badge.DownstreamError(w, e)
		return
	}
	berr, data := badge.CreateBadge("Latest Version", v, "green")
	if berr != nil {
		log.Error(berr)
		badge.DownstreamError(w, berr)
		return
	} else {
		w.Header().Add("Content-Type", "image/svg+xml")
		w.Write(data)
		return
	}
}

func init() {
	flag.StringVar(&nexusURLFlag, "nexusUrl", defNexusURL, "configure url for nexus")
}

func Configure(r *mux.Router) {
	if os.Getenv("NEXUS_URL") != "" {
		nexusURL = os.Getenv("NEXUS_URL")
	}
	if nexusURLFlag != defNexusURL {
		nexusURL = nexusURLFlag
	}
	log.Infof("Using baseUrl: %s", nexusURL)
	r.PathPrefix("/latestVersion").Handler(http.HandlerFunc(latestVersion))
}

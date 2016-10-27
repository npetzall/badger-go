package nexus

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/npetzall/badger-go/lib/badge"
	"github.com/npetzall/badger-go/lib/requests"
)

const defNexusURL = "http://localhost:8081/nexus"

var nexusURL string

const PathPrefix = "/nexus"

func latestVersion(w http.ResponseWriter, r *http.Request) {
	req, rerr := http.NewRequest("GET", nexusURL+"/service/local/artifact/maven/resolve?v=LATEST&"+r.URL.RawQuery, nil)
	if rerr != nil {
		badge.DownstreamError(w, rerr)
		return
	}
	req.Header.Add("Accept", "application/json")
	v, e := requests.Request(req, requests.JsonFirstMatch("$.data.version+"))
	if e != nil {
		badge.DownstreamError(w, e)
		return
	}
	berr, data := badge.CreateBadge("Latest Version", v, "green")
	if berr != nil {
		badge.DownstreamError(w, berr)
		return
	} else {
		w.Header().Add("Content-Type", "image/svg+xml")
		w.Write(data)
		return
	}
}

func Configure(r *mux.Router) {
	if os.Getenv("NEXUS_URL") != "" {
		nexusURL = os.Getenv("NEXUS_URL")
	} else {
		nexusURL = defNexusURL
	}
	r.PathPrefix("/latestVersion").Handler(http.HandlerFunc(latestVersion))
}

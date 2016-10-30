package nexus

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/npetzall/badger-go/lib/badge"
	"github.com/npetzall/badger-go/lib/metric"
	"github.com/npetzall/badger-go/lib/requests"
	logging "github.com/op/go-logging"
)

const PathPrefix = "/nexus"
const defNexusURL = "http://localhost:8081/nexus"

var nexusURLFlag = flag.String("nexusUrl", defNexusURL, "configure url for nexus")

var log = logging.MustGetLogger("nexus")

var nexusURL = defNexusURL

func latestVersion(w http.ResponseWriter, r *http.Request) {
	defer metric.Timed(time.Now(), "latestVersion", log)
	req, rerr := http.NewRequest("GET", nexusURL+"/service/local/artifact/maven/resolve?v=LATEST&"+r.URL.RawQuery, nil)
	if rerr != nil {
		log.Error(rerr)
		badge.DownstreamError(w)
		return
	}
	req.Header.Add("Accept", "application/json")
	v, e := requests.Request(req, requests.JsonFirstMatch("$.data.version+"))
	if e != nil {
		log.Error(e)
		badge.DownstreamError(w)
		return
	}
	berr, data := badge.CreateBadge("Latest Version", v, "green")
	if berr != nil {
		log.Error(berr)
		badge.DownstreamError(w)
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
	}
	if defNexusURL != *nexusURLFlag {
		nexusURL = *nexusURLFlag
	}
	log.Infof("Using baseUrl: %s", nexusURL)
	r.PathPrefix("/latestVersion").Handler(http.HandlerFunc(latestVersion))
}

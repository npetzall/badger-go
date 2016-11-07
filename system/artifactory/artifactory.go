package artifactory

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

const PathPrefix = "/artifactory"

const defArtifactoryURL = "http://localhost:8081"

var artifactoryURLFlag = flag.String("artifactoryUrl", defArtifactoryURL, "configure url for artifactory")

var log = logging.MustGetLogger("artifactory")

var artifactoryURL = defArtifactoryURL

func Configure(r *mux.Router) {
	if os.Getenv("ARTIFACTORY_URL") != "" {
		artifactoryURL = os.Getenv("ARTIFACTORY_URL")
	}
	if defArtifactoryURL != *artifactoryURLFlag {
		artifactoryURL = *artifactoryURLFlag
	}
	log.Infof("Using baseUrl: %s", artifactoryURL)
	r.PathPrefix("/latestVersion").Methods("GET").Handler(http.HandlerFunc(getLatestVersion))
}

func getLatestVersion(w http.ResponseWriter, r *http.Request) {
	defer metric.Timed(time.Now(), "getLatestVersion", log)
	req, rerr := http.NewRequest("GET", artifactoryURL+"/api/search/latestVersion?"+r.URL.RawQuery, nil)
	if rerr != nil {
		log.Error(rerr)
		badge.DownstreamError(w)
		return
	}
	v, e := requests.Request(req, requests.Text())
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

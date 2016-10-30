package sonarqube

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

const PathPrefix = "/sonarqube"
const defSonarqubeURL = "http://localhost:9000"

var sonarqubeURLFlag = flag.String("sonarqubeUrl", defSonarqubeURL, "configure url for sonarqube")

var log = logging.MustGetLogger("sonarqube")

var sonarqubeURL = defSonarqubeURL

func Configure(r *mux.Router) {
	if os.Getenv("SONARQUBE_URL") != "" {
		sonarqubeURL = os.Getenv("SONARQUBE_URL")
	}
	if defSonarqubeURL != *sonarqubeURLFlag {
		sonarqubeURL = *sonarqubeURLFlag
	}
	log.Infof("Using baseUrl: %s", sonarqubeURL)
	r.PathPrefix("/version").Methods("GET").Handler(http.HandlerFunc(getVersion))
	r.PathPrefix("/overall-coverage").Methods("GET").Handler(http.HandlerFunc(overallCoverage))
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	defer metric.Timed(time.Now(), "getVersion", log)
	v, e := version()
	if e != nil {
		log.Error(e)
		badge.DownstreamError(w)
		return
	}
	berr, data := badge.CreateBadge("Sonarqube Version", v, "green")
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

func version() (string, error) {
	defer metric.Timed(time.Now(), "version", log)
	req, rerr := http.NewRequest("GET", sonarqubeURL+"/api/system/status", nil)
	if rerr != nil {
		return "", rerr
	}
	v, e := requests.Request(req, requests.JsonFirstMatch("$.version+"))
	if e != nil {
		return "", e
	}
	return v, nil
}

func overallCoverage(w http.ResponseWriter, r *http.Request) {
	defer metric.Timed(time.Now(), "overallCoverage", log)
	req, rerr := http.NewRequest("GET", sonarqubeURL+"/api/resources/index?metrics=overall_coverage&resource="+r.URL.Query().Get("id"), nil)
	if rerr != nil {
		log.Error(rerr)
		badge.DownstreamError(w)
		return
	}
	val, err := requests.Request(req, requests.JsonFirstMatch("$[0].msr[0].frmt_val+"))
	if err != nil {
		log.Error(err)
		badge.DownstreamError(w)
		return
	}
	berr, data := badge.CreateBadge("Overall Coverage", val, "green")
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

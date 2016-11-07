package requests

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/NodePrime/jsonpath"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("requests")

func Request(r *http.Request, fn func(body io.ReadCloser) (string, error)) (string, error) {
	log.Debugf("%s: %s", r.Method, r.URL.String())
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Error(err)
		return "", err
	}
	v, e := fn(res.Body)
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	return v, e
}

func JsonFirstMatch(path string) func(body io.ReadCloser) (string, error) {
	return func(body io.ReadCloser) (string, error) {
		paths, err := jsonpath.ParsePaths(path)
		if err != nil {
			log.Error(err)
			return "", err
		}
		eval, err := jsonpath.EvalPathsInReader(body, paths)
		if err != nil {
			log.Error(err)
			return "", err
		}
		for {
			result, ok := eval.Next()
			if ok {
				return strings.Trim(string(result.Value), "\""), nil
			}
			if eval.Error != nil {
				log.Error(eval.Error)
				return "", eval.Error
			}
			return "", fmt.Errorf("0 results under path: %s", path)
		}
	}
}

func Text() func(body io.ReadCloser) (string, error) {
	return func(body io.ReadCloser) (string, error) {
		data, err := ioutil.ReadAll(body)
		if err != nil {
			log.Error(err)
			return "", err
		}
		return string(data), nil
	}
}

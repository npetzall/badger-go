package requests

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/NodePrime/jsonpath"
)

func Request(r *http.Request, fn func(body io.ReadCloser) (string, error)) (string, error) {
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
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
			return "", err
		}
		eval, err := jsonpath.EvalPathsInReader(body, paths)
		if err != nil {
			return "", err
		}
		for {
			result, ok := eval.Next()
			if ok {
				return strings.Trim(string(result.Value), "\""), nil
			}
			if eval.Error != nil {
				return "", eval.Error
			}
			return "", fmt.Errorf("0 results under path: %s", path)
		}
	}
}

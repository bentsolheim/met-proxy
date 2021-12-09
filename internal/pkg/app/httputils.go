package app

import (
	"github.com/palantir/stacktrace"
	"io/ioutil"
	"net/http"
)

func HttpGetWithResponse(url string) (*http.Response, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "error while getting %s", url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "error while reading response body")
	}
	return resp, body, nil
}

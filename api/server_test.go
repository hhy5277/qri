package api

import (
	"bytes"
	// "encoding/json"
	"github.com/qri-io/qri/repo/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerRoutes(t *testing.T) {
	// TODO: refactor cases struct:
	// cases := []struct{
	// 	method, endpoint string
	// 	body string // json? or should this be map[string]interface{}?
	// 	resBody string // json? or should this be map[string]interface{}?
	// 	resStatus int
	// }
	cases := []struct {
		method, endpoint string
		body             []byte
		resStatus        int
	}{
		// {"GET", "/", nil, 200},
		{"GET", "/status", nil, 200},
		{"OPTIONS", "/add/", nil, 200},
		{"POST", "/add/", nil, 400},
		{"PUT", "/add/", nil, 400},
		// TODO: more tests for /add/ endpoint:
		// {"POST", "/add/", {data to add dataset}, {response body}, 200}
		// {"POST", "/add/", {badly formed body}, {response body}, 400}
		// {"PUT", "/add/", {data to add dataset}, {response body}, 200}
		// {"PUT", "/add/", {badly formed body}, {response body}, 400}
	}

	client := &http.Client{}

	r, err := test.NewTestRepo()
	if err != nil {
		t.Errorf("error allocating test repo: %s", err.Error())
		return
	}

	s, err := New(r, func(opt *Config) {
		opt.Online = false
		opt.MemOnly = true
	})
	if err != nil {
		t.Error(err.Error())
		return
	}

	server := httptest.NewServer(NewServerRoutes(s))

	for i, c := range cases {
		req, err := http.NewRequest(c.method, server.URL+c.endpoint, bytes.NewReader(c.body))
		if err != nil {
			t.Errorf("case %d error creating request: %s", i, err.Error())
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("case %d error performing request: %s", i, err.Error())
			continue
		}

		if res.StatusCode != c.resStatus {
			t.Errorf("case %d: %s - %s status code mismatch. expected: %d, got: %d", i, c.method, c.endpoint, c.resStatus, res.StatusCode)
			continue
		}
	}
}

package plugin_blockuseragent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc    string
		regexps []string
		expErr  bool
	}{
		{
			desc:    "should return no error",
			regexps: []string{`\bagent1\b`},
			expErr:  false,
		},
		{
			desc:    "should return an error",
			regexps: []string{"*"},
			expErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				Regex: test.regexps,
			}

			if _, err := New(context.Background(), nil, cfg, "name"); test.expErr && err == nil {
				t.Errorf("expected error on bad regexp format")
			}
		})
	}
}

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		desc          string
		regexps       []string
		reqUserAgent  string
		expNextCall   bool
		expStatusCode int
	}{
		{
			desc:          "should return forbidden status",
			regexps:       []string{"\\bagent1\\b"},
			reqUserAgent:  "agent1",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return forbidden status",
			regexps:       []string{"\\bagent1\\b", "\\bagent2\\b"},
			reqUserAgent:  "agent2",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			regexps:       []string{"\\bagent1\\b", "\\bagent2\\b"},
			reqUserAgent:  "agentok",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return ok status",
			reqUserAgent:  "agentok",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return forbidden status",
			regexps:       []string{"\\bagent.*"},
			reqUserAgent:  "agent1",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return forbidden status",
			regexps:       []string{"^agent.*"},
			reqUserAgent:  "agent1",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				Regex: test.regexps,
			}

			nextCall := false
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				nextCall = true
			})

			handler, err := New(context.Background(), next, cfg, "blockpath")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)
			req.Header.Add("User-Agent", test.reqUserAgent)

			handler.ServeHTTP(recorder, req)

			if nextCall != test.expNextCall {
				t.Errorf("next handler should not be called")
			}

			if recorder.Result().StatusCode != test.expStatusCode {
				t.Errorf("got status code %d, want %d", recorder.Code, test.expStatusCode)
			}
		})
	}
}

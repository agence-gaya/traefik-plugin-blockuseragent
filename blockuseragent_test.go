package traefik_plugin_blockuseragent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc       string
		regexAllow []string
		regexDeny  []string
		expErr     bool
	}{
		{
			desc:       "should return no error",
			regexAllow: []string{`\bagent1\b`},
			regexDeny:  []string{`\bagent2\b`},
			expErr:     false,
		},
		{
			desc:       "should return an error",
			regexAllow: []string{"*"},
			regexDeny:  nil,
			expErr:     true,
		},
		{
			desc:       "should return an error",
			regexAllow: nil,
			regexDeny:  []string{"*"},
			expErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				RegexAllow: test.regexAllow,
				Regex:      test.regexDeny,
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
		regexAllow    []string
		regexDeny     []string
		reqUserAgent  string
		reqURI        string
		expNextCall   bool
		expStatusCode int
	}{
		{
			desc:          "should return forbidden status",
			regexAllow:    nil,
			regexDeny:     []string{"\\bagent1\\b"},
			reqUserAgent:  "agent1",
			reqURI:        "http://localhost/",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return forbidden status",
			regexAllow:    nil,
			regexDeny:     []string{"\\bagent1\\b", "\\bagent2\\b"},
			reqUserAgent:  "agent2",
			reqURI:        "http://localhost/test",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			regexAllow:    nil,
			regexDeny:     []string{"\\bagent1\\b", "\\bagent2\\b"},
			reqUserAgent:  "agentok",
			reqURI:        "http://localhost/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return ok status",
			regexAllow:    nil,
			regexDeny:     nil,
			reqUserAgent:  "agentok",
			reqURI:        "http://localhost/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return forbidden status",
			regexAllow:    nil,
			regexDeny:     []string{"\\bagent.*"},
			reqUserAgent:  "agent1",
			reqURI:        "http://localhost/test",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return forbidden status",
			regexAllow:    nil,
			regexDeny:     []string{"^agent.*"},
			reqUserAgent:  "agent1",
			reqURI:        "http://localhost/test",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return forbidden status",
			regexAllow:    nil,
			regexDeny:     []string{"^$"},
			reqUserAgent:  "",
			reqURI:        "http://localhost/test",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			regexAllow:    nil,
			regexDeny:     []string{"^$"},
			reqUserAgent:  "agentok",
			reqURI:        "http://localhost/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return ok status",
			regexAllow:    []string{"^agent1$"},
			regexDeny:     []string{".*"},
			reqUserAgent:  "agent1",
			reqURI:        "http://localhost/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return ok status",
			regexAllow:    []string{"^agent1$", "^agent2$"},
			regexDeny:     []string{".*"},
			reqUserAgent:  "agent2",
			reqURI:        "http://localhost/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return forbidden status",
			regexAllow:    []string{"^agent1$", "^agent2$"},
			regexDeny:     []string{".*"},
			reqUserAgent:  "agent3",
			reqURI:        "http://localhost/test",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			regexAllow:    []string{"^agent1$", "^agent2$"},
			regexDeny:     nil,
			reqUserAgent:  "agent3",
			reqURI:        "http://localhost/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return forbidden status",
			regexAllow:    []string{"\\bAllowed\\b"},
			regexDeny:     []string{"\\bTheAgent\\b"},
			reqUserAgent:  "This is TheAgent-1.0",
			reqURI:        "http://localhost/test",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			regexAllow:    []string{"\\bAllowed\\b"},
            regexDeny:     []string{"\\bTheAgent\\b"},
			reqUserAgent:  "This is TheAgent-Allowed-1.0",
			reqURI:        "http://localhost/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				RegexAllow: test.regexAllow,
				Regex:      test.regexDeny,
			}

			nextCall := false
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				nextCall = true
			})

			handler, err := New(context.Background(), next, cfg, "blockuseragent")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, test.reqURI, nil)
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

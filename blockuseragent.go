// Package traefik_plugin_blockuseragent a plugin to block User-Agent.
package traefik_plugin_blockuseragent

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// Config holds the plugin configuration.
type Config struct {
	Regex []string `json:"regex,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{Regex: make([]string, 0)}
}

type blockUserAgent struct {
	name    string
	next    http.Handler
	regexps []*regexp.Regexp
}

// New creates and returns a plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	regexps := make([]*regexp.Regexp, len(config.Regex))

	for i, regex := range config.Regex {
		re, err := regexp.Compile(regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", regex, err)
		}

		regexps[i] = re
	}

	return &blockUserAgent{
		name:    name,
		next:    next,
		regexps: regexps,
	}, nil
}

func (b *blockUserAgent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req != nil {
		userAgent := req.UserAgent()

		for _, re := range b.regexps {
			if re.MatchString(userAgent) {
				log.Printf("Block User-Agent: '%s'", userAgent)
				rw.WriteHeader(http.StatusForbidden)

				return
			}
		}
	}

	b.next.ServeHTTP(rw, req)
}

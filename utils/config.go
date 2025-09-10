package utils

import (
	"crypto/tls"
	"time"

	"github.com/valyala/fasthttp"
)

type Config struct {
	URL         string
	Wordlist    string
	Threads     int
	Method      string
	StatusCodes string
	Extensions  string
	Headers     []string
	Delay       int
	Retries     int
	Timeout     int
	Recursion   bool
	MaxDepth    int
	FilterSize  string
	FilterLines string
	FilterRegex string
	FollowRedirects bool
	NoTLS       bool
}

func NewFastHTTPClient(cfg *Config) *fasthttp.Client {
	return &fasthttp.Client{
		ReadTimeout:  time.Duration(cfg.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: cfg.NoTLS,
		},
	}
}
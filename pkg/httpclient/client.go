package httpclient

import (
	"net"
	"net/http"
	"time"
)

type HTTPClientOpts struct {
	ConnTimeout time.Duration
	ReadTimeout time.Duration
}

func GetClient(opts *HTTPClientOpts) *http.Client {
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: opts.ConnTimeout,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &http.Client{
		Timeout:   opts.ReadTimeout,
		Transport: netTransport,
	}
}

package common

import (
	"crypto/tls"
	"net/http"
)

func NewHttpClient() *http.Client {
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: httpTransport}
}

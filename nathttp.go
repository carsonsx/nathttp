package nathttp

import (
	"log"
	"net/url"
)

const (
	nat_http_request = "nat_http_request"
	nat_http_response = "nat_http_response"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type NatHttpRequestMessage struct {
	Method   string                 `json:"method"`
	Url      string                 `json:"url"`
	JsonData string                 `json:"jsonData"`
	FormData url.Values             `json:"formData"`
	Sync     bool                   `json:"sync"`
	Next     *NatHttpRequestMessage `json:"next"`
}

func CreateNatHttpConnection(qurl string) *NatHttpConnection {
	return &NatHttpConnection{
		qurl:  qurl,
		qreq: nat_http_request,
		qres: nat_http_response,
	}
}

func CreateNatAsyncHttpConnection(qurl string) *NatAsyncHttpConnection {
	return &NatAsyncHttpConnection{
		qurl: qurl,
		qreq: nat_http_request,
	}
}

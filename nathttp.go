package nathttp

import (
	"log"
	"net/url"
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

func CreateNatHttpConnection(qurl, qname string) *NatHttpConnection {
	return &NatHttpConnection{
		qurl:  qurl,
		qname: qname,
	}
}

func CreateNatAsyncHttpConnection(qurl, qname string) *NatAsyncHttpConnection {
	return &NatAsyncHttpConnection{
		qurl: qurl,
		qname: qname,
	}
}

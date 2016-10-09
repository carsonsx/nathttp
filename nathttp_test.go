package nathttpclient

import (
	"testing"
	"fmt"
)

func TestGet(t *testing.T) {
	conn := CreateNatHttpConnection("amqp://121.14.28.232:5672", "http_nat_queue")
	res, err := conn.Get("http://tako.im/service/openapi/app/uri?uri=tako")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
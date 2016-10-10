package nathttp

import (
	"testing"
	"fmt"
)

const URL  = "amqp://121.14.28.232:5672"
const QUEUE  = "http_nat_queue"

func TestGet(t *testing.T) {
	conn := CreateNatHttpConnection(URL, QUEUE)
	res, err := conn.Get("https://baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestAsyncGet(t *testing.T) {
	conn := CreateNatAsyncHttpConnection(URL, QUEUE)
	conn.AddGet("https://baidu.com")
	conn.AddGet("http://www.qq.com/")
	err := conn.Request()
	if err != nil {
		panic(err)
	}
}
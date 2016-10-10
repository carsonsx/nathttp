package nathttp

import (
	"testing"
	"fmt"
)

const URL  = "amqp://121.14.28.232:5672"

func TestGet(t *testing.T) {
	conn := CreateNatHttpConnection(URL)
	res, err := conn.Get("https://tako.im/tako")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestAsyncGet(t *testing.T) {
	conn := CreateNatAsyncHttpConnection(URL)
	conn.AddGet("https://tako.im/tako")
	conn.AddGet("https://tako.im/takoios")
	err := conn.Request()
	if err != nil {
		panic(err)
	}
}

func TestPost(t *testing.T) {
	conn := CreateNatHttpConnection(URL)
	res, err := conn.Post("https://tako.im/service/login")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
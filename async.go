package nathttp

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"net/http"
)

type NatAsyncHttpConnection struct {
	qurl    string
	qreq    string
	message *NatHttpRequestMessage
}

func (c *NatAsyncHttpConnection) SetRequestQueueName(name string) {
	c.qreq = name
}

func (c *NatAsyncHttpConnection) getNextMessage() *NatHttpRequestMessage {
	if c.message == nil {
		c.message = new(NatHttpRequestMessage)
		c.message.Sync = false
	}
	m := c.message
	for m.Url != "" {
		if m.Next == nil {
			m.Next = new(NatHttpRequestMessage)
			m.Next.Sync = false
		}
		m = m.Next
	}
	return m
}

func (c *NatAsyncHttpConnection) AddGet(url string) {
	if url == "" {
		return
	}
	c.getNextMessage().Url = url
}

func (c *NatAsyncHttpConnection) AddPost(url string) error {
	return c.addJson(http.MethodPost, url, "")
}

func (c *NatAsyncHttpConnection) AddPostJson(url string, data interface{}) error {
	return c.addJson(http.MethodPost, url, data)
}

func (c *NatAsyncHttpConnection) AddPostForm(url string, data map[string][]string) {
	c.addForm(http.MethodPost, url, data)
}

func (c *NatAsyncHttpConnection) AddPut(url string) error {
	return c.addJson(http.MethodPut, url, "")
}

func (c *NatAsyncHttpConnection) AddPutJson(url string, data interface{}) error {
	return c.addJson(http.MethodPut, url, data)
}

func (c *NatAsyncHttpConnection) AddPutForm(url string, data map[string][]string) {
	c.addForm(http.MethodPut, url, data)
}

func (c *NatAsyncHttpConnection) AddDelete(url string) error {
	return c.addJson(http.MethodDelete, url, "")
}

func (c *NatAsyncHttpConnection) AddDeleteJson(url string, data interface{}) error {
	return c.addJson(http.MethodDelete, url, data)
}

func (c *NatAsyncHttpConnection) AddDeleteForm(url string, data map[string][]string) {
	c.addForm(http.MethodDelete, url, data)
}

func (c *NatAsyncHttpConnection) addJson(method, url string, data interface{}) error {
	if url == "" {
		return nil
	}
	m := c.getNextMessage()
	m.Method = method
	m.Url = url
	var err error
	m.JsonData, err = GetJsonString(data)
	return err
}

func (c *NatAsyncHttpConnection) addForm(method, url string, data map[string][]string) {
	if url == "" {
		return
	}
	message := c.getNextMessage()
	message.Method = method
	message.Url = url
	message.FormData = data
}

func (c *NatAsyncHttpConnection) Request() error {
	conn, err := amqp.Dial(c.qurl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	body, err := json.Marshal(c.message)
	if err != nil {
		return err
	}
	return ch.Publish(
		"",      // exchange
		c.qreq, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          body,
		})
}

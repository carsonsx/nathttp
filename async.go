package nathttp

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"net/http"
)

type NatAsyncHttpConnection struct {
	qurl    string
	qname   string
	message *NatHttpRequestMessage
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

func (c *NatAsyncHttpConnection) AddPostJson(url string, data interface{}) error {
	if url == "" {
		return nil
	}
	m := c.getNextMessage()
	m.Method = http.MethodPost
	m.Url = url
	var err error
	m.JsonData, err = GetJsonString(data)
	return err
}

func (c *NatAsyncHttpConnection) AddPostForm(url string, data map[string][]string) {
	if url == "" {
		return
	}
	message := c.getNextMessage()
	message.Method = http.MethodPost
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

	var body []byte
	body, err = json.Marshal(c.message)
	if err != nil {
		return err
	}
	return ch.Publish(
		"",      // exchange
		c.qname, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          body,
		})
}

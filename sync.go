package nathttp

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"net/http"
)

type NatHttpConnection struct {
	qurl string
	qreq string
	qres string
}

func (c *NatHttpConnection) SetRequestQueueName(name string) {
	c.qreq = name
}

func (c *NatHttpConnection) SetResponseQueueName(name string) {
	c.qres = name
}

func (c *NatHttpConnection) Get(url string) (string, error) {
	var message NatHttpRequestMessage
	message.Sync = true
	message.Url = url
	return c.Request(message)
}

func (c *NatHttpConnection) Post(url string) (string, error) {
	return c.sendJson(http.MethodPost, url, "")
}

func (c *NatHttpConnection) PostJson(url string, data interface{}) (string, error) {
	return c.sendJson(http.MethodPost, url, data)
}

func (c *NatHttpConnection) PostForm(url string, data map[string][]string) (string, error) {
	return c.sendForm(http.MethodPost, url, data)
}

func (c *NatHttpConnection) Put(url string) (string, error) {
	return c.sendJson(http.MethodPut, url, "")
}

func (c *NatHttpConnection) PutJson(url string, data interface{}) (string, error) {
	return c.sendJson(http.MethodPut, url, data)
}

func (c *NatHttpConnection) PutForm(url string, data map[string][]string) (string, error) {
	return c.sendForm(http.MethodPut, url, data)
}

func (c *NatHttpConnection) Delete(url string) (string, error) {
	return c.sendJson(http.MethodDelete, url, "")
}

func (c *NatHttpConnection) DeleteJson(url string, data interface{}) (string, error) {
	return c.sendJson(http.MethodDelete, url, data)
}

func (c *NatHttpConnection) DeleteForm(url string, data map[string][]string) (string, error) {
	return c.sendForm(http.MethodDelete, url, data)
}

func (c *NatHttpConnection) sendJson(method, url string, data interface{}) (string, error) {
	var message NatHttpRequestMessage
	message.Sync = true
	message.Url = url
	message.Method = method
	jsonData, err := GetJsonString(data)
	if err != nil {
		return "", err
	}
	message.JsonData = jsonData
	return c.Request(message)
}

func (c *NatHttpConnection) sendForm(method, url string, data map[string][]string) (string, error) {
	var message NatHttpRequestMessage
	message.Sync = true
	message.Url = url
	message.Method = method
	message.FormData = data
	return c.Request(message)
}



func (c *NatHttpConnection) Request(message NatHttpRequestMessage) (res string, err error) {
	conn, err := amqp.Dial(c.qurl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	var body []byte
	body, err = json.Marshal(message)
	if err != nil {
		return
	}

	corrId := uuid.NewV4().String()

	msgs, err := ch.Consume(
		c.qres, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failOnError(err, "Failed to register a consumer")

	err = ch.Publish(
		"",      // exchange
		c.qreq, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       c.qres,
			Body:          body,
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res = string(d.Body)
			failOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

package nathttp

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"net/http"
)

type NatHttpConnection struct {
	qurl  string
	qname string
}

func (c *NatHttpConnection) Get(url string) (string, error) {
	var message NatHttpRequestMessage
	message.Sync = true
	message.Url = url
	return c.Request(message)
}

func (c *NatHttpConnection) PostJson(url string, data interface{}) (string, error) {
	var message NatHttpRequestMessage
	message.Sync = true
	message.Url = url
	message.Method = http.MethodPost
	jsonData, err := GetJsonString(data)
	if err != nil {
		return "", err
	}
	message.JsonData = jsonData
	return c.Request(message)
}

func (c *NatHttpConnection) PostForm(url string, data map[string][]string) (string, error) {
	var message NatHttpRequestMessage
	message.Sync = true
	message.Url = url
	message.Method = http.MethodPost
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
	callback := c.qname + "_callback"

	msgs, err := ch.Consume(
		callback, // queue
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
		c.qname, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       callback,
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

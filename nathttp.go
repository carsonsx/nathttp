package nathttpclient

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"log"
	"net/url"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type natHttpRequestMessage struct {
	Method   string              `json:"method"`
	Url      string              `json:"url"`
	JsonData string              `json:"jsonData"`
	FormData url.Values          `json:"formData"`
	Sync     bool                `json:"sync"`
	Next     *natHttpRequestMessage `json:"next"`
}

type NatHttpConnection struct {
	qurl  string
	qname string
}

func CreateNatHttpConnection(qurl, qname string) *NatHttpConnection {
	return &NatHttpConnection{qurl, qname}
}

func (c *NatHttpConnection) Get (url string) (string, error)  {
	return c.get(url, true)
}

func (c *NatHttpConnection) AsyncGet (url string) (string, error)  {
	return c.get(url, false)
}

func (c *NatHttpConnection) get (url string, sync bool) (string, error) {
	var message natHttpRequestMessage
	message.Url = url
	message.Sync = sync
	return c.invokeHttp(message)
}


func (c *NatHttpConnection) invokeHttp(message natHttpRequestMessage) (res string, err error) {
	conn, err := amqp.Dial(c.qurl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	corrId := uuid.NewV4().String()

	var callback_qname string

	if message.Sync {

		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when usused
			true,  // exclusive
			false, // noWait
			nil,   // arguments
		)
		failOnError(err, "Failed to declare a queue")
		callback_qname = q.Name
	}

	var b []byte
	b, err = json.Marshal(message)
	if err != nil {
		return
	}
	err = ch.Publish(
		"",    // exchange
		c.qname, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       callback_qname,
			Body:          b,
		})
	failOnError(err, "Failed to publish a message")

	if message.Sync {
		msgs, err := ch.Consume(
			callback_qname, // queue
			"",             // consumer
			true,           // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		failOnError(err, "Failed to register a consumer")

		for d := range msgs {
			if corrId == d.CorrelationId {
				res = string(d.Body)
				failOnError(err, "Failed to convert body to integer")
				break
			}
		}
	}

	return
}

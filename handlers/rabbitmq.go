package handlers

import (
	"fmt"
	"os"

	syslog "github.com/RackSec/srslog"
	v2syslog "github.com/influxdata/go-syslog/v2"
	"github.com/influxdata/go-syslog/v2/rfc5424"
	"github.com/loafoe/go-rabbitmq"
	"github.com/streadway/amqp"
)

type RabbitMQHandler struct {
	debug  bool
	writer *syslog.Writer
	parser v2syslog.Machine
}

func NewRabbitMQHandler(promtailAddr string) (*RabbitMQHandler, error) {
	if promtailAddr == "" {
		return nil, fmt.Errorf("missing promtail address")
	}
	handler := &RabbitMQHandler{}

	parser := rfc5424.NewParser()

	if os.Getenv("DEBUG") == "true" {
		handler.debug = true
	}
	writer, err := syslog.Dial("tcp", promtailAddr,
		syslog.LOG_WARNING|syslog.LOG_DAEMON, "rabbitmq-logdrain")
	if err != nil {
		return nil, fmt.Errorf("promtail: %w", err)
	}
	writer.SetFramer(syslog.RFC5425MessageLengthFramer)
	writer.SetFormatter(RFC5424PassThroughFormatter)
	handler.writer = writer
	handler.parser = parser
	return handler, nil
}

func (h *RabbitMQHandler) CreateWorker(exchange, routingKey, queueName, consumerTag string) (chan bool, error) {
	doneChannel := make(chan bool)
	// Consumer
	consumer, err := rabbitmq.NewConsumer(rabbitmq.Config{
		RoutingKey:   routingKey,
		Exchange:     exchange,
		ExchangeType: "topic",
		Durable:      true,
		AutoDelete:   false,
		QueueName:    queueName,
		CTag:         consumerTag,
		HandlerFunc:  h.RabbitMQRFC5424Worker(doneChannel),
	})
	if err != nil {
		return nil, err
	}
	if err := consumer.Start(); err != nil {
		return nil, err
	}
	return doneChannel, nil
}

func (h *RabbitMQHandler) RabbitMQRFC5424Worker(doneChannel <-chan bool) rabbitmq.ConsumerHandlerFunc {
	return func(deliveries <-chan amqp.Delivery, done <-chan bool) {
		for {
			select {
			case d := <-deliveries:
				// TODO: PARSE message here
				ackDelivery(d)
				syslogMessage, err := h.parser.Parse(d.Body)
				fmt.Printf("version=%d\n", syslogMessage.Version())
				if err != nil {
					fmt.Printf("Error processing syslog message: %v\n", err)
					continue
				}
				_, _ = h.writer.Write(d.Body)
			case <-done:
				fmt.Printf("Worker received done message (server)...\n")
				return
			case <-doneChannel:
				fmt.Printf("Worker received done message (main)...\n")
				return
			}
		}
	}
}

func ackDelivery(d amqp.Delivery) {
	err := d.Ack(true)
	if err != nil {
		fmt.Printf("Error Acking delivery: %v\n", err)
	}
}

package worker

import (
	"github.com/streadway/amqp"

	"github.com/K-Phoen/flatrise/flatrise/emitter"
	"github.com/pkg/errors"
)

type rabbitMqWorker struct {
	channel *amqp.Channel

	searchRequestsQueue string

	crawler CrawlerFunc
	emitter emitter.Emitter
}

func NewRabbitMqWorker(channel *amqp.Channel, searchRequestsQueue string, crawler CrawlerFunc, emitter emitter.Emitter) (*rabbitMqWorker, error) {
	_, err := channel.QueueDeclare(
		searchRequestsQueue, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "Could not declare RabbitMq search requests queue")
	}

	return &rabbitMqWorker{
		channel:             channel,
		searchRequestsQueue: searchRequestsQueue,
		crawler:             crawler,
		emitter:             emitter,
	}, nil
}

func (w rabbitMqWorker) Run() error {
	msgs, err := w.channel.Consume(
		w.searchRequestsQueue, // queue
		"",                    // consumer
		true,                  // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		return errors.Wrap(err, "Could not register search offers consumer")
	}

	for _ = range msgs {
		go func() {
			w.crawler(w.emitter)
		}()
	}

	return nil
}

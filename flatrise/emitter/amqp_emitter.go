package emitter

import (
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/K-Phoen/flatrise/flatrise/model"
	"github.com/streadway/amqp"
)

const OffersQueue = "offers"

type amqpEmitter struct {
	channel *amqp.Channel
}

func NewAmqpEmitter(channel *amqp.Channel) (*amqpEmitter, error) {
	_, err := channel.QueueDeclare(
		OffersQueue, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "Could not declare RabbitMq offers queue")
	}

	emitter := &amqpEmitter{
		channel: channel,
	}

	return emitter, nil
}

func (emitter amqpEmitter) Emit(offer model.Offer) error {
	jsonOffer, err := json.Marshal(offer)
	if err != nil {
		return errors.Wrap(err, "Could not marshal offer")
	}

	err = emitter.channel.Publish(
		"",          // exchange
		OffersQueue, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonOffer,
		},
	)
	if err != nil {
		errors.Wrap(err, "Could not publish offer")
	}

	return nil
}

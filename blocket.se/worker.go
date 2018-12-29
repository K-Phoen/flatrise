package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const OffersQueue = "offers"
const BlocketSearchsQueue = "blocket_searchs"

func crawlOffers(channel *amqp.Channel) {
	resultChan, err := Crawl()
	if err != nil {
		log.Fatalf("Could not crawl offers: %s", err)
	}

	for offer := range resultChan {
		jsonOffer, err := json.Marshal(offer)
		if err != nil {
			log.Fatalf("Could not marshal offer: %s", err)
		}

		log.Infof("Found offer \"%s\" -- %s\n", offer.Title, offer.Identifier)

		err = channel.Publish(
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
			log.Fatalf("Could not publish offer: %s", err)
		}
	}
}

func main() {
	rabbitMqUrl, ok := os.LookupEnv("RABBITMQ_URL")
	if !ok {
		log.Fatal("The DSN to use to connect to RabbitMq must be specified by the RABBITMQ_URL environment variable.")
	}

	rabbitConn, err := amqp.Dial(rabbitMqUrl)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMq: %s", err)
	}
	defer rabbitConn.Close()

	channel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Could not connect to RabbitMq channel: %s", err)
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		OffersQueue, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalf("Could not declare RabbitMq offers queue: %s", err)
	}

	_, err = channel.QueueDeclare(
		BlocketSearchsQueue, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Fatalf("Could not declare RabbitMq search requests queue: %s", err)
	}

	msgs, err := channel.Consume(
		BlocketSearchsQueue, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		log.Fatalf("Could not register search offers consumer: %s", err)
	}

	log.Printf("Waiting for search requests. To exit press CTRL+C")

	for _ = range msgs {
		go func() {
			log.Printf("Received search request!")

			crawlOffers(channel)
		}()
	}
}

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/K-Phoen/flatrise/flatrise/emitter"
	"github.com/K-Phoen/flatrise/flatrise/worker"
)

const BlocketSearchsQueue = "blocket_searchs"

func newEmitter(channel *amqp.Channel) emitter.Emitter {
	stdoutEmitter := emitter.NewWriterEmitter(os.Stdout)
	amqpEmitter, err := emitter.NewAmqpEmitter(channel)
	if err != nil {
		log.Fatalf("Could not create emitter: %s", err)
	}

	return emitter.NewMultiEmitter(stdoutEmitter, amqpEmitter)
}

func crawlOffers(emitter emitter.Emitter) {
	log.Printf("Received search request!")

	resultChan, err := Crawl()
	if err != nil {
		log.Fatalf("Could not crawl offers: %s", err)
	}

	for offer := range resultChan {
		err := emitter.Emit(offer)
		if err != nil {
			log.Fatalf("Could not emit offer: %s", err)
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

	emitter := newEmitter(channel)
	worker, err := worker.NewRabbitMqWorker(channel, BlocketSearchsQueue, crawlOffers, emitter)
	if err != nil {
		log.Fatalf("Could not create RabbitMq worker: %s", err)
	}

	log.Printf("Waiting for search requests. To exit press CTRL+C")

	worker.Run()
}

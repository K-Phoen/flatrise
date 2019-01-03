package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/K-Phoen/flatrise/flatrise/emitter"
	"github.com/K-Phoen/flatrise/flatrise/worker"
)

const BoligPortalSearchsQueue = "boligportal_searchs"

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

	offers, err := Crawl("https://www.boligportal.dk/RAP/ads?placeIds[]=15&housingTypes[]=3&listViewResults=true&limitRecords=60&sort=paid&tid=5c0283e94d35b")
	if err != nil {
		log.Fatalf("Could not crawl offers: %s", err)
	}

	for _, offer := range offers {
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
	worker, err := worker.NewRabbitMqWorker(channel, BoligPortalSearchsQueue, crawlOffers, emitter)
	if err != nil {
		log.Fatalf("Could not create RabbitMq worker: %s", err)
	}

	log.Printf("Waiting for search requests. To exit press CTRL+C")

	worker.Run()
}

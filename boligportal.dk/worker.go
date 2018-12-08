package main

import (
  "encoding/json"
  "fmt"
  "log"
  "os"

  "github.com/streadway/amqp"
)

const OffersQueue = "offers"
const BoligPortalSearchsQueue = "boligportal_searchs"

func crawlOffers(channel *amqp.Channel) {
  offers, err := Crawl("https://www.boligportal.dk/RAP/ads?placeIds[]=15&housingTypes[]=3&listViewResults=true&limitRecords=60&sort=paid&tid=5c0283e94d35b")
  if err != nil {
    log.Fatalf("Could not crawl offers: %s", err)
  }

  for _, offer := range offers {
    jsonOffer, err := json.Marshal(offer)
    if err != nil {
      log.Fatalf("Could not marshal offer: %s", err)
    }

    fmt.Printf("%s -- %s\n", offer.Title, offer.Identifier)

    err = channel.Publish(
      "",     // exchange
      OffersQueue, // routing key
      false,  // mandatory
      false,  // immediate
      amqp.Publishing {
        ContentType: "application/json",
        Body: jsonOffer,
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
    true,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  if err != nil {
    log.Fatalf("Could not declare RabbitMq offers queue: %s", err)
  }

  _, err = channel.QueueDeclare(
    BoligPortalSearchsQueue, // name
    true,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  if err != nil {
    log.Fatalf("Could not declare RabbitMq search requests queue: %s", err)
  }

  msgs, err := channel.Consume(
    BoligPortalSearchsQueue, // queue
    "",     // consumer
    true,   // auto-ack
    false,  // exclusive
    false,  // no-local
    false,  // no-wait
    nil,    // args
  )
  if err != nil {
    log.Fatalf("Could not register search offers consumer: %s", err)
  }

  forever := make(chan bool)

  go func() {
    for _ = range msgs {
      log.Printf("Received search request!")

      crawlOffers(channel)
    }
  }()

  log.Printf("[*] Waiting for search requests. To exit press CTRL+C")
  <-forever
}

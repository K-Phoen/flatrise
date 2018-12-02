package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "time"

  "github.com/streadway/amqp"
)

type Location struct {
	Lat float64 `json:"lat"`
  Lon float64 `json:"lon"`
}

type Offer struct {
	Identifier string `json:"identifier"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price int `json:"price"`
	Area int `json:"area"`
	Rooms int `json:"rooms"`
	Location Location `json:"location"`
}

func Crawl(url string) (offers []Offer, err error) {
  var httpClient = &http.Client{Timeout: 10 * time.Second}

  response, err := httpClient.Get(url)
  if err != nil {
    return offers, err
  }

  defer response.Body.Close()

  jsonPayload, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return offers, err
  }

  var payload map[string]interface{}
  json.Unmarshal(jsonPayload, &payload)

  data := payload["data"].(map[string]interface{})
  properties := data["properties"].(map[string]interface{})
  propertiesCollection := properties["collection"].([]interface{})

  for _, propertyData := range propertiesCollection {
    offerData := propertyData.(map[string]interface{})

    offer := Offer {
      Identifier: "https://www.boligportal.dk/en" + offerData["url"].(string),
      Title: offerData["title"].(string),
      Description: offerData["description"].(string),
      Price: int(offerData["monthlyPrice"].(float64)),
      Area: int(offerData["sizeM2"].(float64)),
      Rooms: int(offerData["numRooms"].(float64)),
      Location: Location { Lat: offerData["lat"].(float64), Lon: offerData["lng"].(float64) },
    }

    offers = append(offers, offer)
  }

  return offers, err
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

  queue, err := channel.QueueDeclare(
    "offers", // name
    true,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  if err != nil {
    log.Fatalf("Could not declare RabbitMq queue: %s", err)
	}

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
      queue.Name, // routing key
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

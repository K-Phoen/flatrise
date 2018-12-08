package main

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "time"
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

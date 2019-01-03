package main

import (
	"encoding/json"
	"github.com/K-Phoen/flatrise/flatrise/model"
	"io/ioutil"
	"net/http"
	"time"
)

// divide a DKK amount by this rate to get the price in euros
const DKKToEuroRate = 7.46635

func Crawl(url string) (offers []model.Offer, err error) {
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

		price := int(offerData["monthlyPrice"].(float64))

		offer := model.Offer{
			Identifier:  "https://www.boligportal.dk/en" + offerData["url"].(string),
			Title:       offerData["title"].(string),
			Description: offerData["description"].(string),
			Price:       price,
			Currency:    "DKK",
			PriceEur:    int(float64(price) / DKKToEuroRate),
			Area:        int(offerData["sizeM2"].(float64)),
			Rooms:       int(offerData["numRooms"].(float64)),
			Location:    model.Location{Lat: offerData["lat"].(float64), Lon: offerData["lng"].(float64)},
		}

		offers = append(offers, offer)
	}

	return offers, err
}

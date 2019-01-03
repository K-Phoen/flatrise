package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/K-Phoen/flatrise/flatrise/model"
	"github.com/PuerkitoBio/goquery"
)

// divide a SEK amount by this rate to get the price in euros
const SEKToEuroRate = 10.2704

func stringToFloat(input string) float64 {
	f, err := strconv.ParseFloat(input, 64)

	if err != nil {
		panic(fmt.Sprintf("Unable to convert '%v' to an int", input))
	}

	return f
}

func stringToInt(input string) int {
	if len(input) == 0 {
		return 0
	}

	cleanInput := input

	if index := strings.Index(input, ","); index != -1 {
		cleanInput = input[:index]
	}

	i, err := strconv.Atoi(cleanInput)

	if err != nil {
		panic(fmt.Sprintf("Unable to convert '%v' to an int", input))
	}

	return i
}

func httpGet(url string) (io.Reader, error) {
	var httpClient = &http.Client{Timeout: 10 * time.Second}

	response, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func extractLocation(identifier string) model.Location {
	source, err := httpGet(identifier)
	if err != nil {
		return model.Location{}
	}

	doc, err := goquery.NewDocumentFromReader(source)
	if err != nil {
		return model.Location{}
	}

	mapLinks := doc.Find("*[id=hitta-map-broker]")

	if mapLinks.Length() == 0 {
		return model.Location{}
	}

	src, _ := mapLinks.First().Attr("src")
	data := strings.Split(strings.Split(strings.Split(src, "/")[7], "?")[0], ":")

	return model.Location{
		Lat: stringToFloat(data[0]),
		Lon: stringToFloat(data[1]),
	}
}

func buildOffer(offerData map[string]interface{}) model.Offer {
	identifier := fmt.Sprintf("https://www.blocket.se/stockholm/seo-friendly-slug_%s.htm", offerData["id"].(string))

	price := stringToInt(offerData["monthly_rent"].(string))

	return model.Offer{
		Identifier:  identifier,
		Title:       offerData["address"].(string),
		Description: "",
		Price:       price,
		Currency:    "SEK",
		PriceEur:    int(float64(price) / SEKToEuroRate),
		Area:        stringToInt(offerData["sqm"].(string)),
		Rooms:       stringToInt(offerData["rooms"].(string)),
		Location:    extractLocation(identifier),
	}
}

func Crawl() (offersChan chan model.Offer, err error) {
	url := "https://www.blocket.se/karta/items?ca=11&ca=11&st=s&cg=3020&sort=&ps=&pe=&ss=&se=&ros=&roe=&mre=&q=&is=1&f=b&w=3&ac=0MNXXY7CTORXWG23IN5WG2000&zl=12&ne=59.39389826993069%2C18.441925048828125&sw=59.2802650449542%2C17.865142822265625"

	response, err := httpGet(url)
	if err != nil {
		return offersChan, err
	}

	jsonPayload, _ := ioutil.ReadAll(response)
	if err != nil {
		return offersChan, err
	}

	var payload map[string]interface{}
	json.Unmarshal(jsonPayload, &payload)

	listItems := payload["list_items"].([]interface{})

	offersChan = make(chan model.Offer, len(listItems))
	var wg sync.WaitGroup

	go func() {
		for _, item := range listItems {
			wg.Add(1)

			go func(offerData map[string]interface{}) {
				offersChan <- buildOffer(offerData)

				wg.Done()
			}(item.(map[string]interface{}))
		}

		wg.Wait()

		close(offersChan)
	}()

	return offersChan, err
}

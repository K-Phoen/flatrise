package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
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

// ex: 1,5 rum
func extractRooms(input string) int {
	if len(input) == 0 {
		return 0
	}

	cleanInput := input

	if index := strings.Index(input, " "); index != -1 {
		cleanInput = input[:index]
	}

	return stringToInt(cleanInput)
}

// ex: 15 000 kr/mån
func extractRent(input string) int {
	if len(input) == 0 {
		return 0
	}

	cleanInput := input

	if index := strings.Index(input, "kr"); index != -1 {
		cleanInput = input[:index]
	}

	cleanInput = strings.Replace(cleanInput, " ", "", -1)

	return stringToInt(cleanInput)
}

// ex: 35 m²
func extractArea(input string) int {
	if len(input) == 0 {
		return 0
	}

	cleanInput := input

	if index := strings.Index(input, " "); index != -1 {
		cleanInput = input[:index]
	}

	return stringToInt(cleanInput)
}

func httpGet(url string) (io.ReadCloser, error) {
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

func buildOffer(s *goquery.Selection) model.Offer {
	linkSelector := s.Find(".item_link")
	roomsSelector := s.Find(".rooms")
	rentSelector := s.Find(".monthly_rent")
	areaSelector := s.Find(".size")

	identifier, _ := linkSelector.Attr("href")
	rent := extractRent(strings.Trim(rentSelector.Text(), " \t\n"))

	return model.Offer{
		Identifier:  identifier,
		Title:       strings.Trim(linkSelector.Text(), " \t\n"),
		Description: "",
		Price:       rent,
		Currency:    "SEK",
		PriceEur:    int(float64(rent) / SEKToEuroRate),
		Area:        extractArea(strings.Trim(areaSelector.Text(), " \t\n")),
		Rooms:       extractRooms(strings.Trim(roomsSelector.Text(), " \t\n")),
		Location:    extractLocation(identifier),
	}
}

func Crawl() (chan model.Offer, error) {
	url := "https://www.blocket.se/bostad/uthyres/stockholm?cg_multi=3020&sort=&ss=&se=&ros=&roe=&bs=&be=&mre=&q=&q=&q=&is=1&save_search=1&l=0&md=th&f=p&f=c&f=b"

	response, err := httpGet(url)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch remote offers listing")
	}
	defer response.Close()

	doc, err := goquery.NewDocumentFromReader(response)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse remote offers listing")
	}

	offersSelector := doc.Find(".item_row")
	offersChan := make(chan model.Offer, offersSelector.Length())

	var wg sync.WaitGroup

	go func() {
		offersSelector.Each(func(i int, s *goquery.Selection) {
			wg.Add(1)

			go func(selection *goquery.Selection) {
				offersChan <- buildOffer(selection)

				wg.Done()
			}(s)
		})

		wg.Wait()

		close(offersChan)
	}()

	return offersChan, err
}

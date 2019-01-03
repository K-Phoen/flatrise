package model

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Offer struct {
	Identifier  string   `json:"identifier"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Currency    string   `json:"currency"`
	PriceEur    int      `json:"price_eur"`
	Area        int      `json:"area"`
	Rooms       int      `json:"rooms"`
	Location    Location `json:"location"`
}

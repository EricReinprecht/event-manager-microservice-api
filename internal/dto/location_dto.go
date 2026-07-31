package dto

type PartyLocation struct {
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
	City        string `json:"city"`
	Country     string `json:"country"`
	PostalCode  string `json:"postalCode"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Timezone string `json:"timezone"`
}

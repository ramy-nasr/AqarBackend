package domain

type Transaction struct {
	ID           string       `json:"id"`
	City         string       `json:"city"`
	Latitude     float64      `json:"latitude"`
	Longitude    float64      `json:"longitude"`
	Time         string       `json:"time"`
	Price        int          `json:"price"`
	PropertyType PropertyType `json:"type"`
}

package operation

import "time"

type Operations []Operation

type Operation struct {
	Id             string    `json:"id"`
	Figi           string    `json:"figi"`
	InstrumentType string    `json:"instrument_type"`
	Description    string    `json:"description"`
	Quantity       int64     `json:"quantity"`
	Payment        float64   `json:"payment"`
	Currency       string    `json:"currency"`
	Time           time.Time `json:"time"`
}

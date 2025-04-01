package dto

import "time"

type Operation struct {
	ID             string    `json:"id"`
	Figi           string    `json:"figi"`
	InstrumentType string    `json:"instrumentType"`
	Quantity       int64     `json:"quantity"`
	Payment        float64   `json:"payment"`
	Currency       string    `json:"currency"`
	Date           time.Time `json:"date"`
}

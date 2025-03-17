package account

import "time"

type Account struct {
	Id          string    `json:"id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	OpenedDate  time.Time `json:"opened_date"`
	ClosedDate  time.Time `json:"closed_date"`
	AccessLevel string    `json:"access_level"`
}

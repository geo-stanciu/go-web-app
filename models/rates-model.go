package models

import "time"

// Rate - exchange rate struct
type Rate struct {
	ReferenceCurrency string    `json:"reference_currency" sql:"reference_currency"`
	Currency          string    `json:"currency" sql:"currency"`
	Date              time.Time `json:"date" sql:"exchange_date"`
	Value             float64   `json:"value" sql:"rate"`
}

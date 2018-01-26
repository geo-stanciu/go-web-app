package models

// ExchangeRatesResponseModel - Exchange Rates Response Model
type ExchangeRatesResponseModel struct {
	GenericResponseModel
	Rates []*Rate `json:"rates"`
}

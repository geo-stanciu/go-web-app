package models

// UsersResponseModel - Exchange Rates Response Model
type UsersResponseModel struct {
	GenericResponseModel
	UserModel []*UserModel `json:"users"`
}

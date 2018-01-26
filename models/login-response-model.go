package models

type LoginResponseModel struct {
	GenericResponseModel
	TemporaryPassword bool
}

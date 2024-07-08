package models

type Address struct {
	ID           int32
	UserID       int32
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Country      string
	PostalCode   string
	IsDefault    bool
}

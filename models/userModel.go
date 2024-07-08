package models

type User struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    int64  `json:"phone"`
}

type UserAddress struct {
	UserID    int32
	AddressID int32
	IsDefault bool
}

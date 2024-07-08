package models

type Products struct {
	ID           int32
	CategoryID   int32
	Category     Categories
	ProductName  string
	Description  string
	Price        int32
	Stock        int32
	ProductImage string
}

type Categories struct {
	ID   int32
	Name string
}

type Cart struct {
	ID        int32
	UserID    int32
	ProductID int32
}

type CartItems struct {
	ID        int32
	CartID    int32
	Product   Products
	UserID    int32
	ProductID int32
	Quantity  int32
	Price     int32
}

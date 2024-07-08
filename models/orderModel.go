package models

type Orders struct {
	ID                int32
	RefNumber         string
	UserID            int32
	PaymentIntentID   string
	ShippingAddressID int32
	OrderStatusID     int32
	OrderTotal        int32
	PaymentStatus     string
	UserName          string
	UserEmail         string
	UserPhone         int
	OrderStatus       string
}

type OrderStatus struct {
	ID     int32
	Status string
}

type OrderItems struct {
	ID        int32
	OrderID   int32
	Product   Products
	ProductID int32
	Quantity  int32
	Price     int32
}

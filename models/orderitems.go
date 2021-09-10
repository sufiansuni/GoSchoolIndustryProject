package models

// Operations for order_items database: Insert(Create), Select(Read), Update, Delete

type OrderItem struct {
	ID       int
	OrderID  int
	FoodID   int
	Quantity int
	Subtotal float64
}
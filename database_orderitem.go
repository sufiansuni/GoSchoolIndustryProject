package main

// Operations for order_items database: Insert(Create), Select(Read), Update, Delete

type order_item struct {
	ID       int
	OrderID  int
	FoodID   int
	Quantity int
	Subtotal float64
}

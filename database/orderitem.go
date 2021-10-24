package database

import (
	"GoIndustryProject/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

// Operations for order_items database: Insert(Create), Select(Read), Update, Delete

// Insert a new order_item entry into database
func InsertOrderItem(db *sql.DB, myOrderItem models.OrderItem) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO order_items (OrderID, FoodID, FoodName, Quantity, SubtotalPrice, SubTotalCalories) " +
		"VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		myOrderItem.OrderID,
		myOrderItem.FoodID,
		myOrderItem.FoodName,
		myOrderItem.Quantity,
		myOrderItem.SubtotalPrice,
		myOrderItem.SubtotalCalories,
	)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = errors.New("no rows updated")
	}
	return
}

// Select/Read a order_item entry from database with a ID input
func SelectOrderItemByID(db *sql.DB, ID int) (myOrderItem models.OrderItem, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM order_items WHERE ID=?"
	err = db.QueryRowContext(ctx, query, ID).Scan(
		&myOrderItem.ID,
		&myOrderItem.OrderID,
		&myOrderItem.FoodID,
		&myOrderItem.FoodName,
		&myOrderItem.Quantity,
		&myOrderItem.SubtotalPrice,
		&myOrderItem.SubtotalCalories,
	)
	return
}

// Select/Read order_items from database with a orderID input
func SelectOrderItemsByOrderID(db *sql.DB, orderID int) (myOrderItems []models.OrderItem, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM order_items WHERE orderID=?"
	rows, err := db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var myOrderItem models.OrderItem
		err = rows.Scan(
			&myOrderItem.ID,
			&myOrderItem.OrderID,
			&myOrderItem.FoodID,
			&myOrderItem.FoodName,
			&myOrderItem.Quantity,
			&myOrderItem.SubtotalPrice,
			&myOrderItem.SubtotalCalories,
		)
		if err != nil {
			return nil, err
		}
		myOrderItems = append(myOrderItems, myOrderItem)
	}

	return
}

// Update a order_item entry in database
func UpdateOrderItem(db *sql.DB, myOrderItem models.OrderItem) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE order_items SET OrderID=?, FoodID=?, FoodName=?, Quantity=?, SubtotalPrice=?, SubTotalCalories=? " +
		"WHERE ID=?"

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		myOrderItem.OrderID,
		myOrderItem.FoodID,
		myOrderItem.FoodName,
		myOrderItem.Quantity,
		myOrderItem.SubtotalPrice,
		myOrderItem.SubtotalCalories,
		myOrderItem.ID,
	)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = errors.New("no rows updated")
	}
	return
}

// Delete a order_item entry in database
func DeleteOrderItem(db *sql.DB, orderItemID int) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM order_items WHERE ID=?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, orderItemID)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = errors.New("no rows updated")
	}
	return
}

// Select/Read order_items from database with a orderID and foodID input
func SelectOrderItemsByOrderIDAndFoodID(db *sql.DB, orderID int, foodID int) (myOrderItems []models.OrderItem, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM order_items WHERE orderID=? AND foodID=?"
	rows, err := db.QueryContext(ctx, query, orderID, foodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var myOrderItem models.OrderItem
		err = rows.Scan(
			&myOrderItem.ID,
			&myOrderItem.OrderID,
			&myOrderItem.FoodID,
			&myOrderItem.FoodName,
			&myOrderItem.Quantity,
			&myOrderItem.SubtotalPrice,
			&myOrderItem.SubtotalCalories,
		)
		if err != nil {
			return nil, err
		}
		myOrderItems = append(myOrderItems, myOrderItem)
	}

	return
}

package database

import (
	"GoIndustryProject/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

// Operations for orders database: Insert(Create), Select(Read), Update, Delete

// Insert a new order entry into database
func InsertOrder(db *sql.DB, myOrder models.Order) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO orders (Username, RestaurantID, RestaurantName, Status, Collection, Date, " +
		"UserAddress, UserUnit, UserLat, UserLng, " +
		"RestaurantAddress, RestaurantUnit, RestaurantLat, RestaurantLng, " +
		"TotalPrice, TotalCalories, BurnCalories) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		myOrder.Username,
		myOrder.RestaurantID,
		myOrder.RestaurantName,
		myOrder.Status,
		myOrder.Collection,
		myOrder.Date,
		myOrder.UserAddress,
		myOrder.UserUnit,
		myOrder.UserLat,
		myOrder.UserLng,
		myOrder.RestaurantAddress,
		myOrder.RestaurantUnit,
		myOrder.RestaurantLat,
		myOrder.RestaurantLng,
		myOrder.TotalPrice,
		myOrder.TotalCalories,
		myOrder.BurnCalories,
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

// Select/Read order entries from database with a username and status input
func SelectOrdersByUsernameAndStatus(db *sql.DB, username string, status string) ([]models.Order, error) {
	var myOrders []models.Order

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM orders WHERE Username=? AND Status=?"

	rows, err := db.QueryContext(ctx, query, username, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var myOrder models.Order
		err := rows.Scan(
			&myOrder.ID,
			&myOrder.Username,
			&myOrder.RestaurantID,
			&myOrder.RestaurantName,
			&myOrder.Status,
			&myOrder.Collection,
			&myOrder.Date,
			&myOrder.UserAddress,
			&myOrder.UserUnit,
			&myOrder.UserLat,
			&myOrder.UserLng,
			&myOrder.RestaurantAddress,
			&myOrder.RestaurantUnit,
			&myOrder.RestaurantLat,
			&myOrder.RestaurantLng,
			&myOrder.TotalPrice,
			&myOrder.TotalCalories,
			&myOrder.BurnCalories,
		)
		if err != nil {
			return nil, err
		}
		myOrders = append(myOrders, myOrder)
	}
	return myOrders, err
}

// Select/Read order entries from database with a restaurantID and status input
func SelectOrdersByRestaurantIDAndStatus(db *sql.DB, restaurantID int, status string) ([]models.Order, error) {
	var myOrders []models.Order

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM orders WHERE RestaurantID=? AND Status=?"

	rows, err := db.QueryContext(ctx, query, restaurantID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var myOrder models.Order
		err := rows.Scan(
			&myOrder.ID,
			&myOrder.Username,
			&myOrder.RestaurantID,
			&myOrder.RestaurantName,
			&myOrder.Status,
			&myOrder.Collection,
			&myOrder.Date,
			&myOrder.UserAddress,
			&myOrder.UserUnit,
			&myOrder.UserLat,
			&myOrder.UserLng,
			&myOrder.RestaurantAddress,
			&myOrder.RestaurantUnit,
			&myOrder.RestaurantLat,
			&myOrder.RestaurantLng,
			&myOrder.TotalPrice,
			&myOrder.TotalCalories,
			&myOrder.BurnCalories,
		)
		if err != nil {
			return nil, err
		}
		myOrders = append(myOrders, myOrder)
	}
	return myOrders, err
}

// Update an order entry in database
func UpdateOrder(db *sql.DB, myOrder models.Order) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE orders SET Username=?, RestaurantID=?, RestaurantName=?, Status=?, Collection=?, Date=?, " +
		"UserAddress=?, UserUnit=?, UserLat=?, UserLng=?, " +
		"RestaurantAddress=?, RestaurantUnit=?, RestaurantLat=?, RestaurantLng=?, " +
		"TotalPrice=?, TotalCalories=?, BurnCalories=? " +
		"WHERE ID=?"

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		myOrder.Username,
		myOrder.RestaurantID,
		myOrder.RestaurantName,
		myOrder.Status,
		myOrder.Collection,
		myOrder.Date,
		myOrder.UserAddress,
		myOrder.UserUnit,
		myOrder.UserLat,
		myOrder.UserLng,
		myOrder.RestaurantAddress,
		myOrder.RestaurantUnit,
		myOrder.RestaurantLat,
		myOrder.RestaurantLng,
		myOrder.TotalPrice,
		myOrder.TotalCalories,
		myOrder.BurnCalories,
		myOrder.ID,
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

// Delete an order entry in database
func DeleteOrder(db *sql.DB, orderID int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM orders WHERE ID=?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, orderID)
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

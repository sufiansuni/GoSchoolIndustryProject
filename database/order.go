package database

import "GoIndustryProject/models"

// Operations for orders database: Insert(Create), Select(Read), Update, Delete

// Insert a new order entry into database
func InsertOrder(myOrder models.Order) error {

	statement := "INSERT INTO orders (Username, Status, Date, Address, PostalCode, Lat, Lng) VALUES (?,?,?,?,?,?,?)"
	_, err := DB.Exec(statement,
		// myOrder.Username,
		// myOrder.Status,
		// myOrder.Date,
		// myOrder.Address,
		// myOrder.PostalCode,
		// myOrder.Lat,
		// myOrder.Lng,
	)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read order entries from database with a username and status input
func SelectOrders(username string, status string) ([]models.Order, error) {
	var myOrder models.Order
	var myOrders []models.Order

	query := "SELECT * FROM orders WHERE Username=? AND Status=?"
	results, err := DB.Query(query, username, status)
	if err != nil {
		return myOrders, err
	}
	defer results.Close()
	for results.Next() {
		err := results.Scan(
			// &myOrder.ID,
			// &myOrder.Username,
			// &myOrder.Status,
			// &myOrder.Date,
			// &myOrder.Address,
			// &myOrder.PostalCode,
			// &myOrder.Lat,
			// &myOrder.Lng,
		)
		if err != nil {
			return myOrders, err
		}
		myOrders = append(myOrders, myOrder)
	}
	return myOrders, err
}

// Update an order entry in database
func UpdateOrder(myOrder models.Order) error {
	statement := "UPDATE orders SET Username=?, Status=?, Date=?, " +
		"Address=?, PostalCode=?,  Lat=?, Lng=? " +
		"WHERE ID=?"

	_, err := DB.Exec(statement,
		// myOrder.Username,
		// myOrder.Status,
		// myOrder.Date,
		// myOrder.Address,
		// myOrder.PostalCode,
		// myOrder.Lat,
		// myOrder.Lng,
		// myOrder.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete an order entry in database
func DeleteOrder(orderID int) error {
	_, err := DB.Exec("DELETE FROM orders WHERE ID=?",
		orderID)
	if err != nil {
		return err
	}
	return nil
}

package main

// Operations for orders database: Insert(Create), Select(Read), Update, Delete

type order struct {
	ID         int
	Username   string
	Status     string
	Date       string
	Address    string
	PostalCode int
	Lat        float64
	Lng        float64
}

// Insert a new order entry into database
func insertOrder(myOrder order) error {

	statement := "INSERT INTO orders (Username, Status, Date, Address, PostalCode, Lat, Lng) VALUES (?,?,?,?,?,?,?)"
	_, err := db.Exec(statement,
		myOrder.Username,
		myOrder.Status,
		myOrder.Date,
		myOrder.Address,
		myOrder.PostalCode,
		myOrder.Lat,
		myOrder.Lng,
	)
	if err != nil {
		return err
	}
	return nil
}

// Select/Read order entries from database with a username and status input
func selectOrders(username string, status string) ([]order, error) {
	var myOrder order
	var myOrders []order

	query := "SELECT * FROM orders WHERE Username=? AND Status=?"
	results, err := db.Query(query, username, status)
	if err != nil {
		return myOrders, err
	}
	defer results.Close()
	for results.Next() {
		err := results.Scan(
			&myOrder.ID,
			&myOrder.Username,
			&myOrder.Status,
			&myOrder.Date,
			&myOrder.Address,
			&myOrder.PostalCode,
			&myOrder.Lat,
			&myOrder.Lng,
		)
		if err != nil {
			return myOrders, err
		}
		myOrders = append(myOrders, myOrder)
	}
	return myOrders, err
}

// Update an order entry in database
func updateOrder(myOrder order) error {
	statement := "UPDATE orders SET Username=?, Status=?, Date=?, " +
		"Address=?, PostalCode=?,  Lat=?, Lng=? " +
		"WHERE ID=?"

	_, err := db.Exec(statement,
		myOrder.Username,
		myOrder.Status,
		myOrder.Date,
		myOrder.Address,
		myOrder.PostalCode,
		myOrder.Lat,
		myOrder.Lng,
		myOrder.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete an order entry in database
func deleteOrder(orderID int) error {
	_, err := db.Exec("DELETE FROM orders WHERE ID=?",
		orderID)
	if err != nil {
		return err
	}
	return nil
}

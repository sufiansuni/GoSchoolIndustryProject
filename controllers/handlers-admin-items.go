package controllers

import (
	"GoIndustryProject/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handles request of "/admin/orders/{orderID}/items/{itemID}/add" page
func adminOrderItemAdd(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(req)

	itemID, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrderItem, err := database.SelectOrderItemByID(database.DB, itemID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	pricePerItem := myOrderItem.SubtotalPrice / float64(myOrderItem.Quantity)
	caloriesPerItem := myOrderItem.SubtotalCalories / myOrderItem.Quantity

	myOrderItem.Quantity++
	myOrderItem.SubtotalPrice += pricePerItem
	myOrderItem.SubtotalCalories += caloriesPerItem

	err = database.UpdateOrderItem(database.DB, myOrderItem)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrder, err := database.SelectOrderByID(database.DB, myOrderItem.OrderID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrder.TotalItems++
	myOrder.TotalPrice += pricePerItem
	myOrder.TotalCalories += caloriesPerItem

	err = database.UpdateOrder(database.DB, myOrder)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, "/admin/orders/"+vars["orderID"]+"/items", http.StatusSeeOther)
}

// Handles request of "/admin/orders/{orderID}/items/{itemID}/subtract" page
func adminOrderItemSubtract(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(req)

	itemID, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrderItem, err := database.SelectOrderItemByID(database.DB, itemID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if myOrderItem.Quantity > 1 {
		pricePerItem := myOrderItem.SubtotalPrice / float64(myOrderItem.Quantity)
		caloriesPerItem := myOrderItem.SubtotalCalories / myOrderItem.Quantity

		myOrderItem.Quantity--
		myOrderItem.SubtotalPrice -= pricePerItem
		myOrderItem.SubtotalCalories -= caloriesPerItem

		err = database.UpdateOrderItem(database.DB, myOrderItem)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		myOrder, err := database.SelectOrderByID(database.DB, myOrderItem.OrderID)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		myOrder.TotalItems--
		myOrder.TotalPrice -= pricePerItem
		myOrder.TotalCalories -= caloriesPerItem

		err = database.UpdateOrder(database.DB, myOrder)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

	}

	http.Redirect(res, req, "/admin/orders/"+vars["orderID"]+"/items", http.StatusSeeOther)
}

// Handles request of "/admin/orders/{orderID}/items/{itemID}/delete" page
func adminOrderItemDelete(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(req)
	itemID, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrderItem, err := database.SelectOrderItemByID(database.DB, itemID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrder, err := database.SelectOrderByID(database.DB, myOrderItem.OrderID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = database.DeleteOrderItem(database.DB, itemID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrder.TotalItems -= myOrderItem.Quantity
	myOrder.TotalPrice -= myOrderItem.SubtotalPrice
	myOrder.TotalCalories -= myOrderItem.SubtotalCalories

	err = database.UpdateOrder(database.DB, myOrder)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, "/admin/orders/"+vars["orderID"]+"/items", http.StatusSeeOther)
}

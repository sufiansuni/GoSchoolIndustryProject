package controllers

import (
	"GoIndustryProject/api"
	"GoIndustryProject/database"
	"GoIndustryProject/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handles request of "/admin/orders" page
func adminOrders(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}

	myOrders, err := database.SelectAllOrders(database.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		User   models.User
		Orders []models.Order
	}{
		myUser,
		myOrders,
	}

	tpl.ExecuteTemplate(res, "admin-orders.html", data)
}

// Handles request of "/admin/orders/{orderID}/items" page
func adminOrderItems(res http.ResponseWriter, req *http.Request) {
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
	targetID, err := strconv.Atoi(vars["orderID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrder, err := database.SelectOrderByID(database.DB, targetID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrderItems, err := database.SelectOrderItemsByOrderID(database.DB, myOrder.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		User       models.User
		Order      models.Order
		OrderItems []models.OrderItem
	}{
		myUser,
		myOrder,
		myOrderItems,
	}

	tpl.ExecuteTemplate(res, "admin-orders-items.html", data)
}

// Handles request of "/admin/orders/{orderID}/details" page
func adminOrderDetails(res http.ResponseWriter, req *http.Request) {
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
	targetID, err := strconv.Atoi(vars["orderID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrder, err := database.SelectOrderByID(database.DB, targetID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	userLat := fmt.Sprintf("%f", myOrder.UserLat)
	userLng := fmt.Sprintf("%f", myOrder.UserLng)

	userMapLink := api.OneMapGenerateMapPNGSingle(userLat, userLng)

	restaurantLat := fmt.Sprintf("%f", myOrder.RestaurantLat)
	restaurantLng := fmt.Sprintf("%f", myOrder.RestaurantLng)

	restaurantMapLink := api.OneMapGenerateMapPNGSingle(restaurantLat, restaurantLng)

	data := struct {
		User              models.User
		Order             models.Order
		UserMapLink       string
		RestaurantMapLink string
	}{
		myUser,
		myOrder,
		userMapLink,
		restaurantMapLink,
	}

	tpl.ExecuteTemplate(res, "admin-orders-details.html", data)
}

// Handles request of "/admin/orders/{orderID}/delete" page
func adminOrderDelete(res http.ResponseWriter, req *http.Request) {
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
	targetID, err := strconv.Atoi(vars["orderID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = database.DeleteOrder(database.DB, targetID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	http.Redirect(res, req, "/admin/orders", http.StatusSeeOther)
}

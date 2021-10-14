package controllers

import (
	"GoIndustryProject/database"
	"GoIndustryProject/models"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

// Handles request of "/admin/restaurants" page
func adminRestaurants(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}

	var myRestaurants []models.Restaurant
	myRestaurants, err := database.SelectAllRestaurants(database.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		User        models.User
		Restaurants []models.Restaurant
	}{
		myUser,
		myRestaurants,
	}
	tpl.ExecuteTemplate(res, "admin-restaurants.html", data)
}

// Handles request of "/admin/restaurants/new" page
func adminRestaurantNew(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		name := req.FormValue("name")
		description := req.FormValue("description")

		if name != "" {
			myRestaurant := models.Restaurant{
				Name:        name,
				Description: description,
			}

			err := database.InsertRestaurant(myRestaurant) // previouslymapUsers[username] = myUser
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			} else {
				fmt.Println("Restaurant Created:", name)
			}
		}

		// redirect to admin page (restaurants)
		http.Redirect(res, req, "/admin/restaurants", http.StatusSeeOther)
		return

	}

	data := struct {
		User models.User
	}{
		myUser,
	}

	tpl.ExecuteTemplate(res, "newrestaurant.html", data)
}

// Handles request of "/admin/restaurants/{restaurantID}/profile" page
func adminRestaurantProfile(res http.ResponseWriter, req *http.Request) {
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
	targetID, err := strconv.Atoi(vars["restaurantID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	targetRestaurant, err := database.SelectRestaurant(targetID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodPost {
		unchangedRestaurant := targetRestaurant

		if req.FormValue("name") != "" {
			targetRestaurant.Name = req.FormValue("name")
		}

		if req.FormValue("description") != "" {
			targetRestaurant.Description = req.FormValue("description")
		}

		if req.FormValue("halal") != "" {
			switch req.FormValue("halal") {
			case "true":
				targetRestaurant.Halal = true
			case "false":
				targetRestaurant.Halal = false
			default:
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			targetRestaurant.Halal = false
		}

		if req.FormValue("vegan") != "" {
			switch req.FormValue("vegan") {
			case "true":
				targetRestaurant.Vegan = true
			case "false":
				targetRestaurant.Vegan = false
			default:
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			targetRestaurant.Vegan = false
		}

		if !reflect.DeepEqual(targetRestaurant, unchangedRestaurant) {
			err := database.UpdateRestaurant(targetRestaurant)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(res, req, "/admin/restaurants/"+strconv.Itoa(targetRestaurant.ID)+"/profile", http.StatusSeeOther)
		return
	}

	data := struct {
		User       models.User
		Restaurant models.Restaurant
	}{
		myUser,
		targetRestaurant,
	}

	tpl.ExecuteTemplate(res, "admin-restaurants-profile.gohtml", data)
}

// Handles request of "/admin/restaurants/{restaurantID}/orders" page
func adminRestaurantOrders(res http.ResponseWriter, req *http.Request) {
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
	targetID, err := strconv.Atoi(vars["restaurantID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	targetRestaurant, err := database.SelectRestaurant(targetID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	started, err := database.SelectOrdersByRestaurantIDAndStatus(database.DB, targetID, "started")
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	awaitingCollection, err := database.SelectOrdersByRestaurantIDAndStatus(database.DB, targetID, "awaiting collection")
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	completed, err := database.SelectOrdersByRestaurantIDAndStatus(database.DB, targetID, "completed")
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		User               models.User
		Restaurant         models.Restaurant
		Started            []models.Order
		AwaitingCollection []models.Order
		Completed          []models.Order
	}{
		myUser,
		targetRestaurant,
		started,
		awaitingCollection,
		completed,
	}
	tpl.ExecuteTemplate(res, "admin-restaurants-orders.html", data)
}

// Handles request of "/admin/restaurants/{restaurantID}/delete" page
func adminRestaurantDelete(res http.ResponseWriter, req *http.Request) {
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
	switch vars["restaurantID"] {
	case "":
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return

	default:
		targetID, err := strconv.Atoi(vars["restaurantID"])
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
		err = database.DeleteRestaurant(targetID)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(res, req, "/admin/restaurants", http.StatusSeeOther)
}

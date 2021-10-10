package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template

// Pre-Database: var mapUsers = map[string]user{}
// Pre-Database: var mapSessions = map[string]string{}

// Init Function for HTTP Server Functionality. Init templates and admin account.
func HTTPServerInit() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

// Map handlers and start the http server
func StartHTTPServer() {
	HTTPServerInit()
	r := mux.NewRouter() //New Router Instance
	r.HandleFunc("/", index)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.HandleFunc("/setlocation", setlocation)
	r.HandleFunc("/confirmlocation", confirmlocation)
	r.HandleFunc("/profile", profile)
	r.HandleFunc("/changepassword", changepassword)
	// r.HandleFunc("/cart", cart)

	// r.HandleFunc("/restaurants", restaurants) // Restaurant Listing
	// r.HandleFunc("/restaurants/{restaurantID}", restaurantPage) // Individual Restaurant Page, Food Listing
	// r.HandleFunc("/restaurants/{restaurantID}/{foodID}", foodPage) // User will set quantity here
	// r.HandleFunc("/restaurants/{restaurantID}/{foodID}/add", addToCart) // Initial check, start new order? add go over calories?
	// r.HandleFunc("/restaurants/{restaurantID}/{foodID}/addConfirm", addToCartConfirm) // Confirm Add

	r.HandleFunc("/admin", admin)

	r.HandleFunc("/admin/users", adminUsers)
	r.HandleFunc("/admin/users/new", adminUserNew)
	r.HandleFunc("/admin/users/{username}/profile", adminUserProfile)
	r.HandleFunc("/admin/users/{username}/location", adminUserLocation)
	r.HandleFunc("/admin/users/{username}/location/set", adminUserLocationSet)
	r.HandleFunc("/admin/users/{username}/location/confirm", adminUserLocationConfirm)
	r.HandleFunc("/admin/users/{username}/delete", adminUserDelete)

	r.HandleFunc("/admin/restaurants", adminRestaurants)
	r.HandleFunc("/admin/restaurants/new", adminRestaurantNew)
	r.HandleFunc("/admin/restaurants/{restaurantID}/profile", adminRestaurantProfile)
	// r.HandleFunc("/admin/restaurants/{restaurantID}/location", adminRestaurantLocation)
	// r.HandleFunc("/admin/restaurants/{restaurantID}/location/set", adminRestaurantLocationSet)
	// r.HandleFunc("/admin/restaurants/{restaurantID}/location/confirm", adminRestaurantLocationConfirm)
	r.HandleFunc("/admin/restaurants/{restaurantID}/delete", adminRestaurantDelete)

	r.HandleFunc("/admin/restaurants/{restaurantID}", adminRestaurantFoods)
	r.HandleFunc("/admin/restaurants/{restaurantID}/newfood", adminRestaurantFoodNew)
	r.HandleFunc("/admin/restaurants/{restaurantID}/{foodID}", adminRestaurantFoodEdit)
	r.HandleFunc("/admin/restaurants/{restaurantID}/{foodID}/delete", adminRestaurantFoodDelete)

	// Sample Handle Func
	r.HandleFunc("/sample", sample)

	log.Fatal(http.ListenAndServe(":8080", r))
}

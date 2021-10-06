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
	r.HandleFunc("/restricted", restricted)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.HandleFunc("/testmap", testmap)
	r.HandleFunc("/setlocation", setlocation)
	r.HandleFunc("/confirmlocation", confirmlocation)
	r.HandleFunc("/profile", profile)
	r.HandleFunc("/changepassword", changepassword)
	// r.HandleFunc("/cart", cart)

	// r.HandleFunc("/restaurants", restaurants)
	// r.HandleFunc("/restaurants/{restaurantID}", restaurantPage)
	// r.HandleFunc("/restaurants/{restaurantID}/{foodID}", foodPage) // User will set quantity here
	// r.HandleFunc("/restaurants/{restaurantID}/{foodID}/check", addToCartCheck) // Check if adding will go over calories
	// r.HandleFunc("/restaurants/{restaurantID}/{foodID}/add", addToCart) // Add to Cart

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
	// r.HandleFunc("/admin/restaurants/{restaurantID}/foods", adminRestaurantFoods)
	// r.HandleFunc("/admin/restaurants/{restaurantID}/foods/{foodsID}", adminRestaurantFood)
	// r.HandleFunc("/admin/restaurants/{restaurantID}/foods/{foodsID}/delete", adminRestaurantFoodDelete)
	r.HandleFunc("/admin/restaurants/{restaurantID}/delete", adminRestaurantDelete)

	// Sample Handle Func
	r.HandleFunc("/sample", sample)

	log.Fatal(http.ListenAndServe(":8080", r))
}

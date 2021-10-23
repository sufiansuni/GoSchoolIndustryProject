package controllers

import (
	"GoIndustryProject/api"
	"GoIndustryProject/database"
	"GoIndustryProject/models"
	"encoding/json"
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

// Handles request of "/admin/restaurants/{restaurantID}/location" page
func adminRestaurantLocation(res http.ResponseWriter, req *http.Request) {
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

	var mapLink string

	if targetRestaurant.Address != "" {
		mapLink = api.OneMapGenerateMapPNGSingle(fmt.Sprintf("%f", targetRestaurant.Lat), fmt.Sprintf("%f", targetRestaurant.Lng))
	}

	// Prepare data to be sent to template
	// Sample Data can be of any type. Use Arrays or Maps for 'group' data.
	data := struct {
		User       models.User
		Restaurant models.Restaurant
		MapLink    string
	}{
		myUser,
		targetRestaurant,
		mapLink,
	}

	tpl.ExecuteTemplate(res, "admin-restaurants-location.gohtml", data)
}

// Handles request of "/admin/restaurants/{restaurantID}/location/set" page
func adminRestaurantLocationSet(res http.ResponseWriter, req *http.Request) {
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

	var searchResults api.OneMapSearchResult
	var newResults []map[string]string
	var locationQuery string

	if req.Method == http.MethodPost {
		locationQuery = req.FormValue("locationQuery")
		var err error
		searchResults, err = api.OneMapSearch(locationQuery)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		for _, v := range searchResults.Results {
			var newResultItem map[string]string
			jV, err := json.Marshal(v)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(jV, &newResultItem)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			var mapLinkStruct struct {
				MapLink string
			}
			mapLinkStruct.MapLink = api.OneMapGenerateMapPNGSingle(v.Latitude, v.Longitude)

			jMapLinkStruct, err := json.Marshal(mapLinkStruct)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(jMapLinkStruct, &newResultItem)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			newResults = append(newResults, newResultItem)

		}
	}

	// Prepare data to be sent to template
	data := struct {
		User          models.User
		Restaurant    models.Restaurant
		Locations     []map[string]string
		LocationQuery string
	}{
		myUser,
		targetRestaurant,
		newResults,
		locationQuery,
	}
	tpl.ExecuteTemplate(res, "admin-restaurants-location-set.gohtml", data)
}

// Handles request of "/admin/users/{username}/location/confirm" page
func adminRestaurantLocationConfirm(res http.ResponseWriter, req *http.Request) {
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
		locationQuery := req.FormValue("locationQuery")
		locationNumberString := req.FormValue("locationNumber")
		locationNumber, err := strconv.Atoi(locationNumberString)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		searchResult, err := api.OneMapSearch(locationQuery)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		targetRestaurant.Address = searchResult.Results[locationNumber].Address

		targetRestaurant.Unit = req.FormValue("unit")

		if searchResult.Results[locationNumber].Latitude != "NIL" {
			targetRestaurant.Lat, err = strconv.ParseFloat(searchResult.Results[locationNumber].Latitude, 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if searchResult.Results[locationNumber].Longitude != "NIL" {
			targetRestaurant.Lng, err = strconv.ParseFloat(searchResult.Results[locationNumber].Longitude, 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		err = database.UpdateRestaurant(targetRestaurant)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(res, req, "/admin/restaurants/"+vars["restaurantID"]+"/location", http.StatusSeeOther)
}

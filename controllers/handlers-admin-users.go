package controllers

import (
	"GoIndustryProject/api"
	"GoIndustryProject/database"
	"GoIndustryProject/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Handles request of "/admin/users" page
func adminUsers(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}

	myUsers, err := database.SelectAllUsers(database.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		User  models.User
		Users []models.User
	}{
		myUser,
		myUsers,
	}
	tpl.ExecuteTemplate(res, "admin-users.html", data)
}

// Handles request of "/admin/users/new" page
func adminUserNew(res http.ResponseWriter, req *http.Request) {
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
		username := req.FormValue("username")
		password := req.FormValue("password")
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")
		if username != "" {
			//check if client tried to create "admin"
			if username == "admin" {
				http.Error(res, "Forbidden", http.StatusForbidden)
				return
			}
			// check if username exist/ taken
			_, err := database.SelectUserByUsername(database.DB, username)
			if err != nil {
				if err != sql.ErrNoRows {
					fmt.Println(err)
					http.Error(res, "Internal server error", http.StatusInternalServerError)
					return
				} else {
					fmt.Println("User '", username, "' not found. ", err.Error())
				}
			} else {
				fmt.Println(err)
				http.Error(res, "Username already taken", http.StatusForbidden)
				return
			}

			//encrypt password
			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			myUser = models.User{
				Username: username,
				Password: bPassword,
				First:    firstname,
				Last:     lastname,
			}
			myUser.AdjustStrings()

			err = database.InsertUser(database.DB, myUser) // previouslymapUsers[username] = myUser
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			} else {
				fmt.Println("User Created:", username)
			}

		}

		// redirect to admin page (users)
		http.Redirect(res, req, "/admin/users", http.StatusSeeOther)
		return

	}

	data := struct {
		User models.User
	}{
		myUser,
	}

	tpl.ExecuteTemplate(res, "signup.html", data)
}

// Handles request of "/admin/users/{username}/profile" page
func adminUserProfile(res http.ResponseWriter, req *http.Request) {
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
	targetUser, err := database.SelectUserByUsername(database.DB, vars["username"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodPost {
		unchangedUser := targetUser

		if req.FormValue("firstName") != "" {
			targetUser.First = req.FormValue("firstName")
		}

		if req.FormValue("lastName") != "" {
			targetUser.First = req.FormValue("lastName")
		}

		if req.FormValue("gender") != "" {
			targetUser.Gender = req.FormValue("gender")
		}

		if req.FormValue("birthday") != "" {
			targetUser.Birthday = req.FormValue("birthday")
		}

		if req.FormValue("height") != "" {
			var err error
			targetUser.Height, err = strconv.Atoi(req.FormValue("height"))
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if req.FormValue("weight") != "" {
			var err error
			targetUser.Weight, err = strconv.ParseFloat(req.FormValue("weight"), 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if req.FormValue("activityLevel") != "" {
			var err error
			targetUser.ActivityLevel, err = strconv.Atoi(req.FormValue("activityLevel"))
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if req.FormValue("caloriesPerDay") != "" {
			var err error
			targetUser.CaloriesPerDay, err = strconv.Atoi(req.FormValue("caloriesPerDay"))
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if req.FormValue("halal") != "" {
			switch req.FormValue("halal") {
			case "true":
				targetUser.Halal = true
			case "false":
				targetUser.Halal = false
			default:
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			targetUser.Halal = false
		}

		if req.FormValue("vegan") != "" {
			switch req.FormValue("vegan") {
			case "true":
				targetUser.Vegan = true
			case "false":
				targetUser.Vegan = false
			default:
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			targetUser.Vegan = false
		}

		if !reflect.DeepEqual(targetUser, unchangedUser) {
			err := database.UpdateUserProfile(database.DB, targetUser)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(res, req, "/admin/users/"+targetUser.Username+"/profile", http.StatusSeeOther)
	}

	if req.Method == http.MethodGet {
		recommendedCaloriesPerDay, err := userRecommendedCaloriesPerDay(targetUser)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
		data := struct {
			User                      models.User
			Target                    models.User
			RecommendedCaloriesPerDay int
		}{
			myUser,
			targetUser,
			recommendedCaloriesPerDay,
		}

		tpl.ExecuteTemplate(res, "admin-users-profile.html", data)
	}
}

// Handles request of "/admin/users/{username}/location" page
func adminUserLocation(res http.ResponseWriter, req *http.Request) {
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
	targetUser, err := database.SelectUserByUsername(database.DB, vars["username"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	var mapLink string

	if targetUser.Address != "" {
		mapLink = api.OneMapGenerateMapPNGSingle(fmt.Sprintf("%f", targetUser.Lat), fmt.Sprintf("%f", targetUser.Lng))
	}

	// Prepare data to be sent to template
	// Sample Data can be of any type. Use Arrays or Maps for 'group' data.
	data := struct {
		User    models.User
		Target  models.User
		MapLink string
	}{
		myUser,
		targetUser,
		mapLink,
	}

	tpl.ExecuteTemplate(res, "admin-users-location.html", data)
}

// Handles request of "/admin/users/{username}/location/set" page
func adminUserLocationSet(res http.ResponseWriter, req *http.Request) {
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
	targetUser, err := database.SelectUserByUsername(database.DB, vars["username"])
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
		Target        models.User
		Locations     []map[string]string
		LocationQuery string
	}{
		myUser,
		targetUser,
		newResults,
		locationQuery,
	}
	tpl.ExecuteTemplate(res, "admin-users-location-set.html", data)
}

// Handles request of "/admin/users/{username}/location/confirm" page
func adminUserLocationConfirm(res http.ResponseWriter, req *http.Request) {
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
	targetUser, err := database.SelectUserByUsername(database.DB, vars["username"])
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

		targetUser.Address = searchResult.Results[locationNumber].Address

		targetUser.Unit = req.FormValue("unit")

		if searchResult.Results[locationNumber].Latitude != "NIL" {
			targetUser.Lat, err = strconv.ParseFloat(searchResult.Results[locationNumber].Latitude, 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if searchResult.Results[locationNumber].Longitude != "NIL" {
			targetUser.Lng, err = strconv.ParseFloat(searchResult.Results[locationNumber].Longitude, 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		err = database.UpdateUserProfile(database.DB, targetUser)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(res, req, "/admin/users/"+targetUser.Username+"/location", http.StatusSeeOther)
}

// Handles request of "/admin/users/{username}/orders" page
func adminUserOrders(res http.ResponseWriter, req *http.Request) {
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
	targetUser, err := database.SelectUserByUsername(database.DB, vars["username"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	searchCart, err := database.SelectOrdersByUsernameAndStatus(database.DB, targetUser.Username, "started")
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}
	cart := searchCart[0]

	awaitingCollection, err := database.SelectOrdersByUsernameAndStatus(database.DB, targetUser.Username, "awaiting collection")
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	completed, err := database.SelectOrdersByUsernameAndStatus(database.DB, targetUser.Username, "completed")
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		User               models.User
		Target             models.User
		Cart               models.Order
		AwaitingCollection []models.Order
		Completed          []models.Order
	}{
		myUser,
		targetUser,
		cart,
		awaitingCollection,
		completed,
	}
	tpl.ExecuteTemplate(res, "admin-users-orders.html", data)
}

// Handles request of "/admin/users/{username}/delete" page
func adminUserDelete(res http.ResponseWriter, req *http.Request) {
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
	switch vars["username"] {
	case "admin":
		http.Redirect(res, req, "/admin/users", http.StatusUnauthorized)
		return

	case "":
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return

	default:
		err := database.DeleteUser(database.DB, vars["username"])
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(res, req, "/admin/users", http.StatusSeeOther)
}

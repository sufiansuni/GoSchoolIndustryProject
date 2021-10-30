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

// Handles request of "/admin/restaurants/{restaurantID}/foods" page
func adminRestaurantFoods(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}

	var targetRestaurant models.Restaurant
	var targetFoods []models.Food

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
		targetRestaurant, err = database.SelectRestaurant(targetID)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
		targetFoods, err = database.SelectAllFoodsByRestaurantID(database.DB, targetID)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	if req.Method == http.MethodGet {
		data := struct {
			User       models.User
			Restaurant models.Restaurant
			Foods      []models.Food
		}{
			myUser,
			targetRestaurant,
			targetFoods,
		}

		tpl.ExecuteTemplate(res, "admin-restaurants-foods.html", data)
	}
}

// Handles request of "/admin/restaurants/{restaurantID}/foods/new" page
func adminRestaurantFoodNew(res http.ResponseWriter, req *http.Request) {
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
	var targetRestaurant models.Restaurant
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
		targetRestaurant, err = database.SelectRestaurant(targetID)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		name := req.FormValue("name")
		description := req.FormValue("description")
		price := req.FormValue("price")
		calories := req.FormValue("calories")

		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		caloriesInt, err := strconv.Atoi(calories)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		if name != "" {
			myFood := models.Food{
				RestaurantID: targetRestaurant.ID,
				Name:         name,
				Description:  description,
				Price:        priceFloat,
				Calories:     caloriesInt,
			}

			err := database.InsertFood(myFood) // previouslymapUsers[username] = myUser
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			} else {
				fmt.Println("Food Created:", name)
			}
		}

		// redirect to admin page (restaurants)
		http.Redirect(res, req, "/admin/restaurants/"+vars["restaurantID"]+"/foods", http.StatusSeeOther)
		return
	}

	data := struct {
		User       models.User
		Restaurant models.Restaurant
	}{
		myUser,
		targetRestaurant,
	}

	tpl.ExecuteTemplate(res, "newrestaurantfood.html", data)
}

// Handles request of "/admin/restaurants/{restaurantID}/foods/{foodID}" page
func adminRestaurantFoodEdit(res http.ResponseWriter, req *http.Request) {
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
	targetID, err := strconv.Atoi(vars["foodID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	targetFood, err := database.SelectFood(targetID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodPost {
		unchangedFood := targetFood
		name := req.FormValue("name")
		description := req.FormValue("description")
		price := req.FormValue("price")
		calories := req.FormValue("calories")

		if req.FormValue("name") != "" {
			targetFood.Name = name
		}

		if req.FormValue("description") != "" {
			targetFood.Description = description
		}

		if req.FormValue("price") != "" {
			priceFloat, err := strconv.ParseFloat(price, 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
			targetFood.Price = priceFloat
		}

		if req.FormValue("calories") != "" {
			caloriesInt, err := strconv.Atoi(calories)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
			targetFood.Calories = caloriesInt
		}

		if !reflect.DeepEqual(targetFood, unchangedFood) {
			err := database.UpdateFood(targetFood)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(res, req, "/admin/restaurants/"+vars["restaurantID"]+"/foods/"+vars["foodID"], http.StatusSeeOther)
		return
	}

	data := struct {
		User models.User
		Food models.Food
	}{
		myUser,
		targetFood,
	}

	tpl.ExecuteTemplate(res, "admin-restaurants-foodedit.gohtml", data)
}

// Handles request of "/admin/restaurants/{restaurantID}/foods/{foodID}/delete" page
func adminRestaurantFoodDelete(res http.ResponseWriter, req *http.Request) {
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
	switch vars["foodID"] {
	case "":
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return

	default:
		targetID, err := strconv.Atoi(vars["foodID"])
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
		err = database.DeleteFood(targetID)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(res, req, "/admin/restaurants/"+vars["restaurantID"], http.StatusSeeOther)
}

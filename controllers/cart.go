package controllers

import (
	"GoIndustryProject/api"
	"GoIndustryProject/database"
	"GoIndustryProject/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Handles request of "/cart" page
func userCart(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	var myOrders []models.Order
	var myOrderItems []models.OrderItem
	myOrders, err := database.SelectOrdersByUsernameAndStatus(database.DB, myUser.Username, "Started")
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}
	var myOrder models.Order
	if len(myOrders) != 0 {
		myOrder = myOrders[0]
		myOrderItems, err = database.SelectOrderItemsByOrderID(database.DB, myOrders[0].ID)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	userLat := fmt.Sprintf("%f", myOrder.UserLat)
	userLng := fmt.Sprintf("%f", myOrder.UserLng)
	userMapLink := api.OneMapGenerateMapPNGSingle(userLat, userLng)

	restaurantLat := fmt.Sprintf("%f", myOrder.RestaurantLat)
	restaurantLng := fmt.Sprintf("%f", myOrder.RestaurantLng)
	restaurantMapLink := api.OneMapGenerateMapPNGSingle(restaurantLat, restaurantLng)

	var todayCalories int

	todayOrders, err := database.SelectOrdersByUsernameDateStatus(database.DB, myUser.Username, time.Now().Format("2006-01-02"), "Completed")
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	for _, order := range todayOrders {
		todayCalories += order.TotalCalories
	}

	caloriesToTarget := myUser.CaloriesPerDay - todayCalories

	var routeMapLink string
	var walkingDistance int
	var averageWalkingCalorieBurn int

	if myOrder.UserAddress != "" && myOrder.RestaurantAddress != "" {
		routeMapLink, walkingDistance = api.OneMapGenerateMapPNGTwoPoints(userLat, userLng, restaurantLat, restaurantLng)
		averageWalkingCalorieBurn = int(float64(walkingDistance) / 1000 * 63)
	}

	data := struct {
		User                      models.User
		Order                     models.Order
		OrderItems                []models.OrderItem
		UserMapLink               string
		RestaurantMapLink         string
		TodayCalories             int
		CaloriesToTarget          int
		WalkingDistance           int
		RouteMapLink              string
		AverageWalkingCalorieBurn int
	}{
		myUser,
		myOrder,
		myOrderItems,
		userMapLink,
		restaurantMapLink,
		todayCalories,
		caloriesToTarget,
		walkingDistance,
		routeMapLink,
		averageWalkingCalorieBurn,
	}
	tpl.ExecuteTemplate(res, "cart.gohtml", data)
}

// Handles request of "/addtocart/{foodID}" page
func addToCart(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		quantityString := req.FormValue("quantity")
		quantity, err := strconv.Atoi(quantityString)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		var targetFood models.Food
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
			targetFood, err = database.SelectFood(targetID)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		targetRestaurant, err := database.SelectRestaurant(targetFood.RestaurantID)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		userCarts, err := database.SelectOrdersByUsernameAndStatus(database.DB, myUser.Username, "Started")
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		if myUser.Address == "" || targetRestaurant.Address == "" {
			fmt.Println("Either User or Restaurant Address Not Set")
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		var targetOrder models.Order

		switch len(userCarts) {
		case 0:
			newOrder := models.Order{
				Username:       myUser.Username,
				RestaurantID:   targetRestaurant.ID,
				RestaurantName: targetRestaurant.Name,

				Status:     "Started",
				Collection: "Delivery",
				Date:       time.Now().Format("2006-01-02"),

				UserAddress: myUser.Address,
				UserUnit:    myUser.Unit,
				UserLat:     myUser.Lat,
				UserLng:     myUser.Lng,

				RestaurantAddress: targetRestaurant.Address,
				RestaurantUnit:    targetRestaurant.Unit,
				RestaurantLat:     targetRestaurant.Lat,
				RestaurantLng:     targetRestaurant.Lng,

				TotalItems: quantity,
				TotalPrice: targetFood.Price * float64(quantity),

				TotalCalories: quantity * int(targetFood.Calories),
				BurnCalories:  0,
			}
			database.InsertOrder(database.DB, newOrder)
			targetOrders, err := database.SelectOrdersByUsernameAndStatus(database.DB, myUser.Username, "Started")
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
			targetOrder = targetOrders[0]
		case 1:
			if userCarts[0].RestaurantID != targetFood.RestaurantID {
				http.Error(res, "Your current order is for a different restaurant", http.StatusInternalServerError)
				return
			}
			targetOrder = userCarts[0]
			targetOrder.TotalItems += quantity
			targetOrder.TotalPrice += targetFood.Price * float64(quantity)
			targetOrder.TotalCalories += quantity * int(targetFood.Calories)
			database.UpdateOrder(database.DB, targetOrder)
		default:
			fmt.Println("User has too many orders started")
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		var targetOrderItem models.OrderItem
		targetOrderItems, err := database.SelectOrderItemsByOrderIDAndFoodID(database.DB, targetOrder.ID, targetFood.ID)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
		if len(targetOrderItems) == 0 {
			newOrderItem := models.OrderItem{
				OrderID:          targetOrder.ID,
				FoodID:           targetFood.ID,
				FoodName:         targetFood.Name,
				Quantity:         quantity,
				SubtotalPrice:    targetFood.Price * float64(quantity),
				SubtotalCalories: targetFood.Calories * quantity,
			}
			err = database.InsertOrderItem(database.DB, newOrderItem)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
			// targetOrderItems, err = database.SelectOrderItemsByOrderIDAndFoodID(database.DB, targetOrder.ID, targetFood.ID)
			// if err != nil {
			// 	fmt.Println(err)
			// 	http.Error(res, "Internal server error", http.StatusInternalServerError)
			// 	return
			// }
			// targetOrderItem = targetOrderItems[0]
		} else {
			targetOrderItem = targetOrderItems[0]
			targetOrderItem.Quantity += quantity
			targetOrderItem.SubtotalPrice += targetFood.Price * float64(quantity)
			targetOrderItem.SubtotalCalories += targetFood.Calories * quantity
			err = database.UpdateOrderItem(database.DB, targetOrderItem)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	}
	http.Redirect(res, req, "/cart", http.StatusSeeOther)
}

// Handles request of "/cart/{itemID}/add" page
func cartItemAdd(res http.ResponseWriter, req *http.Request) {
	// myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
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

	http.Redirect(res, req, "/cart", http.StatusSeeOther)
}

// Handles request of "/cart/{itemID}/subtract" page
func cartItemSubtract(res http.ResponseWriter, req *http.Request) {
	// myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
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

	http.Redirect(res, req, "/cart", http.StatusSeeOther)
}

// Handles request of "/cart/{itemID}/delete" page
func cartItemDelete(res http.ResponseWriter, req *http.Request) {
	// myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
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

	http.Redirect(res, req, "/cart", http.StatusSeeOther)
}

// Handles request of "/cart/delivery/{orderID}" page
func cartDelivery(res http.ResponseWriter, req *http.Request) {
	// myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(req)
	orderID, err := strconv.Atoi(vars["orderID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrder, err := database.SelectOrderByID(database.DB, orderID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if myOrder.Collection != "Delivery" {
		myOrder.Collection = "Delivery"
		err = database.UpdateOrder(database.DB, myOrder)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(res, req, "/cart", http.StatusSeeOther)
}

// Handles request of "/cart/delivery/{orderID}" page
func cartSelfCollect(res http.ResponseWriter, req *http.Request) {
	// myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(req)
	orderID, err := strconv.Atoi(vars["orderID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrder, err := database.SelectOrderByID(database.DB, orderID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if myOrder.Collection != "Self-Collect" {
		myOrder.Collection = "Self-Collect"
		err = database.UpdateOrder(database.DB, myOrder)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(res, req, "/cart", http.StatusSeeOther)
}

// Handles request of "/cart/confirm/{orderID}" page
func cartConfirm(res http.ResponseWriter, req *http.Request) {
	// myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(req)
	orderID, err := strconv.Atoi(vars["orderID"])
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	myOrder, err := database.SelectOrderByID(database.DB, orderID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	if myOrder.Status != "Awaiting Collection" {
		myOrder.Status = "Awaiting Collection"
		err = database.UpdateOrder(database.DB, myOrder)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(res, req, "/cart", http.StatusSeeOther)
}

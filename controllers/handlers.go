package controllers

import (
	"GoIndustryProject/api"
	"GoIndustryProject/database"
	"GoIndustryProject/health"
	"GoIndustryProject/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Sample of handler
func sample(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// Do whatever you need to here, these lines are just examples
	var myArray []string
	myArray = append(myArray, "string1", "string2")

	myMap := make(map[string]int)
	myMap["Age"] = 100
	myMap["Birth Year"] = 1911

	// Prepare data to be sent to template
	// Sample Data can be of any type. Use Arrays or Maps for 'group' data.
	data := struct {
		User        models.User
		SampleData  string
		SampleArray []string
		SampleMap   map[string]int
	}{
		myUser,
		"A Sample Data",
		myArray,
		myMap,
	}
	tpl.ExecuteTemplate(res, "sample.html", data)
}

// Handles request of index/homepage
func index(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	data := struct {
		User models.User
	}{
		myUser,
	}
	tpl.ExecuteTemplate(res, "index.html", data)
}

// Handles request of sign-up page. Also login the user on success.
func signup(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	var myUser models.User

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

			// create session
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:  "myCookie",
				Value: id.String(),
			}

			http.SetCookie(res, myCookie)

			mySession := models.Session{
				UUID:     myCookie.Value,
				Username: myUser.Username,
			}

			err = database.InsertSession(database.DB, mySession) // previously: mapSessions[myCookie.Value] = username
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			} else {
				fmt.Println("Session Created")
			}

		}

		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	data := struct {
		User models.User
	}{
		myUser,
	}
	tpl.ExecuteTemplate(res, "signup.html", data)
}

// Handles request of login page. Login user on successful POST.
func login(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		// check if user exist with username
		checker, err := database.SelectUserByUsername(database.DB, username)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}

		// Matching of password entered
		err = bcrypt.CompareHashAndPassword(checker.Password, []byte(password))
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}

		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
		http.SetCookie(res, myCookie)

		mySession := models.Session{
			UUID:     myCookie.Value,
			Username: username,
		}

		err = database.InsertSession(database.DB, mySession) // previously: mapSessions[myCookie.Value] = username
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			fmt.Println("Session Created")
		}

		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// Execute Template when Method not POST.
	tpl.ExecuteTemplate(res, "login.html", nil)
}

// Handles request of logout page
func logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session

	err := database.DeleteSession(database.DB, myCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// Locates user's cookie and check against session data. Creates cookie if not present.
// If user is found, returns the user data.
func checkUser(res http.ResponseWriter, req *http.Request) (myUser models.User) {
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
	}
	http.SetCookie(res, myCookie)

	// if the user exists already, get user

	mySession, err := database.SelectSession(database.DB, myCookie.Value)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
		} else {
			fmt.Println("No Entry Found in Database for UUID:" + myCookie.Value)
		}
	} else {
		myUser, err = database.SelectUserByUsername(database.DB, mySession.Username)
		if err != nil {
			if err != sql.ErrNoRows {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
			} else {
				fmt.Println("No Entry Found in Database for User:" + mySession.Username)
			}
		}
	}
	return myUser
}

// Locates user's cookie and check against session data.
// Returns true if user found(logged in), else return false.
// Function DOES NOT issue cookie if not found.
func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}

	mySession, err := database.SelectSession(database.DB, myCookie.Value)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Print(err)
		} else {
			fmt.Println("No Entry Found in Database for UUID:" + myCookie.Value)
		}
	} else {
		_, err = database.SelectUserByUsername(database.DB, mySession.Username)
		if err != nil {
			fmt.Print(err)
		} else {
			return true
		}
	}
	return false
}

// Handles request of admin page
func admin(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusUnauthorized)
		return
	}

	data := struct {
		User models.User
	}{
		myUser,
	}
	tpl.ExecuteTemplate(res, "admin.html", data)
}

// Handles request of setlocation page
func setlocation(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// Do whatever you need to here

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
		Locations     []map[string]string
		LocationQuery string
	}{
		myUser,
		newResults,
		locationQuery,
	}
	tpl.ExecuteTemplate(res, "setlocation.html", data)
}

// Handles request of confirmlocation page
func confirmlocation(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
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

		myUser.Address = searchResult.Results[locationNumber].Address

		myUser.Unit = req.FormValue("unit")

		if searchResult.Results[locationNumber].Latitude != "NIL" {
			myUser.Lat, err = strconv.ParseFloat(searchResult.Results[locationNumber].Latitude, 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if searchResult.Results[locationNumber].Longitude != "NIL" {
			myUser.Lng, err = strconv.ParseFloat(searchResult.Results[locationNumber].Longitude, 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		err = database.UpdateUserProfile(database.DB, myUser)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// Handles request of profile page
func profile(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		unchangedUser := myUser

		if req.FormValue("firstName") != "" {
			myUser.First = req.FormValue("firstName")
		}

		if req.FormValue("lastName") != "" {
			myUser.First = req.FormValue("lastName")
		}

		if req.FormValue("gender") != "" {
			myUser.Gender = req.FormValue("gender")
		}

		if req.FormValue("birthday") != "" {
			myUser.Birthday = req.FormValue("birthday")
		}

		if req.FormValue("height") != "" {
			var err error
			myUser.Height, err = strconv.Atoi(req.FormValue("height"))
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if req.FormValue("weight") != "" {
			var err error
			myUser.Weight, err = strconv.ParseFloat(req.FormValue("weight"), 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if req.FormValue("activityLevel") != "" {
			var err error
			myUser.ActivityLevel, err = strconv.Atoi(req.FormValue("activityLevel"))
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if req.FormValue("caloriesPerDay") != "" {
			var err error
			myUser.CaloriesPerDay, err = strconv.Atoi(req.FormValue("caloriesPerDay"))
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if req.FormValue("halal") != "" {
			switch req.FormValue("halal") {
			case "true":
				myUser.Halal = true
			case "false":
				myUser.Halal = false
			default:
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			myUser.Halal = false
		}

		if req.FormValue("vegan") != "" {
			switch req.FormValue("vegan") {
			case "true":
				myUser.Vegan = true
			case "false":
				myUser.Vegan = false
			default:
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			myUser.Vegan = false
		}

		if !reflect.DeepEqual(myUser, unchangedUser) {
			err := database.UpdateUserProfile(database.DB, myUser)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(res, req, "/profile", http.StatusSeeOther)
	}

	if req.Method == http.MethodGet {
		recommendedCaloriesPerDay, err := userRecommendedCaloriesPerDay(myUser)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Prepare data to be sent to template
		data := struct {
			User                      models.User
			RecommendedCaloriesPerDay int
		}{
			myUser,
			recommendedCaloriesPerDay,
		}
		tpl.ExecuteTemplate(res, "profile.html", data)
	}
}

// Helper function to calculate a user's recommended daily calories
// Returns 0 if required parameters are missing
func userRecommendedCaloriesPerDay(myUser models.User) (int, error) {
	if (myUser.Birthday != "") && (myUser.Gender != "") && (myUser.Height != 0) && (myUser.Weight != 0) && (myUser.ActivityLevel != 0) {
		userBirthday, err := time.Parse("2006-01-02", myUser.Birthday)
		if err != nil {
			return 0, err
		}
		userAge := time.Now().Year() - userBirthday.Year()
		userBMR := health.CalculateBMR(myUser.Gender, myUser.Height, myUser.Weight, userAge)
		return health.CalculateDailyCalories(userBMR, myUser.ActivityLevel), nil
	}
	return 0, nil
}

// Handles request of changepassword page
func changepassword(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		currentPassword := req.FormValue("currentPassword")
		newPassword := req.FormValue("newPassword")

		// Matching of current password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(currentPassword))
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}

		bNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		err = database.UpdateUserPassword(database.DB, myUser.Username, bNewPassword)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(res, req, "/changepassword", http.StatusSeeOther)
	}

	if req.Method == http.MethodGet {
		// Prepare data to be sent to template
		data := struct {
			User models.User
		}{
			myUser,
		}
		tpl.ExecuteTemplate(res, "changepassword.html", data)
	}
}

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

	data := struct {
		User       models.User
		Order      models.Order
		OrderItems []models.OrderItem
	}{
		myUser,
		myOrder,
		myOrderItems,
	}
	tpl.ExecuteTemplate(res, "cart.html", data)
}

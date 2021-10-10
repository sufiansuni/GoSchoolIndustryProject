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

	if req.Method == http.MethodPost {
		locationQuery := req.FormValue("locationQuery")
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
		User      models.User
		Locations []map[string]string
	}{
		myUser,
		newResults,
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
		currentLocation := req.FormValue("currentLocation")
		searchResult, err := api.OneMapSearch(currentLocation)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		myUser.Address = searchResult.Results[0].Address

		myUser.Unit = req.FormValue("unit")

		if searchResult.Results[0].Latitude != "NIL" {
			myUser.Lat, err = strconv.ParseFloat(searchResult.Results[0].Latitude, 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if searchResult.Results[0].Longitude != "NIL" {
			myUser.Lng, err = strconv.ParseFloat(searchResult.Results[0].Longitude, 64)
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
		// Prepare data to be sent to template
		data := struct {
			User models.User
		}{
			myUser,
		}
		tpl.ExecuteTemplate(res, "profile.html", data)
	}
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
		data := struct {
			User   models.User
			Target models.User
		}{
			myUser,
			targetUser,
		}

		tpl.ExecuteTemplate(res, "admin-users-profile.html", data)
	}
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
		searchResults, err := api.OneMapSearch(targetUser.Address)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		mapLink = api.OneMapGenerateMapPNGSingle(searchResults.Results[0].Latitude, searchResults.Results[0].Longitude)
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

	if req.Method == http.MethodPost {
		locationQuery := req.FormValue("locationQuery")
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
		User      models.User
		Target    models.User
		Locations []map[string]string
	}{
		myUser,
		targetUser,
		newResults,
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
		currentLocation := req.FormValue("currentLocation")
		searchResult, err := api.OneMapSearch(currentLocation)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		targetUser.Address = searchResult.Results[0].Address

		targetUser.Unit = req.FormValue("unit")

		if searchResult.Results[0].Latitude != "NIL" {
			targetUser.Lat, err = strconv.ParseFloat(searchResult.Results[0].Latitude, 64)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if searchResult.Results[0].Longitude != "NIL" {
			targetUser.Lng, err = strconv.ParseFloat(searchResult.Results[0].Longitude, 64)
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

// Handles request of "/admin/restaurants/{restaurantID}" page
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

// Handles request of "/admin/restaurants/{restaurantID}/newfood" page
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
		price := req.FormValue("price")
		calories := req.FormValue("calories")

		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		caloriesFloat, err := strconv.ParseFloat(calories, 64)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		if name != "" {
			myFood := models.Food{
				RestaurantID: targetRestaurant.ID,
				Name:         name,
				Price:        priceFloat,
				Calories:     caloriesFloat,
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
		http.Redirect(res, req, "/admin/restaurants/"+vars["restaurantID"], http.StatusSeeOther)
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

// Handles request of "/admin/restaurants/{restaurantID}/{foodID}" page
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
		price := req.FormValue("price")
		calories := req.FormValue("calories")

		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		caloriesFloat, err := strconv.ParseFloat(calories, 64)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}

		if req.FormValue("name") != "" {
			targetFood.Name = name
		}

		if req.FormValue("price") != "" {
			targetFood.Price = priceFloat
		}

		if req.FormValue("calories") != "" {
			targetFood.Calories = caloriesFloat
		}

		if !reflect.DeepEqual(targetFood, unchangedFood) {
			err := database.UpdateFood(targetFood)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(res, req, "/admin/restaurants/"+strconv.Itoa(targetFood.RestaurantID)+"/"+strconv.Itoa(targetFood.ID), http.StatusSeeOther)
		return
	}

	data := struct {
		User models.User
		Food models.Food
	}{
		myUser,
		targetFood,
	}

	tpl.ExecuteTemplate(res, "admin-restaurants-foodedit.html", data)
}

// Handles request of "/admin/restaurants/{restaurantID}/{foodID}/delete" page
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

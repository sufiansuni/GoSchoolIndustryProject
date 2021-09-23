package controllers

import (
	"GoIndustryProject/apis"
	"GoIndustryProject/database"
	"GoIndustryProject/models"
	"database/sql"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

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

// Handles request of restricted page
func restricted(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	data := struct {
		User models.User
	}{
		myUser,
	}
	tpl.ExecuteTemplate(res, "restricted.html", data)
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

			err = database.InsertUser(database.DB, myUser) // previouslymapUsers[username] = myUser
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			} else {
				fmt.Println("User Created:", username)
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

func testmap(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	search1, err := apis.OneMapSearch("13 Marsiling Lane")
	if err != nil {
		fmt.Println(err)
	}
	search2, err := apis.OneMapSearch("Woodlands MRT NS9")
	if err != nil {
		fmt.Println(err)
	}

	start_lat := search1.Results[0].Latitude
	start_lng := search1.Results[0].Longitude

	end_lat := search2.Results[0].Latitude
	end_lng := search2.Results[0].Longitude

	MapPNG := apis.OneMapGenerateMapPNG(start_lat, start_lng, end_lat, end_lng)
	data := struct{
		User	models.User
		MapPNG string
	}{
		myUser,
		MapPNG,
	}
	tpl.ExecuteTemplate(res, "testmap.html", data)
}

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

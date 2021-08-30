package main

// Basic HTTP Server Framework with the following functionality:
// Login/Logout
// Signup
// "Restricted Page" is a test page, to be removed

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username       string
	Password       []byte
	First          string
	Last           string
	Gender         string
	Birthday       string
	Height         int
	Weight         float64
	CaloriesPerDay int
	Halal          bool
	Vegan          bool
	Address        string
	PostalCode     int
	Lat            float64
	Lng            float64
}

type session struct {
	UUID     string //primary key
	Username string //foreign key
}

var tpl *template.Template

//Pre-Database: var mapUsers = map[string]user{}
//Pre-Database: var mapSessions = map[string]string{}

//Function creates initial admin account. If account already exist, error will be printed.
func createAdminAccount() {
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	myUser := user{
		Username: "admin",
		Password: bPassword,
		First:    "first",
		Last:     "last",
	}
	err := insertUser(myUser) //previously mapUsers["admin"] = myUser
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Admin Account Created")
	}
}

//Init Function for HTTP Server Functionality. Init templates and admin account.
func HTTPServerInit() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	createAdminAccount() // Create Admin Account, Previously: mapUsers["admin"] = user{"admin", bPassword, "admin", "admin"}
}

//Function will map handlers and start the http server
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

	log.Fatal(http.ListenAndServe(":8080", r))
}

//Function handles request of index/homepage
func index(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	tpl.ExecuteTemplate(res, "index.html", myUser)
}

//Function handles request of restricted page
func restricted(res http.ResponseWriter, req *http.Request) {
	myUser := checkUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "restricted.html", myUser)
}

//Function handles request of sign-up page. Also login the user on success.
func signup(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	var myUser user
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		username := req.FormValue("username")
		password := req.FormValue("password")
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")
		if username != "" {
			// check if username exist/ taken
			var checker string

			query := "SELECT Username FROM users WHERE Username=?"

			err := db.QueryRow(query, username).Scan(&checker)
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

			mySession := session{myCookie.Value, username}

			err = insertSession(mySession) // previously: mapSessions[myCookie.Value] = username
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			} else {
				fmt.Println("Session Created")
			}

			//check password
			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			myUser = user{
				Username: username,
				Password: bPassword,
				First:    firstname,
				Last:     lastname,
			}

			err = insertUser(myUser) // previouslymapUsers[username] = myUser
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
	tpl.ExecuteTemplate(res, "signup.html", myUser)
}

//Function handles request of login page. Login user on successful POST
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
		var checker user

		query := "SELECT Username, Password FROM users WHERE Username=?"
		err := db.QueryRow(query, username).Scan(
			&checker.Username,
			&checker.Password,
		)
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

		mySession := session{myCookie.Value, username}
		err = insertSession(mySession) // previously: mapSessions[myCookie.Value] = username
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

	//Execute Template when Method not POST.
	tpl.ExecuteTemplate(res, "login.html", nil)
}

//Function handles request of logout page
func logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session

	err := deleteSession(myCookie.Value)
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

//function locates user's cookie and check against session data. Creates cookie if not present.
//if user is found, returns the user data
func checkUser(res http.ResponseWriter, req *http.Request) user {
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
	var checker string
	var myUser user

	query := "SELECT Username FROM sessions WHERE UUID=?"
	err = db.QueryRow(query, myCookie.Value).Scan(&checker)

	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
		} else {
			fmt.Println("No Entry Found in Database for UUID:" + myCookie.Value)
		}
	} else {
		query = "SELECT * FROM users WHERE Username=?"

		err = db.QueryRow(query, checker).Scan(
			&myUser.Username,
			&myUser.Password,
			&myUser.First,
			&myUser.Last,
			&myUser.Birthday,
			&myUser.Height,
			&myUser.Weight,
			&myUser.CaloriesPerDay,
			&myUser.Halal,
			&myUser.Vegan,
			&myUser.Address,
			&myUser.PostalCode,
			&myUser.Lat,
			&myUser.Lng,
		)
		if err != nil {
			if err != sql.ErrNoRows {
				fmt.Println(err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
			} else {
				fmt.Println("No Entry Found in Database for User:" + checker)
			}
		}
	}
	return myUser
}

//function locates user's cookie and check against session data.
//returns true if user found(logged in), else return false
//function DOES NOT issue cookie if not found
func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	var checker string

	query := "SELECT Username FROM sessions WHERE UUID=?"

	err = db.QueryRow(query, myCookie.Value).Scan(&checker)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Print(err)
		} else {
			fmt.Println("No Entry Found in Database for UUID:" + myCookie.Value)
		}
	} else {
		query = "SELECT Username FROM users WHERE Username=?"

		err = db.QueryRow(query, checker).Scan(&checker)
		if err != nil {
			fmt.Print(err)
		} else {
			return true
		}
	}
	return false
}

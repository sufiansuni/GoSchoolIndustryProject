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
	r.HandleFunc("/admin", admin)
	r.HandleFunc("/setlocation", setlocation)
	r.HandleFunc("/confirmlocation", confirmlocation)

	// Sample Handle Func
	r.HandleFunc("/sample", sample)

	log.Fatal(http.ListenAndServe(":8080", r))
}

package models

import "strings"

type User struct {
	Username       string //primary key
	Password       []byte
	First          string
	Last           string
	Gender         string
	Birthday       string
	Height         int
	Weight         float64
	ActivityLevel  int
	CaloriesPerDay int
	Halal          bool
	Vegan          bool
	Address        string
	Unit           string
	Lat            float64
	Lng            float64
}

func (myUser *User) FillDefaults() {
	// set default value for activity level if 0
	if myUser.ActivityLevel == 0 {
		myUser.ActivityLevel = 1
	}
}

func (myUser *User) AdjustStrings() {
	//titlecase first name
	myUser.First = strings.ToLower(myUser.First)
	myUser.First = strings.Title(myUser.First)
	//titlecase last name
	myUser.Last = strings.ToLower(myUser.Last)
	myUser.Last = strings.Title(myUser.Last)
	//titlecase address
	myUser.Address = strings.ToUpper(myUser.Address)
	//titlecase unit
	myUser.Unit = strings.ToUpper(myUser.Unit)
}

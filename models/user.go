package models

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
	PostalCode     int
	Lat            float64
	Lng            float64
}

func (myUser *User) FillDefaults() {
	// set default value for birthday if blank
	if myUser.Birthday == "" {
		myUser.Birthday = "1000-01-01"
	}

	// set default value for activity level if 0
	if myUser.ActivityLevel == 0 {
		myUser.ActivityLevel = 1
	}
}

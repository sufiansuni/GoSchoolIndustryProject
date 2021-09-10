package models

type Session struct {
	UUID     string //primary key
	Username string //foreign key
}
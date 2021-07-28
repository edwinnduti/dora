package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// student info
type Student struct {
	ID       primitive.ObjectID `bson:"_id" json:"studentID"`
	AdmNo    string             `json:"admNo"`
	FullName string             `json:"studentFullname"`
	Email    string             `json:"studentEmail"`
	Password string             `json:"studentPassword"`
}

// course details
type Course struct {
	CourseID   primitive.ObjectID `bson:"_id" json:"courseID"`
	CourseName string             `json:"courseName"`
	Units      []Unit             `json:"units"`
}

// unit info
type Unit struct {
	UnitID   primitive.ObjectID `bson:"_id" json:"unitID"`
	UnitName string             `json:"unitname"`
	UnitCode string             `json:"unitcode"`
}

// authentication struct
type Auth struct {
	Userid   string `json:"userid"`
	Password string `json:"password"`
}

// response struct
type Response struct {
	ID      primitive.ObjectID `json:"ID"`
	Message string             `json:"message"`
}

// describe config model
type Config struct {
	Host       string
	Dbport     string
	Dbusername string
	Dbname     string
	Passwd     string
	Key        string
}

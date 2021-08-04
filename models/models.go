package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// student info
type Student struct {
	ID          primitive.ObjectID `bson:"_id" json:"studentID"`
	AdmNo       string             `json:"admNo"`
	FullName    string             `json:"studentFullname"`
	Email       string             `json:"studentEmail"`
	Password    string             `json:"studentPassword"`
	YearOfStudy string             `json:"yearOfStudy"`
	Course      string             `json:"course"`
	Semester    string             `json:"currentSemester"`
}

type Details struct {
	DetailsID primitive.ObjectID `bson:"_id" json:"detailsID"`
	Unit      string             `json:"unit"`
	Lecturer  string             `json:"lecturer"`
}

// questions
type Questions struct {
	QuestionsID                        primitive.ObjectID `bson:"_id" json:"questionsID"`
	ClarityOfCourseUnitObjective       string             `json:"clarityOfCourseUnitObjective"`
	AchievementOfCourseUnitObjective   string             `json:"achievementOfCourseUnitObjective"`
	ValuableCourseOutline              string             `json:"valuableCourseOutline"`
	InterpretationOfConcepts           string             `json:"interpretationOfConcepts"`
	ExtentOfCoverage                   string             `json:"extentOfCoverage"`
	ClarityOfPresentation              string             `json:"clarityOfPresentation"`
	SufficiencyOfHandouts              string             `json:"sufficiencyOfHandouts"`
	GuidanceOnUse                      string             `json:"guidanceOnUse"`
	AdequancyOfReadings                string             `json:"adequancyOfReadings"`
	ExhibitsHighLevel                  string             `json:"exhibitsHighLevel"`
	OrganizedNotes                     string             `json:"organizedNotes"`
	RelevantAssignment                 string             `json:"relevantAssignment"`
	MakesAssignments                   string             `json:"makesAssignments"`
	GivesFeedback                      string             `json:"givesFeedback"`
	AttendsToLessons                   string             `json:"attendsToLessons"`
	KeepsTimetable                     string             `json:"keepsTimetable"`
	Punctual                           string             `json:"punctual"`
	TeachesFullSession                 string             `json:"teachesFullSession"`
	UseOfClassTime                     string             `json:"useOfClassTime"`
	PresentCourseConceptsInterestingly string             `json:"presentCourseConceptsInterestingly"`
	PresentCourseConceptsClearly       string             `json:"presentCourseConceptsClearly"`
	FacilitatesClassParticipation      string             `json:"facilitatesClassParticipation"`
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

// pass student with string id
type IdDetail struct {
	ID string `json:"ID"`
}

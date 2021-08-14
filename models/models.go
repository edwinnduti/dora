package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// student info
type Student struct {
	ID          primitive.ObjectID `bson:"_id" json:"studentID"`
	AdmNo       string             `json:"admno"`
	FullName    string             `json:"fullname"`
	Email       string             `json:"email"`
	Password    string             `json:"password"`
	YearOfStudy string             `json:"yearofstudy"`
	Course      string             `json:"course"`
	Semester    string             `json:"semester"`
}

type Details struct {
	DetailsID primitive.ObjectID `bson:"_id" json:"detailsID"`
	Student   string             `json:"studentId"`
	Unit      string             `json:"unit"`
	Lecturer  string             `json:"lecturer"`
	Rating    int                `json:"rating"`
}

// questions
type Questions struct {
	QuestionsID                        primitive.ObjectID `bson:"_id" json:"questionsID"`
	Userid                             primitive.ObjectID `json:"userid"`
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
	CourseID      primitive.ObjectID `bson:"_id" json:"courseID"`
	CourseName    string             `json:"coursename"`
	NumberOfUnits int                `json:"numberofunits"`
	Units         []Unit             `json:"units"`
}

// unit info
type Unit struct {
	UnitID   primitive.ObjectID `bson:"_id" json:"unitID"`
	UnitName string             `json:"unitname"`
	UnitCode string             `json:"unitcode"`
	Lecturer string             `json:"lecturer"`
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
	ID        string   `json:"ID"`
	Courses   []Course `json:"courses"`
	Units     []Unit   `json:"units"`
	Course    Course   `json:"course"`
	DetailsID string   `json:"detailsid"`
}

// Lecturer struct
type Lecturer struct {
	ID         primitive.ObjectID `bson:"_id" json:"Id"`
	LecturerID string             `json:"lecturerid"`
	FullName   string             `json:"lecturerfullname"`
	Email      string             `json:"lectureremail"`
	Password   string             `json:"lecturerpassword"`
}

// StringIdCourses
type StringIdCourse struct {
	CourseID      string      `json:"courseID"`
	CourseName    string      `json:"coursename"`
	NumberOfUnits int         `json:"numberofunits"`
	Units         []StrIdUnit `json:"units"`
}

// unit info
type StrIdUnit struct {
	UnitID   string `json:"unitID"`
	UnitName string `json:"unitname"`
	UnitCode string `json:"unitcode"`
	Lecturer string `json:"lecturer"`
}

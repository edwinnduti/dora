/*
[*] Copyright © 2021
[*] Dev/Author ->  Edwin Nduti
[*] Description:
	The code stores names in a mongodb file.
    Written in pure Golang.
*/

package main

// libraries to use
import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/edwinnduti/dora/helpers"
	"github.com/edwinnduti/dora/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/urfave/negroni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// templ
var (
	dir = "assets/"
)

// set session key
var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte(os.Getenv("SECRET_KEY"))
	store = sessions.NewCookieStore(key)
)

// match templates
var templates map[string]*template.Template

//Compile view templates
func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["welcomehandler"] = template.Must(template.ParseFiles("templates/welcomehandler.html", "templates/base.html"))
	templates["studentsignuphandler"] = template.Must(template.ParseFiles("templates/studentsignuphandler.html", "templates/base.html"))
	templates["yearofstudyhandler"] = template.Must(template.ParseFiles("templates/yearofstudyform.html", "templates/navbar.html", "templates/base.html"))
	templates["unitandlechandler"] = template.Must(template.ParseFiles("templates/unitandlecturerform.html", "templates/navbar.html", "templates/base.html"))

	templates["questionpagehandler"] = template.Must(template.ParseFiles("templates/questionspage.html", "templates/base.html"))
	templates["submitresponsehandler"] = template.Must(template.ParseFiles("templates/submitresponse.html", "templates/base.html"))

	// admin
	templates["adminloginhandler"] = template.Must(template.ParseFiles("templates/adminloginhandler.html", "templates/base.html"))
	templates["allcourseshandler"] = template.Must(template.ParseFiles("templates/allcourseshandler.html", "templates/base.html"))
	templates["dashboardhandler"] = template.Must(template.ParseFiles("templates/dashboardhandler.html", "templates/base.html"))
	templates["allunithandler"] = template.Must(template.ParseFiles("templates/allunithandler.html", "templates/base.html"))
	templates["addunithandler"] = template.Must(template.ParseFiles("templates/addunithandler.html", "templates/base.html"))
	templates["addcourseshandler"] = template.Must(template.ParseFiles("templates/addcourseshandler.html", "templates/base.html"))

}

// database and collection names are statically declared
const (
	database            = "lecture-progress"
	studentCollection   = "studentDetails"
	detailsCollection   = "details"
	questionsCollection = "questions"
	adminCollection     = "adminsDetails"
	courseCollection    = "courseAndUnit"
)

// create connection to mongodb
func CreateConnection() (*mongo.Client, error) {
	// connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set MONGOURI
	MongoURI := os.Getenv("MONGOURI")
	// connect to mongodb
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		MongoURI,
	))
	Check(err)

	// return client and error
	return client, nil
}

// hash password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// match two passwords
func Match(passwd, confirmPasswd string) (string, error) {
	if passwd == confirmPasswd {
		// hash password
		hash, err := HashPassword(passwd)

		return hash, err
	} else {
		var err error = fmt.Errorf("password not matching")
		return "", err
	}
}

/* save new members */
func PostSaveStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student

	client, err := CreateConnection()
	Check(err)

	c := client.Database(database).Collection(studentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//create the new member
	student.ID = primitive.NewObjectID()

	if r.Method == "POST" {
		r.ParseForm()
		// decode incoming values
		student.AdmNo = r.FormValue("userid")
		student.FullName = r.FormValue("fullname")
		student.Email = r.FormValue("email")

		// password confirmation
		passwd := r.FormValue("password")
		confirmPasswd := r.FormValue("confirmPasswd")

		// hash password
		hash, err := Match(passwd, confirmPasswd)
		if err != nil {
			fmt.Println(fmt.Errorf("match password error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// store hash
		student.Password = hash

		result, err := c.InsertOne(ctx, student)
		Check(err)
		fmt.Println("added new object of Id ", result.InsertedID.(primitive.ObjectID))

		err = helpers.SendMailTo(&student)
		Check(err)

		// set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "POST")

		http.Redirect(w, r, "/", http.StatusFound)
	}
	if r.Method == "GET" {
		// set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "GET")
		w.WriteHeader(http.StatusOK)
		RenderTemp(w, "studentsignuphandler", "base", nil)
	}
}

/* login page view */
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "welcomehandler", "base", nil)
}

/* sign up view */
func StudentSignUpHandler(w http.ResponseWriter, r *http.Request) {
	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "studentsignuphandler", "base", nil)
}

/* login handler */
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// connect to database
	client, err := CreateConnection()
	Check(err)

	inCollection := client.Database(database).Collection(studentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	// create an empty student struct
	var student models.Student

	//  allow parsing form
	r.ParseForm()

	// decode incoming values
	userid := r.FormValue("userid")
	passwd := r.FormValue("password")

	// Authentication goes here
	// find table document
	err = inCollection.FindOne(ctx, bson.M{"admno": userid}).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println(fmt.Errorf("no documents error: %v", err))
			w.WriteHeader(http.StatusOK)
			//render template
			RenderTemp(w, "studentsignuphandler", "base", nil)

		} else {
			w.WriteHeader(http.StatusOK)
			//render template
			RenderTemp(w, "studentsignuphandler", "base", nil)

		}
	}

	// hash password
	hash, err := HashPassword(passwd)
	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// compare hashes
	if ok := CheckPasswordHash(passwd, hash); !ok {
		log.Fatalln("Wrong password!")
		w.WriteHeader(http.StatusOK)
		//render template
		RenderTemp(w, "studentsignuphandler", "base", nil)
	}

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)

	// Redirect to long url
	// id := Between(student.ID.Hex(), "ObjectID(\"", "\")")
	log.Println(student.ID.Hex())
	uri := fmt.Sprintf("/year-of-study/%v", student.ID.Hex())
	http.Redirect(w, r, uri, http.StatusFound)
}

/* yearofstudyhandler view */
func YearofstudyhandlerHandler(w http.ResponseWriter, r *http.Request) {

	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Year of study Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["userid"]

	// empty slice of Courses
	var courses []models.Course

	// connect to database
	client, err := CreateConnection()
	Check(err)

	inCourseCollection := client.Database(database).Collection(courseCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	// get all documents
	cursor, err := inCourseCollection.Find(ctx, bson.M{})
	Check(err)

	err = cursor.All(ctx, &courses)
	Check(err)

	// student struct with id and issue url id to student.ID
	studentDetail := models.IdDetail{
		ID:      objId,
		Courses: courses,
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "yearofstudyhandler", "base", studentDetail)
}

/* get-year-and-course */
func GetYearAndCourse(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Save year and course Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["userid"]
	userid, err := primitive.ObjectIDFromHex(objId)
	Check(err)

	// connect to database
	client, err := CreateConnection()
	Check(err)

	inStudentCollection := client.Database(database).Collection(studentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	//  allow parsing form
	r.ParseForm()

	// decode incoming values
	yearOfStudy := r.FormValue("yearOfStudy")
	course := r.FormValue("course")
	semester := r.FormValue("currentSemester")

	// find table document
	filter := bson.M{"_id": userid}

	// define update
	update := bson.D{
		{Key: "$set", Value: bson.M{"yearofstudy": yearOfStudy}},
		{Key: "$set", Value: bson.M{"course": course}},
		{Key: "$set", Value: bson.M{"semester": semester}},
	}
	_, err = inStudentCollection.UpdateOne(ctx, filter, update)
	Check(err)

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST")

	//redirect to profile
	uri := fmt.Sprintf("/unit-and-lecture/%s", userid.Hex())
	http.Redirect(w, r, uri, http.StatusFound)
}

/* unit-and-lecturer view */
func UnitAndLecturerHandler(w http.ResponseWriter, r *http.Request) {

	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "unitandlecturer PAGE Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["userid"]
	userid, err := primitive.ObjectIDFromHex(objId)
	Check(err)

	// empty slice of Courses
	var course models.Course

	// connect to database
	client, err := CreateConnection()
	Check(err)
	fmt.Print("database connected\n")

	// select student collection
	inStudentCollection := client.Database(database).Collection(studentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create an empty student struct
	var student models.Student

	// find student table document
	err = inStudentCollection.FindOne(ctx, bson.M{"_id": userid}).Decode(&student)
	Check(err)

	// get course selected in previous page
	inCourseCollection := client.Database(database).Collection(courseCollection)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	// find course table document
	err = inCourseCollection.FindOne(ctx, bson.M{"coursename": student.Course}).Decode(&course)
	Check(err)

	// student struct with id and issue url id to student.ID
	studentDetail := models.IdDetail{
		ID:     objId,
		Course: course,
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "unitandlechandler", "base", studentDetail)
}

/* get-unit-and-lec */
func GetUnitAndLec(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Save unit and lecturer Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["userid"]
	userid, err := primitive.ObjectIDFromHex(objId)
	Check(err)

	// connect to database
	client, err := CreateConnection()
	Check(err)

	inCollection := client.Database(database).Collection(detailsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	//  allow parsing form
	r.ParseForm()

	// empty details struct
	var details models.Details

	// decode incoming values
	details.DetailsID = primitive.NewObjectID()
	details.Student = userid.Hex()
	details.Unit = r.FormValue("unitOfStudy")
	details.Lecturer = r.FormValue("lecturer")
	details.Rating = 0

	// insert detail info
	result, err := inCollection.InsertOne(ctx, details)
	Check(err)
	fmt.Println("added new object of Id ", result.InsertedID.(primitive.ObjectID))

	// // define update
	// update := bson.D{
	// 	{Key: "$set", Value: bson.M{"unit": unitOfStudy}},
	// 	{Key: "$set", Value: bson.M{"lecturer": lecturer}},
	// }
	// _, err = inCollection.UpdateOne(ctx, filter, update)
	// Check(err)

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST")

	//redirect to profile
	uri := fmt.Sprintf("/questions/%s/%s", userid.Hex(), details.DetailsID.Hex())
	http.Redirect(w, r, uri, http.StatusFound)

}

/* questionpage view */
func QuestionsHandler(w http.ResponseWriter, r *http.Request) {

	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Questions page is Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["userid"]
	detailsId := vars["detailsid"]

	// student struct with id and issue url id to student.ID
	studentID := models.IdDetail{
		ID:        objId,
		DetailsID: detailsId,
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "questionpagehandler", "base", studentID)
}

/* evaluate answers */
func GetEvaluationAnswers(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Questions Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["userid"]
	detailsID := vars["detailsid"]
	userid, err := primitive.ObjectIDFromHex(objId)
	Check(err)
	detailsId, err := primitive.ObjectIDFromHex(detailsID)
	Check(err)

	// connect to database
	client, err := CreateConnection()
	Check(err)

	// access questions collection
	inQuestionsCollection := client.Database(database).Collection(questionsCollection)

	// access details collection
	inDetailsCollection := client.Database(database).Collection(detailsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	//  allow parsing form
	r.ParseForm()

	// empty details struct
	var questions models.Questions

	// set empty slice
	var box []int

	// decode incoming values
	questions.QuestionsID = primitive.NewObjectID()
	questions.Userid = userid
	questions.ClarityOfCourseUnitObjective = r.FormValue("clarityOfCourseUnitObjective")
	helpers.Add(&box, &questions.ClarityOfCourseUnitObjective)
	questions.AchievementOfCourseUnitObjective = r.FormValue("achievementOfCourseUnitObjective")
	helpers.Add(&box, &questions.AchievementOfCourseUnitObjective)

	questions.ValuableCourseOutline = r.FormValue("valuableCourseOutline")
	helpers.Add(&box, &questions.ValuableCourseOutline)

	questions.InterpretationOfConcepts = r.FormValue("interpretationOfConcepts")
	helpers.Add(&box, &questions.InterpretationOfConcepts)

	questions.ExtentOfCoverage = r.FormValue("extentOfCoverage")
	helpers.Add(&box, &questions.ExtentOfCoverage)

	questions.ClarityOfPresentation = r.FormValue("clarityOfPresentation")
	helpers.Add(&box, &questions.ClarityOfPresentation)

	questions.SufficiencyOfHandouts = r.FormValue("sufficiencyOfHandouts")
	helpers.Add(&box, &questions.SufficiencyOfHandouts)

	questions.GuidanceOnUse = r.FormValue("guidanceOnUse")
	helpers.Add(&box, &questions.GuidanceOnUse)

	questions.AdequancyOfReadings = r.FormValue("adequancyOfReadings")
	helpers.Add(&box, &questions.AdequancyOfReadings)

	questions.ExhibitsHighLevel = r.FormValue("exhibitsHighLevel")
	helpers.Add(&box, &questions.ExhibitsHighLevel)

	questions.OrganizedNotes = r.FormValue("organizedNotes")
	helpers.Add(&box, &questions.OrganizedNotes)

	questions.RelevantAssignment = r.FormValue("relevantAssignment")
	helpers.Add(&box, &questions.RelevantAssignment)

	questions.MakesAssignments = r.FormValue("makesAssignments")
	helpers.Add(&box, &questions.MakesAssignments)

	questions.GivesFeedback = r.FormValue("givesFeedback")
	helpers.Add(&box, &questions.GivesFeedback)

	questions.AttendsToLessons = r.FormValue("attendsToLessons")
	helpers.Add(&box, &questions.AttendsToLessons)

	questions.KeepsTimetable = r.FormValue("keepsTimetable")
	helpers.Add(&box, &questions.KeepsTimetable)

	questions.Punctual = r.FormValue("punctual")
	helpers.Add(&box, &questions.Punctual)

	questions.TeachesFullSession = r.FormValue("teachesFullSession")
	helpers.Add(&box, &questions.TeachesFullSession)

	questions.UseOfClassTime = r.FormValue("useOfClassTime")
	helpers.Add(&box, &questions.UseOfClassTime)

	questions.PresentCourseConceptsInterestingly = r.FormValue("presentCourseConceptsInterestingly")
	helpers.Add(&box, &questions.PresentCourseConceptsInterestingly)

	questions.PresentCourseConceptsClearly = r.FormValue("presentCourseConceptsClearly")
	helpers.Add(&box, &questions.PresentCourseConceptsClearly)

	questions.FacilitatesClassParticipation = r.FormValue("facilitatesClassParticipation")
	helpers.Add(&box, &questions.FacilitatesClassParticipation)

	// calculate summation and get percentage
	var sum = 0
	for _, num := range box {
		sum = sum + num
	}
	percentageValue := (sum / 110) * 100

	// find table document
	filter := bson.M{"_id": detailsId}

	// update var
	update := bson.D{
		{Key: "$set", Value: bson.M{"rating": percentageValue}},
	}

	// update in detail collection
	_, err = inDetailsCollection.UpdateOne(ctx, filter, update)
	Check(err)

	// insert questions info
	result, err := inQuestionsCollection.InsertOne(ctx, questions)
	Check(err)
	fmt.Println("added new object of Id ", result.InsertedID.(primitive.ObjectID))

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST")

	//redirect to profile
	uri := fmt.Sprintf("/thank-you-response/%s", userid.Hex())
	http.Redirect(w, r, uri, http.StatusFound)
}

/* submitresponse view */
func SubmitResponseHandler(w http.ResponseWriter, r *http.Request) {

	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Submit response Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["userid"]

	// student struct with id and issue url id to student.ID
	studentID := models.IdDetail{
		ID: objId,
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "submitresponsehandler", "base", studentID)
}

/* logout handler */
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)

	// Redirect to long url
	http.Redirect(w, r, "/", http.StatusFound)
}

// admin view
func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "adminloginhandler", "base", nil)
}

/* admin  sign-in handler */
func AdminSignInHandler(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// connect to database
	client, err := CreateConnection()
	Check(err)

	inCollection := client.Database(database).Collection(adminCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	// create an empty student struct
	var lecturer models.Lecturer

	//  allow parsing form
	r.ParseForm()

	// decode incoming values
	userid := r.FormValue("userid")
	passwd := r.FormValue("password")

	// find table document
	err = inCollection.FindOne(ctx, bson.M{"userid": userid}).Decode(&lecturer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println(fmt.Errorf("no documents error: %v", err))
			w.WriteHeader(http.StatusOK)
			//render template
			RenderTemp(w, "adminloginhandler", "base", nil)

		} else {
			w.WriteHeader(http.StatusOK)
			//render template
			RenderTemp(w, "adminloginhandler", "base", nil)

		}
	}

	// hash password
	hash, err := HashPassword(passwd)
	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// compare hashes
	if ok := CheckPasswordHash(passwd, hash); !ok {
		log.Fatalln("Wrong password!")
		w.WriteHeader(http.StatusOK)
		//render template
		RenderTemp(w, "adminloginhandler", "base", nil)
	}

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)

	// redirect user
	log.Println(lecturer.ID.Hex())
	uri := fmt.Sprintln("/dashboard")
	http.Redirect(w, r, uri, http.StatusFound)
}

// dashboard view
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Dashboard Page Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// empty slice of details
	var details []models.Details

	// connect to database
	client, err := CreateConnection()
	Check(err)

	// access details collection
	inDetailsCollection := client.Database(database).Collection(detailsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	// get all documents
	cursor, err := inDetailsCollection.Find(ctx, bson.M{})
	Check(err)

	err = cursor.All(ctx, &details)
	Check(err)

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "dashboardhandler", "base", details)
}

// retrieve courses
func AllCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "All courses Page Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// empty slice of Courses
	var courses []models.Course

	// connect to database
	client, err := CreateConnection()
	Check(err)

	inCourseCollection := client.Database(database).Collection(courseCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	// get all documents
	cursor, err := inCourseCollection.Find(ctx, bson.M{})
	Check(err)

	err = cursor.All(ctx, &courses)
	Check(err)

	var strIdCourses []models.StringIdCourse
	var strIdCourse models.StringIdCourse

	for _, course := range courses {
		strIdCourse.CourseID = course.CourseID.Hex()
		strIdCourse.CourseName = course.CourseName
		strIdCourse.NumberOfUnits = len(course.Units)

		strIdCourses = append(strIdCourses, strIdCourse)
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "allcourseshandler", "base", strIdCourses)
}

// add courses view
func AddCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Add Courses Page Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "addcourseshandler", "base", nil)
}

// save new unit
func SaveCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Save new unit Page Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get values
	if r.Method == "POST" {
		// connect to database
		client, err := CreateConnection()
		Check(err)

		// empty Course struct
		var course models.Course

		// set course id
		course.CourseID = primitive.NewObjectID()

		// select collection
		inCourseCollection := client.Database(database).Collection(courseCollection)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Print("database connected\n")

		r.ParseForm()
		// decode incoming values
		course.CourseName = r.FormValue("coursename")
		if course.CourseName == "" {
			log.Fatal("NO COURSENAME ENTERED")
			// set headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Method", "POST")

			//redirect to profile
			uri := fmt.Sprintln("/dashboard/courses")
			http.Redirect(w, r, uri, http.StatusFound)
		}

		var unit models.Unit

		unit.UnitID = primitive.NewObjectID()
		unit.UnitName = r.FormValue("unitname")
		unit.Lecturer = r.FormValue("lecturer")
		unit.UnitCode = r.FormValue("unitcode")

		if unit.UnitCode == "" || unit.UnitName == "" {
			log.Fatal("NO COURSENAME ENTERED")
			// set headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Method", "POST")

			//redirect to profile
			uri := fmt.Sprintln("/dashboard/courses")
			http.Redirect(w, r, uri, http.StatusFound)
		}

		// push to struct
		course.NumberOfUnits = 0

		course.Units = append(course.Units, unit)

		// insert in collection
		_, err = inCourseCollection.InsertOne(ctx, course)
		Check(err)

		// set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "POST")

		//redirect to profile
		uri := fmt.Sprintln("/dashboard/courses")
		http.Redirect(w, r, uri, http.StatusFound)
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")

	//redirect to profile
	uri := fmt.Sprintln("/dashboard/courses")
	http.Redirect(w, r, uri, http.StatusFound)

}

// show units
func AllUnitsHandler(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "All units Page Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["courseid"]
	courseid, err := primitive.ObjectIDFromHex(objId)
	Check(err)

	// connect to database
	client, err := CreateConnection()
	Check(err)

	// empty Course struct
	var course models.Course

	inCourseCollection := client.Database(database).Collection(courseCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	// find table document
	err = inCourseCollection.FindOne(ctx, bson.M{"_id": courseid}).Decode(&course)
	Check(err)

	// empty string id course struct
	var strIdCourse models.StringIdCourse
	var strIdUnit models.StrIdUnit

	// match data
	strIdCourse.CourseID = course.CourseID.Hex()
	strIdCourse.CourseName = course.CourseName
	strIdCourse.NumberOfUnits = course.NumberOfUnits

	// range over units
	for _, units := range course.Units {
		strIdUnit.UnitID = units.UnitID.Hex()
		strIdUnit.UnitCode = units.UnitCode
		strIdUnit.UnitName = units.UnitName
		strIdUnit.Lecturer = units.Lecturer

		// append unit to units
		strIdCourse.Units = append(strIdCourse.Units, strIdUnit)
	}

	// number of units
	strIdCourse.NumberOfUnits = len(strIdCourse.Units)

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "allunithandler", "base", strIdCourse)
}

// add unit view
func AddNewUnitHandler(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Add unit Page Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["courseid"]
	courseid, err := primitive.ObjectIDFromHex(objId)
	Check(err)

	// connect to database
	client, err := CreateConnection()
	Check(err)

	// empty Course struct
	var course models.Course

	inCourseCollection := client.Database(database).Collection(courseCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	// find table document
	err = inCourseCollection.FindOne(ctx, bson.M{"_id": courseid}).Decode(&course)
	Check(err)

	// empty string id course struct
	var strIdCourse models.StringIdCourse
	var strIdUnit models.StrIdUnit

	// match data
	strIdCourse.CourseID = course.CourseID.Hex()
	strIdCourse.CourseName = course.CourseName
	strIdCourse.NumberOfUnits = course.NumberOfUnits

	// range over units
	for _, units := range course.Units {
		strIdUnit.UnitID = units.UnitID.Hex()
		strIdUnit.UnitCode = units.UnitCode
		strIdUnit.UnitName = units.UnitName
		strIdUnit.Lecturer = units.Lecturer

		// append unit to units
		strIdCourse.Units = append(strIdCourse.Units, strIdUnit)
	}

	// number of units
	strIdCourse.NumberOfUnits = len(strIdCourse.Units)

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "addunithandler", "base", strIdCourse)
}

// save new unit
func SaveNewUnitHandler(w http.ResponseWriter, r *http.Request) {
	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Save new unit Page Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// get id
	vars := mux.Vars(r)
	objId := vars["courseid"]
	courseid, err := primitive.ObjectIDFromHex(objId)
	Check(err)

	// get values
	if r.Method == "POST" {
		// connect to database
		client, err := CreateConnection()
		Check(err)

		// empty Course struct
		var unit models.Unit

		inCourseCollection := client.Database(database).Collection(courseCollection)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Print("database connected\n")

		r.ParseForm()
		// decode incoming values
		unit.UnitCode = r.FormValue("unitcode")
		unit.UnitName = r.FormValue("nameofunit")
		unit.Lecturer = r.FormValue("lecturer")

		// find table document
		filter := bson.M{"_id": courseid}

		// update var
		update := bson.D{
			{Key: "$push", Value: bson.M{"units": unit}},
		}

		// update in collection
		_, err = inCourseCollection.UpdateOne(ctx, filter, update)
		Check(err)

		// set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "POST")

		//redirect to profile
		uri := fmt.Sprintf("/dashboard/courses/%s", courseid.Hex())
		http.Redirect(w, r, uri, http.StatusFound)
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")

	//redirect to profile
	uri := fmt.Sprintf("/dashboard/courses/%s/addUnit", courseid.Hex())
	http.Redirect(w, r, uri, http.StatusFound)

}

/* logout admin handler */
func LogoutAdminHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)

	// Redirect to long url
	http.Redirect(w, r, "/admin", http.StatusFound)
}

/* function render template */
//Render templates for the given name, template definition and data object
func RenderTemp(w http.ResponseWriter, name string, template string, viewModel interface{}) {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "The template does not exist.", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, template, viewModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* log errors */
func Check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func Checkf(line string, err error) {
	if err != nil {
		log.Fatalln(line, " : ", err)
	}
}

// Main function
func main() {
	/*
	   mgo.SetDebug(true)
	   mgo.SetLogger(log.New(os.Stdout,"err",6))

	   The above two lines are for debugging errors
	   that occur straight from accessing the mongo db
	*/

	//Register router{}
	r := mux.NewRouter().StrictSlash(false)

	// API routes,handlers and methods
	r.HandleFunc("/", WelcomeHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/signUp", StudentSignUpHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/year-of-study/{userid}", YearofstudyhandlerHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/unit-and-lecture/{userid}", UnitAndLecturerHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/questions/{userid}/{detailsid}", QuestionsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/thank-you-response/{userid}", SubmitResponseHandler).Methods("GET", "OPTIONS")

	// admin
	r.HandleFunc("/admin", AdminLoginHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/adminsignIn", AdminSignInHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/dashboard", DashboardHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/dashboard/courses", AllCoursesHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/dashboard/courses/addCourse", AddCoursesHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/savenewcourse", SaveCoursesHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/dashboard/courses/{courseid}", AllUnitsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/dashboard/courses/{courseid}/addUnit", AddNewUnitHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/savenewunit/{courseid}", SaveNewUnitHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/adminLogout", LogoutAdminHandler).Methods("GET", "OPTIONS")

	// route action links
	r.HandleFunc("/register", PostSaveStudent).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", LoginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/get-year-and-course/{userid}", GetYearAndCourse).Methods("POST", "OPTIONS")
	r.HandleFunc("/get-unit-and-lec/{userid}", GetUnitAndLec).Methods("POST", "OPTIONS")
	r.HandleFunc("/get-evaluation-answers/{userid}/{detailsid}", GetEvaluationAnswers).Methods("POST", "OPTIONS")
	r.HandleFunc("/logout", LogoutHandler).Methods("GET", "OPTIONS")

	// route assets e.g images, css, javascript
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(dir))))

	//Get port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}

	// establish logger
	n := negroni.Classic()
	n.UseHandler(r)
	server := &http.Server{
		Handler: n,
		Addr:    ":" + Port,
	}
	log.Printf("Listening on PORT: %s", Port)
	server.ListenAndServe()
}

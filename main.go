/*
[*] Copyright © 2020
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
	templates["dashboardhandler"] = template.Must(template.ParseFiles("templates/dashboardhandler.html", "templates/yearofstudyform.html", "templates/navbar.html", "templates/base.html"))
	templates["unitandlechandler"] = template.Must(template.ParseFiles("templates/dashboardhandler.html", "templates/unitandlecturerform.html", "templates/navbar.html", "templates/base.html"))

	templates["questionpagehandler"] = template.Must(template.ParseFiles("templates/questionspage.html", "templates/base.html"))
	templates["submitresponsehandler"] = template.Must(template.ParseFiles("templates/submitresponse.html", "templates/base.html"))
}

// database and collection names are statically declared
const database, collection = "lecture-progress", "studentDetails"

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

	c := client.Database(database).Collection(collection)
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

	inCollection := client.Database(database).Collection(collection)
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
	uri := fmt.Sprintf("/dashboard/%v", student.ID.String())
	http.Redirect(w, r, uri, http.StatusFound)
}

/* dashboad view */
func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	// instance of sessions
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Fatal("session error: ", err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, "Dashboard Forbidden!")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "dashboardhandler", "base", nil)
}

/* get year and course */
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

	inCollection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	//  allow parsing form
	r.ParseForm()

	// decode incoming values
	yearOfStudy := r.FormValue("yearOfStudy")
	course := r.FormValue("course")

	// find table document
	filter := bson.M{"_id": userid}

	// define update
	update := bson.D{
		{Key: "$set", Value: bson.M{"yearOfStudy": yearOfStudy}},
		{Key: "$set", Value: bson.M{"course": course}},
	}
	_, err = inCollection.UpdateOne(ctx, filter, update)
	Check(err)

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST")

	//redirect to profile
	uri := fmt.Sprintf("/unit-and-lecture/%s", userid.String())
	http.Redirect(w, r, uri, http.StatusFound)
}

/* unitandlecturer view */
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

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "unitandlechandler", "base", nil)
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

	inCollection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	//  allow parsing form
	r.ParseForm()

	// empty details struct
	var details models.Details

	// decode incoming values
	details.DetailsID = userid
	details.Unit = r.FormValue("unitOfStudy")
	details.Lecturer = r.FormValue("lecturer")

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
	uri := fmt.Sprintf("/questions/%s", userid.String())
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

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "questionpagehandler", "base", nil)
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
	userid, err := primitive.ObjectIDFromHex(objId)
	Check(err)

	// connect to database
	client, err := CreateConnection()
	Check(err)

	inCollection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("database connected\n")

	//  allow parsing form
	r.ParseForm()

	// empty details struct
	var questions models.Questions

	// decode incoming values
	questions.QuestionsID = userid
	questions.ClarityOfCourseUnitObjective = r.FormValue("clarityOfCourseUnitObjective")
	questions.AchievementOfCourseUnitObjective = r.FormValue("achievementOfCourseUnitObjective")
	questions.ValuableCourseOutline = r.FormValue("valuableCourseOutline")
	questions.InterpretationOfConcepts = r.FormValue("interpretationOfConcepts")
	questions.ExtentOfCoverage = r.FormValue("extentOfCoverage")
	questions.ClarityOfPresentation = r.FormValue("clarityOfPresentation")
	questions.SufficiencyOfHandouts = r.FormValue("sufficiencyOfHandouts")
	questions.GuidanceOnUse = r.FormValue("guidanceOnUse")
	questions.AdequancyOfReadings = r.FormValue("adequancyOfReadings")
	questions.ExhibitsHighLevel = r.FormValue("exhibitsHighLevel")
	questions.OrganizedNotes = r.FormValue("organizedNotes")
	questions.RelevantAssignment = r.FormValue("relevantAssignment")
	questions.MakesAssignments = r.FormValue("makesAssignments")
	questions.GivesFeedback = r.FormValue("givesFeedback")
	questions.AttendsToLessons = r.FormValue("attendsToLessons")
	questions.KeepsTimetable = r.FormValue("keepsTimetable")
	questions.Punctual = r.FormValue("punctual")
	questions.TeachesFullSession = r.FormValue("teachesFullSession")
	questions.UseOfClassTime = r.FormValue("useOfClassTime")
	questions.PresentCourseConceptsInterestingly = r.FormValue("presentCourseConceptsInterestingly")
	questions.PresentCourseConceptsClearly = r.FormValue("presentCourseConceptsClearly")
	questions.FacilitatesClassParticipation = r.FormValue("facilitatesClassParticipation")

	// insert detail info
	result, err := inCollection.InsertOne(ctx, questions)
	Check(err)
	fmt.Println("added new object of Id ", result.InsertedID.(primitive.ObjectID))

	// // find table document
	// filter := bson.M{"_id": userid}

	// // define update
	// update := bson.D{
	// 	{Key: "$set", Value: bson.M{"yearOfStudy": unitOfStudy}},
	// 	{Key: "$set", Value: bson.M{"course": lecturer}},
	// }
	// _, err = inCollection.UpdateOne(ctx, filter, update)
	// Check(err)

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST")

	//redirect to profile
	uri := fmt.Sprintf("/thank-you-response/%s", userid)
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

	// set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET")
	w.WriteHeader(http.StatusOK)

	//render template
	RenderTemp(w, "submitresponsehandler", "base", nil)
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
	r.HandleFunc("/dashboard/{userid}", DashboardHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/unit-and-lecture/{userid}", UnitAndLecturerHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/questions/{userid}", QuestionsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/thank-you-response/{userid}", SubmitResponseHandler).Methods("GET", "OPTIONS")

	// route action links
	r.HandleFunc("/register", PostSaveStudent).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", LoginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/get-year-and-course/{userid}", GetYearAndCourse).Methods("POST", "OPTIONS")
	r.HandleFunc("/get-unit-and-lec/{userid}", GetUnitAndLec).Methods("POST", "OPTIONS")
	r.HandleFunc("/get-evaluation-answers/{userid}", GetEvaluationAnswers).Methods("POST", "OPTIONS")
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

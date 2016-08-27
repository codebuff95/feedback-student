package database

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var DatabaseName string
var DatabaseSession *mgo.Session

var SectionCollection *mgo.Collection
var FacultyCollection *mgo.Collection
var SubjectCollection *mgo.Collection
var QuestionCollection *mgo.Collection
var CourseCollection *mgo.Collection
var FeedbackCollection *mgo.Collection
var TextQuestionCollection *mgo.Collection

type DatabaseDetails struct {
	Url    string `json:"url"`
	Dbname string `json:"dbname"`
}

type InitDatabaseError int

func (err InitDatabaseError) Error() string {
	if err == PROBLEMOPENINGFILE {
		return "Could not open file faconfig.json"
	}
	return "Miscellaneous Issues in database"
}

const (
	PROBLEMOPENINGFILE = 0
	PROBLEMDECODING    = 1
)

func InitDatabaseSession() error {
	var err error
	//Database details should match with feedback-admin, so that both processes may operate in the same database.
	myFile, err := os.Open("feedbackadminres/faconfig.json")
	defer myFile.Close()

	if err != nil {
		return InitDatabaseError(PROBLEMOPENINGFILE)
	}
	var myDBDetails DatabaseDetails
	err = json.NewDecoder(myFile).Decode(&myDBDetails)
	if err != nil {
		return InitDatabaseError(PROBLEMDECODING)
	}
	DatabaseSession, err = mgo.Dial(myDBDetails.Url)
	DatabaseName = myDBDetails.Dbname
	log.Println("feedback-student: Initialised new database session to url:", myDBDetails.Url, "and Dbname:", myDBDetails.Dbname)
	return nil
}

func InitCollections() {
	//log.Println("**Initialising Essential Collections with Database",DatabaseName,"**")
	SectionCollection = DatabaseSession.DB(DatabaseName).C("section")
	FacultyCollection = DatabaseSession.DB(DatabaseName).C("faculty")
	SubjectCollection = DatabaseSession.DB(DatabaseName).C("subject")
	QuestionCollection = DatabaseSession.DB(DatabaseName).C("question")
	CourseCollection = DatabaseSession.DB(DatabaseName).C("course")
	FeedbackCollection = DatabaseSession.DB(DatabaseName).C("feedback")
	TextQuestionCollection = DatabaseSession.DB(DatabaseName).C("textquestion")
}

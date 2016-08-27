package question

import (
	"errors"
	"feedback-student/database"
	"log"
)

type Question struct {
	Questionid string `bson:"questionid"`
	Text       string `bson:"text"`
	Weightage  int    `bson:"weightage"`
}

var GlobalQuestions []Question
var GlobalTextQuestions []Question

func InitQuestions() error {
	database.QuestionCollection.Find(nil).All(&GlobalQuestions)
	if len(GlobalQuestions) == 0 {
		return errors.New("No questions found")
	}
	database.TextQuestionCollection.Find(nil).All(&GlobalTextQuestions)
	if len(GlobalTextQuestions) == 0 {
		log.Println("Could not find any text qquestions")
	}
	return nil
}

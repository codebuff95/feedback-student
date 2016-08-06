package question

import(
  "errors"
  "feedback-student/database"
)

type Question struct{
  Questionid string `bson:"questionid"`
  Text string `bson:"text"`
  Weightage int `bson:"weightage"`
}

var GlobalQuestions []Question

func InitQuestions() error{
  err := database.QuestionCollection.Find(nil).All(&GlobalQuestions)
  if len(GlobalQuestions) == 0{
    return errors.New("No questions found")
  }
  return err
}

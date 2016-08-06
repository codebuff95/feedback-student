package subject

import(
  "feedback-student/database"
  "log"
  "gopkg.in/mgo.v2/bson"
)

type Subject struct{
  Subjectid string `bson:"subjectid"`
  Subjectname string `bson:"subjectname"`
  Addedon string  `bson:"addedon"`
}

func GetSubject(sid string) (*Subject,error){
  log.Println("Getting subject with subjectid:",sid)
  var mySubject *Subject = &Subject{}
  err := database.SubjectCollection.Find(bson.M{"subjectid" : sid}).Limit(1).One(mySubject)
  if err != nil{
    log.Println("Could not get subject.")
    return nil,err
  }
  return mySubject,err
}

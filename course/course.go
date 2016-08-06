package course

import(
  "feedback-student/database"
  "log"
  "gopkg.in/mgo.v2/bson"
)

type Course struct{
  Courseid string `bson:"courseid"`
  Coursename string `bson:"coursename"`
  Addedon string  `bson:"addedon"`
}

func GetCourse(cid string) (*Course,error){
  log.Println("Getting course with courseid:",cid)
  var myCourse *Course = &Course{}
  err := database.CourseCollection.Find(bson.M{"courseid" : cid}).Limit(1).One(myCourse)
  if err != nil{
    log.Println("Could not get course.")
    return nil,err
  }
  if myCourse.Courseid == ""{
    log.Println("Course does not exist. Returning nil course.")
    return nil,nil
  }
  log.Println("Successfully got Course")
  return myCourse,err
}

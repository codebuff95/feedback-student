package course

import(
  "feedback-admin/database"
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
  return myCourse,err
}

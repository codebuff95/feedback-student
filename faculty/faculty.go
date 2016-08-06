package faculty

import(
  "feedback-student/database"
  "log"
  "gopkg.in/mgo.v2/bson"
)

type Faculty struct{
  Facultyid string `bson:"facultyid"`
  Facultyname string `bson:"facultyname"`
  Addedon string  `bson:"addedon"`
}

func GetFaculty(fid string) (*Faculty,error){
  log.Println("Getting faculty with facultyid:",fid)
  var myFaculty *Faculty = &Faculty{}
  err := database.FacultyCollection.Find(bson.M{"facultyid" : fid}).Limit(1).One(myFaculty)
  if err != nil{
    log.Println("Could not get faculty.")
    return nil,err
  }
  return myFaculty,err
}

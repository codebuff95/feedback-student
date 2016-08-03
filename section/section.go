package section

import(
  "feedback-student/database"
  "log"
  "gopkg.in/mgo.v2/bson"
)

type Teacher struct{
  Facultyid string `bson:"facultyid"`
  Subjectid string `bson:"subjectid"`
}

type Section struct{
  Sectionid string `bson:"sectionid"`
  Sectionname string `bson:"sectionname"`
  Year int `bson:"year"`
  Session int `bson:"session"`
  Courseid string `bson:"courseid"`
  Teachers *[]Teacher `bson:"teachers,omitempty"`
  Password string `bson:"password"`
  Students int `bson:"students"`
  Addedon string  `bson:"addedon"`
}

func GetSection(sectionid string) (*Section,error){
  log.Println("Getting section with sectionid:",sectionid)
  var mySection *Section = &Section{}
  err := database.SectionCollection.Find(bson.M{"sectionid" : sectionid}).Limit(1).One(mySection)
  if err != nil{
    log.Println("Could not get section.")
    return nil,err
  }
  return mySection,nil
}

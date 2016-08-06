package section

import(
  "feedback-student/database"
  "feedback-student/faculty"
  "feedback-student/subject"
  "log"
  "gopkg.in/mgo.v2/bson"
  "errors"
)

type Teacher struct{
  Facultyid string `bson:"facultyid"`
  Subjectid string `bson:"subjectid"`
}

type DetailedTeacher struct{
  Facultyid string `bson:"facultyid"`
  Facultyname string `bson:"facultyname"`
  Subjectid string `bson:"subjectid"`
  Subjectname string `bson:"subjectname"`
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

func GetDetailedTeachers(myTeachers *[]Teacher) (*[]DetailedTeacher,error){
  log.Println("*Getting Detailed Teachers*")
  if myTeachers == nil || len(*myTeachers) == 0{
    log.Println("myTeachers slice passed to GetDetailedTeachers is nil. Returning nil DetailedTeacher slice.")
    return nil,errors.New("No Teachers")
  }
  myDetailedTeachers := make([]DetailedTeacher,len(*myTeachers))
  for i,myTeacher := range *myTeachers{
    myDetailedTeachers[i].Facultyid = myTeacher.Facultyid
    myDetailedTeachers[i].Subjectid = myTeacher.Subjectid
    myFaculty,err := faculty.GetFaculty(myTeacher.Facultyid)
    if err != nil{
      return nil,errors.New("Could not find faculty with facultyid:" + myTeacher.Facultyid)
    }
    myDetailedTeachers[i].Facultyname = myFaculty.Facultyname

    mySubject,err := subject.GetSubject(myTeacher.Subjectid)
    if err != nil{
      return nil,errors.New("Could not find subject with subjectid:" + mySubject.Subjectid)
    }
    myDetailedTeachers[i].Subjectname = mySubject.Subjectname
  }
  log.Println("Success getting DetailedTeachers from Teachers. Returning slice of DetailedTeachers")
  return &myDetailedTeachers,nil
}

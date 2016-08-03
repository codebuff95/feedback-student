package user

import(
  "net/http"
  "github.com/codebuff95/uafm/usersession"
  "feedback-student/database"
  "log"
  "html/template"
  //"gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Student struct{
  Sectionid *string `bson:"sectionid" json:"sectionid"`
  Password string `bson:"password" json:"password"`
}

func AuthenticateRequest(r *http.Request) (*string,error){
  log.Println("Authenticating incoming request from client:",r.RemoteAddr,"for authentic user")
  userSidCookie,err := r.Cookie("usersid")
  if err != nil{
    log.Println("User Authentication failed")
    return nil,err
  }
  userSid := userSidCookie.Value
  log.Println("Got Cookie Value")
  log.Println("Authenticating UserSession with Sid:\"",userSid,"\"")
  userRid, err := usersession.ValidateSession(userSid)
  if err != nil{
    log.Println("Error authenticating userSid:",err)
  }
  return userRid,err
}

func AuthenticateLoginAttempt(r *http.Request) (*Student,error){
  log.Println("*Authenticating Login Attempt*")
  attemptSectionid := template.HTMLEscapeString(r.Form.Get("sectionid"))       //Escape special characters for security.
  attemptPassword := template.HTMLEscapeString(r.Form.Get("password"))       //Escape special characters for security.
  log.Println("Attempt sectionid :", attemptSectionid, ", Attempt Password:", attemptPassword)
  var myStudent Student = Student{}
  err := database.SectionCollection.
            Find(bson.M{"sectionid":attemptSectionid, "password":attemptPassword}).
                Limit(1).One(&myStudent)
  if err != nil{
    log.Println("Error finding requested section in collection:",err)
    return nil,err
  }
  if myStudent.Sectionid == nil{
    log.Println("Could not find section with supplied credentials. Returning nil Student")
  }
  return &myStudent,err
}

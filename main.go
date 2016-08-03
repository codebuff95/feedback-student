package main

import(
  "feedback-student/login"
  "feedback-student/database"
  "feedback-student/feedback"
  "feedback-student/templates"
  "net/http"
  "github.com/codebuff95/uafm"
  "log"
)

func handlefatalerror(err error){
  if err != nil{
    log.Fatal("*_*_* Fatal Error:",err,"*_*_*")
  }
}

func main(){

  err := uafm.Init("feedbackadminres","studentsession","formsession") //make sure that
  // usersession collection of admin and student account types are different.
  handlefatalerror(err)

  err = database.InitDatabaseSession()
  handlefatalerror(err)

  database.InitCollections()

  err = templates.InitEssentialTemplates()
  handlefatalerror(err)

  http.HandleFunc("/login",login.LoginHandler)
  http.HandleFunc("/feedback",feedback.FeedbackHandler)
  http.ListenAndServe(":8080",nil)
}

package main

import(
  "feedback-student/login"
  "feedback-student/database"
  "feedback-student/feedback"
  "feedback-student/question"
  "feedback-student/templates"
  "feedback-student/logout"
  "net/http"
  "github.com/codebuff95/uafm"
  "log"
)

func handlefatalerror(err error){
  if err != nil{
    log.Fatal("*_*_* Fatal Error:",err,"*_*_*")
  }
}

func MyHandler(w http.ResponseWriter, r *http.Request){
  http.Redirect(w, r, "/login", http.StatusSeeOther)
  return
}

func main(){

  err := uafm.Init("feedbackadminres","studentsession","formsession") //make sure that
  // usersession collection of admin and student account types are different.
  handlefatalerror(err)

  err = database.InitDatabaseSession()
  handlefatalerror(err)

  database.InitCollections()

  log.Println("Initialised Database Collections")

  err = templates.InitEssentialTemplates()
  handlefatalerror(err)

  log.Println("Initialised Essential Templates")

  err = question.InitQuestions()
  handlefatalerror(err)

  log.Println("Initialised Questions")

  http.HandleFunc("/",MyHandler)
  http.HandleFunc("/login",login.LoginHandler)
  http.HandleFunc("/feedback",feedback.FeedbackHandler)

  http.HandleFunc("/logout",logout.LogoutHandler)

  //Start file server for public files.
  http.Handle("/resources/",http.StripPrefix("/resources/",http.FileServer(http.Dir("feedbackstudentres/publicres"))))

  http.ListenAndServe(":8080",nil)
}

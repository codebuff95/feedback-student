package templates

import(
  "html/template"
  "log"
)

var LoginFormTemplate *template.Template
var BadPageTemplate *template.Template
var FeedbackFormTemplate *template.Template
/*var HomePageTemplate *template.Template
var CoursePageTemplate *template.Template
var SubjectPageTemplate *template.Template
var FacultyPageTemplate *template.Template
var SectionPageTemplate *template.Template*/

func InitEssentialTemplates() error{
  var err error
  LoginFormTemplate,err = template.ParseFiles("feedbackstudentres/login.html")
  if err != nil{
    log.Println("Error parsing LoginFormTemplate:",err)
    return err
  }
  BadPageTemplate,err = template.ParseFiles("feedbackstudentres/badpage.html")
  if err != nil{
    log.Println("Error parsing BadPageTemplate:",err)
    return err
  }
  FeedbackFormTemplate,err = template.ParseFiles("feedbackstudentres/feedback.html")
  if err != nil{
    log.Println("Error parsing FeedbackFormTemplate:",err)
    return err
  }
  return nil
}

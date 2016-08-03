package feedback

import(
  //"github.com/codebuff95/uafm"
  "github.com/codebuff95/uafm/usersession"
  "github.com/codebuff95/uafm/formsession"
  "feedback-student/course"
  "feedback-student/user"
  "feedback-student/section"
  "feedback-student/templates"
  "net/http"
  "time"
  "log"
  "errors"
  //"html/template"
)

type FeedbackPage struct{
  Sectionid *string `bson:"sectionid"`
  Sectionname *string `bson:"sectionname"`
  Coursename *string `bson:"coursename"`
}

func displayFeedbackForm(w http.ResponseWriter, r *http.Request,sectionId string){
  log.Println("Displaying feedback form to user.")
  formSid, err := formsession.CreateSession("0",time.Minute*30) //Form created will be valid for 30 minutes.
  if err != nil{
    log.Println("Error creating new session for feedback form:",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Creating new feedback form to client",r.RemoteAddr,"with formSid:",*formSid)  //Enter client ip address and new form SID.
  mySection,err := section.GetSection(sectionId)
  if err != nil{
    displayBadPage(w,r,errors.New("Bad Section ID"))
    return
  }
  var myFeedbackPage FeedbackPage
  myFeedbackPage.Sectionid = &mySection.Sectionid
  myFeedbackPage.Sectionname = &mySection.Sectionname

  myCourse,err := course.GetCourse(mySection.Courseid)

  if err != nil{
    displayBadPage(w,r,errors.New("Bad Course ID for given Section ID"))
    return
  }

  myFeedbackPage.Coursename = &myCourse.Coursename

  templates.FeedbackFormTemplate.Execute(w,myFeedbackPage)
}

func displayBadPage(w http.ResponseWriter, r *http.Request, err error){
  templates.BadPageTemplate.Execute(w,err.Error())
}

func FeedbackHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***FEEDBACK HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  sectionId, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, redirecting to login page:",err)
    http.Redirect(w, r, "/login", http.StatusSeeOther)
    return
  }
  log.Println("User Session authentic:",err)
  if r.Method == "GET"{
    log.Println("Request method is GET")
    displayFeedbackForm(w,r,*sectionId)
    return
  }
  //method = POST
  log.Println("Request method is POST")
  r.ParseForm()

  //Validate entered form session.
  _,err = formsession.ValidateSession(r.Form.Get("formsid"))
  if err != nil{
    log.Println("Error validating formSid:",err)
    displayBadPage(w,r,err)
    return
  }

  //Form is valid. Continue with deleting the form session.
  formSidDeleted,err := formsession.DeleteSession(r.Form.Get("formsid"))
  log.Println("Deleted",formSidDeleted,"form sessions for sid",r.Form.Get("formsid"),".")

  //Login Attempt Authentication begin.
  attemptStudent,err := user.AuthenticateLoginAttempt(r)
  if err != nil || attemptStudent == nil{
    displayBadPage(w,r,errors.New("Section with given credentials does not exist"))
    return
  }

  //Authentic Login Attempt. Set cookie on client, and redirect to homepage.
  newUserSid,err := usersession.CreateSession(*attemptStudent.Sectionid,time.Hour*24*3)
  if err != nil{
    log.Println("Error creating new UserSid:",err)
    displayBadPage(w,r,err)
    return
  }
  myCookie := &http.Cookie{Name:"usersid",Value:*newUserSid}  //Cookie is not persistent for security purposes.
  http.SetCookie(w, myCookie)
  log.Println("usersid Cookie successfully set on client. Redirecting to feedbackpage.")
  http.Redirect(w, r, "/feedback", http.StatusSeeOther)
  return
}

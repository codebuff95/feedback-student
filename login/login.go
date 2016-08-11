package login

import(
  //"github.com/codebuff95/uafm"
  "github.com/codebuff95/uafm/usersession"
  "github.com/codebuff95/uafm/formsession"
  "feedback-student/user"
  "feedback-student/college"
  "feedback-student/templates"
  "net/http"
  "time"
  "log"
  "errors"
  //"html/template"
)

type LoginPage struct{
  Formsid *string
  Collegename *string
}

func displayLoginForm(w http.ResponseWriter, r *http.Request){
  log.Println("Displaying login form to user.")
  formSid, err := formsession.CreateSession("0",time.Minute*10) //Form created will be valid for 10 minutes.
  if err != nil{
    log.Println("Error creating new session for login form:",err)
    displayBadPage(w,r,err)
    return
  }
  var myLoginPage LoginPage
  myLoginPage.Formsid = formSid
  myLoginPage.Collegename = &college.GlobalDetails.Collegename
  log.Println("Creating new login form to client",r.RemoteAddr,"with formSid:",*formSid)  //Enter client ip address and new form SID.
  templates.LoginFormTemplate.Execute(w,myLoginPage)
}

func displayBadPage(w http.ResponseWriter, r *http.Request, err error){
  templates.BadPageTemplate.Execute(w,err.Error())
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***STUDENT LOGIN HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  userId, err := user.AuthenticateRequest(r)
  if err == nil{ //user session authentic. Redirect to home page.
    log.Println("User session is authentic with UserId:",*userId,", redirecting to homepage.")
    http.Redirect(w, r, "/feedback", http.StatusSeeOther)
    return
  }
  log.Println("User Session not authentic:",err)
  if r.Method == "GET"{
    log.Println("Request method is GET")
    displayLoginForm(w,r)
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
    displayBadPage(w,r,errors.New("Given credentials do not exist"))
    return
  }

  //Authentic Login Attempt. Set cookie on client, and redirect to homepage.
  newUserSid,err := usersession.CreateSession(*attemptStudent.Sectionid,time.Hour*24*3)
  if err != nil{
    log.Println("Error creating new UserSid:",err)
    displayBadPage(w,r,err)
    return
  }
  myCookie := &http.Cookie{Name:"usersid",Value:*newUserSid, Expires: time.Now().Add(time.Minute * 60)}  //Cookie is being made persistent so that the password of the section can be deleted.
  http.SetCookie(w, myCookie)
  log.Println("usersid Cookie successfully set on client. Deleting password from sectionid")
  err = user.RemovePassword(r)
  if err != nil{
    log.Println("Error removing password for logging in user.")
    displayBadPage(w,r,errors.New("Error completing log in sequence"))
    return
  }
  log.Println("Success removing password. Redirecting to feedback page.")
  http.Redirect(w, r, "/feedback", http.StatusSeeOther)
  return
}

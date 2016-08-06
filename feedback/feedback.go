package feedback

import(
  //"github.com/codebuff95/uafm"
  //"github.com/codebuff95/uafm/usersession"
  "github.com/codebuff95/uafm/formsession"
  "feedback-student/course"
  "feedback-student/question"
  "feedback-student/user"
  "feedback-student/section"
  "feedback-student/templates"
  "feedback-student/database"
  "net/http"
  "html/template"
  "strings"
  "time"
  "log"
  "errors"
  "strconv"
  //"html/template"
)

type Feedback struct{
  Sectionid string  `bson:"sectionid"`
  Ratings []Rating  `bson:"ratings"`
  Addedon string  `bson:"addedon"`
}

type Rating struct{
  Facultyid string  `bson:"facultyid"`
  Subjectid string `bson:"subjectid"`
  Points []Point  `bson:"points"`
}

type Point struct{
  Questionid string `bson:"questionid"`
  Marks int `bson:"marks"`
}

type FeedbackPage struct{
  Sectionid *string `bson:"sectionid"`
  Sectionname *string `bson:"sectionname"`
  Coursename *string `bson:"coursename"`
  DetailedTeachers *[]section.DetailedTeacher `bson:"detailedteachers"`
  Questions *[]question.Question
  FormSid *string
}

func parseFeedbackForm(w http.ResponseWriter, r *http.Request, sectionId string){
  log.Println("**PARSE FEEDBACK FORM**")
  r.ParseForm()
  submittedSectionId := template.HTMLEscapeString(r.Form.Get("sectionid"))
  if strings.Compare(sectionId,submittedSectionId) != 0{ //SectionID from form and from user authenti
    //cation are not same. Display bad page
    log.Println("Section IDs don't match in usersession and formsession.")
    displayBadPage(w,r,errors.New("SectionIDs don't match. Please submit form properly"))
    return
  }
  //Validate entered form session.
  formSessionSectionId,err := formsession.ValidateSession(template.HTMLEscapeString(r.Form.Get("formsid")))
  if err != nil{
    log.Println("Error validating formSid:",err)
    displayBadPage(w,r,err)
    return
  }

  if strings.Compare(*formSessionSectionId,sectionId) != 0{ //SectionID from form session and from user authenti
    //cation are not same. Display bad page
    displayBadPage(w,r,errors.New("Bad Form used to submit feedback"))
  }
  //Form is valid. Continue with deleting the form session.
  log.Println("Form session and section IDs are valid.")
  formSidDeleted,err := formsession.DeleteSession(template.HTMLEscapeString(r.Form.Get("formsid")))
  log.Println("Deleted",formSidDeleted,"form sessions for sid",r.Form.Get("formsid"),".")

  mySection,_ := section.GetSection(sectionId)
  myTeachers := mySection.Teachers
  myQuestions := &question.GlobalQuestions

  var myFeedback Feedback
  myFeedback.Sectionid = sectionId
  myFeedback.Ratings = make([]Rating,len(*myTeachers))

  for i,myTeacher := range *myTeachers{
    myFeedback.Ratings[i].Facultyid = myTeacher.Facultyid
    myFeedback.Ratings[i].Subjectid = myTeacher.Subjectid
    myFeedback.Ratings[i].Points = make([]Point,len(*myQuestions))
    for qno,myQuestion := range *myQuestions{
      myFeedback.Ratings[i].Points[qno].Questionid = myQuestion.Questionid
      if template.HTMLEscapeString(r.Form.Get(myTeacher.Facultyid + "~" + myTeacher.Subjectid + "~" + myQuestion.Questionid)) == ""{
        displayBadPage(w,r,errors.New("No feedback found for teacher ID: "+myTeacher.Facultyid+",subject ID: "+myTeacher.Subjectid+",question ID: "+myQuestion.Questionid))
        return
      }
      myFeedback.Ratings[i].Points[qno].Marks,err = strconv.Atoi(template.HTMLEscapeString(r.Form.Get(myTeacher.Facultyid + "~" + myTeacher.Subjectid + "~" + myQuestion.Questionid)))
      if err != nil{
        displayBadPage(w,r,err)
        return
      }
    }
  }
  myFeedback.Addedon = time.Now().Format("2006-01-02 15:04:05")
  log.Println("Successfully Parsed Feedback Form. Inserting Feedback.")
  err = database.FeedbackCollection.Insert(myFeedback)
  if err != nil{
    log.Println("Error inserting to database:",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Successfully inserted to database. Logging user out.")
  http.Redirect(w, r, "/logout", http.StatusSeeOther)
}

func displayFeedbackForm(w http.ResponseWriter, r *http.Request,sectionId string){
  log.Println("**DISPLAY FEEDBACK FORM**")
  log.Println("Creating feedback form to user:",r.RemoteAddr)
  log.Println("Fetching feedback form information for client",r.RemoteAddr)  //Enter client ip address and new form SID.
  mySection,err := section.GetSection(sectionId)
  if err != nil{
    displayBadPage(w,r,errors.New("Bad Section ID"))
    return
  }
  //New session being created with rid = mySection.Sectionid, so that when validating form submission,
  //we can also validate that the request received is for the correct section.
  formSid, err := formsession.CreateSession(mySection.Sectionid,time.Minute*30) //Form created will be valid for 30 minutes.
  if err != nil{
    log.Println("Error creating new session for feedback form:",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Successfully created feedback formSession for section ID:",mySection.Sectionid)
  var myFeedbackPage FeedbackPage
  myFeedbackPage.FormSid = formSid
  myFeedbackPage.Sectionid = &mySection.Sectionid
  myFeedbackPage.Sectionname = &mySection.Sectionname

  myCourse,err := course.GetCourse(mySection.Courseid)

  if err != nil || myCourse == nil{
    displayBadPage(w,r,errors.New("Bad Course ID " + mySection.Courseid + " for given Section ID:" + err.Error()))
    return
  }
  log.Println("Getting CourseName")
  myFeedbackPage.Coursename = &myCourse.Coursename
  log.Println("Successfully got CourseName")

  myFeedbackPage.DetailedTeachers,err = section.GetDetailedTeachers(mySection.Teachers)

  if err != nil{
    displayBadPage(w,r,err )
    return
  }

  myFeedbackPage.Questions = &question.GlobalQuestions

  log.Println("Success generating feedback information for client:",r.RemoteAddr,". Now displaying Feedback Form.")
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
  log.Println("User Session authentic")
  if r.Method == "GET"{
    log.Println("Request method is GET")
    displayFeedbackForm(w,r,*sectionId)
    return
  }
  //method = POST
  log.Println("Request method is POST")
  parseFeedbackForm(w,r,*sectionId)
  /*r.ParseForm()

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
  return*/
}

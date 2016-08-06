package college

import(
  "os"
  "encoding/json"
)

type CollegeDetails struct{
  Collegename string  `bson:"collegename"`
}

var GlobalDetails CollegeDetails

func InitCollegeDetails() error{
  //Sharing common resources with Feedback-admin.
  myFile, err := os.Open("feedbackadminres/faconfig.json")
  defer myFile.Close()

  if err != nil{
    return err
  }
  err = json.NewDecoder(myFile).Decode(&GlobalDetails)
  if err != nil{
    return err
  }
  return nil
}

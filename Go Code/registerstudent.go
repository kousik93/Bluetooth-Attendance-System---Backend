package main

import (
	"github.com/drone/routes"
	"log"
    "encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type UUID struct{
	Uuids []string `json:"uuids"`
}


type RowReturn struct{
	Id string `json:"id"`
	Key int `json:"key"`
	Value string `json:"value"`
}

type HasStudentRegistered struct{
	Total_rows int `json:"total_rows"`
	Offset int `json:"offset"`
	Rows []struct{
			Id string `json:"id"`
			Key int `json:"key"`
			Value string `json:"value"`
		}`json:"rows"`
}

type ValidStudent struct{

	Ans string `json:"ans"`
}


type Student struct{
	StudentId int `json:"studentid"`
	Password string `json:"password"`
}

type StudentFull struct{
	StudentId int `json:"studentid"`
	Password string `json:"password"`
	Rev string `json:"_rev"`
	Id string `json:"_id"`
}

var uniqueid UUID
var hasStudentRegistered HasStudentRegistered
var validStudent ValidStudent
var BaseUrl string
var student Student
var studentFull StudentFull
var studentPasswordData HasStudentRegistered

func getUUID() string{
	response,_ := http.Get("https://couchdb-80f683.smileupps.com/_uuids")
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&uniqueid)
	if err != nil {
		panic(err)
	}
	return uniqueid.Uuids[0]
}



func doGetRegister(id string){
	var Url string
	Url=BaseUrl+"/studentprofile/_design/studentdetails/_view/studentregistered?key=+"+string(id)
	response,_ := http.Get(Url)
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&hasStudentRegistered)
	if err != nil {
		panic(err)
	}
}


func doGetPassword(id string){
	var Url string
	Url=BaseUrl+"/studentprofile/_design/studentdetails/_view/studentpassword?key="+string(id)
	response,_ := http.Get(Url)
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&studentPasswordData)
	if err != nil {
		panic(err)
	}
}


func doGetEnrolled(id string){
	response,_ := http.Get("http://localhost:3000/checkstudentvalid/"+id)
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&validStudent)
	if err != nil {
		panic(err)
	}

}

func doPut(body string,uuid string){
    Url:= BaseUrl+"/studentprofile/"+uuid
	request, _ := http.NewRequest("PUT", Url, strings.NewReader(body))
	client := &http.Client{}
	client.Do(request)
}

func doPutAttendance(classid string,body string,uuid string){
    Url:= BaseUrl+"/class"+classid+"/"+uuid
	request, _ := http.NewRequest("PUT", Url, strings.NewReader(body))
	client := &http.Client{}
	log.Println(Url)
	client.Do(request)
}

func doDelete(id string){
	response,_ := http.Get("https://couchdb-80f683.smileupps.com/studentprofile/"+id)
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&studentFull)

	Url:= BaseUrl+"/studentprofile/"+id+"?rev="+studentFull.Rev
	request, _ := http.NewRequest("DELETE", Url, nil)
	client := &http.Client{}
	client.Do(request)
}

func RegisterStudent(rw http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	pass:= params.Get(":pass")
	doGetRegister(id)
	if len(hasStudentRegistered.Rows)>0 {
		a:=`{"error":"Already Registered"}`
		rw.Write([]byte(a))
	}

	if len(hasStudentRegistered.Rows)==0 {
		doGetEnrolled(id)
		if validStudent.Ans=="yes"{
			//To-Do Do JSON Unmarshall
			student.StudentId,_=strconv.Atoi(id)
			student.Password=pass
			a, _ := json.Marshal(student)
			uuid:=getUUID()
			doPut(string([]byte(a)),uuid)
			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte(`{"deviceid":"`+uuid+`"}`))
		}
		if validStudent.Ans=="no" {
			a:=`{"error":"Not Valid Student"}`
			rw.Write([]byte(a))
		}	
	}
}

func DeleteStudent(rw http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	pass:= params.Get(":pass")
	doGetPassword(id)
	if len(studentPasswordData.Rows)==0 {
		a:=`{"error":"Not Exist"}`
		rw.Write([]byte(a))
	}
	if len(studentPasswordData.Rows)>0 {
		if studentPasswordData.Rows[0].Value==pass {
			doDelete(studentPasswordData.Rows[0].Id)
		}
	}

}


func MarkPresent(rw http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	deviceid:= params.Get(":deviceid")
	classid:= params.Get(":classid")
	doGetRegister(id)
	if len(hasStudentRegistered.Rows)>0 {
		if deviceid==hasStudentRegistered.Rows[0].Id {
			//Do duplication Checks
			log.Println("comes here")
			doPutAttendance(classid,`{"studentid":`+id+`}`,getUUID())
		}
	}
	if len(hasStudentRegistered.Rows)==0 {
		a:=`{"error":"Student not registered"}`
		rw.Write([]byte(a))		
	}
}


func main(){

	BaseUrl="https://admin:9631aa6374e6@couchdb-80f683.smileupps.com"

	//REST Config begins
			mux := routes.New()
			//mux.Get("/studentname/:id", GetStudentName)
			mux.Post("/registerstudent/:id/:pass",RegisterStudent)
			mux.Post("/markpresent/:id/:deviceid/:classid",MarkPresent)
			mux.Del("/deletestudent/:id/:pass",DeleteStudent)
			http.Handle("/", mux)
			log.Println("REST has been set up: "+strconv.Itoa(3001))
			log.Println("Listening...")
			http.ListenAndServe(":"+strconv.Itoa(3001), nil)
	//REST Config end
}
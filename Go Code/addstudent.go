package main

import (
	"github.com/drone/routes"
	"log"
    "encoding/json"
	"net/http"
 	"strings"
	"strconv"
)

type UUID struct{
	Uuids []string `json:"uuids"`
}

type Student struct{
	StudentId int `json:"studentid"`
	RegClasses []int `json:"regclasses"`
	StudentName string `json:"studentname"`
}

type RowReturn struct{
	Id string `json:"id"`
	Key int `json:"key"`
	Value string `json:"value"`
}

type RowReturnEnrolled struct{
	Id string `json:"id"`
	Key int `json:"key"`
	Value []int `json:"value"`
}

type StudentNameAPI struct{
	Total_rows int `json:"total_rows"`
	Offset int `json:"offset"`
	Rows []struct{
			Id string `json:"id"`
			Key int `json:"key"`
			Value string `json:"value"`
		}`json:"rows"`
}

type StudentEnrolledAPI struct{
	Total_rows int `json:"total_rows"`
	Offset int `json:"offset"`
	Rows []struct{
			Id string `json:"id"`
			Key int `json:"key"`
			Value []int `json:"value"`
		}`json:"rows"`
}


var uniqueid UUID
var returnStudent StudentNameAPI
var returnStudentenrolled StudentEnrolledAPI
var BaseUrl string
var rowreturn RowReturn
var rowreturnenrolled RowReturnEnrolled

//Get Unique UUID
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


//Helper Put funtion for new student insert
func doPut(body string,uuid string){
    Url:= BaseUrl+"/studentlist/"+uuid
	request, _ := http.NewRequest("PUT", Url, strings.NewReader(body))
	client := &http.Client{}
	client.Do(request)
}

//Helper GET Function for Student Names
func doGetName(id string){
	var Url string
	if id=="" {
		Url=BaseUrl+"/studentlist/_design/getlistdata/_view/studentname"
	}
	if id!="" {
	Url=BaseUrl+"/studentlist/_design/getlistdata/_view/studentname?key=+"+string(id)
	}
	response,_ := http.Get(Url)
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&returnStudent)
	if err != nil {
		panic(err)
	}
}


// Helper GET Function for Student Enrolled
func doGetEnroll(id string){
	var Url string
	if id=="" {
		Url=BaseUrl+"/studentlist/_design/getlistdata/_view/studentenrolled"
	}
	if id!="" {
	Url=BaseUrl+"/studentlist/_design/getlistdata/_view/studentenrolled?key=+"+string(id)
	}
	response,_ := http.Get(Url)
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&returnStudentenrolled)
	if err != nil {
		panic(err)
	}
}


//Main ADD Student
func AddStudent(rw http.ResponseWriter, r *http.Request) {
	var student Student
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&student)
	if err != nil {
		panic(err)
	}
	a, _ := json.Marshal(student)
	doPut(string([]byte(a)),getUUID())
	log.Println("Done")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(""))
	
}


//Main StudentName Handle Func
func GetStudentName(rw http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	doGetName(id)
	rowreturn.Value= returnStudent.Rows[0].Value
	rowreturn.Key= returnStudent.Rows[0].Key
	rowreturn.Id= returnStudent.Rows[0].Id
	a, _ := json.Marshal(rowreturn)
	log.Println("Done")
	rw.Write([]byte(a))
}

//Main Give all student name function
func GetAllStudent(rw http.ResponseWriter, r *http.Request) {
	doGetName("")
	a, _ := json.Marshal(returnStudent)
	log.Println("Done")
	rw.Write([]byte(a))
}


//Main Get Enrolled Details for StudentId
func GetStudentEnrolled(rw http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	doGetEnroll(id)
	rowreturnenrolled.Value= returnStudentenrolled.Rows[0].Value
	rowreturnenrolled.Key= returnStudentenrolled.Rows[0].Key
	rowreturnenrolled.Id= returnStudentenrolled.Rows[0].Id
	a, _ := json.Marshal(rowreturnenrolled)
	log.Println("Done")
	rw.Write([]byte(a))
}


//Main Check Student Valid Function
func CheckStudentValid(rw http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	doGetName(id)
	var a string
	if len(returnStudent.Rows)>0 {
		a=`{"ans":"yes"}`
	}

	if len(returnStudent.Rows)==0 {
		a=`{"ans":"no"}`
	}
	log.Println("Done")
	rw.Write([]byte(a))
}


func main(){

	BaseUrl="https://admin:9631aa6374e6@couchdb-80f683.smileupps.com"

	//REST Config begins
			mux := routes.New()
			mux.Get("/studentname/:id", GetStudentName)
			mux.Get("/checkstudentvalid/:id", CheckStudentValid)
			mux.Get("/studentenrolled/:id", GetStudentEnrolled)
			mux.Get("/allstudent", GetAllStudent)
			mux.Post("/addstudent",AddStudent)
			http.Handle("/", mux)
			log.Println("REST has been set up: "+strconv.Itoa(3000))
			log.Println("Listening...")
			http.ListenAndServe(":"+strconv.Itoa(3000), nil)
	//REST Config end
}
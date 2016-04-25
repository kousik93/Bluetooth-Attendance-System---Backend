package main

import (
	"github.com/drone/routes"
	"log"
    "encoding/json"
	"net/http"
	"strconv"
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


var uniqueid UUID
var hasStudentRegistered HasStudentRegistered
var validStudent ValidStudent
var BaseUrl string

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

func doGetEnrolled(id string){
	response,_ := http.Get("http://localhost:3000/checkstudentvalid/"+id)
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&validStudent)
	if err != nil {
		panic(err)
	}

}



func RegisterStudent(rw http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	//pass:= params.Get(":pass")
	doGetRegister(id)
	if len(hasStudentRegistered.Rows)>0 {
		a:=`{"error":"Already Registered"}`
		rw.Write([]byte(a))
	}

	if len(hasStudentRegistered.Rows)==0 {
		doGetEnrolled(id)
		if validStudent.Ans=="yes"{
			log.Println("Student is valid")
		}
		if validStudent.Ans=="no" {
			a:=`{"error":"Not Valid Student"}`
			rw.Write([]byte(a))
		}

	//log.Println("Done")	
	}
}


func main(){

	BaseUrl="https://admin:9631aa6374e6@couchdb-80f683.smileupps.com"

	//REST Config begins
			mux := routes.New()
			//mux.Get("/studentname/:id", GetStudentName)
			mux.Post("/registerstudent/:id/:pass",RegisterStudent)
			http.Handle("/", mux)
			log.Println("REST has been set up: "+strconv.Itoa(3001))
			log.Println("Listening...")
			http.ListenAndServe(":"+strconv.Itoa(3001), nil)
	//REST Config end
}
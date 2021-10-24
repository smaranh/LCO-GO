package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//model for course - file
type Course struct {
	CourseId    string  `json:"courseid`
	CourseName  string  `json:"coursename`
	CoursePrice int     `json:"courseprice"`
	Author      *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname`
	Website  string `json:"website`
}

//fake DB
var courses []Course

//middleware
func (c *Course) IsEmpty() bool {
	return c.CourseName == ""
}

func main() {
	r := mux.NewRouter()

	//seeding
	courses = append(courses, Course{
		CourseId:    "2",
		CourseName:  "ProtoBuf",
		CoursePrice: 140,
		Author: &Author{
			Fullname: "LCO",
			Website:  "lco.dev",
		},
	})
	courses = append(courses, Course{
		CourseId:    "4",
		CourseName:  "gRPC",
		CoursePrice: 240,
		Author: &Author{
			Fullname: "Sam",
			Website:  "smaranharihar.dev",
		},
	})

	//routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{courseid}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/course/{courseid}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/course/{courseid}", deleteOneCourse).Methods("DELETE")
	r.HandleFunc("/courses", deleteAllCourses).Methods("DELETE")

	fmt.Println("Listening on port 4000")
	log.Fatal(http.ListenAndServe(":4000", r))
}

//controllers

// serve home route

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to Go API</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one Course")
	w.Header().Set("Content-Type", "application/json")

	//grab id from request
	params := mux.Vars(r)
	requestedId := params["courseid"]

	// loop through courses, find matching id and return the response
	for _, course := range courses {
		if course.CourseId == requestedId {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	responseMsg := fmt.Sprintf("No course found for CourseId: %v\n", string(requestedId))

	json.NewEncoder(w).Encode(responseMsg)
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create one Course")
	w.Header().Set("Content-Type", "application/json")

	// what if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("No data sent")
		return
	}

	// what if data is {}
	var newCourse Course
	json.NewDecoder(r.Body).Decode(&newCourse)
	if newCourse.IsEmpty() {
		json.NewEncoder(w).Encode("Data is empty")
		return
	}

	for _, course := range courses {
		if course.CourseName == newCourse.CourseName {
			json.NewEncoder(w).Encode("This course already exists")
			return
		}
	}

	//generate unique id, string
	// append course into courses

	rand.Seed(time.Now().UnixNano())
	newCourse.CourseId = strconv.Itoa(rand.Intn(100))
	courses = append(courses, newCourse)
	json.NewEncoder(w).Encode(fmt.Sprintf("Course: %s has been created", newCourse.CourseName))
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one Course")
	w.Header().Set("Content-Type", "application/json")

	// first grad the id
	params := mux.Vars(r)

	//loop through courses, remove from courses and add the updated course value to courses
	for index, course := range courses {
		if course.CourseId == params["courseid"] {
			courses = append(courses[:index], courses[index+1:]...)
			var updatedCourse Course
			json.NewDecoder(r.Body).Decode(&updatedCourse)
			updatedCourse.CourseId = params["courseid"]
			courses = append(courses, updatedCourse)
			json.NewEncoder(w).Encode("Course has been updated")
			return
		}
	}
	json.NewEncoder(w).Encode("Course does not exist")
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one Course")
	w.Header().Set("Content-Type", "application/json")

	// first grad the id
	params := mux.Vars(r)

	//loop through courses, remove from courses and add the updated course value to courses
	for index, course := range courses {
		if course.CourseId == params["courseid"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode(fmt.Sprintf("Course id: %s has been deleted", params["courseid"]))
			return
		}
	}
	json.NewEncoder(w).Encode("Course does not exist")
}

func deleteAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete all Courses")
	w.Header().Set("Content-Type", "application/json")

	courses = nil
	json.NewEncoder(w).Encode("All the courses have been deleted")
}

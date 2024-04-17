package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Student struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	ContactNumber string `json:"contact_number"`
}

var students []Student

var db *sql.DB

func init() {
	// Open a database connection
	var err error
	db, err = sql.Open("mysql", "dev:dev123@tcp(43.204.101.76:3306)/go_project")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	// Check if the database connection is valid
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging the database:", err)
	}
	fmt.Println("Pinging to MySQL is successful")
}
func main() {
	http.HandleFunc("/users", getUSerHandler)
	http.HandleFunc("/users/", getUserByIDHandler)
	http.HandleFunc("/create", createUser)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}
func getUSerHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT * FROM Users")
	if err != nil {
		fmt.Println("Error querying to fetch data", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.ContactNumber)
		if err != nil {
			fmt.Println("Error in scanning", err)
		}
		students = append(students, student)
	}
	defer rows.Close()
	jsonBytes, err := json.Marshal(students)
	if err != nil {
		fmt.Println("Error in marshalling to JSON", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
	students = nil
}
func createUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var student Student
	w.Header().Set("Content-Type", "application/json")
	// err := json.NewDecoder(r.Body).Decode(&student)
	stmt, err := db.Prepare("insert into Users(id,name, email, contact_number) values(?,?,?,?)")
	if err != nil {
		fmt.Println("error in executing query", err)
	}
	result, err := stmt.Exec(student.ID, student.Name, student.Email, student.ContactNumber)
	if err != nil {
		log.Fatal("error in executing the statement", err)
	}
	resultId, err := result.LastInsertId()
	if err != nil {
		log.Fatal("last inseret id error", err)
	}
	fmt.Println(resultId)
	fmt.Println(student.ID, " student id and student name", student.Name)
	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		fmt.Println("error in encoding json", err)
		http.Error(w, "internal server errror", http.StatusInternalServerError)
		return
	}
}
func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	// Parse ID from URL parameter
	id := r.URL.Path[len("/users/"):]
	studentID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// Query the database for the student with the provided ID
	row := db.QueryRow("SELECT * FROM Users WHERE id=?", studentID)
	var student Student
	err = row.Scan(&student.ID, &student.Name, &student.Email, &student.ContactNumber)
	if err == sql.ErrNoRows {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Marshal student data to JSON
	jsonBytes, err := json.Marshal(student)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

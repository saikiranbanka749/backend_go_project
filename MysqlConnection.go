package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var students []Student
var db *sql.DB

func init() {
	// Open a database connection
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/go_project")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	// Check if the database connection is valid
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging the database:", err)
	}
	fmt.Println("pinging to mysql is successfull")
}
func main() {
	http.HandleFunc("/users", getUSerHandler)
	http.HandleFunc("/create", createUser)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
func getUSerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(students)
	rows, err := db.Query("select * from Student")
	if err != nil {
		fmt.Println("error quering in fetching data", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.ID, &student.Name)
		if err != nil {
			fmt.Println("error in scanning", err)
		}
		students = append(students, student)
	}
	jsonBytes, err := json.Marshal(students)
	fmt.Println("data at get handler", students)
	if err != nil {
		fmt.Println("error in marshalling  to json", err)
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}
func createUser(w http.ResponseWriter, r *http.Request) {
	var student Student
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&student)
	stmt, err := db.Prepare("insert into student(id,name) values(?,?)")
	if err != nil {
		fmt.Println("error in executing query", err)
	}
	result, err := stmt.Exec(student.ID, student.Name)
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

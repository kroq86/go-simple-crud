package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Posts []Post
}

type Post struct {
	gorm.Model

	Title      string
	Body       string
	UserID     int
}

var db *gorm.DB
var err error

func main() {
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)

	db, err = gorm.Open(dialect, dbURI)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to database successfully")
	}

	defer db.Close()

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})

	router := mux.NewRouter()

	router.HandleFunc("/posts", GetPosts).Methods("GET")
	router.HandleFunc("/post/{id}", GetPost).Methods("GET")
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/user/{id}", GetUser).Methods("GET")

	router.HandleFunc("/create/user", CreateUser).Methods("POST")
	router.HandleFunc("/create/post", CreatePost).Methods("POST")

	router.HandleFunc("/delete/user/{id}", DeleteUser).Methods("DELETE")
	router.HandleFunc("/delete/post/{id}", DeletePost).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}


func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	var posts []Post

	db.First(&user, params["id"])
	db.Model(&user).Related(&posts)

	user.Posts = posts

	json.NewEncoder(w).Encode(&user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var user []User

	db.Find(&user)

	json.NewEncoder(w).Encode(&user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	createdUser := db.Create(&user)
	err = createdUser.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user User

	db.First(&user, params["id"])
	db.Delete(&user)

	json.NewEncoder(w).Encode(&user)
}


func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var post Post

	db.First(&post, params["id"])

	json.NewEncoder(w).Encode(&post)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []Post

	db.Find(&posts)

	json.NewEncoder(w).Encode(&posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	json.NewDecoder(r.Body).Decode(&post)

	createdPost := db.Create(&post)
	err = createdPost.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdPost)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var post Post

	db.First(&post, params["id"])
	db.Delete(&post)

	json.NewEncoder(w).Encode(&post)
}

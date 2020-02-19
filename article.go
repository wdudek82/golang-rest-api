package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	"net/http"
)

type Article struct {
	gorm.Model
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type JsonData struct {
	Message string  `json:"message"`
	Article Article `json:"article"`
}

var db *gorm.DB
var err error
var dbConnectionString = "host=localhost port=5432 user=learning_go dbname=learning_go password=learning_go sslmode=disable"

func InitialArticleMigration() {
	db, err = gorm.Open("postgres", dbConnectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Article{})
}

func handleDbConnectionError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to the database")
	}
}

func AllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: all articles")

	db, err = gorm.Open("postgres", dbConnectionString)
	handleDbConnectionError(err)
	defer db.Close()

	var articles []Article
	db.Find(&articles)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func AddArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: add article")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	article := &Article{}

	err = json.Unmarshal([]byte(data), article)
	if err != nil {
		fmt.Println(err)
	}

	db, err = gorm.Open("postgres", dbConnectionString)
	handleDbConnectionError(err)
	defer db.Close()

	db.Create(article)

	fmt.Println("New user has been added successfully: ", article)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JsonData{Message: "New article has been added.", Article: *article})
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: add article")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	updArticle := &Article{}

	err = json.Unmarshal([]byte(data), updArticle)
	if err != nil {
		fmt.Println(err)
	}

	db, err = gorm.Open("postgres", dbConnectionString)
	handleDbConnectionError(err)
	defer db.Close()

	vars := mux.Vars(r)
	articleId := vars["articleId"]

	var article Article
	db.Where("id = ?", articleId).Find(&article)
	article.Title = updArticle.Title
	article.Desc = updArticle.Desc
	article.Content = updArticle.Content

	db.Save(&article)

	fmt.Println("Article %s has been successfully updated:", article)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		JsonData{Message: fmt.Sprintf("Article id %s has been successfully updated", articleId)},
	)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: delete article")

	db, err = gorm.Open("postgres", dbConnectionString)
	handleDbConnectionError(err)
	defer db.Close()

	vars := mux.Vars(r)
	articleId := vars["articleId"]

	var article Article
	db.Where("id = ?", articleId).Delete(&article)

	fmt.Println("Article %s has been successfully deleted:", article)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		JsonData{Message: fmt.Sprintf("Article id %s has been deleted", articleId)},
	)
}

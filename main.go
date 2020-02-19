package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello, World!</h1>")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	// Generic routes
	myRouter.HandleFunc("/", helloWorld).Methods("GET")

	// Article routing
	myRouter.HandleFunc("/articles", AllArticles).Methods("GET")
	myRouter.HandleFunc("/articles", AddArticles).Methods("POST")
	myRouter.HandleFunc("/articles/{articleId}", UpdateArticle).Methods("PUT")
	myRouter.HandleFunc("/articles/{articleId}", DeleteArticle).Methods("DELETE")

	// Server
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	fmt.Println("Server is running...")

	InitialArticleMigration()
	handleRequests()
	//JsonTutorial()
}

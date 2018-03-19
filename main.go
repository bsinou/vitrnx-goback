package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bsinou/vitrnx-goback/model"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/posts", model.GetArticleByTag)
	myRouter.HandleFunc("/posts/{id}", model.GetSingleArticle)
	myRouter.HandleFunc("/by-tag/{tag}", model.GetArticleByTag)
	log.Fatal(http.ListenAndServe(":7777", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}

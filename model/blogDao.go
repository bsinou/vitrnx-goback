package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
)

type Article struct {
	Id     string    `json:"id"`
	Date   time.Time `json:"date"`
	Author string    `json:"author"`
	Title  string    `json:"title"`
	Tags   string    `json:"tags"`
	Desc   string    `json:"desc"`
	Body   string    `json:"body"`
}

type Articles []Article

var (
	fakeRepo map[string]*Article
)

func init() {
	fakeRepo = *populateFakeRepo()
}

func GetArticleByTag(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GetArticleByTag")

	vars := mux.Vars(r)
	slug := vars["tag"]
	fmt.Println("Get for tag: " + slug)

	var tag *regexp.Regexp
	if slug != "" {
		tag = regexp.MustCompile(slug)
	}

	var results []Article
	for _, v := range fakeRepo {
		if tag == nil {
			results = append(results, *v)
		} else if tag.MatchString(v.Tags) {
			results = append(results, *v)
		}
	}

	// writeCrossDomainHeaders(w, r)
	json.NewEncoder(w).Encode(results)
}

func GetSingleArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetSingleArticle")

	vars := mux.Vars(r)
	slug := vars["id"]
	fmt.Println("Get article for id " + slug)

	art := Article{Id: slug, Title: "Oops", Tags: "Actualités Réflexions", Author: "Marie-Madeleine", Desc: "L'article que vous cherckez n'est pas disponible", Body: "Veuillez revoir votre recherche"}
	if val, ok := fakeRepo[slug]; ok {
		art = *val
	}

	// writeCrossDomainHeaders(w, r)
	json.NewEncoder(w).Encode(art)
}

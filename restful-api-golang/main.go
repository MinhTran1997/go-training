// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// youtube source: https://www.youtube.com/watch?v=SonwZ6MF5BE
// morioh source: https://morioh.com/p/9a0e53da7908?fbclid=IwAR05Tp1czeMZCp1s4MHAkRIqABqd-12cRsb3o71894vTGVlT13G91t9kJFk

// Article - Our struct for all articles
type Article struct {
	Id      string  `json:"Id"`
	Title   string  `json:"title"`
	Desc    string  `json:"desc"`
	Content string  `json:"content"`
	Author  *Author `json:"author"`
}

type Author struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

var Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	article.Id = strconv.Itoa(rand.Intn(1000000))
	// update our global Articles array to include
	// our new Article
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article) // bien cai request body thanh 1 article co dang struct la Article

	//xoa article cu di va thay bang article moi (delete then add new)
	for index, article := range Articles {
		if article.Id == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(Articles)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, article := range Articles {
		if article.Id == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	Articles = []Article{
		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content", Author: &Author{Name: "Chris Nolan", Email: "nolan@gmail.com"}},
		Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content", Author: &Author{Name: "Dave Fincher", Email: "davef@gmail.com"}},
	}
	handleRequests()
}

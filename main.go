package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Repo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stargazers  int    `json:"startgazers_count"`
	Watchers    int    `json:"watchers_count"`
	Forks       int    `json:"forks_count"`
}

func GetUserRepos(userName string) []Repo {
	response, err := http.Get(fmt.Sprintf("https://api.github.com/users/%v/repos", userName))
	if err != nil {
		log.Fatal(err.Error())
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var repos []Repo
	err = json.Unmarshal(responseData, &repos)
	if err != nil {
		log.Fatal(err)
	}
	return repos
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/users/"):]
	if len(path) != 0 {
		tmpl := template.Must(template.ParseFiles("template/user.html"))
		tmpl.Execute(w, struct {
			User string
			Data []Repo
		}{path, GetUserRepos(path)})
	} else {
		fmt.Fprintf(w, "<h1>Not found!</h1><a href='/'>Go back</a>")
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/users/", handleUsers)

	fmt.Println("App listening at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

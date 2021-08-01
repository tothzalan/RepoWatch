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

func GetUserRepos(userName string) ([]Repo, error) {
	response, err := http.Get(fmt.Sprintf("https://api.github.com/users/%v/repos", userName))
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var repos []Repo
	err = json.Unmarshal(responseData, &repos)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return repos, nil
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/users/"):]
	if len(path) != 0 {
		tmpl := template.Must(template.ParseFiles("template/user.html"))
		data, err := GetUserRepos(path)
		if err != nil {
			fmt.Fprint(w, "<h1>An error occured</h1>\n<a href='/'>Go back</a>")
			return
		}
		tmpl.Execute(w, struct {
			User string
			Data []Repo
		}{path, data})
	} else {
		fmt.Fprintf(w, "<h1>Not found!</h1>\n<a href='/'>Go back</a>")
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/users/", handleUsers)

	fmt.Println("App listening at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type NameResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type JokeEntry struct {
	Joke string `json:"joke"`
}

type JokeResponse struct {
	Type  string `json:"type"`
	Value JokeEntry
}

func handler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	name_res, err := http.Get("https://names.mcquay.me/api/v0")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if name_res.StatusCode != 200 {
		http.Error(res, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}

	var name NameResponse
	err = json.NewDecoder(name_res.Body).Decode(&name)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	joke_res, err := http.Get(fmt.Sprintf(
		"http://api.icndb.com/jokes/random?firstName=%s&lastName=%s&limitTo=nerdy",
		name.FirstName,
		name.LastName))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if joke_res.StatusCode != 200 {
		http.Error(res, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}

	var joke JokeResponse
	err = json.NewDecoder(joke_res.Body).Decode(&joke)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(res, joke.Value.Joke)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

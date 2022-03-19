package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Name struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type NameResult struct {
	name *Name
	err  error
}

type Joke struct {
	Joke string `json:"joke"`
}

type JokeResponse struct {
	Type  string `json:"type"`
	Value Joke
}

type JokeResult struct {
	joke string
	err  error
}

func fetchName(c chan NameResult) {
	// TODO: Name API is rate limited. Consider caching previous results and reusing cached values in case of errors
	res, err := http.Get("https://names.mcquay.me/api/v0")
	if err != nil {
		c <- NameResult{nil, err}
		return
	}
	if res.StatusCode != 200 {
		c <- NameResult{nil, fmt.Errorf("Failed to fetch random name. Code: %d", res.StatusCode)}
		return
	}

	var name Name
	err = json.NewDecoder(res.Body).Decode(&name)
	if err != nil {
		c <- NameResult{nil, err}
		return
	}

	c <- NameResult{&name, nil}
}

func fetchJoke(c chan JokeResult) {
	res, err := http.Get("http://api.icndb.com/jokes/random?firstName=%s&lastName=%s&limitTo=nerdy")

	if err != nil {
		c <- JokeResult{"", err}
		return
	}
	if res.StatusCode != 200 {
		c <- JokeResult{"", fmt.Errorf("Failed to fetch random joke. Code: %d", res.StatusCode)}
		return
	}

	var joke JokeResponse
	err = json.NewDecoder(res.Body).Decode(&joke)
	if err != nil {
		c <- JokeResult{"", err}
		return
	}

	c <- JokeResult{joke.Value.Joke, nil}
}

func handler(c *gin.Context) {
	name_chan := make(chan NameResult)
	joke_chan := make(chan JokeResult)

	go fetchName(name_chan)
	go fetchJoke(joke_chan)

	nres := <-name_chan
	if nres.err != nil {
		logrus.Error(nres.err.Error())
		c.String(http.StatusBadGateway, nres.err.Error())
		return
	}
	jres := <-joke_chan
	if jres.err != nil {
		logrus.Error(jres.err.Error())
		c.String(http.StatusBadGateway, jres.err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf(jres.joke, nres.name.FirstName, nres.name.LastName))
}

func main() {
	router := gin.Default()
	router.GET("/", handler)
	router.Run(":8080")
}

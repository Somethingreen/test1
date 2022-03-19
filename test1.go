package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func errorCode(c *gin.Context, code int) {
	c.String(code, http.StatusText(code))
}

func handler(c *gin.Context) {

	name_res, err := http.Get("https://names.mcquay.me/api/v0")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if name_res.StatusCode != 200 {
		logrus.Error("Failed to fetch name. code: ", name_res.StatusCode)
		errorCode(c, http.StatusBadGateway)
		return
	}

	var name NameResponse
	err = json.NewDecoder(name_res.Body).Decode(&name)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	joke_res, err := http.Get(fmt.Sprintf(
		"http://api.icndb.com/jokes/random?firstName=%s&lastName=%s&limitTo=nerdy",
		name.FirstName,
		name.LastName))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if joke_res.StatusCode != 200 {
		logrus.Error("Failed to fetch joke. code: ", joke_res.StatusCode)
		errorCode(c, http.StatusBadGateway)
		return
	}

	var joke JokeResponse
	err = json.NewDecoder(joke_res.Body).Decode(&joke)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, joke.Value.Joke)
}

func main() {
	router := gin.Default()
	router.GET("/", handler)
	router.Run(":8080")
}

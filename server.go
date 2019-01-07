package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/json-iterator/go"
)

// Coder represents a programmer who may
// or may not be coding at this very moment.
type Coder struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	LatestHeartbeat string `json:"last_heartbeat"`
	LatestProject   string `json:"last_project"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("WAKATIME_API_KEY")
	url := "https://wakatime.com/api/v1/users/current/teams/0d49a7ce-bbc6-4ca9-916c-57ed9d2b65dd/members?api_key=" + apiKey

	client := http.Client{
		Timeout: time.Second * 60,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "iscoding.live")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	textBytes := []byte(body)
	coders := jsoniter.Get(textBytes, "data", 0).ToString()
	fmt.Println(coders)

	coder := Coder{}
	parseErr := jsoniter.Unmarshal(textBytes, &coder)
	if parseErr != nil {
		fmt.Println(parseErr)
		return
	}

	// TODO: Iterate over output and print each.Coder
	//
	// for c := range coders {
	// 	fmt.Printf("'%s' last seen on '%s'\n", coders[c].Email, coders[c].LastHeartbeat)
	// }

	fmt.Println(coder.Email)
}

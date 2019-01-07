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

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	teamGUID := os.Getenv("WAKATIME_TEAM_GUID")
	url := "https://wakatime.com/api/v1/users/current/teams/" + teamGUID + "/members?api_key=" + apiKey

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

	teamJSON := json.Get([]byte(body)).ToString()

	// Creating the maps for JSON
	m := map[string]interface{}{}

	//coder := Coder{}
	parseErr := json.Unmarshal([]byte(teamJSON), &m)
	if parseErr != nil {
		fmt.Println(parseErr)
		return
	}

	fmt.Println(m)
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/joho/godotenv"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Coder represents a programmer who may
// or may not be coding at this very moment.
type Coder struct {
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	LatestProject   string    `json:"last_project"`
	Timezone        string    `json:"timezone"`
	LatestHeartbeat time.Time `json:"last_heartbeat"`
	Active          bool
}

// Team contains an array of Coder to represent
// each team member.
type Team struct {
	Coders []Coder `json:"data"`
}

// timeIn converts UTC time to local time.
func timeIn(t time.Time, tzName string) time.Time {
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		panic(err)
	}
	return t.In(loc)
}

func main() {
	// Initialize godotenv for reading secrets stored in .env files.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Grab secrets from .env using godotenv.
	apiKey := os.Getenv("WAKATIME_API_KEY")
	teamGUID := os.Getenv("WAKATIME_TEAM_GUID")
	url := "https://wakatime.com/api/v1/users/current/teams/" + teamGUID + "/members?api_key=" + apiKey

	client := http.Client{
		// If you have a large team in Wakatime, it might take a while to return the data.
		// Set the timeout higher for these requests.
		Timeout: time.Second * 60,
	}

	// Use http.NewRequest when you need to specify attributes
	// of the request. Example: setting custom headers.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set a custom header for User-Agent so the
	// Wakatime API knows who made the request.
	req.Header.Set("User-Agent", "iscoding.live")

	// Execute the GET request.
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// Response data is available in response.Body
	// Read until error or EOF.
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// Create a new Team to store the data we're interested in.
	team := &Team{}

	// Unmarshal the JSON data to populate the team object.
	parseErr := json.Unmarshal([]byte(body), &team)

	if parseErr != nil {
		fmt.Println(parseErr)
		return
	}

	// Set up the CLI output to look nice.
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 15, 8, 8, '\t', 0)

	separator := strings.Repeat("-", 15)
	fmt.Println()
	fmt.Fprintln(w, "Coder\t", "Current Project\t", "Last Seen")
	fmt.Fprintln(w, separator, "\t", separator, "\t", separator)

	// Print the resulting object.
	// Convert LatestHeartbeat to the current user's timezone.
	for _, coder := range team.Coders {
		// Provide a default value for TZ if none is provided.
		if coder.Timezone == "" {
			coder.Timezone = "America/Los_Angeles"
		}

		// Convert UTC heartbeat into the user's local time.
		coder.LatestHeartbeat = timeIn(coder.LatestHeartbeat, coder.Timezone)
		now := timeIn(time.Now(), coder.Timezone)

		// Diff the time from localized time.Now() and update the Coder's status.
		secondsDiff := math.Abs(now.Sub(coder.LatestHeartbeat).Seconds())
		activityTimeout := 120.0
		isActive := secondsDiff <= activityTimeout
		coder.Active = isActive

		// Set a 2 minute timeout window for activity.
		if isActive {
			fmt.Fprintln(w, coder.Email, "\t", coder.LatestProject, "\t", int(secondsDiff), "seconds ago")
		}
	}

	fmt.Fprintln(w, separator, "\t", separator, "\t", separator)
	w.Flush()
}


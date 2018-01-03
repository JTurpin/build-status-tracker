package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

// BuildArtifact Struct used for keeping track of build entity
type BuildArtifact struct {
	Name            string `validate:"nonzero"`
	BuildVersion    string `validate:"nonzero"`
	BuildPromoted   bool
	Description     string
	LastBuildStatus string `validate:"nonzero"`
	LastBuild       time.Time
}

// ArtifactsList Used to contain a slice of build artifacts
type ArtifactsList struct {
	Builds []BuildArtifact
}

var db *bolt.DB

func init() {

}

func main() {
	var err error
	db, err = setupDB()
	if err != nil {
		fmt.Println("DB setup failed")
	}
	//t := time.Now()

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Here's all the enpoints!
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/update", updateArtifactHandler)
	http.HandleFunc("/delete", deleteArtifactHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/backup", backupHandler)

	http.ListenAndServe(":7080", nil)
}

// Helper Functions

// GetFormattedTimeFromEpoch gets time in a nice format from int64
func GetFormattedTimeFromEpoch(z int64, zone string, style int) string {
	var x string
	secondsSinceEpoch := z
	unixTime := time.Unix(secondsSinceEpoch, 0)
	timeZoneLocation, err := time.LoadLocation(zone)
	if err != nil {
		fmt.Println("Error loading timezone:", err)
	}

	timeInZone := unixTime.In(timeZoneLocation)

	switch style {
	case 1:
		timeInZoneStyleOne := timeInZone.Format("Mon Jan 2 15:04:05")
		//Mon Aug 14 13:36:02
		return timeInZoneStyleOne
	case 2:
		timeInZoneStyleTwo := timeInZone.Format("02-01-2006 15:04:05")
		//14-08-2017 13:36:02
		return timeInZoneStyleTwo
	case 3:
		timeInZoneStyleThree := timeInZone.Format("2006-02-01 15:04:05")
		//2017-14-08 13:36:02
		return timeInZoneStyleThree
	}
	return x
}

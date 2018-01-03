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
	BuildPromoted   bool   `json:"BuildPromoted"`
	Description     string
	LastBuildStatus string `validate:"nonzero"`
	LastBuild       time.Time
	CallbackURL     string
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

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Here's all the enpoints!
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/update", updateArtifactHandler)
	http.HandleFunc("/delete", deleteArtifactHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/help", helpHandler)
	http.HandleFunc("/backup", backupHandler)

	http.ListenAndServe(":7080", nil)
}

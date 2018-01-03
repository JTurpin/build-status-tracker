package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	validator "gopkg.in/validator.v2"
)

func rootHandler(rw http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("assets/html/index.html"))
	rw.WriteHeader(http.StatusOK)
	artifactslist := getAllArtifacts(db)
	tmpl.Execute(rw, artifactslist)
}

func updateArtifactHandler(rw http.ResponseWriter, request *http.Request) {
	// Example Request:
	// curl -vvv -X POST -d '{"Name":"autobot","Description":"","BuildVersion": "1","LastBuildStatus":"Fail","LastBuild":"1514002902"}' http://localhost:7080/update
	rw.WriteHeader(http.StatusCreated)
	decoder := json.NewDecoder(request.Body)
	var t BuildArtifact
	//now := time.Now()

	err := decoder.Decode(&t)
	fmt.Println(t.LastBuild)

	// Set build time as when it gets posted
	now := time.Now()
	t.LastBuild = now

	// Get name to lower
	t.Name = strings.ToLower(t.Name)

	if errs := validator.Validate(t); errs != nil {
		// values not valid, deal with errors here
		log.Println("There were errors in validating the request.")
	} else {
		updateDBArtifact(db, t)
		if err != nil {
			panic(err)
		}
	}

}

func deleteArtifactHandler(rw http.ResponseWriter, request *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("200 - Things seem to be running!\n"))
}

func healthHandler(rw http.ResponseWriter, request *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Health OK\n"))
}

func backupHandler(w http.ResponseWriter, req *http.Request) {
	// This is invoked like this: curl http://localhost/backup > my.db
	err := db.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

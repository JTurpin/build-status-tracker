package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	validator "gopkg.in/validator.v2"
)

func rootHandler(rw http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("assets/html/index.html"))
	rw.WriteHeader(http.StatusOK)
	artifactslist := getAllArtifacts(db)
	/*
		for l := range artifactslist {
			log.Printf("\n Name: %v \n Description: %v \n BuildStatus: %v \n LastBuild: %v", artifactslist[l].Name, artifactslist[l].Description, artifactslist[l].LastBuildStatus, artifactslist[l].LastBuild)
		}
	*/
	tmpl.Execute(rw, artifactslist)
}

func updateArtifactHandler(rw http.ResponseWriter, request *http.Request) {
	// Example Request:
	// curl -vvv -X POST -d '{"Name":"autobot","Description":"","BuildVersion": "1","BuildPromoted": "true","LastBuildStatus":"Fail","LastBuild":"1514002902"}' http://localhost:7080/update
	// LastBuild can be left blank and it will use current system time
	// otherwise use epoch time in seconds
	rw.WriteHeader(http.StatusCreated)
	decoder := json.NewDecoder(request.Body)
	var t BuildArtifact
	if t.LastBuild == 0 {
		now := time.Now()
		t.LastBuild = now.Unix()
	}
	err := decoder.Decode(&t)
	if errs := validator.Validate(t); errs != nil {
		// values not valid, deal with errors here
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

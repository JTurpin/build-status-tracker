package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func getAllArtifacts(db *bolt.DB) []BuildArtifact {
	var artifactList []BuildArtifact

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("BuildArtifacts"))
		bucket.ForEach(func(k, v []byte) error {
			var build BuildArtifact
			// Debug whast's in the bucket
			//log.Printf("%s", v)

			// unmarshal json into struct
			json.Unmarshal(v, &build)

			// add each item to the struct
			artifactList = append(artifactList, build)

			return nil
		})
		return nil
	})
	return artifactList
}

func updateDBArtifact(db *bolt.DB, buildArt BuildArtifact) error {
	// Store the artifact model in the BuildArtifacts bucket using the name as the key.
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("BuildArtifacts"))
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(buildArt)
		if err != nil {
			return err
		}
		return b.Put([]byte(buildArt.Name), encoded)
	})
	return err
}

func deleteDBArtifact(db *bolt.DB, buildArt BuildArtifact) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("BuildArtifacts"))
		err := bucket.Delete([]byte(buildArt.Name))
		if err != nil {
			return err
		}
		log.Println("Deleted Artifact: " + buildArt.Name)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("BuildArtifacts"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		if root == nil {
			return fmt.Errorf("Bucket %q not found", []byte("BuildArtifacts"))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	fmt.Println("DB Setup Done")
	return db, nil
}

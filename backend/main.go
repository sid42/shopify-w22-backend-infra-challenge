package main

import (
	"net/http"
	"log"
	// "fmt"
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main () { 
	log.Print("attempting to start image server...")
	
	server := Server{}

	// s3 initialization 

	// db initialization
	opt, err := pg.ParseURL("postgres://postgres:very-secret-db-password@image-db:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	server.db = pg.Connect(opt)

	if err := server.db.Ping(context.Background()); err != nil {
		log.Fatalf("failed to connect to db %s", err)
	} else {
		log.Print("db connection established!")
	}
	
	// routes
	server.r = mux.NewRouter()

	server.r.HandleFunc("/login", server.Login).Methods("POST")
	server.r.HandleFunc("/signup", server.Signup).Methods("POST")
	server.r.HandleFunc("/addImages", server.AddImages).Methods("POST")
	server.r.HandleFunc("/deleteImage", server.DeleteImage).Methods("POST")
	server.r.HandleFunc("/search", server.SearchImage).Methods("GET")

	log.Print("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", server.r))
}
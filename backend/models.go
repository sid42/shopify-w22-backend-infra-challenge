package main

import (
	"time"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Server struct {
	r				*mux.Router
	db				*pg.DB
	s3_session		*s3.S3
	s3_bucket	 	string
}

type User struct {
	tableName struct{} `pg:"users"`

	Id 			string `json:"-" pg:",pk" `
	Email		string `json:"email"`	
	Password	string `json:"password"`
}

type Image struct {
	tableName struct{} `pg:"images"`

	Id 			string `json:"id" pg:",pk" `
	Key		 	string `json:"-"`
	UserEmail	string `json:"-"`	
	CreatedAt 	time.Time `json:"-"`
}

type ImageIds struct {
	Ids		[]string `json:"ids"`
}

type SearchImageRequest struct {
	ByEmails	[]string `json:"by_emails"`
}

type Token struct {
	Token	string `json:"token"`
}

package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Server struct {
	r	*mux.Router
	db	*pg.DB
	s	*session.Session
}

type User struct {
	tableName struct{} `pg:"users,alias:u"`

	Id 			string `json:"-" pg:",pk" `
	Email		string `json:"email"`	
	Password	string `json:"password"`
}

type Token struct {
	Token	string `json:"token"`
}

// type Error struct {
// 	Error 	string `json:"error"`
// }
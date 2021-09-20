package main

import (
	"fmt"
	"log"
	"strings"
	"errors"
	"net/http"
	"encoding/json"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/go-pg/pg/v10"
)

const (
	signing_key = "VerySecureJWTPassphrase"
)

func generateTokenString(user User) string {
	//new token object with user data
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email" : user.Email,
	})
	
	//signing the token to get token string to be sent to user 
	tokenString, tokenErr := token.SignedString([]byte(signing_key))
	if tokenErr != nil {
		panic(tokenErr)
	}

	return tokenString
}

func getTokenFomHeader (r http.Request) (string, error) {
	header := r.Header.Get("x-authentication-token")
	
	if len(header) == 0{
		//no auth token provided
		return "", fmt.Errorf("No auth token provided")
	} else {
		headerSplit := strings.Split(header, " ")
		if len(headerSplit) == 1 {
			// either 'bearer' or token provided, invalid usage
			return "", fmt.Errorf("Invalid usage")
		} else if headerSplit[0] != "Bearer" {
			return "", fmt.Errorf("Invalid usage")
		} else {
			return headerSplit[1], nil
		}
	}
}

func (s *Server) Signup (w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL)

	var u User
	_ = json.NewDecoder(r.Body).Decode(&u)
	
	//check if user with provided email already exists
	var temp_user User
	err := s.db.Model(&temp_user).
	Where("email = ?", u.Email).
	Select()

	if errors.Is(err, pg.ErrNoRows) && u.Email != "" {
		_, err = s.db.Model(&u).Insert()
		if err != nil {
			log.Fatalf("error inserting user in database %s", err)
		}

		resp := Token{Token: generateTokenString(u)}
		json.NewEncoder(w).Encode(resp)
	
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User with provided e-mail already exists or is invalid"))
	}
}

func (s *Server) Login (w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL)

	var u User
	_ = json.NewDecoder(r.Body).Decode(&u)

	var res User
	err := s.db.Model(&res).
	Where("email = ?", u.Email).
	Where("password = ?", u.Password).
	Select()

	if err != nil || res == (User{}){
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect login credentials"))
	} else {
		tokenString := generateTokenString(res)
		json.NewEncoder(w).Encode(Token{Token: tokenString})
	}

}

func (s *Server) AddImages (w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL)
}


func (s *Server) DeleteImage (w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL)
}

func (s *Server) SearchImage (w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL)
}


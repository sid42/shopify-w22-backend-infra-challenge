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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dchest/uniuri"
)

const (
	signing_key = "very_secure_jwt_passphrase"
)

func generateTokenString(user User) string {
	//new token object with user data
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email" : user.Email,
	})
	
	//signing the token to get token string to be sent to user 
	encoded_token, tokenErr := token.SignedString([]byte(signing_key))
	if tokenErr != nil {
		panic(tokenErr)
	}

	return encoded_token
}

func getTokenFomHeader(r http.Request) (string, error) {
	header := r.Header.Get("x-authentication-token")
	
	if len(header) == 0{
		//no auth token provided
		return "", fmt.Errorf("No auth token provided")
	} else {
		header_split := strings.Split(header, " ")
		if len(header_split) == 1 {
			// either 'bearer' or token provided, invalid usage
			return "", fmt.Errorf("Invalid usage")
		} else if header_split[0] != "Bearer" {
			return "", fmt.Errorf("Invalid usage")
		} else {
			return header_split[1], nil
		}
	}
}

func checkSigningMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Error occured")
	}
	return []byte(signing_key), nil
}

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	var u User
	_ = json.NewDecoder(r.Body).Decode(&u)
	
	//check if user with provided email already exists
	var res User
	err := s.db.Model(&res).
	Where("email = ?", u.Email).
	Select()

	if errors.Is(err, pg.ErrNoRows) && u.Email != "" {
		_, err = s.db.Model(&u).Insert()
		if err != nil {
			log.Fatalf("error inserting user in database %s", err)
		}

		resp := Token{Token: generateTokenString(u)}
		log.Print("successfully created user account")
		json.NewEncoder(w).Encode(resp)
	
	} else {
		http.Error(w, "User with provided e-mail already exists or is invalid", http.StatusConflict)
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var u User
	_ = json.NewDecoder(r.Body).Decode(&u)

	var res User
	err := s.db.Model(&res).
	Where("email = ?", u.Email).
	Where("password = ?", u.Password).
	Select()

	if err != nil || res == (User{}){
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
	} else {
		tokenString := generateTokenString(res)
		log.Print("successfully logged in")
		json.NewEncoder(w).Encode(Token{Token: tokenString})
	}

}

func (s *Server) AddImage(w http.ResponseWriter, r *http.Request) {
	var encoded_token string
	var err error
	if encoded_token, err = getTokenFomHeader(*r); err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}

	var token *jwt.Token
	if token, err = jwt.Parse(encoded_token, checkSigningMethod); err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}
	email := token.Claims.(jwt.MapClaims)["Email"].(string)

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// split filename by '.' to obtain file type 
	file_type := strings.Split(header.Filename, ".")
	// generate random url for image to allow multiple uploads of same file
	key := fmt.Sprintf("%s.%s", uniuri.New(), file_type[len(file_type) - 1])
	log.Printf("filename: %s size: %d KB", key, header.Size/1024)
	
	// persist image to bucket
	_, err = s.s3_session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.s3_bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}

	// insert into database
	i := Image {
		Key : key,
		UserEmail : email, 
	}
	_, err = s.db.Model(&i).Returning("id").Insert()
	if err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}

	log.Print("successfully added image")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(i)
}

func (s *Server) DeleteImage(w http.ResponseWriter, r *http.Request) {
	var i Image
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var encoded_token string
	var err error
	if encoded_token, err = getTokenFomHeader(*r); err != nil {
		http.Error(w, "Error deleting image", http.StatusInternalServerError)
		return
	}

	var token *jwt.Token
	if token, err = jwt.Parse(encoded_token, checkSigningMethod); err != nil {
		http.Error(w, "Error deleting image", http.StatusInternalServerError)
		return
	}

	email := token.Claims.(jwt.MapClaims)["Email"].(string)
	if err := s.db.Model(&i).WherePK().Select(); err != nil {
		http.Error(w, "Error deleting image", http.StatusInternalServerError)
		return
	} 

	// no permission to delete another user's image
	if i.UserEmail != email {
		http.Error(w, "Unable to delete another user's image", http.StatusUnauthorized)
		return
	}

	log.Printf("deleting %s", i.Key)

	// delete from s3 bucket
	_, err = s.s3_session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.s3_bucket),
		Key: aws.String(i.Key),
	})
	if err != nil {
		http.Error(w, "Error deleting image", http.StatusInternalServerError)
		return
	}

	// delete from db 
	_, err = s.db.Model(&i).WherePK().Delete()	
	if err != nil {
		http.Error(w, "Error deleting image", http.StatusInternalServerError)
		return
	}

	log.Print("Succesfully deleted image")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Succesfully deleted image")
}

func (s *Server) SearchImages(w http.ResponseWriter, r *http.Request) {
	// by user
	// by latest
	// by tags? 
}

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Print(r.Method, " ", r.RequestURI)

		if r.RequestURI == "/login" || r.RequestURI == "/signup" {
			next.ServeHTTP(w, r)
		} else {
			// authentication to check if user is valid
			encoded_token, err := getTokenFomHeader(*r)

			if err != nil {
				log.Print("error occured in auth middleware ", err)
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				// parse token from jwt encoding
				token, err := jwt.Parse(encoded_token, checkSigningMethod)
			
				if err != nil{	
					log.Print("unauthorized request %s", err)
					http.Error(w, "Forbidden", http.StatusForbidden)
				
				} else if !token.Valid {	
					log.Print("invalid jwt token")
					http.Error(w, "Forbidden", http.StatusForbidden)
					
				} else if token.Valid {
					next.ServeHTTP(w, r)
				}
			}
		}
    })
}

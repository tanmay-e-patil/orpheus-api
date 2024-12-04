package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	///
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Middleware called")
		authorizationHeader := r.Header.Get("Authorization")
		accessToken, found := strings.CutPrefix(authorizationHeader, "Bearer ")
		if !found {
			log.Printf("Token not found")
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		customClaims := jwt.RegisteredClaims{}
		jwtParsed, err := jwt.ParseWithClaims(accessToken, &customClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			log.Printf("jwt parsing failed")
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if jwtParsed == nil {
			log.Printf("jwt parsing failed %v", jwtParsed)
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		idString, err := customClaims.GetSubject()

		if err != nil {
			log.Printf("Id not in token")
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		tokenExpiry, err := customClaims.GetExpirationTime()
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
		}
		if tokenExpiry.Before(time.Now().UTC()) {
			respondWithError(w, http.StatusUnauthorized, "Token is expired")
		}

		id, err := uuid.Parse(idString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		log.Printf("Checking if user exists")

		fetchedUser, err := cfg.DB.GetUserByID(r.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		handler(w, r, fetchedUser)
	}

}

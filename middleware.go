package main

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
	"net/http"
	"strings"
	"time"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	///
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		accessToken, found := strings.CutPrefix(authorizationHeader, "Bearer ")
		if !found {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		customClaims := jwt.RegisteredClaims{}
		jwtParsed, err := jwt.ParseWithClaims(accessToken, &customClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if jwtParsed == nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		idString, err := customClaims.GetSubject()

		if err != nil {
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

		fetchedUser, err := cfg.DB.GetUserByID(r.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		handler(w, r, fetchedUser)
	}

}

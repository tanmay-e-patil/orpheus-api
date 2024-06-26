package main

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

func (cfg *apiConfig) handlerSignUp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return

	}

	id, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jwtClaimsRefresh := jwt.RegisteredClaims{
		Issuer:    "orpheus",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(24) * time.Hour)),
		Subject:   id.String(),
	}

	jwtClaimsAccess := jwt.RegisteredClaims{
		Issuer:    "orpheus",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(24) * time.Minute)),
		Subject:   id.String(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaimsRefresh)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaimsAccess)

	refreshSignedString, err := refreshToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Error while signing jwt: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	accessSignedString, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Error while signing jwt: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           id,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Username:     params.Username,
		Email:        params.Email,
		Password:     string(passwordHash),
		RefreshToken: refreshSignedString,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user, accessSignedString))

}

func (cfg *apiConfig) handlerSignIn(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	defer r.Body.Close()

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	jwtClaimsRefresh := jwt.RegisteredClaims{
		Issuer:    "orpheus",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(24) * time.Hour)),
		Subject:   user.ID.String(),
	}

	jwtClaimsAccess := jwt.RegisteredClaims{
		Issuer:    "orpheus",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(24) * time.Minute)),
		Subject:   user.ID.String(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaimsRefresh)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaimsAccess)

	refreshSignedString, err := refreshToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Error while signing jwt: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	accessSignedString, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Error while signing jwt: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err = cfg.DB.UpdateUserWithRefreshToken(r.Context(), database.UpdateUserWithRefreshTokenParams{
		ID:           user.ID,
		RefreshToken: refreshSignedString,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user, accessSignedString))

}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	authorizationHeader := r.Header.Get("Authorization")
	refreshToken := authorizationHeader[7:]

	fetchedUser, err := cfg.DB.GetUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	customClaims := jwt.RegisteredClaims{}
	jwtParsed, err := jwt.ParseWithClaims(refreshToken, &customClaims, func(token *jwt.Token) (interface{}, error) {
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
	id, err := customClaims.GetSubject()

	log.Printf("%v", id)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshTokenExpiry := customClaims.ExpiresAt

	if time.Now().UTC().After(refreshTokenExpiry.Time) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is expired")
		return
	}

	jwtClaims := jwt.RegisteredClaims{
		Issuer:    "orpheus",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(24) * time.Minute)),
		Subject:   fetchedUser.ID.String(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedString, err := jwtToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Error while signing jwt: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		AccessToken string `json:"access_token"`
	}{AccessToken: signedString})
}

func (cfg *apiConfig) handlerMe(w http.ResponseWriter, r *http.Request, user database.User) {
	accessToken, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !found {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user, accessToken))
}

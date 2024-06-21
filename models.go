package main

import (
	"github.com/google/uuid"
	"github.com/tanmay-e-patil/orpheus-api/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

func databaseUserToUser(user database.User, accessToken string) User {
	return User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: user.RefreshToken,
	}
}

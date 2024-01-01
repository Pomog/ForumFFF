package dto

import (
	"forum-authentication/types"
)

type UserDTO struct {
	Id        int
	Username  string
	FirstName string
	LastName  string
	Email     string
}

// NewUserDTO creates a UserDTO from a User, excluding sensitive fields.
func NewUserDTO(user types.User) UserDTO {
	return UserDTO{
		Id:        user.Id,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

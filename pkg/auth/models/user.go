package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    Password  string `json:"password"`
    UserType  string `json:"user_type" gorm:"enum:passenger,driver"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(provided string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(provided))
	if err != nil {
		return err
	}
	return nil
}
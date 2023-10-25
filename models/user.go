package models

import (
	"html"
	"strings"

	"github.com/acanoe/newsbytes-api-go/utils/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;size:255" json:"username"`
	Password string `gorm:"not null;size:255" json:"password"`
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CheckLogin(username, password string) (string, error) {
	var err error
	var user User

	result := DB.Model(&User{}).Where("username = ?", username).Take(&user)
	if result.Error != nil {
		return "", err
	}

	err = verifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *User) SaveUser() (*User, error) {
	result := DB.Create(&u)
	if result.Error != nil {
		return &User{}, result.Error
	}
	return u, nil
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	// hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	// sanitize the user's username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}

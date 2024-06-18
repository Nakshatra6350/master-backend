package data

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type UserCredentials struct {
	Email     string    `json:"email" bson:"email" validate:"required"`
	Password  string    `json:"password" bson:"password" validate:"required"`
	Name      string    `json:"name" bson:"name" validate:"required"`
	Username  string    `json:"username" bson:"username" validate:"required"`
	CreatedAt time.Time `json:"-" bson:"created_at"`
}

func (u *UserCredentials) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type LoginUser struct {
	Email    string `json:"email" bson:"email" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

func (u *LoginUser) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

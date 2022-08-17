package models

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson: "name" json: "name"`
	Email            string `bson:"email" json:"email"`
	Password         string `bson:"password" json:"password"`
}

func NewUser(name string, email string, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

type RegisterResponse struct {
	Name     string `bson: "name" json: "name" validate:"required"`
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required,min=3"`
}

func (register *RegisterResponse) ValidateResponse() error {
	validator := validator.New()
	err := validator.Struct(register)
	if err != nil {
		return err
	}

	user := &User{}

	col := mgm.Coll(user)

	_ = col.First(bson.M{"email": register.Email}, user)
	if user.Name != "" {
		return errors.New("User with this email already exist")
	}
	return nil

}

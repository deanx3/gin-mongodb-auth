package models

import (
	"errors"
	"math/rand"
	"time"

	"github.com/deanx3/gin-mongodb-auth/helpers"
	"github.com/go-playground/validator/v10"
	"github.com/kamva/mgm/v3"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson: "name" json: "name"`
	Email            string `bson:"email" json:"email"`
	Password         string `bson:"password" json:"password"`
	Token            string `bson:"token" json:"token"`
}

func NewUser(name string, email string, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

func (user *User) GenerateToken() error {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	res := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	user.Token = res.String()
	err := user.UpdateUser()
	if err != nil {
		return err
	}
	return nil
}

func (user *User) UpdateUser() error {
	err := mgm.Coll(user).Update(user)
	if err != nil {
		return err
	}
	return nil
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

type LoginResponse struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required,min=3"`
}

func (login *LoginResponse) ValidateResponse() error {
	validator := validator.New()
	err := validator.Struct(login)
	if err != nil {
		return err
	}
	user := &User{}

	col := mgm.Coll(user)
	_ = col.First(bson.M{"email": login.Email}, user)
	if user.Email == "" {
		return errors.New("no user found againts this email")
	}
	err = helpers.VerifyPassword(user.Password, login.Password)
	if err != nil {
		return err
	}

	return nil
}

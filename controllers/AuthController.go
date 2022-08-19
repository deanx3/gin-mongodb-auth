package controllers

import (
	"net/http"

	"github.com/deanx3/gin-mongodb-auth/helpers"
	"github.com/deanx3/gin-mongodb-auth/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (auth *AuthController) login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	login := models.LoginResponse{Email: email, Password: password}

	err := login.ValidateResponse()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	user := &models.User{}

	_ = mgm.Coll(user).First(bson.M{"email": login.Email}, user)

	err = user.GenerateToken()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "user": user})
}

func (auth *AuthController) register(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	name := ctx.PostForm("name")
	requestRegister := models.RegisterResponse{Name: name, Email: email, Password: password}

	// validate request
	err := requestRegister.ValidateResponse()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	hash, err := helpers.HashPassword(requestRegister.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	// err, hash = helpers.hashedPassword(requestRegister.Password)
	user := models.NewUser(requestRegister.Name, requestRegister.Email, hash)
	err = mgm.Coll(user).Create(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "user": user})
}

func (auth *AuthController) AuthRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.GET("/login", auth.login)
	router.POST("/register", auth.register)
}

package controllers

import (
	"net/http"

	"github.com/deanx3/gin-mongodb-auth/middlewares"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (home *HomeController) StatusCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}

func (home *HomeController) HomeRoutes(r *gin.RouterGroup) {
	router := r.Group("healthCheck")
	router.Use(middlewares.Auth())
	router.GET("/", home.StatusCheck)
}
